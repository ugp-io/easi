package easi

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/jszwec/csvutil"
)

type Standard940V4 struct {
	EnvelopeHeaderV3  EnvelopeHeaderV3
	Transaction       Standard940V4Transaction
	LineItems         []Standard940V4LineItem
	OtherCharges      []Standard940V4OtherCharge
	Trailer           Standard940V4Trailer
	EnvelopeTrailerV3 EnvelopeTrailerV3
}

type Standard940V4Transaction struct {
	Header                                     string
	TransactionType                            string
	TransactionSetPurpose                      string
	VersionNumber                              string
	PurchaseOrderTypeCode                      string
	PurchaseOrderNumber                        string
	PurchaseOrderDate                          string
	PurchaseOrderTime                          string
	CustomerPONumber                           string
	CurrencyCode                               string
	PurchaserAccountID                         string
	DistributionCenterID                       string
	VendorID                                   string
	ContactNameNumber                          string
	FOBPaymentInstructions                     string
	SalesRequirementCodeShipment               string
	SalesRequirementCodeTruckLoad              string
	SalesRequirementCodeShipDate               string
	SalesRequirementCodeConsignmentOrShipBlind string
	CODForMerchandise                          string
	PaymentTermsDiscountOffered                string
	PaymentTermsDiscountDays                   int
	PaymentDueInDaysWithoutDiscount            int
	SpecificPaymentDate                        string
	LiteralOfPaymentTerms                      string
	RequestedShipDate                          string
	DeliveryServiceLevel                       string
	CommonCarrier                              string
	CancelDate                                 string
	DeliverToCompanyName                       string
	DeliverToContactName                       string
	ReceiversPhoneNumber                       string
	ReceiversEmailAddress                      string
	AccountNumber                              string
	NameOfAccount                              string
	TrackingID                                 string
	PurchasersAccountID                        string
	DeliverToAddress1                          string
	DeliverToAddress2                          string
	DeliverToCityName                          string
	DeliverToStateCode                         string
	DeliverToPostalCode                        string
	DeliverToCountryCode                       string
	DropShipCode                               string
	SpecialDeliveryInstructions                string
	SpecialOrderInstructions                   string
	DeliverToCommercialOrResidentialSite       string
	CODTagsIndicator                           string
	ThirdPartyAccountNumber                    string
	DeliverToCountyProvinceTownTerritory       string
	PromotionalCode                            string
}

type Standard940V4LineItem struct {
	LineItemRecord                         string
	LineItemNumber                         int
	GTIN                                   string
	MasterStyle                            string
	ColorCode                              string
	SizeCode                               string
	QuantityOrdered                        int
	UnitOrBasisForMeasurementCode          string
	SellingUnitPrice                       int `csv:"-"`
	SellingUnitPriceFormatted              string
	TotalMonetaryAmountOfLineItem          int `csv:"-"`
	TotalMonetaryAmountOfLineItemFormatted string
}

type Standard940V4OtherCharge struct {
	OtherChargeRecord             string
	LineItemNumberForOtherCharges int
	OtherChargeDescription        string
	OtherChargeAmount             int `csv:"-"`
	OtherChargeAmountFormatted    string
}

type Standard940V4Trailer struct {
	TrailerRecord        string
	RecordCount          int
	TotalQuantityOrdered int
	NumberOfCases        int
	TotalAmountOfPO      string
}

