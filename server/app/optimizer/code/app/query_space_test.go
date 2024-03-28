package app

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"tail.server/app/optimizer/code/space"
	"testing"
)

func Test_spaceHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/optimize", nil)
	assert.Nil(t, err)
	q := req.URL.Query()
	q.Add("ctx", "13677617117914323147")
	req.URL.RawQuery = q.Encode()
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(spaceHandler)
	handler.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)
	learnedEstimations := &space.LearnedEstimation{}
	err = json.NewDecoder(res.Body).Decode(learnedEstimations)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(learnedEstimations.Level))
	assert.Equal(t, len(learnedEstimations.Level[0].Price), len(learnedEstimations.Level[0].Pr))
	assert.Equal(t, len(learnedEstimations.Level[1].Price), len(learnedEstimations.Level[1].Pr))
	assert.Equal(t, len(learnedEstimations.Level[2].Price), len(learnedEstimations.Level[2].Pr))
}
