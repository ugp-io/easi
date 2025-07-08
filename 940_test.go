package easi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	standard940v4 = Standard940V4{
		EnvelopeHeaderV3: EnvelopeHeaderV3{
			InterchangeID: "202102268484912",
			ReceiverID:    "123456789",
			SenderID:      "383601069",
		},
		Transaction: Standard940V4Transaction{
			TransactionSetPurpose: "00",
			PurchaseOrderTypeCode: "SA",
			PurchaseOrderNumber:   "12345678",
			CurrencyCode:          "USD",
			PurchaserAccountID:    "12345",
			VendorID:              "707738",
			// StoreID:                      "1234",
			FOBPaymentInstructions:               "PP",
			SalesRequirementCodeShipment:         "SC",
			DropShipCode:                         "N",
			DeliverToCompanyName:                 "Overlook Hotel",
			DeliverToContactName:                 "Jack Torrance",
			DeliverToAddress1:                    "333 E Wonderview Ave",
			DeliverToCityName:                    "Estes Park",
			DeliverToStateCode:                   "CO",
			DeliverToPostalCode:                  "80517",
			DistributionCenterID:                 "1234",
			CODForMerchandise:                    "N",
			DeliveryServiceLevel:                 "01",
			CommonCarrier:                        "UPS",
			AccountNumber:                        "123456789",
			NameOfAccount:                        "Jack Torrance",
			DeliverToCommercialOrResidentialSite: "C",
			CODTagsIndicator:                     "N",
		},
		LineItems: []Standard940V4LineItem{
			{
				GTIN:                          "00821780002660",
				MasterStyle:                   "2002",
				ColorCode:                     "BLK",
				SizeCode:                      "L",
				QuantityOrdered:               12,
				UnitOrBasisForMeasurementCode: "EA",
				SellingUnitPrice:              185,
				TotalMonetaryAmountOfLineItem: 2220,
			},
			{
				GTIN:                          "00821780002799",
				MasterStyle:                   "2002",
				ColorCode:                     "BLK",
				SizeCode:                      "L",
				QuantityOrdered:               6,
				UnitOrBasisForMeasurementCode: "EA",
				SellingUnitPrice:              185,
				TotalMonetaryAmountOfLineItem: 2220,
			},
		},
		OtherCharges: []Standard940V4OtherCharge{
			{
				OtherChargeDescription: "Shipping and Handling",
				OtherChargeAmount:      200,
			},
		},
		EnvelopeTrailerV3: EnvelopeTrailerV3{
			InterchangeID: "202102268484912",
		},
	}
)

func TestConvertToFile(t *testing.T) {

	ctx := context.Background()

	byteArrayPointer, err := standard940v4.ToBytes(ctx)
	if byteArrayPointer != nil {
		byteArray := *byteArrayPointer
		err := ioutil.WriteFile("./examples/940v4.txt", byteArray, 0644)
		if err != nil {
			assert.Nil(t, err)
		}
	}
	// file, _ := json.MarshalIndent(product, "", "	")
	// _ = ioutil.WriteFile("./test/" +  gtin +".json", file, 0644)
	if err != nil {
		t.Logf("Validation returned errors:\n%v", err)
	}
}

func TestConvertFromFile(t *testing.T) {

	ctx := context.Background()

	bytes, readErr := ioutil.ReadFile("./examples/940v4-input.txt")
	if readErr != nil {
		assert.Nil(t, readErr)
	}

	var standard940v4 Standard940V4
	err := standard940v4.FromBytes(ctx, bytes)
	intB, _ := json.Marshal(standard940v4)
	fmt.Println(string(intB))
	assert.Nil(t, err)

	byteArrayPointer, err := standard940v4.ToBytes(ctx)
	if byteArrayPointer != nil {
		byteArray := *byteArrayPointer
		err := ioutil.WriteFile("./examples/940v4-output.txt", byteArray, 0644)
		if err != nil {
			assert.Nil(t, err)
		}
	}
	// file, _ := json.MarshalIndent(product, "", "	")
	// _ = ioutil.WriteFile("./test/" +  gtin +".json", file, 0644)

	assert.Nil(t, err)

}
