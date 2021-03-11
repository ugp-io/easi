package easi

import(
	// "fmt"
	"context"
	"io"
	"time"
	"bytes"
	"strconv"
	"encoding/csv"
	"github.com/jszwec/csvutil"
)

type Standard856V5 struct {
	EnvelopeHeaderV2 EnvelopeHeaderV2
	Transactions []Standard856V5Transaction
	EnvelopeTrailerV2 EnvelopeTrailerV2
}

type Standard856V5TransactionHeader struct {
	Header string
	TransactionType string
	TransactionSetPurpose string
	VersionNumber string
	ShipmentNumber string
	ASNDate string
	ASNTime string
	VendorID string
	PurchaserAccountID string
	StoreID string
	DeliverToCompanyName string
	DeliverToAddress1 string
	DeliverToAddress2 string
	DeliverToCityName string
	DeliverToStateCode string
	DeliverToPostalCode string
	DeliverToCountryCode string
	BOLNumber string
	CarrierRoutingDetails string
	TrailerID string
	CarrierTrackingNumber string
	ShipmentDate string
	
}

type Standard856V5TransactionTrailer struct {
	TrailerRecord string
	TotalCaseCount int
	TotalQtyShipped int
	TotalGrossWeight int
	RecordCount int
	TotalPalletCount int
}

type Standard856V5Transaction struct {
	Header Standard856V5TransactionHeader
	Pallets []Standard856V5Pallet `csv:"-"`
	Trailer Standard856V5TransactionTrailer
	
}

type Standard856V5Pallet struct {
	PalletRecord string
	PalletID string
	LineItems []Standard856V5LineItem `csv:"-"`
}

type Standard856V5LineItem struct {
	DetailSectionLoopB string
	LineItemNumber int
	ManufacturersSerialCaseNumber string
	BuyersPurchaseOrderNumber string
	ItemIdentificationGTIN string
	MasterStyle string
	DetailStyle string
	ColorCode string
	SizeCode string
	RevisionCode string
	UnitOrBasisForMeasurementCode string
	QuantityShipped int
	CountryOfOrigin string
	ManufacturersOrderNumber string
	ManufacturersLotID string
}



func (s *Standard856V5) Prep(ctx context.Context) (error){

	// Header
	errHeader := s.EnvelopeHeaderV2.Prep(ctx)
	if errHeader != nil {
		return errHeader
	}
	s.EnvelopeHeaderV2.TransactionType = "856"

	// Transactions
	for transactionKey, transaction := range s.Transactions {

		s.Transactions[transactionKey].Header.Header = "01"
		s.Transactions[transactionKey].Header.TransactionType = "856"
		s.Transactions[transactionKey].Header.TransactionSetPurpose = "00"
		
		s.Transactions[transactionKey].Header.VersionNumber = "7.0"
		s.Transactions[transactionKey].Header.ASNDate = time.Now().Format("20060102")
		s.Transactions[transactionKey].Header.ASNTime = time.Now().Format("150405")
		s.Transactions[transactionKey].Header.ShipmentDate = time.Now().Format("20060102")

		// Pallets
		for palletKey, pallet := range transaction.Pallets {
			s.Transactions[transactionKey].Pallets[palletKey].PalletRecord = "05"

			// Line Items
			for palletLineItemKey, _ := range pallet.LineItems {
				s.Transactions[transactionKey].Pallets[palletKey].LineItems[palletLineItemKey].DetailSectionLoopB = "02"
			}

		}

		s.Transactions[transactionKey].Trailer.TrailerRecord = "09"
	}
	

	// Trailer
	errTrailer := s.EnvelopeTrailerV2.Prep(ctx)
	if errTrailer != nil {
		return errTrailer
	}

	return nil
}

