package space

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func Test_writeToCSV(t *testing.T) {
	estimations := Estimations{Estimation{0.1, 0.2}, Estimation{0.2, 0.3}, Estimation{0.3, 0.4}}
	f, err := writeToCSV(estimations)
	assert.Nil(t, err)
	assert.NotNil(t, f)
	defer os.Remove(f.Name())
}

func Test_generateTrainConfig(t *testing.T) {
	f, n, err := generateTrainConfig("in.txt", "model.txt")
	defer os.Remove(f.Name())
	assert.Nil(t, err)
	assert.Greater(t, n, 0)
	fReader, err := os.Open(f.Name())
	assert.Nil(t, err)
	b := make([]byte, n)
	r, err := fReader.Read(b)
	assert.Nil(t, err)
	assert.Equal(t, n, r)
	str := string(b)
	assert.True(t, strings.Contains(str, "train"))
	assert.True(t, strings.Contains(str, "learning_rate"))
	fReader.Close()
}

func Test_execute(t *testing.T) {
	got, err := execute("test_data/train.conf", "test_data/predict.conf", "test_data/test_prediction.gout")
	assert.Nil(t, err)
	assert.Equal(t, 100, len(got))
	assert.Equal(t, 4.098817615583548, got[0])
	assert.Equal(t, 55.29056714574496, got[99])
	os.Remove("prediction.txt")
	os.Remove("LightGBM_model.txt")
}

func Test_readPredictions(t *testing.T) {
	got, err := readPredictions("test_data/test_prediction.gout")
	assert.Nil(t, err)
	assert.Equal(t, 100, len(got))
	assert.Equal(t, 4.098817615583548, got[0])
	assert.Equal(t, 55.29056714574496, got[99])
}
