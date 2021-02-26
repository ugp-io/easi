package easi

import(
	"context"
	"io"
	"fmt"
	"time"
	"bytes"
	"strconv"
	"encoding/csv"
	"github.com/jszwec/csvutil"
)

type Standard856V7 struct {
	EnvelopeHeaderV3 EnvelopeHeaderV3
	Transaction Standard856V7Transaction
	Pallets []Standard856V7Pallet
	Trailer Standard856V7Trailer
	EnvelopeTrailerV3 EnvelopeTrailerV3
}

type Standard856V7Transaction struct {
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
	DistributionCenterID string
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
	ShipmentDate string
	DeliverToContactName string
	DropShipCode string
}

type Standard856V7Pallet struct {
	PalletRecord string
	PalletID string
	Shipments []Standard856V7Shipment `csv:"-"`
}

type Standard856V7Shipment struct {
	DetailSectionLoopA string
	CarrierTrackingNumber string
	ManufacturersSerialCaseNumber string
	PurchaseOrderTypeCode string
	BuyersPurchaseOrderNumber string
	PODate string
	POTime string
	TrackingID string
	ManufacturersOrderNumber string
	CaseWeight float64 `csv:"-"`
	CaseWeightFormatted string
	FreightCharge int `csv:"-"`
	FreightChargeFormatted string
	LineItems []Standard856V7LineItem `csv:"-"`
}

type Standard856V7LineItem struct {
	DetailSectionLoopB string
	LineItemNumber int
	ItemIdentificationGTIN string
	MasterStyle string
	DetailStyle string
	ColorCode string
	SizeCode string
	RevisionCode string
	UnitOrBasisForMeasurementCode string
	QuantityShipped int
	CountryOfOrigin string
	ManufacturersLotID string
	BuyersPurchaseOrderNumber string
}

type Standard856V7Trailer struct {
	TrailerRecord string
	TotalCaseCount int
	TotalQtyShipped int
	TotalGrossWeight int
	TotalFreightCharges int `csv:"-"`
	TotalFreightChargesFormatted string
	RecordCount int
	TotalPalletCount int
}

func (s *Standard856V7) Prep(ctx context.Context) (error){

	// Header
	errHeader := s.EnvelopeHeaderV3.Prep(ctx)
	if errHeader != nil {
		return errHeader
	}
	s.EnvelopeHeaderV3.TransactionType = "856"

	// Transaction
	s.Transaction.Header = "01"
	s.Transaction.TransactionType = "856"
	s.Transaction.TransactionSetPurpose = "00"
	
	s.Transaction.VersionNumber = "7.0"
	s.Transaction.ASNDate = time.Now().Format("20060102")
	s.Transaction.ASNTime = time.Now().Format("150405")
	s.Transaction.ShipmentDate = time.Now().Format("20060102")

	// Pallets
	var totalFreightCharges int
	for palletKey, pallet := range s.Pallets {
		s.Pallets[palletKey].PalletRecord = "05"

		// Shipments
		for palletShipmentKey, palletShipment := range pallet.Shipments {
			s.Pallets[palletKey].Shipments[palletShipmentKey].DetailSectionLoopA = "02"
			s.Pallets[palletKey].Shipments[palletShipmentKey].CaseWeightFormatted = fmt.Sprintf("%.4f", palletShipment.CaseWeight)
			s.Pallets[palletKey].Shipments[palletShipmentKey].FreightChargeFormatted = fmt.Sprintf("%.4f", float64(palletShipment.FreightCharge) / 100)
			totalFreightCharges += palletShipment.FreightCharge

			// Line Items
			for palletShipmentLineItemKey, _ := range palletShipment.LineItems {
				s.Pallets[palletKey].Shipments[palletShipmentKey].LineItems[palletShipmentLineItemKey].DetailSectionLoopB = "03"
			}

		}

	}

	// Trailer
	s.Trailer.TrailerRecord = "09"
	s.Trailer.TotalFreightChargesFormatted = fmt.Sprintf("%.4f", float64(totalFreightCharges) / 100)

	// Trailer
	errTrailer := s.EnvelopeTrailerV3.Prep(ctx)
	if errTrailer != nil {
		return errTrailer
	}

	return nil
}