func (s *Standard940V4) Prep(ctx context.Context) error {

	// Header
	errHeader := s.EnvelopeHeaderV3.Prep(ctx)
	if errHeader != nil {
		return errHeader
	}
	s.EnvelopeHeaderV3.TransactionType = "940"

	// Transaction
	s.Transaction.Header = "01"
	s.Transaction.TransactionType = "940"
	s.Transaction.VersionNumber = "4.0"
	s.Transaction.TransactionSetPurpose = "00"
	s.Transaction.PurchaseOrderDate = time.Now().Format("20060102")

	// Line Items
	var totalQuantityOrdered, totalMonetaryValue int
	for lineItemKey, lineItem := range s.LineItems {
		s.LineItems[lineItemKey].LineItemRecord = "02"
		s.LineItems[lineItemKey].LineItemNumber = lineItemKey + 1
		s.LineItems[lineItemKey].UnitOrBasisForMeasurementCode = "EA"
		s.LineItems[lineItemKey].SellingUnitPriceFormatted = fmt.Sprintf("%.4f", float64(lineItem.SellingUnitPrice)/100)
		s.LineItems[lineItemKey].TotalMonetaryAmountOfLineItemFormatted = fmt.Sprintf("%.4f", float64(lineItem.TotalMonetaryAmountOfLineItem)/100)
		totalQuantityOrdered += lineItem.QuantityOrdered
		totalMonetaryValue += lineItem.SellingUnitPrice * lineItem.QuantityOrdered
	}

	// Other Charges
	var totalMonetaryValueOfOtherCharges int
	for otherChargeKey, otherCharge := range s.OtherCharges {
		s.OtherCharges[otherChargeKey].OtherChargeRecord = "06"
		s.OtherCharges[otherChargeKey].LineItemNumberForOtherCharges = otherChargeKey + 1
		s.OtherCharges[otherChargeKey].OtherChargeAmountFormatted = fmt.Sprintf("%.4f", float64(otherCharge.OtherChargeAmount)/100)
		totalMonetaryValueOfOtherCharges += otherCharge.OtherChargeAmount
	}

	// Trailer
	s.Trailer.TrailerRecord = "09"
	s.Trailer.RecordCount = len(s.LineItems)
	s.Trailer.TotalQuantityOrdered = totalQuantityOrdered

	// s.Trailer.TotalMonetaryValueFormatted = fmt.Sprintf("%.4f", float64(totalMonetaryValue)/100)
	// s.Trailer.TotalMonetaryValueOfOtherChargesFormatted = fmt.Sprintf("%.4f", float64(totalMonetaryValueOfOtherCharges)/100)
	// s.Trailer.PurchaseOrderTotalAmountFormatted = fmt.Sprintf("%.4f", float64(totalMonetaryValue+totalMonetaryValueOfOtherCharges)/100)

	// Trailer
	errTrailer := s.EnvelopeTrailerV3.Prep(ctx)
	if errTrailer != nil {
		return errTrailer
	}

	return nil
}

func (s *Standard940V4) ToBytes(ctx context.Context) (*[]byte, error) {

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	w.Comma = '\t'
	enc := csvutil.NewEncoder(w)
	enc.AutoHeader = false

	// Prep
	errPrep := s.Prep(ctx)
	if errPrep != nil {
		return nil, errPrep
	}

	if err := s.ValidateStandard940V4(ctx); err != nil {
		return nil, err
	}

	// Envelope Header)
	errEnvelopeHeaderV3 := enc.Encode(s.EnvelopeHeaderV3)
	if errEnvelopeHeaderV3 != nil {
		return nil, errEnvelopeHeaderV3
	}

	// Convert Transaction
	errTransaction := enc.Encode(s.Transaction)
	if errTransaction != nil {
		return nil, errTransaction
	}

	// Convert Line Items
	for _, lineItem := range s.LineItems {
		errLineItem := enc.Encode(lineItem)
		if errLineItem != nil {
			return nil, errLineItem
		}
	}

	// Convert Other Charges
	for _, otherCharge := range s.OtherCharges {
		errOtherCharges := enc.Encode(otherCharge)
		if errOtherCharges != nil {
			return nil, errOtherCharges
		}
	}

	// Convert Trailer
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

func (s *Standard940V4) FromBytes(ctx context.Context, req []byte) error {

	r := csv.NewReader(bytes.NewReader(req))
	r.Comma = '\t'

	// Headerless
	blankHeader, errHeader := csvutil.Header(Header{}, "txt")
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

		record := dec.Record()
		if len(record) > 0 {
			// Build
			switch record[0] {
			case "EASI":
				var x EnvelopeHeaderV3
				err := x.FromSlice(ctx, record)
				if err != nil {
					return err
				}
				s.EnvelopeHeaderV3 = x
			case "01":
				var x Standard940V4Transaction
				err := x.ConvertFromFile(ctx, record)
				if err != nil {
					return err
				}
				s.Transaction = x
			case "02":
				var x Standard940V4LineItem
				err := x.ConvertFromFile(ctx, record)
				if err != nil {
					return err
				}
				s.LineItems = append(s.LineItems, x)
			case "06":
				var x Standard940V4OtherCharge
				err := x.ConvertFromFile(ctx, record)
				if err != nil {
					return err
				}
				s.OtherCharges = append(s.OtherCharges, x)
			case "09":
				var x Standard940V4Trailer
				err := x.ConvertFromFile(ctx, record)
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
			}
		}
	}

	return nil
}

