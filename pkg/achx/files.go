// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package achx

import (
	"fmt"
	"time"

	"github.com/moov-io/ach"
	customers "github.com/moov-io/customers/client"
	"github.com/moov-io/paygate/pkg/client"
	"github.com/moov-io/paygate/pkg/config"
)

type Source struct {
	Customer customers.Customer
	Account  customers.Account
}

type Destination struct {
	Customer customers.Customer
	Account  customers.Account

	// AccountNumber contains the decrypted account number from the customers service
	AccountNumber string
}

func ConstrctFile(id string, odfi config.ODFI, companyID string, xfer *client.Transfer, source Source, destination Destination) (*ach.File, error) {
	file, now := ach.NewFile(), time.Now()
	file.ID = id
	file.Control = ach.NewFileControl()

	// File Header
	file.Header.ID = id
	file.Header.ImmediateOrigin = odfi.Gateway.Origin
	file.Header.ImmediateOriginName = odfi.Gateway.OriginName
	file.Header.ImmediateDestination = odfi.Gateway.Destination
	file.Header.ImmediateDestinationName = odfi.Gateway.DestinationName
	file.Header.FileCreationDate = now.Format("060102") // YYMMDD
	file.Header.FileCreationTime = now.Format("1504")   // HHMM

	// Right now we only support creating PPD files
	batch, err := createPPDBatch(id, odfi, companyID, xfer, source, destination)
	if err != nil {
		return nil, fmt.Errorf("constructACHFile: PPD: %v", err)
	}
	file.AddBatch(batch)

	return file, file.Validate()
}
