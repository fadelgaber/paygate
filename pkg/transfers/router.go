// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package transfers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/moov-io/base"
	"github.com/moov-io/paygate/pkg/client"
	"github.com/moov-io/paygate/pkg/customers"
	"github.com/moov-io/paygate/pkg/customers/accounts"
	"github.com/moov-io/paygate/pkg/model"
	"github.com/moov-io/paygate/pkg/tenants"
	"github.com/moov-io/paygate/pkg/transfers/fundflow"
	"github.com/moov-io/paygate/pkg/transfers/pipeline"
	"github.com/moov-io/paygate/pkg/util"
	"github.com/moov-io/paygate/x/route"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
)

type Router struct {
	Logger log.Logger
	Repo   Repository

	Publisher pipeline.XferPublisher

	GetUserTransfers   http.HandlerFunc
	CreateUserTransfer http.HandlerFunc
	GetUserTransfer    http.HandlerFunc
	DeleteUserTransfer http.HandlerFunc
}

func NewRouter(
	logger log.Logger,
	repo Repository,
	tenantRepo tenants.Repository,
	customersClient customers.Client,
	accountDecryptor accounts.Decryptor,
	fundStrategy fundflow.Strategy,
	pub pipeline.XferPublisher,
) *Router {
	return &Router{
		Logger:             logger,
		Repo:               repo,
		Publisher:          pub,
		GetUserTransfers:   GetUserTransfers(logger, repo),
		CreateUserTransfer: CreateUserTransfer(logger, repo, tenantRepo, customersClient, accountDecryptor, fundStrategy, pub),
		GetUserTransfer:    GetUserTransfer(logger, repo),
		DeleteUserTransfer: DeleteUserTransfer(logger, repo, pub),
	}
}

func (c *Router) RegisterRoutes(r *mux.Router) {
	r.Methods("GET").Path("/transfers").HandlerFunc(c.GetUserTransfers)
	r.Methods("POST").Path("/transfers").HandlerFunc(c.CreateUserTransfer)
	r.Methods("GET").Path("/transfers/{transferID}").HandlerFunc(c.GetUserTransfer)
	r.Methods("DELETE").Path("/transfers/{transferID}").HandlerFunc(c.DeleteUserTransfer)
}

func getTransferID(r *http.Request) string {
	return route.ReadPathID("transferID", r)
}

type transferFilterParams struct {
	Status    client.TransferStatus
	StartDate time.Time
	EndDate   time.Time
	Limit     int64
	Offset    int64
}

func readTransferFilterParams(r *http.Request) transferFilterParams {
	params := transferFilterParams{
		StartDate: time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Now().Add(24 * time.Hour),
		Limit:     100,
		Offset:    0,
	}
	if r == nil {
		return params
	}
	if r.URL != nil {
		q := r.URL.Query()
		if v := q.Get("startDate"); v != "" {
			params.StartDate = util.FirstParsedTime(v, base.ISO8601Format, util.YYMMDDTimeFormat)
		}
		if v := q.Get("endDate"); v != "" {
			params.EndDate, _ = time.Parse(base.ISO8601Format, v)
			fmt.Printf("params.EndDate=%v\n", params.EndDate)
		}
		if s := strings.TrimSpace(q.Get("status")); s != "" {
			params.Status = client.TransferStatus(s)
		}
	}
	if limit := route.ReadLimit(r); limit != 0 {
		params.Limit = limit
	}
	if offset := route.ReadOffset(r); offset != 0 {
		params.Offset = offset
	}
	return params
}

func GetUserTransfers(logger log.Logger, repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responder := route.NewResponder(logger, w, r)

		params := readTransferFilterParams(r)
		xfers, err := repo.getUserTransfers(responder.XUserID, params)
		if err != nil {
			responder.Problem(err)
			return
		}

		responder.Respond(func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(xfers)
		})
	}
}

