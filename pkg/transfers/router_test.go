// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package transfers

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/moov-io/base"
	moovcustomers "github.com/moov-io/customers/client"
	"github.com/moov-io/paygate/pkg/client"
	"github.com/moov-io/paygate/pkg/customers"
	"github.com/moov-io/paygate/pkg/customers/accounts"
	"github.com/moov-io/paygate/pkg/tenants"
	"github.com/moov-io/paygate/pkg/testclient"
	"github.com/moov-io/paygate/pkg/transfers/fundflow"
	"github.com/moov-io/paygate/pkg/transfers/pipeline"
	"github.com/moov-io/paygate/pkg/util"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
)

var (
	repoWithTransfer = &MockRepository{
		Transfers: []*client.Transfer{
			{
				TransferID: base.ID(),
				Amount:     "USD 12.44",
				Source: client.Source{
					CustomerID: base.ID(),
					AccountID:  base.ID(),
				},
				Destination: client.Destination{
					CustomerID: base.ID(),
					AccountID:  base.ID(),
				},
				Description: "test transfer",
				Status:      client.PENDING,
				Created:     time.Now(),
			},
		},
	}

	tenantRepo = &tenants.MockRepository{}

	fakePublisher = &pipeline.MockPublisher{}

	mockStrategy = &fundflow.MockStrategy{}

	mockDecryptor = &accounts.MockDecryptor{Number: "12345"}
)

func mockCustomersClient() *customers.MockClient {
	return &customers.MockClient{
		Account: &moovcustomers.Account{
			AccountID:           base.ID(),
			MaskedAccountNumber: "****34",
			RoutingNumber:       "987654320",
			Status:              "validated",
			Type:                moovcustomers.CHECKING,
		},
		Customer: &moovcustomers.Customer{
			CustomerID: base.ID(),
			FirstName:  "John",
			LastName:   "Doe",
			Email:      "john.doe@example.com",
		},
	}
}

func TestTransfers__readTransferFilterParams(t *testing.T) {
	u, _ := url.Parse("http://localhost:8082/transfers?startDate=2020-04-06&limit=10&status=failed")
	req := &http.Request{URL: u}
	params := readTransferFilterParams(req)

	if params.StartDate.Format(util.YYMMDDTimeFormat) != "2020-04-06" {
		t.Errorf("unexpected StartDate: %v", params.StartDate)
	}
	if !params.EndDate.After(time.Now()) {
		t.Errorf("unexpected EndDate: %v", params.EndDate)
	}
	if params.Status != client.FAILED {
		t.Errorf("expected status: %q", params.Status)
	}
	if params.Limit != 10 {
		t.Errorf("unexpected limit: %d", params.Limit)
	}
	if params.Offset != 0 {
		t.Errorf("unexpected offset: %d", params.Offset)
	}
}