func (s *Standard940V4Transaction) ConvertFromFile(ctx context.Context, transactions []string) error {
	for transactionIndex, transactionValue := range transactions {
		switch transactionIndex {
		case 0:
			s.Header = transactionValue
		case 1:
			s.TransactionType = transactionValue
		case 2:
			s.TransactionSetPurpose = transactionValue
		case 3:
			s.VersionNumber = transactionValue
		case 4:
			s.PurchaseOrderTypeCode = transactionValue
		case 5:
			s.PurchaseOrderNumber = transactionValue
		case 6:
			s.PurchaseOrderDate = transactionValue
		case 7:
			s.PurchaseOrderTime = transactionValue
		case 8:
			s.CustomerPONumber = transactionValue
		case 9:
			s.CurrencyCode = transactionValue
		case 10:
			s.PurchaserAccountID = transactionValue
		case 11:
			s.DistributionCenterID = transactionValue
		case 12:
			s.VendorID = transactionValue
		case 13:
			s.ContactNameNumber = transactionValue
		case 14:
			s.FOBPaymentInstructions = transactionValue
		case 15:
			s.SalesRequirementCodeShipment = transactionValue
		case 16:
			s.SalesRequirementCodeTruckLoad = transactionValue
		case 17:
			s.SalesRequirementCodeShipDate = transactionValue
		case 18:
			s.SalesRequirementCodeConsignmentOrShipBlind = transactionValue
		case 19:
			s.CODForMerchandise = transactionValue
		case 20:
			s.PaymentTermsDiscountOffered = transactionValue
		case 21:
			paymentTermsDiscountDays, _ := strconv.Atoi(transactionValue)
			s.PaymentTermsDiscountDays = paymentTermsDiscountDays
		case 22:
			paymentDueInDaysWithoutDiscount, _ := strconv.Atoi(transactionValue)
			s.PaymentDueInDaysWithoutDiscount = paymentDueInDaysWithoutDiscount
		case 23:
			s.SpecificPaymentDate = transactionValue
		case 24:
			s.LiteralOfPaymentTerms = transactionValue
		case 25:
			s.RequestedShipDate = transactionValue
		case 26:
			s.DeliveryServiceLevel = transactionValue
		case 27:
			s.CommonCarrier = transactionValue
		case 28:
			s.CancelDate = transactionValue
		case 29:
			s.DeliverToCompanyName = transactionValue
		case 30:
			s.DeliverToContactName = transactionValue
		case 31:
			s.ReceiversPhoneNumber = transactionValue
		case 32:
			s.ReceiversEmailAddress = transactionValue
		case 33:
			s.AccountNumber = transactionValue
		case 34:
			s.NameOfAccount = transactionValue
		case 35:
			s.TrackingID = transactionValue
		case 36:
			s.PurchasersAccountID = transactionValue
		case 37:
			s.DeliverToAddress1 = transactionValue
		case 38:
			s.DeliverToAddress2 = transactionValue
		case 39:
			s.DeliverToCityName = transactionValue
		case 40:
			s.DeliverToStateCode = transactionValue
		case 41:
			s.DeliverToPostalCode = transactionValue
		case 42:
			s.DeliverToCountryCode = transactionValue
		case 43:
			s.DropShipCode = transactionValue
		case 44:
			s.SpecialDeliveryInstructions = transactionValue
		case 45:
			s.SpecialOrderInstructions = transactionValue
		case 46:
			s.DeliverToCommercialOrResidentialSite = transactionValue
		case 47:
			s.CODTagsIndicator = transactionValue
		case 48:
			s.ThirdPartyAccountNumber = transactionValue
		case 49:
			s.DeliverToCountyProvinceTownTerritory = transactionValue
		case 50:
			s.PromotionalCode = transactionValue
		}
	}
	return nil
}