func (s *Standard856V5) ToBytes(ctx context.Context) (*[]byte, error){

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
    w.Comma = '\t'
	enc := csvutil.NewEncoder(w)

	// Prep
	errPrep := s.Prep(ctx)
	if errPrep != nil {
		return nil, errPrep
	}
	
	// Headerless
	type header struct{}
	errHeader := enc.EncodeHeader(header{})
	if errHeader != nil {
		return nil, errHeader
	}

	// Envelope Header
	errEnvelopeHeaderV2 := enc.Encode(s.EnvelopeHeaderV2)
	if errEnvelopeHeaderV2 != nil {
		return nil, errEnvelopeHeaderV2
	}

	// Transactions
	for _, transaction := range s.Transactions {

		// Header
		errTransaction := enc.Encode(transaction.Header)
		if errTransaction != nil {
			return nil, errTransaction
		}

		// Pallets
		for _, pallet := range transaction.Pallets {
			errPallet := enc.Encode(pallet)
			if errPallet != nil {
				return nil, errPallet
			}

			// Line Items
			for _, lineItem := range pallet.LineItems {
				errLineItem := enc.Encode(lineItem)
				if errLineItem != nil {
					return nil, errLineItem
				}
			}
		}

		// Trailer
		errTrailer := enc.Encode(transaction.Trailer)
		if errTrailer != nil {
			return nil, errTrailer
		}

	}
	
	// Envelope Trailer
	errEnvelopeTrailerV2 := enc.Encode(s.EnvelopeTrailerV2)
	if errEnvelopeTrailerV2 != nil {
		return nil, errEnvelopeTrailerV2
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}

	byteArray := buf.Bytes()

	return &byteArray, nil
}

func (s *Standard856V5) FromBytes(ctx context.Context, req []byte) (error){

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
	
	var transactionCount, palletCount int
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
			var x EnvelopeHeaderV2
			err := x.FromSlice(ctx, record)
			if err != nil {
				return err
			}
			s.EnvelopeHeaderV2 = x
		case "01":
			var x Standard856V5TransactionHeader
			err := x.FromSlice(ctx, record)
			if err != nil {
				return err
			}
			s.Transactions = append(s.Transactions, Standard856V5Transaction{})
			palletCount = 0
			transactionCount++
			s.Transactions[transactionCount - 1].Header = x
		case "05":
			var x Standard856V5Pallet
			err := x.FromSlice(ctx, record)
			if err != nil {
				return err
			}
			s.Transactions[transactionCount - 1].Pallets = append(s.Transactions[transactionCount - 1].Pallets, x)
			palletCount++
		case "02":
			var x Standard856V5LineItem
			err := x.FromSlice(ctx, record)
			if err != nil {
				return err
			}
			if palletCount <= 0 {
				s.Transactions[transactionCount - 1].Pallets = append(s.Transactions[transactionCount - 1].Pallets, Standard856V5Pallet{})
				palletCount++
			}
			s.Transactions[transactionCount - 1].Pallets[palletCount - 1].LineItems = append(s.Transactions[transactionCount - 1].Pallets[palletCount - 1].LineItems, x)
		case "09":
			var x Standard856V5TransactionTrailer
			err := x.FromSlice(ctx, record)
			if err != nil {
				return err
			}
			s.Transactions[transactionCount - 1].Trailer = x
		case "EASX":
			var x EnvelopeTrailerV2
			err := x.FromSlice(ctx, record)
			if err != nil {
				return err
			}
			s.EnvelopeTrailerV2 = x
		default:
			
		}

	}

	return nil
}