func TestRouter__getUserTransfers(t *testing.T) {
	customersClient := mockCustomersClient()

	r := mux.NewRouter()
	router := NewRouter(log.NewNopLogger(), repoWithTransfer, tenantRepo, customersClient, mockDecryptor, mockStrategy, fakePublisher)
	router.RegisterRoutes(r)

	c := testclient.New(t, r)

	xfers, resp, err := c.TransfersApi.GetTransfers(context.TODO(), "userID", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if n := len(xfers); n != 1 {
		t.Errorf("got %d transfers: %#v", n, xfers)
	}
}

func TestRouter__createUserTransfer(t *testing.T) {
	customersClient := mockCustomersClient()

	r := mux.NewRouter()
	router := NewRouter(log.NewNopLogger(), repoWithTransfer, tenantRepo, customersClient, mockDecryptor, mockStrategy, fakePublisher)
	router.RegisterRoutes(r)

	c := testclient.New(t, r)

	opts := client.CreateTransfer{
		Amount: "USD 12.44",
		Source: client.Source{
			CustomerID: base.ID(),
			AccountID:  base.ID(),
		},
		Destination: client.Destination{
			CustomerID: base.ID(),
			AccountID:  base.ID(),
		},
		Description: "test transfer",
		SameDay:     true,
	}
	xfer, resp, err := c.TransfersApi.AddTransfer(context.TODO(), "userID", opts, nil)
	if err != nil {
		bs, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("error=%v \n body=%s", err, string(bs))
	}
	defer resp.Body.Close()

	if xfer.TransferID == "" {
		t.Errorf("missing Transfer=%#v", xfer)
	}
}

func TestRouter__createUserTransfersInvalidAmount(t *testing.T) {
	customersClient := mockCustomersClient()

	r := mux.NewRouter()
	router := NewRouter(log.NewNopLogger(), repoWithTransfer, tenantRepo, customersClient, mockDecryptor, mockStrategy, fakePublisher)
	router.RegisterRoutes(r)

	c := testclient.New(t, r)

	opts := client.CreateTransfer{
		Amount: "USD YY.44",
	}
	xfer, resp, err := c.TransfersApi.AddTransfer(context.TODO(), "userID", opts, nil)
	if err == nil {
		t.Error("expected error")
	}
	defer resp.Body.Close()

	if xfer.TransferID != "" {
		t.Errorf("unexpected transfer: %#v", xfer)
	}
}

func TestRouter__MissingSource(t *testing.T) {
	customersClient := mockCustomersClient()

	r := mux.NewRouter()
	router := NewRouter(log.NewNopLogger(), repoWithTransfer, tenantRepo, customersClient, mockDecryptor, mockStrategy, fakePublisher)
	router.RegisterRoutes(r)

	c := testclient.New(t, r)

	opts := client.CreateTransfer{
		Amount: "USD 12.54",
		Source: client.Source{
			AccountID: base.ID(), // missing CustomerID
		},
	}
	xfer, resp, err := c.TransfersApi.AddTransfer(context.TODO(), "userID", opts, nil)
	if err == nil {
		t.Error("expected error")
	}
	defer resp.Body.Close()

	if xfer.TransferID != "" {
		t.Errorf("unexpected transfer: %#v", xfer)
	}
}

func TestRouter__MissingDestination(t *testing.T) {
	customersClient := mockCustomersClient()

	r := mux.NewRouter()
	router := NewRouter(log.NewNopLogger(), repoWithTransfer, tenantRepo, customersClient, mockDecryptor, mockStrategy, fakePublisher)
	router.RegisterRoutes(r)

	c := testclient.New(t, r)

	opts := client.CreateTransfer{
		Amount: "USD 12.54",
		Source: client.Source{
			CustomerID: base.ID(),
			AccountID:  base.ID(),
		},
		Destination: client.Destination{
			CustomerID: base.ID(), // missing AccountID
		},
	}
	xfer, resp, err := c.TransfersApi.AddTransfer(context.TODO(), "userID", opts, nil)
	if err == nil {
		t.Error("expected error")
	}
	defer resp.Body.Close()

	if xfer.TransferID != "" {
		t.Errorf("unexpected transfer: %#v", xfer)
	}
}

func TestRouter__validateAmount(t *testing.T) {
	if err := validateAmount(""); err == nil {
		t.Error("expected error")
	}
}

func TestRouter__getUserTransfer(t *testing.T) {
	customersClient := mockCustomersClient()

	r := mux.NewRouter()
	router := NewRouter(log.NewNopLogger(), repoWithTransfer, tenantRepo, customersClient, mockDecryptor, mockStrategy, fakePublisher)
	router.RegisterRoutes(r)

	c := testclient.New(t, r)

	xfer, resp, err := c.TransfersApi.GetTransferByID(context.TODO(), "transferID", "userID", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if xfer.TransferID == "" {
		t.Errorf("missing Transfer=%#v", xfer)
	}
}

func TestRouter__deleteUserTransfer(t *testing.T) {
	customersClient := mockCustomersClient()

	r := mux.NewRouter()
	router := NewRouter(log.NewNopLogger(), repoWithTransfer, tenantRepo, customersClient, mockDecryptor, mockStrategy, fakePublisher)
	router.RegisterRoutes(r)

	c := testclient.New(t, r)

	resp, err := c.TransfersApi.DeleteTransferByID(context.TODO(), "transferID", "userID", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
}
