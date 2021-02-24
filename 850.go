package easi

import(
	"context"
	"fmt"
	"bytes"
	"io"
	"time"
	"strconv"
	"encoding/csv"
	"github.com/jszwec/csvutil"
)

type Standard850 struct{
	Transaction Standard850Transaction
	LineItems []Standard850LineItem
	OtherCharges []Standard850OtherCharge
	Trailer Standard850Trailer
}

type Standard850Transaction struct {
	Header string
	TransactionType string
	TransactionSetPurpose string
	VersionNumber string
	PurchaseOrderTypeCode string
	PurchaseOrderNumber string
	ReleaseNumber string
	PODate string
	POTime string
	ContractNumber string
	CurrencyCode string
	PurchaserAccountID string
	StoreID string
	DistributionCenterID string
	VendorID string
	ContactNameNumber string
	FOBPaymentInstructions string
	SalesRequirementCodeShipment string
	SalesRequirementCodeTruckLoad string
	SalesRequirementCodeShipDate string
	SalesRequirementCodeConsignmentOrShipBlind string
	PaymentTermsDiscountOffered string
	PaymentTermsDiscountDays string
	PaymentDueInNumberOfDaysWithoutDiscount string
	SpecificPaymentDate string
	LiteralOfPaymentTerms string
	RequestedShipDate string
	CancelDate string
	CarrierRoutingDetails string
	DeliverToCompanyName string
	DeliverToContactName string
	DeliverToAddress1 string
	DeliverToAddress2 string
	DeliverToCityName string
	DeliverToStateCode string
	DeliverToPostalCode string
	DeliverToCountryCode string
	DropShipCode string
	SpecialDeliveryInstructions string
	SpecialOrderInstructions string
	DeliverToCountyProvinceTownTerritory string
	PromotionalCode string
	DeliveryServiceLevel string
	DeliverToReceiversPhoneNumber string
	CustomerPONumber string
	CODForMerchandise string
	ReceiversEmailAddress string
	AccountNumber string 
	NameOfAccount string
	TrackingID string
	PurchasersAccountID string
	DeliverToCommercialOrResidentialSite string
	CODTagsIndicator string
	ThirdPartyAccountNumber string
}

type Standard850LineItem struct {
	DetailSectionLoopA string
	LineItemNumber int
	ItemIdentificationGTIN string
	MasterStyle string
	ColorCode string
	SizeCode string
	QuantityOrdered int
	UnitOrBasisForMeasurementCode string
	PurchaseUnitPrice int `csv:"-"`
	PurchaseUnitPriceFormatted string
	TotalMonetaryAmountOfLineItem int `csv:"-"`
	TotalMonetaryAmountOfLineItemFormatted string
}

type Standard850OtherCharge struct {
	OtherChargesRecord string
	LineItemNumberForOtherCharges int
	OtherChargeDescription string
	OtherChargeAmount int
	OtherChargeAmountFormatted string
}

type Standard850Trailer struct {
	TrailerRecord string
	RecordCount int
	TotalQuantityOrdered int
	TotalMonetaryValue int `csv:"-"`
	TotalMonetaryValueFormatted string
	TotalMonetaryValueOfOtherCharges int `csv:"-"`
	TotalMonetaryValueOfOtherChargesFormatted string
	NumberOfCases int
	PurchaseOrderTotalAmount int `csv:"-"`
	PurchaseOrderTotalAmountFormatted string
}