func (s *Standard940V4LineItem) ConvertFromFile(ctx context.Context, lineItems []string) error {
	for indexLineItem, valLineItem := range lineItems {
		switch indexLineItem {
		case 0:
			s.LineItemRecord = valLineItem
		case 1:
			lineItemNumber, _ := strconv.Atoi(valLineItem)
			s.LineItemNumber = lineItemNumber
		case 2:
			s.GTIN = valLineItem
		case 3:
			s.MasterStyle = valLineItem
		case 4:
			s.ColorCode = valLineItem
		case 5:
			s.SizeCode = valLineItem
		case 6:
			quantityOrdered, _ := strconv.Atoi(valLineItem)
			s.QuantityOrdered = quantityOrdered
		case 7:
			s.UnitOrBasisForMeasurementCode = valLineItem
		case 8:
			sellingUnitPrice, _ := strconv.Atoi(valLineItem)
			s.SellingUnitPrice = sellingUnitPrice
		case 9:
			totalMonetaryAmountOfLineItem, _ := strconv.Atoi(valLineItem)
			s.TotalMonetaryAmountOfLineItem = totalMonetaryAmountOfLineItem
		}
	}
	return nil
}

func (s *Standard940V4OtherCharge) ConvertFromFile(ctx context.Context, otherCharges []string) error {
	for indexOtherCharge, valOtherCharge := range otherCharges {
		switch indexOtherCharge {
		case 0:
			s.OtherChargeRecord = valOtherCharge
		case 1:
			lineItemNumberForOtherCharges, _ := strconv.Atoi(valOtherCharge)
			s.LineItemNumberForOtherCharges = lineItemNumberForOtherCharges
		case 2:
			s.OtherChargeDescription = valOtherCharge
		case 3:
			valOtherCharge, _ := strconv.Atoi(valOtherCharge)
			s.OtherChargeAmount = valOtherCharge
		}
	}
	return nil
}

func (s *Standard940V4Trailer) ConvertFromFile(ctx context.Context, trailers []string) error {
	for indexTrailer, valTrailer := range trailers {
		switch indexTrailer {
		case 0:
			s.TrailerRecord = valTrailer
		case 1:
			recordCount, _ := strconv.Atoi(valTrailer)
			s.RecordCount = recordCount
		case 2:
			totalQuantityOrdered, _ := strconv.Atoi(valTrailer)
			s.TotalQuantityOrdered = totalQuantityOrdered
		case 3:
			numberOfCases, _ := strconv.Atoi(valTrailer)
			s.NumberOfCases = numberOfCases
		case 4:
			s.TotalAmountOfPO = valTrailer
		}
	}
	return nil
}

