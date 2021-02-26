package easi

import (
	"context"
	"strconv"
	"testing"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
)

var (
	standard850V1s = []Standard850V1 {
		Standard850V1 {
			EnvelopeHeaderV2 : EnvelopeHeaderV2{
				InterchangeID : "202102268484912",
				ReceiverID : "123456789",
				SenderID : "383601069",
			},
			Transaction : Standard850V1Transaction {
				TransactionSetPurpose : "00",
				PurchaseOrderTypeCode : "SA",
				PurchaseOrderNumber : "12345678",
				CurrencyCode : "USD",
				PurchaserAccountID : "12345",
				VendorID : "707738",
				StoreID : "1234",
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
			LineItems : []Standard850V1LineItem{
				Standard850V1LineItem {
					ItemIdentificationGTIN : "00821780002660",
					QuantityOrdered : 12,
					PurchaseUnitPrice : 185,
				},
				Standard850V1LineItem {
					ItemIdentificationGTIN : "00821780002799",
					QuantityOrdered : 6,
					PurchaseUnitPrice : 185,
				},
			},
			OtherCharges : []Standard850V1OtherCharge {
				Standard850V1OtherCharge{
					OtherChargeAmount : 200,
				},
			},
			EnvelopeTrailerV2 : EnvelopeTrailerV2{
				InterchangeID : "202102268484912",
			},
		},
	}
	Standard850V4s = []Standard850V4 {
		Standard850V4 {
			EnvelopeHeaderV3 : EnvelopeHeaderV3{
				InterchangeID : "202102268484912",
				ReceiverID : "123456789",
				SenderID : "383601069",
			},
			Transaction : Standard850V4Transaction {
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
			LineItems : []Standard850V4LineItem{
				Standard850V4LineItem {
					ItemIdentificationGTIN : "00821780002660",
					QuantityOrdered : 12,
					PurchaseUnitPrice : 185,
				},
				Standard850V4LineItem {
					ItemIdentificationGTIN : "00821780002799",
					QuantityOrdered : 6,
					PurchaseUnitPrice : 185,
				},
			},
			OtherCharges : []Standard850V4OtherCharge {
				Standard850V4OtherCharge{
					OtherChargeAmount : 200,
				},
			},
			EnvelopeTrailerV3 : EnvelopeTrailerV3{
				InterchangeID : "202102268484912",
			},
		},
	}
)

func TestStandard850V1ToBytes(t *testing.T) {

    ctx := context.Background()
	
	for standard850V1Key, standard850V1 := range standard850V1s {
		byteArrayPointer, err := standard850V1.ToBytes(ctx)
		if byteArrayPointer != nil {
			byteArray := *byteArrayPointer
			err := ioutil.WriteFile("./examples/850v1-" + strconv.Itoa(standard850V1Key) + ".txt", byteArray, 0644)
			if err != nil {
				assert.Nil(t, err)
			}
		}
		// file, _ := json.MarshalIndent(product, "", "	")
		// _ = ioutil.WriteFile("./test/" +  gtin +".json", file, 0644)

		assert.Nil(t, err)
	}
	
}

func TestStandard850V4ToBytes(t *testing.T) {

    ctx := context.Background()
	
	for Standard850V4Key, Standard850V4 := range Standard850V4s {
		byteArrayPointer, err := Standard850V4.ToBytes(ctx)
		if byteArrayPointer != nil {
			byteArray := *byteArrayPointer
			err := ioutil.WriteFile("./examples/850v4-" + strconv.Itoa(Standard850V4Key) + ".txt", byteArray, 0644)
			if err != nil {
				assert.Nil(t, err)
			}
		}
		// file, _ := json.MarshalIndent(product, "", "	")
		// _ = ioutil.WriteFile("./test/" +  gtin +".json", file, 0644)

		assert.Nil(t, err)
	}
	
}

func TestStandard850V4FromBytes(t *testing.T) {

    ctx := context.Background()
	
	bytes, readErr := ioutil.ReadFile("./examples/850v4-0.txt")
	if readErr != nil {
		assert.Nil(t, readErr)
	}

	var Standard850V4 Standard850V4
	err := Standard850V4.FromBytes(ctx, bytes)
	assert.Nil(t, err)
	
}