func (s *Standard850) Prep(ctx context.Context) (error){

	// Transaction
	s.Transaction.Header = "01"
	s.Transaction.TransactionType = "850"
	s.Transaction.VersionNumber = "4.0"
	s.Transaction.PODate = time.Now().Format("20060102")

	// Line Items
	var totalQuantityOrdered, totalMonetaryValue int
	for lineItemKey, lineItem := range s.LineItems {
		s.LineItems[lineItemKey].DetailSectionLoopA = "02"
		s.LineItems[lineItemKey].LineItemNumber = lineItemKey + 1
		s.LineItems[lineItemKey].UnitOrBasisForMeasurementCode = "EA"
		s.LineItems[lineItemKey].PurchaseUnitPriceFormatted = fmt.Sprintf("%.4f", float64(lineItem.PurchaseUnitPrice) / 100)
		s.LineItems[lineItemKey].TotalMonetaryAmountOfLineItemFormatted = fmt.Sprintf("%.4f", float64(lineItem.TotalMonetaryAmountOfLineItem) / 100)
		totalQuantityOrdered += lineItem.QuantityOrdered
		totalMonetaryValue += lineItem.PurchaseUnitPrice * lineItem.QuantityOrdered
	}

	// Other Charges
	var totalMonetaryValueOfOtherCharges int
	for otherChargeKey, otherCharge := range s.OtherCharges {
		s.OtherCharges[otherChargeKey].OtherChargesRecord = "06"
		s.OtherCharges[otherChargeKey].LineItemNumberForOtherCharges = otherChargeKey + 1
		s.OtherCharges[otherChargeKey].OtherChargeAmountFormatted = fmt.Sprintf("%.4f", float64(otherCharge.OtherChargeAmount) / 100)
		totalMonetaryValueOfOtherCharges += otherCharge.OtherChargeAmount
	}

	// Trailer
	s.Trailer.TrailerRecord = "09"
	s.Trailer.RecordCount = len(s.LineItems)
	s.Trailer.TotalQuantityOrdered = totalQuantityOrdered
	s.Trailer.TotalMonetaryValueFormatted = fmt.Sprintf("%.4f", float64(totalMonetaryValue) / 100)
	s.Trailer.TotalMonetaryValueOfOtherChargesFormatted = fmt.Sprintf("%.4f", float64(totalMonetaryValueOfOtherCharges) / 100)
	s.Trailer.PurchaseOrderTotalAmountFormatted = fmt.Sprintf("%.4f", float64(totalMonetaryValue + totalMonetaryValueOfOtherCharges) / 100)

	return nil
}



func (s *Standard850) ToBytes(ctx context.Context) (*[]byte, error){

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

	// Transaction
	errTransaction := enc.Encode(s.Transaction)
	if errTransaction != nil {
		return nil, errTransaction
	}

	// Line Items
	for _, lineItem := range s.LineItems {
		errLineItem := enc.Encode(lineItem)
		if errLineItem != nil {
			return nil, errLineItem
		}
	}

	// Other Charges
	errOtherCharges := enc.Encode(s.OtherCharges)
	if errOtherCharges != nil {
		return nil, errOtherCharges
	}

	// Trailer
	errTrailer := enc.Encode(s.Trailer)
	if errTrailer != nil {
		return nil, errTrailer
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}

	byteArray := buf.Bytes()

	return &byteArray, nil
}

func (s *Standard850) FromBytes(ctx context.Context, req []byte) (error){

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
		case "01":
			var x Standard850Transaction
			err := x.FromSlice(ctx, record)
			if err != nil {
				return err
			}
			s.Transaction = x
		case "02":
			var x Standard850LineItem
			err := x.FromSlice(ctx, record)
			if err != nil {
				return err
			}
			s.LineItems = append(s.LineItems, x)
		case "06":
			var x Standard850OtherCharge
			err := x.FromSlice(ctx, record)
			if err != nil {
				return err
			}
			s.OtherCharges = append(s.OtherCharges, x)
		case "09":
			var x Standard850Trailer
			err := x.FromSlice(ctx, record)
			if err != nil {
				return err
			}
			s.Trailer = x
		default:
			
		}

	}

	return nil
}

