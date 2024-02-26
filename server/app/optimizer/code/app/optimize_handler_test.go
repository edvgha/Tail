package app

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var _ = setUpEnv()

func setUpEnv() int {
	os.Setenv("BUFFER_SIZE", "10")
	os.Setenv("DISCOUNT", "0.25")
	os.Setenv("DESIRED_EXPLORATION_SPEED", "2")
	os.Setenv("LOG_LEVEL", "0")
	os.Setenv("LEVEL_SIZE", "3")
	os.Setenv("BUCKET_SIZE", "20")
	os.Setenv("SPACE_DESC_FILE", "../space/spaces_desc.json")
	os.Setenv("CACHE_TTL", "1")
	return 0
}

func Test_explore(t *testing.T) {
	//0.1 - 0.3
	request := &Request{
		ID:          "123456789",
		Price:       0.2,
		FloorPrice:  0.1,
		DC:          "us-east4gcp",
		PublisherID: "1007950",
		BundleID:    "1207472156",
		TagID:       "BANNER",
		GeoCountry:  "USA",
		AdFormat:    "banner",
	}
	errN := 0.0
	noOKN := 0.0
	newPriceN := 0.0
	requestN := 100.0
	delay := 10 * time.Millisecond
	threshold := 15.0
	for i := 0; i < int(requestN); i++ {
		time.Sleep(delay)
		_, OK, err := explore(request)
		if err != nil {
			errN += 1
			continue
		}
		if !OK {
			noOKN += 1
			continue
		}
		newPriceN += 1
	}
	assert.Greater(t, threshold, (100*errN)/requestN, "number of errors above threshold")
	assert.Greater(t, threshold, (100*noOKN)/requestN, "number of exploitations above threshold")
	assert.Greater(t, (100*newPriceN)/requestN, 100-threshold, "number of exploration below threshold")
}
