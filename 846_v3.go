package easi

import (
	// "fmt"
	"bytes"
	"context"
	"encoding/csv"
	"io"
	"strconv"

	"github.com/jszwec/csvutil"
)

type Standard846V3 struct {
	EnvelopeHeaderV3  EnvelopeHeaderV3
	Header            Standard846V3TransactionHeader
	LineItems         []Standard846V3LineItem
	Trailer           Standard846V3TransactionTrailer
	Sections          []Standard846V3Section
	EnvelopeTrailerV3 EnvelopeTrailerV3
}

type Standard846V3Section struct {
	Header    Standard846V3TransactionHeader
	LineItems []Standard846V3LineItem
	Trailer   Standard846V3TransactionTrailer
}

type Standard846V3TransactionHeader struct {
	Header                  string
	TransactionType         string
	TransactionSetPurpose   string
	VersionNumber           string
	VendorID                string
	AsOfDate                string
	AsOfTime                string
	TimeZone                string
	ElapsedTimeToNextUpdate string
	DistributionCenter      string
	DistributionCenterID    string
}

type Standard846V3TransactionTrailer struct {
	TrailerRecord    string
	FileCreationDate string
	FileCreationTime string
	RecordCount      int
}

type Standard846V3LineItem struct {
	DetailSectionLoopA                    string
	LineItemNumber                        int
	ItemIdentificationGTIN                string
	CurrentInventoryLevel                 int
	UnitOfMeasure                         string
	QuantityToArriveWithinTheNextTwoWeeks string
	PurchaseUnitPriceEaches               string
	PurchaseUnitPriceDozens               string
	PurchaseUnitPriceCases                string
	CustomPriceUOMDescription             string
	PurchaseUnitPriceCustom               string
}

func (s *Standard846V3) Prep(ctx context.Context) error {

	return nil
}

func (s *Standard846V3) ToBytes(ctx context.Context) (*[]byte, error) {

	return nil, nil
}

func (s *Standard846V3) FromBytes(ctx context.Context, req []byte) error {

	r := csv.NewReader(bytes.NewReader(req))
	r.Comma = '\t'

	// Headerless
	blankHeader, errHeader := csvutil.Header(Header{}, "csv")
	if errHeader != nil {
		return errHeader
	}

	// Decoder
	dec, errDecoder := csvutil.NewDecoder(r, blankHeader...)
	if errDecoder != nil {
		return errDecoder
	}

	// Section
	var section Standard846V3Section

	for {
		var v struct{}
		if err := dec.Decode(&v); err == io.EOF {
			break
		}

		// Record
		var lineType string
		record := dec.Record()
		if len(record) > 0 {
			lineType = record[0]
		}

		// Build
		switch lineType {
		case "EASI":
			var x EnvelopeHeaderV3
			err := x.FromSlice(ctx, record)
			if err != nil {
				return err
			}
			s.EnvelopeHeaderV3 = x
		case "01":
			var x Standard846V3TransactionHeader
			err := x.FromSlice(ctx, record)
			if err != nil {
				return err
			}
			s.Header = x
			section.Header = x
		case "02":
			var x Standard846V3LineItem
			err := x.FromSlice(ctx, record)
			if err != nil {
				return err
			}
			s.LineItems = append(s.LineItems, x)
			section.LineItems = append(section.LineItems, x)
		case "09":
			var x Standard846V3TransactionTrailer
			err := x.FromSlice(ctx, record)
			if err != nil {
				return err
			}
			s.Trailer = x
			section.Trailer = x

			// Append Section
			s.Sections = append(s.Sections, section)

			// Reset Section
			section = Standard846V3Section{}

		case "EASX":
			var x EnvelopeTrailerV3
			err := x.FromSlice(ctx, record)
			if err != nil {
				return err
			}
			s.EnvelopeTrailerV3 = x
		default:

		}

	}

	return nil
}

func (s *Standard846V3TransactionHeader) FromSlice(ctx context.Context, req []string) error {

	if len(req) > 0 {
		s.Header = req[0]
	}

	if len(req) > 1 {
		s.TransactionType = req[1]
	}

	if len(req) > 2 {
		s.TransactionSetPurpose = req[2]
	}

	if len(req) > 3 {
		s.VersionNumber = req[3]
	}

	if len(req) > 4 {
		s.VendorID = req[4]
	}

	if len(req) > 5 {
		s.AsOfDate = req[5]
	}

	if len(req) > 6 {
		s.AsOfTime = req[6]
	}

	if len(req) > 7 {
		s.TimeZone = req[7]
	}

	if len(req) > 8 {
		s.ElapsedTimeToNextUpdate = req[8]
	}

	if len(req) > 9 {
		s.DistributionCenter = req[9]
	}

	if len(req) > 10 {
		s.DistributionCenterID = req[10]
	}

	return nil
}

func (s *Standard846V3TransactionTrailer) FromSlice(ctx context.Context, req []string) error {

	if len(req) > 0 {
		s.TrailerRecord = req[0]
	}

	if len(req) > 1 {
		s.FileCreationDate = req[1]
	}

	if len(req) > 2 {
		s.FileCreationTime = req[2]
	}

	if len(req) > 3 {
		if req[3] != "" {
			recordCount, err := strconv.Atoi(req[3])
			if err != nil {
				return err
			}
			s.RecordCount = recordCount
		}
	}

	return nil
}

func (s *Standard846V3LineItem) FromSlice(ctx context.Context, req []string) error {

	if len(req) > 0 {
		s.DetailSectionLoopA = req[0]
	}

	if len(req) > 1 {
		if req[1] != "" {
			lineItemNumber, err := strconv.Atoi(req[1])
			if err != nil {
				return err
			}
			s.LineItemNumber = lineItemNumber
		}
	}

	if len(req) > 2 {
		s.ItemIdentificationGTIN = req[2]
	}

	if len(req) > 3 {
		if req[3] != "" {
			currentInventoryLevel, err := strconv.Atoi(req[3])
			if err != nil {
				return err
			}
			s.CurrentInventoryLevel = currentInventoryLevel
		}
	}

	return nil
}
