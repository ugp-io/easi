package easi

import (
	"context"
	"strconv"
	"testing"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
)

var (
	Standard856V7s = []Standard856V7 {
		Standard856V7 {
			EnvelopeHeaderV3 : EnvelopeHeaderV3{
				InterchangeID : "202102268484912",
				ReceiverID : "123456789",
				SenderID : "383601069",
			},
			Transaction : Standard856V7Transaction {
				ShipmentNumber : "987",
				PurchaserAccountID : "12345",
				VendorID : "707738",
				StoreID : "8976",
				DistributionCenterID : "05",
				DeliverToCompanyName : "Overlook Hotel",
				DeliverToContactName : "Jack Torrance",
				DeliverToAddress1 : "333 E Wonderview Ave",
				DeliverToCityName : "Estes Park",
				DeliverToStateCode : "CO",
				DeliverToPostalCode : "80517",
				DropShipCode : "N",
			},	
			Pallets : []Standard856V7Pallet {
				Standard856V7Pallet {
					PalletID : "123456789123456789",
					Shipments : []Standard856V7Shipment {
						Standard856V7Shipment{
							ManufacturersSerialCaseNumber : "345345",
							PurchaseOrderTypeCode : "SA",
							BuyersPurchaseOrderNumber : "34534534",
							CarrierTrackingNumber : "986979879878",
							PODate : "20210224",
							POTime : "150405",
							TrackingID : "987987987987",
							ManufacturersOrderNumber : "79878798798",
							CaseWeight : 12,
							FreightCharge : 12,
							LineItems : []Standard856V7LineItem{
								Standard856V7LineItem {
									ItemIdentificationGTIN : "00821780002660",
									QuantityShipped : 12,
									MasterStyle : "345345",
									ColorCode : "345345",
									SizeCode : "345345",
									UnitOrBasisForMeasurementCode : "23",
									CountryOfOrigin : "US",
								},
								Standard856V7LineItem {
									ItemIdentificationGTIN : "00821780002799",
									QuantityShipped : 6,
									MasterStyle : "345345",
									ColorCode : "345345",
									SizeCode : "345345",
									UnitOrBasisForMeasurementCode : "23",
									CountryOfOrigin : "US",
								},
							},
						},
					},
				},
			},
			EnvelopeTrailerV3 : EnvelopeTrailerV3{
				InterchangeID : "202102268484912",
			},
		},
	}
)


func TestStandard856V7ToBytes(t *testing.T) {

    ctx := context.Background()
	
	for Standard856V7Key, Standard856V7 := range Standard856V7s {
		byteArrayPointer, err := Standard856V7.ToBytes(ctx)
		if byteArrayPointer != nil {
			byteArray := *byteArrayPointer
			err := ioutil.WriteFile("./examples/856-" + strconv.Itoa(Standard856V7Key) + ".txt", byteArray, 0644)
			if err != nil {
				assert.Nil(t, err)
			}
		}

		assert.Nil(t, err)
	}
	
}

func TestStandard856V5FromBytes(t *testing.T) {

    ctx := context.Background()
	
	bytes, readErr := ioutil.ReadFile("./examples/856_173384223_20210130005845.txt")
	if readErr != nil {
		assert.Nil(t, readErr)
	}

	var standard856V5 Standard856V5
	err := standard856V5.FromBytes(ctx, bytes)
	assert.Nil(t, err)

	c, _ := json.Marshal(standard856V5)
	fmt.Println(string(c))
	
}

func TestStandard856V7FromBytes(t *testing.T) {

    ctx := context.Background()
	
	bytes, readErr := ioutil.ReadFile("./examples/856_173384223_20210130005845.txt")
	if readErr != nil {
		assert.Nil(t, readErr)
	}

	var Standard856V7 Standard856V7
	err := Standard856V7.FromBytes(ctx, bytes)
	assert.Nil(t, err)

	c, _ := json.Marshal(Standard856V7)
	fmt.Println(string(c))
	
}