func (s *Standard940V4) ValidateStandard940V4(ctx context.Context) error {

	var errs []error
	transaction := s.Transaction
	lineItems := s.LineItems
	otherCharges := s.OtherCharges
	trailer := s.Trailer

	if transaction.Header == "" {
		errs = append(errs, fmt.Errorf("Transaction Header is required"))
	} else if transaction.Header != "01" {
		errs = append(errs, fmt.Errorf("Transaction Header must be '01'"))
	}

	if transaction.TransactionType == "" {
		errs = append(errs, fmt.Errorf("Transaction Type is required"))
	} else if transaction.TransactionType != "940" {
		errs = append(errs, fmt.Errorf("Transaction Type must be '940'"))
	}

	if transaction.TransactionSetPurpose == "" {
		errs = append(errs, fmt.Errorf("Transaction Set Purpose is required"))
	}

	if transaction.VersionNumber == "" {
		errs = append(errs, fmt.Errorf("Version Number is required"))
	}

	if transaction.PurchaseOrderNumber == "" {
		errs = append(errs, fmt.Errorf("Purchase Order Number is required"))
	}

	if transaction.PurchaseOrderDate == "" {
		errs = append(errs, fmt.Errorf("Purchase Order Date is required"))
	}

	if transaction.CurrencyCode == "" {
		errs = append(errs, fmt.Errorf("Currency Code is required"))
	}

	if transaction.PurchaserAccountID == "" {
		errs = append(errs, fmt.Errorf("Purchaser Account ID is required"))
	}

	if transaction.DistributionCenterID == "" {
		errs = append(errs, fmt.Errorf("Distribution Center ID is required"))
	}

	if transaction.VendorID == "" {
		errs = append(errs, fmt.Errorf("Vendor ID is required"))
	}

	if transaction.FOBPaymentInstructions == "" {
		errs = append(errs, fmt.Errorf("FOB Payment Instructions are required"))
	}

	if transaction.SalesRequirementCodeShipment == "" {
		errs = append(errs, fmt.Errorf("Sales Requirement Code Shipment is required"))
	}

	if transaction.CODForMerchandise == "" {
		errs = append(errs, fmt.Errorf("COD For Merchandise is required"))
	}

	if transaction.PaymentDueInDaysWithoutDiscount < 0 {
		errs = append(errs, fmt.Errorf("Payment Due In Days Without Discount is required"))
	}

	if transaction.DeliveryServiceLevel == "" {
		errs = append(errs, fmt.Errorf("Delivery Service Level is required"))
	}

	if transaction.CommonCarrier == "" {
		errs = append(errs, fmt.Errorf("Common Carrier is required"))
	}

	if transaction.DeliverToCompanyName == "" {
		errs = append(errs, fmt.Errorf("Deliver To Company Name is required"))
	}

	if transaction.AccountNumber == "" {
		errs = append(errs, fmt.Errorf("Account Number is required"))
	}

	if transaction.NameOfAccount == "" {
		errs = append(errs, fmt.Errorf("Name Of Account is required"))
	}

	if transaction.DeliverToAddress1 == "" {
		errs = append(errs, fmt.Errorf("Deliver To Address 1 is required"))
	}

	if transaction.DeliverToCityName == "" {
		errs = append(errs, fmt.Errorf("Deliver To City Name is required"))
	}

	if transaction.DeliverToStateCode == "" {
		errs = append(errs, fmt.Errorf("Deliver To State Code is required"))
	}

	if transaction.DeliverToPostalCode == "" {
		errs = append(errs, fmt.Errorf("Deliver To Postal Code is required"))
	}

	if transaction.DropShipCode == "" {
		errs = append(errs, fmt.Errorf("Drop Ship Code is required"))
	}

	if transaction.DeliverToCommercialOrResidentialSite == "" {
		errs = append(errs, fmt.Errorf("Deliver To Commercial Or Residential Site is required"))
	}

	if transaction.CODTagsIndicator == "" {
		errs = append(errs, fmt.Errorf("COD Tags Indicator is required"))
	}

	if len(lineItems) <= 0 {
		errs = append(errs, fmt.Errorf("Line Items are required"))
	} else {
		for _, lineItem := range lineItems {

			if lineItem.LineItemRecord == "" {
				errs = append(errs, fmt.Errorf("Line Item Record is required"))
			} else if lineItem.LineItemRecord != "02" {
				errs = append(errs, fmt.Errorf("Line Item Record must be '02'"))
			}

			if lineItem.LineItemNumber <= 0 {
				errs = append(errs, fmt.Errorf("Line Item Number is required"))
			}

			if lineItem.GTIN == "" {
				errs = append(errs, fmt.Errorf("GTIN is required"))
			}

			if lineItem.QuantityOrdered <= 0 {
				errs = append(errs, fmt.Errorf("Quantity Ordered is required"))
			}

			if lineItem.UnitOrBasisForMeasurementCode == "" {
				errs = append(errs, fmt.Errorf("Unit Or Basis For Measurement Code is required"))
			}
		}
	}

	if len(otherCharges) <= 0 {
		// errs = append(errs, fmt.Errorf("Other Charges are required"))
	} else {
		for _, otherCharge := range otherCharges {

			if otherCharge.OtherChargeRecord == "" {
				errs = append(errs, fmt.Errorf("Other Charge Record is required"))
			} else if otherCharge.OtherChargeRecord != "06" {
				errs = append(errs, fmt.Errorf("Other Charge Record must be '06'"))
			}

			if otherCharge.LineItemNumberForOtherCharges <= 0 {
				errs = append(errs, fmt.Errorf("Line Item Number For Other Charges is required"))
			}

			if otherCharge.OtherChargeAmount == 0 {
				errs = append(errs, fmt.Errorf("Other Charge Amount is required"))
			}
		}
	}

	if trailer.TrailerRecord == "" {
		errs = append(errs, fmt.Errorf("Trailer Record is required"))
	} else if trailer.TrailerRecord != "09" {
		errs = append(errs, fmt.Errorf("Trailer Record must be '09'"))
	}

	if trailer.RecordCount <= 0 {
		errs = append(errs, fmt.Errorf("Record Count is required"))
	}

	if trailer.TotalQuantityOrdered < 0 {
		errs = append(errs, fmt.Errorf("Total Quantity Ordered is required"))
	}

	if trailer.NumberOfCases < 0 {
		errs = append(errs, fmt.Errorf("Number Of Cases is required"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
