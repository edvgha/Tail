package space

import (
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"os"
	"sort"
	"strings"
	"testing"
)

func Test_writeToCSV(t *testing.T) {
	estimations := Estimations{Estimation{0.1, 0.2}, Estimation{0.2, 0.3}, Estimation{0.3, 0.4}}
	f, err := writeToCSV(estimations, zerolog.New(os.Stdout))
	assert.Nil(t, err)
	assert.NotNil(t, f)
	defer os.Remove(f.Name())
}

func Test_generateTrainConfig(t *testing.T) {
	f, n, err := generateTrainConfig("in.txt", "model.txt", zerolog.New(os.Stdout))
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

func buildEstimations() Estimations {
	pr := []float64{
		0.047054574844036745, 0.09505877181679956, 0.094273,
		0.1264417079834652, 0.15502474802486887, 0.16590212670950855,
		0.2261055918791397, 0.2, 0.2524434976597921,
		0.37808882972711866, 0.44300337394223666, 0.435,
		0.46377841954288446, 0.48058779358617676, 0.5139765783039226,
		0.5420645713502026, 0.6044742224914758, 0.626928832583778,
		0.6, 0.6287616092000385, 0.7071510305046703,
		0.7151211239888838, 0.7685735437278574, 0.806423588176367,
		0.8527733245642592, 0.8616743885235296, 0.84,
		0.9023273492016132, 0.9552039412651554, 0.9704883627880745}
	price := []float64{
		0.23, 0.7013793103448276, 1.1727586206896552,
		1.6441379310344828, 2.1155172413793104, 2.5868965517241382,
		3.0582758620689656, 3.529655172413793, 4.001034482758621,
		4.472413793103449, 4.943793103448277, 5.415172413793104,
		5.886551724137932, 6.3579310344827595, 6.8293103448275865,
		7.300689655172414, 7.772068965517242, 8.24344827586207,
		8.714827586206898, 9.186206896551726, 9.657586206896553,
		10.12896551724138, 10.600344827586207, 11.071724137931035,
		11.543103448275863, 12.01448275862069, 12.485862068965519,
		12.957241379310346, 13.428620689655173, 13.9}
	estimations := make(Estimations, 30)
	for i := 0; i < 30; i++ {
		estimations[i] = Estimation{
			price: price[i],
			pr:    pr[i],
		}
	}
	return estimations
}
func Test_learnNonDecreasing(t *testing.T) {
	estimations := buildEstimations()
	got, err := learnNonDecreasing(estimations, zerolog.New(os.Stdout))
	assert.Nil(t, err)
	assert.Equal(t, 30, len(got))
	sort.SliceIsSorted(got, func(i, j int) bool {
		return got[i] <= got[j]
	})
}
