package easi

import (
	"context"
	"strconv"
	"testing"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
)

var (
	standard850s = []Standard850 {
		Standard850 {
			Transaction : Standard850Transaction {
				TransactionSetPurpose : "00",
				PurchaseOrderTypeCode : "SA",
				PurchaseOrderNumber : "12345678",
				CurrencyCode : "USD",
				PurchaserAccountID : "12345",
				VendorID : "707738",
				DistributionCenterID : "05",
				FOBPaymentInstructions : "PP",
				SalesRequirementCodeShipment : "SC",
				DropShipCode : "N",
				DeliverToCompanyName : "Overlook Hotel",
				DeliverToContactName : "Jack Torrance",
				DeliverToAddress1 : "333 E Wonderview Ave",
				DeliverToCityName : "Estes Park",
				DeliverToStateCode : "CO",
				DeliverToPostalCode : "80517",
			},
			LineItems : []Standard850LineItem{
				Standard850LineItem {
					ItemIdentificationGTIN : "00821780002660",
					QuantityOrdered : 12,
					PurchaseUnitPrice : 185,
				},
				Standard850LineItem {
					ItemIdentificationGTIN : "00821780002799",
					QuantityOrdered : 6,
					PurchaseUnitPrice : 185,
				},
			},
			OtherCharges : []Standard850OtherCharge {
				Standard850OtherCharge{
					OtherChargeAmount : 200,
				},
			},
		},
	}
)


func TestStandard850ToBytes(t *testing.T) {

    ctx := context.Background()
	
	for standard850Key, standard850 := range standard850s {
		byteArrayPointer, err := standard850.ToBytes(ctx)
		if byteArrayPointer != nil {
			byteArray := *byteArrayPointer
			err := ioutil.WriteFile("./examples/850-" + strconv.Itoa(standard850Key) + ".txt", byteArray, 0644)
			if err != nil {
				assert.Nil(t, err)
			}
		}
		// file, _ := json.MarshalIndent(product, "", "	")
		// _ = ioutil.WriteFile("./test/" +  gtin +".json", file, 0644)

		assert.Nil(t, err)
	}
	
}

func TestStandard850FromBytes(t *testing.T) {

    ctx := context.Background()
	
	bytes, readErr := ioutil.ReadFile("./examples/850-0.txt")
	if readErr != nil {
		assert.Nil(t, readErr)
	}

	var standard850 Standard850
	err := standard850.FromBytes(ctx, bytes)
	assert.Nil(t, err)
	
}