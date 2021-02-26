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
	Standard997V2s = []Standard997V2 {
		Standard997V2 {
			EnvelopeHeaderV3 : EnvelopeHeaderV3{
				InterchangeID : "202102268484912",
				ReceiverID : "123456789",
				SenderID : "383601069",
			},
			Body : Standard997V2Body{
				SenderID : "014628093",
				ProductionOrTest : "T",
				InterchangeID : "123456789",
				TransactionSetAcknowledgementCodes : "A",
			},
			EnvelopeTrailerV3 : EnvelopeTrailerV3{
				InterchangeID : "202102268484912",
			},
		},
	}
)


func TestStandard997V2ToBytes(t *testing.T) {

    ctx := context.Background()
	
	for Standard997V2Key, Standard997V2 := range Standard997V2s {
		byteArrayPointer, err := Standard997V2.ToBytes(ctx)
		if byteArrayPointer != nil {
			byteArray := *byteArrayPointer
			err := ioutil.WriteFile("./examples/997-" + strconv.Itoa(Standard997V2Key) + ".txt", byteArray, 0644)
			if err != nil {
				assert.Nil(t, err)
			}
		}
		// file, _ := json.MarshalIndent(product, "", "	")
		// _ = ioutil.WriteFile("./test/" +  gtin +".json", file, 0644)

		assert.Nil(t, err)
	}
	
}

func TestStandard997V1FromBytes(t *testing.T) {

    ctx := context.Background()
	
	bytes, readErr := ioutil.ReadFile("./examples/997_173384223_292101152647.txt")
	if readErr != nil {
		assert.Nil(t, readErr)
	}

	var Standard997V1 Standard997V1
	err := Standard997V1.FromBytes(ctx, bytes)
	assert.Nil(t, err)
	c, _ := json.Marshal(Standard997V1)
	fmt.Println(string(c))
	
}

func TestStandard997V2FromBytes(t *testing.T) {

    ctx := context.Background()
	
	bytes, readErr := ioutil.ReadFile("./examples/997-0.txt")
	if readErr != nil {
		assert.Nil(t, readErr)
	}

	var Standard997V2 Standard997V2
	err := Standard997V2.FromBytes(ctx, bytes)
	assert.Nil(t, err)
	
}