func (s *Standard850Transaction) FromSlice(ctx context.Context, req []string) (error){

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
		s.PurchaseOrderTypeCode = req[4]
	}
	if len(req) > 5 {
		s.PurchaseOrderNumber = req[5]
	}
	if len(req) > 6 {
		s.ReleaseNumber = req[6]
	}
	if len(req) > 7 {
		s.PODate = req[7]
	}
	if len(req) > 8 {
		s.POTime = req[8]
	}
	if len(req) > 9 {
		s.ContractNumber = req[9]
	}
	if len(req) > 10 {
		s.CurrencyCode = req[10]
	}
	if len(req) > 11 {
		s.PurchaserAccountID = req[11]
	}
	if len(req) > 12 {
		s.StoreID = req[12]
	}
	if len(req) > 13 {
		s.DistributionCenterID = req[13]
	}
	if len(req) > 14 {
		s.VendorID = req[14]
	}
	if len(req) > 15 {
		s.ContactNameNumber = req[15]
	}
	if len(req) > 16 {
		s.FOBPaymentInstructions = req[16]
	}
	if len(req) > 17 {
		s.SalesRequirementCodeShipment = req[17]
	}
	if len(req) > 18 {
		s.SalesRequirementCodeTruckLoad = req[18]
	}
	if len(req) > 19 {
		s.SalesRequirementCodeShipDate = req[19]
	}
	if len(req) > 20 {
		s.SalesRequirementCodeConsignmentOrShipBlind = req[20]
	}
	if len(req) > 21 {
		s.PaymentTermsDiscountOffered = req[21]
	}
	if len(req) > 22 {
		s.PaymentTermsDiscountDays = req[22]
	}
	if len(req) > 23 {
		s.PaymentDueInNumberOfDaysWithoutDiscount = req[23]
	}
	if len(req) > 24 {
		s.SpecificPaymentDate = req[24]
	}
	if len(req) > 25 {
		s.LiteralOfPaymentTerms = req[25]
	}
	if len(req) > 26 {
		s.RequestedShipDate = req[26]
	}
	if len(req) > 27 {
		s.CancelDate = req[27]
	}
	if len(req) > 28 {
		s.CarrierRoutingDetails = req[28]
	}
	if len(req) > 29 {
		s.DeliverToCompanyName = req[29]
	}
	if len(req) > 30 {
		s.DeliverToContactName = req[30]
	}
	if len(req) > 31 {
		s.DeliverToAddress1 = req[31]
	}
	if len(req) > 32 {
		s.DeliverToAddress2 = req[32]
	}
	if len(req) > 33 {
		s.DeliverToCityName = req[33]
	}
	if len(req) > 34 {
		s.DeliverToStateCode = req[34]
	}
	if len(req) > 35 {
		s.DeliverToPostalCode = req[35]
	}
	if len(req) > 36 {
		s.DeliverToCountryCode = req[36]
	}
	if len(req) > 37 {
		s.DropShipCode = req[37]
	}
	if len(req) > 38 {
		s.SpecialDeliveryInstructions = req[38]
	}
	if len(req) > 39 {
		s.SpecialOrderInstructions = req[39]
	}
	if len(req) > 40 {
		s.DeliverToCountyProvinceTownTerritory = req[40]
	}
	if len(req) > 41 {
		s.PromotionalCode = req[41]
	}
	if len(req) > 42 {
		s.DeliveryServiceLevel = req[42]
	}
	if len(req) > 43 {
		s.DeliverToReceiversPhoneNumber = req[43]
	}
	if len(req) > 44 {
		s.CustomerPONumber = req[44]
	}
	if len(req) > 45 {
		s.CODForMerchandise = req[45]
	}
	if len(req) > 46 {
		s.ReceiversEmailAddress = req[46]
	}
	if len(req) > 47 {
		s.AccountNumber = req[47]
	}
	if len(req) > 48 {
		s.NameOfAccount = req[48]
	}
	if len(req) > 49 {
		s.TrackingID = req[49]
	}
	if len(req) > 50 {
		s.PurchasersAccountID = req[50]
	}
	if len(req) > 51 {
		s.DeliverToCommercialOrResidentialSite = req[51]
	}
	if len(req) > 52 {
		s.CODTagsIndicator = req[52]
	}
	if len(req) > 53 {
		s.ThirdPartyAccountNumber = req[53]
	}

	return nil
}