func CreateUserTransfer(
	logger log.Logger,
	repo Repository,
	tenantRepo tenants.Repository,
	customersClient customers.Client,
	accountDecryptor accounts.Decryptor,
	fundStrategy fundflow.Strategy,
	pub pipeline.XferPublisher,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responder := route.NewResponder(logger, w, r)

		var req client.CreateTransfer
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			responder.Problem(err)
			return
		}
		if err := validateTransferRequest(req); err != nil {
			responder.Problem(err)
			return
		}

		transfer := &client.Transfer{
			TransferID:  base.ID(),
			Amount:      req.Amount,
			Source:      req.Source,
			Destination: req.Destination,
			Description: req.Description,
			Status:      client.PENDING,
			SameDay:     req.SameDay,
			Created:     time.Now(),
		}

		// TODO(adam): validate both Customer and get get their Accounts
		// get/decrypt destination account number
		//
		// TODO(adam): future: limits checks

		// Save our Transfer to the database
		if err := repo.writeUserTransfers(responder.XUserID, transfer); err != nil {
			responder.Problem(err)
			return
		}

		// According to our strategy create (originate) ACH files to be published somewhere
		if fundStrategy != nil {
			companyID, err := tenantRepo.GetCompanyIdentification("tenantID") // TODO(adam): need to get from auth
			if err != nil {
				responder.Problem(err)
				return
			}
			source, err := getFundflowSource(customersClient, req.Source)
			if err != nil {
				fmt.Printf("error getting source: %v\n", err)
				responder.Problem(err)
				return
			}
			destination, err := getFundflowDestination(customersClient, accountDecryptor, req.Destination)
			if err != nil {
				fmt.Printf("error getting destination: %v\n", err)
				responder.Problem(err)
				return
			}
			files, err := fundStrategy.Originate(companyID, transfer, source, destination)
			if err != nil {
				fmt.Printf("error originating ACH files: %v\n", err)
				responder.Problem(err)
				return
			}
			if err := pipeline.PublishFiles(pub, transfer, files); err != nil {
				fmt.Printf("error publishing ACH files: %v\n", err)
				responder.Problem(err)
				return
			}
		}

		responder.Respond(func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(transfer)
		})
	}
}

func validateTransferRequest(req client.CreateTransfer) error {
	if req.Source.CustomerID == "" || req.Source.AccountID == "" {
		return errors.New("incomplete source")
	}
	if req.Destination.CustomerID == "" || req.Destination.AccountID == "" {
		return errors.New("incomplete destination")
	}
	if err := validateAmount(req.Amount); err != nil {
		return err
	}
	if req.Description == "" {
		return errors.New("missing description")
	}

	return nil
}

func validateAmount(amount string) error {
	if amount == "" {
		return errors.New("missing amount")
	}
	var amt model.Amount
	if err := amt.FromString(amount); err != nil {
		return fmt.Errorf("unable to parse '%s': %v", amount, err)
	}
	return nil
}

func getFundflowSource(client customers.Client, src client.Source) (fundflow.Source, error) {
	var source fundflow.Source

	// Set source Customer
	cust, err := client.Lookup(src.CustomerID, "requestID", "userID")
	if err != nil {
		return source, err
	}
	if cust == nil {
		return source, fmt.Errorf("customerID=%s is not found", src.CustomerID)
	}
	source.Customer = *cust

	// Get customer Account
	if acct, err := client.FindAccount(src.CustomerID, src.AccountID); acct == nil || err != nil {
		return source, fmt.Errorf("accountID=%s not found for customerID=%s error=%v", src.AccountID, src.CustomerID, err)
	} else {
		source.Account = *acct
	}

	return source, nil
}

func getFundflowDestination(client customers.Client, accountDecryptor accounts.Decryptor, dst client.Destination) (fundflow.Destination, error) {
	var destination fundflow.Destination

	// Set destination Customer
	cust, err := client.Lookup(dst.CustomerID, "requestID", "userID")
	if err != nil {
		return destination, err
	}
	if cust == nil {
		return destination, fmt.Errorf("customerID=%s is not found", dst.CustomerID)
	}
	destination.Customer = *cust

	// Get customer Account
	if acct, err := client.FindAccount(dst.CustomerID, dst.AccountID); acct == nil || err != nil {
		return destination, fmt.Errorf("accountID=%s not found for customerID=%s error=%v", dst.AccountID, dst.CustomerID, err)
	} else {
		destination.Account = *acct
	}

	if num, err := accountDecryptor.AccountNumber(dst.CustomerID, dst.AccountID); num == "" || err != nil {
		return destination, fmt.Errorf("unable to decrypt accountID=%s for customerID=%s error=%v", dst.AccountID, dst.CustomerID, err)
	} else {
		destination.AccountNumber = num
	}

	return destination, nil
}

func GetUserTransfer(logger log.Logger, repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responder := route.NewResponder(logger, w, r)

		xfer, err := repo.GetTransfer(getTransferID(r))
		if err != nil {
			responder.Problem(err)
			return
		}

		responder.Respond(func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(xfer)
		})
	}
}

func DeleteUserTransfer(logger log.Logger, repo Repository, pub pipeline.XferPublisher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responder := route.NewResponder(logger, w, r)

		transferID := getTransferID(r)
		if err := repo.deleteUserTransfer(responder.XUserID, transferID); err != nil {
			responder.Problem(err)
			return
		}

		if pub != nil {
			msg := pipeline.CanceledTransfer{
				TransferID: transferID,
			}
			if err := pub.Cancel(msg); err != nil {
				responder.Problem(err)
				return
			}
		}

		responder.Respond(func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusOK)
		})
	}
}
