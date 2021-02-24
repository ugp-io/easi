package easi

import (
	"context"
	"strconv"
	"testing"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
)

var (
	standard997s = []Standard997 {
		Standard997 {
			SenderID : "014628093",
			ProductionOrTest : "T",
			InterchangeID : "123456789",
			TransactionSetAcknowledgementCodes : "A",
		},
	}
)


func TestStandard997ToBytes(t *testing.T) {

    ctx := context.Background()
	
	for standard997Key, standard997 := range standard997s {
		byteArrayPointer, err := standard997.ToBytes(ctx)
		if byteArrayPointer != nil {
			byteArray := *byteArrayPointer
			err := ioutil.WriteFile("./examples/997-" + strconv.Itoa(standard997Key) + ".txt", byteArray, 0644)
			if err != nil {
				assert.Nil(t, err)
			}
		}
		// file, _ := json.MarshalIndent(product, "", "	")
		// _ = ioutil.WriteFile("./test/" +  gtin +".json", file, 0644)

		assert.Nil(t, err)
	}
	
}

func TestStandard997FromBytes(t *testing.T) {

    ctx := context.Background()
	
	bytes, readErr := ioutil.ReadFile("./examples/997-0.txt")
	if readErr != nil {
		assert.Nil(t, readErr)
	}

	var standard997 Standard997
	err := standard997.FromBytes(ctx, bytes)
	assert.Nil(t, err)
	
}