func (s *Standard850LineItem) FromSlice(ctx context.Context, req []string) (error){

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
		s.MasterStyle = req[3]
	}
	if len(req) > 4 {
		s.ColorCode = req[4]
	}
	if len(req) > 5 {
		s.SizeCode = req[5]
	}
	if len(req) > 6 {
		if req[6] != "" {
			quantityOrdered, err := strconv.Atoi(req[6])
			if err != nil {
				return err
			}
			s.QuantityOrdered = quantityOrdered
		}
	}
	if len(req) > 7 {
		s.UnitOrBasisForMeasurementCode = req[7]
	}
	if len(req) > 8 {
		if req[8] != "" {
			purchaseUnitPriceFloat, err := strconv.ParseFloat(req[8], 64)
			if err != nil {
				return err
			}
			s.PurchaseUnitPrice = int((purchaseUnitPriceFloat * float64(100) + 0.5))
		}
	}
	if len(req) > 9 {
		if req[9] != "" {
			totalMonetaryAmountOfLineItemFloat, err := strconv.ParseFloat(req[9], 64)
			if err != nil {
				return err
			}
			s.TotalMonetaryAmountOfLineItem = int((totalMonetaryAmountOfLineItemFloat * float64(100) + 0.5))
		}
		
	}

	return nil
}

func (s *Standard850OtherCharge) FromSlice(ctx context.Context, req []string) (error){

	if len(req) > 0 {
		s.OtherChargesRecord = req[0]
	}
	if len(req) > 1 {
		if req[1] != "" {
			lineItemNumberForOtherCharges, err := strconv.Atoi(req[1])
			if err != nil {
				return err
			}
			s.LineItemNumberForOtherCharges = lineItemNumberForOtherCharges
		}
	}
	if len(req) > 2 {
		s.OtherChargeDescription = req[2]
	}
	if len(req) > 3 {
		if req[3] != "" {
			otherChargeAmountFloat, err := strconv.ParseFloat(req[3], 64)
			if err != nil {
				return err
			}
			s.OtherChargeAmount = int((otherChargeAmountFloat * float64(100) + 0.5))
		}
		
	}
	return nil
}

func (s *Standard850Trailer) FromSlice(ctx context.Context, req []string) (error){

	if len(req) > 0 {
		s.TrailerRecord = req[0]
	}
	if len(req) > 1 {
		if req[1] != "" {
			recordCount, err := strconv.Atoi(req[1])
			if err != nil {
				return err
			}
			s.RecordCount = recordCount
		}
	}
	if len(req) > 2 {
		if req[2] != "" {
			totalQuantityOrdered, err := strconv.Atoi(req[1])
			if err != nil {
				return err
			}
			s.TotalQuantityOrdered = totalQuantityOrdered
		}
	}
	if len(req) > 3 {
		if req[3] != "" {
			totalMonetaryValueFloat, err := strconv.ParseFloat(req[3], 64)
			if err != nil {
				return err
			}
			s.TotalMonetaryValue = int((totalMonetaryValueFloat * float64(100) + 0.5))
		}
		
	}
	if len(req) > 4 {
		if req[4] != "" {
			totalMonetaryValueOfOtherChargesFloat, err := strconv.ParseFloat(req[4], 64)
			if err != nil {
				return err
			}
			s.TotalMonetaryValueOfOtherCharges = int((totalMonetaryValueOfOtherChargesFloat * float64(100) + 0.5))
		}
		
	}
	if len(req) > 5 {
		if req[5] != "" {
			numberOfCases, err := strconv.Atoi(req[1])
			if err != nil {
				return err
			}
			s.NumberOfCases = numberOfCases
		}
	}
	if len(req) > 6 {
		if req[6] != "" {
			purchaseOrderTotalAmountFloat, err := strconv.ParseFloat(req[6], 64)
			if err != nil {
				return err
			}
			s.PurchaseOrderTotalAmount = int((purchaseOrderTotalAmountFloat * float64(100) + 0.5))
		}
		
	}

	return nil
}