func (s *Standard856V7) ToBytes(ctx context.Context) (*[]byte, error){

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
	errEnvelopeHeaderV3 := enc.Encode(s.EnvelopeHeaderV3)
	if errEnvelopeHeaderV3 != nil {
		return nil, errEnvelopeHeaderV3
	}

	// Transaction
	errTransaction := enc.Encode(s.Transaction)
	if errTransaction != nil {
		return nil, errTransaction
	}

	// Pallets
	for _, pallet := range s.Pallets {
		errPallet := enc.Encode(pallet)
		if errPallet != nil {
			return nil, errPallet
		}

		// Shipments
		for _, shipment := range pallet.Shipments {
			errShipment := enc.Encode(shipment)
			if errShipment != nil {
				return nil, errShipment
			}

			// Line Items
			for _, lineItem := range shipment.LineItems {
				errLineItem := enc.Encode(lineItem)
				if errLineItem != nil {
					return nil, errLineItem
				}
			}

		}
	}

	// Trailer
	errTrailer := enc.Encode(s.Trailer)
	if errTrailer != nil {
		return nil, errTrailer
	}

	// Envelope Trailer
	errEnvelopeTrailerV3 := enc.Encode(s.EnvelopeTrailerV3)
	if errEnvelopeTrailerV3 != nil {
		return nil, errEnvelopeTrailerV3
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}

	byteArray := buf.Bytes()

	return &byteArray, nil
}

func (s *Standard856V7) FromBytes(ctx context.Context, req []byte) (error){

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
	
	var palletCount, shipmentCount int
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
			var x Standard856V7Transaction
			err := x.FromSlice(ctx, record)
			if err != nil {
				return err
			}
			s.Transaction = x
		case "05":
			var x Standard856V7Pallet
			err := x.FromSlice(ctx, record)
			if err != nil {
				return err
			}
			s.Pallets = append(s.Pallets, x)
			palletCount++
			shipmentCount = 0
		case "02":
			var x Standard856V7Shipment
			err := x.FromSlice(ctx, record)
			if err != nil {
				return err
			}
			s.Pallets[palletCount - 1].Shipments = append(s.Pallets[palletCount - 1].Shipments, x)
			shipmentCount++
		case "03":
			var x Standard856V7LineItem
			err := x.FromSlice(ctx, record)
			if err != nil {
				return err
			}
			s.Pallets[palletCount - 1].Shipments[shipmentCount - 1].LineItems = append(s.Pallets[palletCount - 1].Shipments[shipmentCount - 1].LineItems, x)
		case "09":
			var x Standard856V7Trailer
			err := x.FromSlice(ctx, record)
			if err != nil {
				return err
			}
			s.Trailer = x
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

func (s *Standard856V7Transaction) FromSlice(ctx context.Context, req []string) (error){

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
		s.DistributionCenterID = req[10]
	}
	if len(req) > 11 {
		s.DeliverToCompanyName = req[11]
	}
	if len(req) > 12 {
		s.DeliverToAddress1 = req[12]
	}
	if len(req) > 13 {
		s.DeliverToAddress2 = req[13]
	}
	if len(req) > 14 {
		s.DeliverToCityName = req[14]
	}
	if len(req) > 15 {
		s.DeliverToStateCode = req[15]
	}
	if len(req) > 16 {
		s.DeliverToPostalCode = req[16]
	}
	if len(req) > 17 {
		s.DeliverToCountryCode = req[17]
	}
	if len(req) > 18 {
		s.BOLNumber = req[18]
	}
	if len(req) > 19 {
		s.CarrierRoutingDetails = req[19]
	}
	if len(req) > 20 {
		s.TrailerID = req[20]
	}
	if len(req) > 21 {
		s.ShipmentDate = req[21]
	}
	if len(req) > 22 {
		s.DeliverToContactName = req[22]
	}
	if len(req) > 23 {
		s.DropShipCode = req[23]
	}

	return nil
}

