package easi

import (
	"context"
	"strconv"
	"testing"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
)

var (
	standard856s = []Standard856 {
		Standard856 {
			Transaction : Standard856Transaction {
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
			Pallets : []Standard856Pallet {
				Standard856Pallet {
					PalletID : "123456789123456789",
					Shipments : []Standard856Shipment {
						Standard856Shipment{
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
							LineItems : []Standard856LineItem{
								Standard856LineItem {
									ItemIdentificationGTIN : "00821780002660",
									QuantityShipped : 12,
									MasterStyle : "345345",
									ColorCode : "345345",
									SizeCode : "345345",
									UnitOrBasisForMeasurementCode : "23",
									CountryOfOrigin : "US",
								},
								Standard856LineItem {
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
		},
	}
)


func TestStandard856ToBytes(t *testing.T) {

    ctx := context.Background()
	
	for standard856Key, standard856 := range standard856s {
		byteArrayPointer, err := standard856.ToBytes(ctx)
		if byteArrayPointer != nil {
			byteArray := *byteArrayPointer
			err := ioutil.WriteFile("./examples/856-" + strconv.Itoa(standard856Key) + ".txt", byteArray, 0644)
			if err != nil {
				assert.Nil(t, err)
			}
		}

		assert.Nil(t, err)
	}
	
}

func TestStandard856FromBytes(t *testing.T) {

    ctx := context.Background()
	
	bytes, readErr := ioutil.ReadFile("./examples/856.txt")
	if readErr != nil {
		assert.Nil(t, readErr)
	}

	var standard856 Standard856
	err := standard856.FromBytes(ctx, bytes)
	assert.Nil(t, err)
	
}