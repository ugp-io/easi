package easi

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStandard846V3FromBytes(t *testing.T) {

	ctx := context.Background()

	bytes, readErr := ioutil.ReadFile("./examples/846.txt")
	if readErr != nil {
		assert.Nil(t, readErr)
	}

	var standard846V3 Standard846V3
	err := standard846V3.FromBytes(ctx, bytes)
	assert.Nil(t, err)

	// c, _ := json.Marshal(standard846V3)
	// fmt.Println(string(c))

}