func (s *Standard856V7Trailer) FromSlice(ctx context.Context, req []string) (error){

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
			totalFreightChargesFloat, err := strconv.ParseFloat(req[4], 64)
			if err != nil {
				return err
			}
			s.TotalFreightCharges = int((totalFreightChargesFloat * float64(100) + 0.5))
		}
	}
	if len(req) > 5 {
		if req[5] != "" {
			recordCount, err := strconv.Atoi(req[5])
			if err != nil {
				return err
			}
			s.RecordCount = recordCount
		}
	}
	if len(req) > 6 {
		if req[6] != "" {
			totalPalletCount, err := strconv.Atoi(req[6])
			if err != nil {
				return err
			}
			s.TotalPalletCount = totalPalletCount
		}
	}

	return nil
}

func (s *Standard856V7Pallet) FromSlice(ctx context.Context, req []string) (error){

	if len(req) > 0 {
		s.PalletRecord = req[0]
	}
	if len(req) > 1 {
		s.PalletID = req[1]
	}

	return nil
}

func (s *Standard856V7Shipment) FromSlice(ctx context.Context, req []string) (error){

	if len(req) > 0 {
		s.DetailSectionLoopA = req[0]
	}
	if len(req) > 1 {
		s.CarrierTrackingNumber = req[1]
	}
	if len(req) > 2 {
		s.ManufacturersSerialCaseNumber = req[2]
	}
	if len(req) > 3 {
		s.PurchaseOrderTypeCode = req[3]
	}
	if len(req) > 4 {
		s.BuyersPurchaseOrderNumber = req[4]
	}
	if len(req) > 5 {
		s.PODate = req[5]
	}
	if len(req) > 6 {
		s.POTime = req[6]
	}
	if len(req) > 7 {
		s.TrackingID = req[7]
	}
	if len(req) > 8 {
		s.ManufacturersOrderNumber = req[8]
	}
	if len(req) > 9 {
		if req[9] != "" {
			caseWeight, err := strconv.ParseFloat(req[9], 64)
			if err != nil {
				return err
			}
			s.CaseWeight = caseWeight
		}
	}
	if len(req) > 10 {
		if req[10] != "" {
			freightCharge, err := strconv.ParseFloat(req[10], 64)
			if err != nil {
				return err
			}
			s.FreightCharge = int((freightCharge * float64(100) + 0.5))
		}
	}

	return nil
}

func (s *Standard856V7LineItem) FromSlice(ctx context.Context, req []string) (error){

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
		s.ItemIdentificationGTIN = req[2]
	}
	if len(req) > 3 {
		s.MasterStyle = req[3]
	}
	if len(req) > 4 {
		s.DetailStyle = req[4]
	}
	if len(req) > 5 {
		s.ColorCode = req[5]
	}
	if len(req) > 6 {
		s.SizeCode = req[6]
	}
	if len(req) > 7 {
		s.RevisionCode = req[7]
	}
	if len(req) > 8 {
		s.UnitOrBasisForMeasurementCode = req[8]
	}
	if len(req) > 9 {
		if req[9] != "" {
			quantityShipped, err := strconv.Atoi(req[9])
			if err != nil {
				return err
			}
			s.QuantityShipped = quantityShipped
		}
	}
	if len(req) > 9 {
		s.CountryOfOrigin = req[9]
	}
	if len(req) > 10 {
		s.ManufacturersLotID = req[10]
	}
	if len(req) > 11 {
		s.BuyersPurchaseOrderNumber = req[11]
	}

	return nil
}