func (s *Standard856V5TransactionHeader) FromSlice(ctx context.Context, req []string) (error){

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
		s.ShipmentNumber = req[4]
	}
	if len(req) > 5 {
		s.ASNDate = req[5]
	}
	if len(req) > 6 {
		s.ASNTime = req[6]
	}
	if len(req) > 7 {
		s.VendorID = req[7]
	}
	if len(req) > 8 {
		s.PurchaserAccountID = req[8]
	}
	if len(req) > 9 {
		s.StoreID = req[9]
	}
	if len(req) > 10 {
		s.DeliverToCompanyName = req[10]
	}
	if len(req) > 11 {
		s.DeliverToAddress1 = req[11]
	}
	if len(req) > 12 {
		s.DeliverToAddress2 = req[12]
	}
	if len(req) > 13 {
		s.DeliverToCityName = req[13]
	}
	if len(req) > 14 {
		s.DeliverToStateCode = req[14]
	}
	if len(req) > 15 {
		s.DeliverToPostalCode = req[15]
	}
	if len(req) > 16 {
		s.DeliverToCountryCode = req[16]
	}
	if len(req) > 17 {
		s.BOLNumber = req[17]
	}
	if len(req) > 18 {
		s.CarrierRoutingDetails = req[18]
	}
	if len(req) > 19 {
		s.TrailerID = req[19]
	}
	if len(req) > 20 {
		s.CarrierTrackingNumber = req[20]
	}
	if len(req) > 21 {
		s.ShipmentDate = req[21]
	}

	return nil
}

func (s *Standard856V5TransactionTrailer) FromSlice(ctx context.Context, req []string) (error){

	if len(req) > 0 {
		s.TrailerRecord = req[0]
	}

	if len(req) > 1 {
		if req[1] != "" {
			totalCaseCount, err := strconv.Atoi(req[1])
			if err != nil {
				return err
			}
			s.TotalCaseCount = totalCaseCount
		}
	}
	if len(req) > 2 {
		if req[2] != "" {
			totalQtyShipped, err := strconv.Atoi(req[2])
			if err != nil {
				return err
			}
			s.TotalQtyShipped = totalQtyShipped
		}
	}
	if len(req) > 3 {
		if req[3] != "" {
			totalGrossWeight, err := strconv.Atoi(req[3])
			if err != nil {
				return err
			}
			s.TotalGrossWeight = totalGrossWeight
		}
	}
	if len(req) > 4 {
		if req[4] != "" {
			recordCount, err := strconv.Atoi(req[4])
			if err != nil {
				return err
			}
			s.RecordCount = recordCount
		}
	}
	if len(req) > 5 {
		if req[5] != "" {
			totalPalletCount, err := strconv.Atoi(req[5])
			if err != nil {
				return err
			}
			s.TotalPalletCount = totalPalletCount
		}
	}

	return nil
}

func (s *Standard856V5Pallet) FromSlice(ctx context.Context, req []string) (error){

	if len(req) > 0 {
		s.PalletRecord = req[0]
	}
	if len(req) > 1 {
		s.PalletID = req[1]
	}

	return nil
}


func (s *Standard856V5LineItem) FromSlice(ctx context.Context, req []string) (error){

	if len(req) > 0 {
		s.DetailSectionLoopB = req[0]
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
		s.ManufacturersSerialCaseNumber = req[2]
	}
	if len(req) > 3 {
		s.BuyersPurchaseOrderNumber = req[3]
	}
	if len(req) > 4 {
		s.ItemIdentificationGTIN = req[4]
	}
	if len(req) > 5 {
		s.MasterStyle = req[5]
	}
	if len(req) > 6 {
		s.DetailStyle = req[6]
	}
	if len(req) > 7 {
		s.ColorCode = req[7]
	}
	if len(req) > 8 {
		s.SizeCode = req[8]
	}
	if len(req) > 9 {
		s.RevisionCode = req[9]
	}
	if len(req) > 10 {
		s.UnitOrBasisForMeasurementCode = req[10]
	}
	if len(req) > 11 {
		if req[11] != "" {
			quantityShipped, err := strconv.Atoi(req[11])
			if err != nil {
				return err
			}
			s.QuantityShipped = quantityShipped
		}
	}
	if len(req) > 12 {
		s.CountryOfOrigin = req[12]
	}
	
	if len(req) > 13 {
		s.ManufacturersOrderNumber = req[13]
	}
	if len(req) > 14 {
		s.ManufacturersLotID = req[14]
	}
	

	return nil
}