package app

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

/*
Spaces:
	{
      "data_center": "us-east4gcp",
      "ext_ad_format": "banner",
      "app_publisher_id": "1007950",
      "bundle_id": "1207472156",
      "tag_id": "BANNER",
      "device_geo_country": "USA",
      "context_hash": "13677617117914323147"
    },
    {
      "data_center":"us-east4gcp",
      "ext_ad_format":"banner",
      "app_publisher_id":"1008199",
      "bundle_id": "com.peoplefun.wordcross",
      "tag_id": "3cf97e2a-0370-4b58-86ef-3174701e3332",
      "device_geo_country": "USA",
      "context_hash": "871397612603656680"
    },
    {
      "data_center":"us-east4gcp",
      "ext_ad_format":"native",
      "app_publisher_id":"0",
      "bundle_id": "591560124",
      "tag_id": "BANNER",
      "device_geo_country": "USA",
      "context_hash": "15508723918213076872"
    },
    {
      "data_center":"us-east4gcp",
      "ext_ad_format":"native",
      "app_publisher_id":"0",
      "bundle_id": "com.pixel.art.coloring.color.number",
      "tag_id": "5d5d2779-c29f-4e11-ab3a-290fc46a93c5",
      "device_geo_country": "USA",
      "context_hash": "8387118903825488324"
    },
    { "data_center":"us-east4gcp",
      "ext_ad_format":"video",
      "app_publisher_id":"0",
      "bundle_id": "74519",
      "tag_id": "400008778",
      "device_geo_country": "USA",
      "context_hash": "18058710746607490439"
    },
    {
      "data_center":"us-east4gcp",
      "ext_ad_format":"video",
      "app_publisher_id":"0",
      "bundle_id": "com.peoplefun.wordcross",
      "tag_id": "INTER",
      "device_geo_country": "USA",
      "context_hash": "3571106062192531396"
    }
*/

func sendToOptimize(body []byte) (*Response, error) {
	req, err := http.NewRequest(http.MethodPost, "/optimize", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(optimizeHandler)
	handler.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		return nil, err
	}
	resBody := &Response{}
	err = json.NewDecoder(res.Body).Decode(resBody)
	if err != nil {
		return nil, err
	}
	return resBody, nil
}

func sendFeedback(body []byte) (*FeedBackResponse, error) {
	req, err := http.NewRequest(http.MethodPost, "/feedback", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(feedbackHandler)
	handler.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		return nil, err
	}
	resBody := &FeedBackResponse{}
	err = json.NewDecoder(res.Body).Decode(resBody)
	if err != nil {
		return nil, err
	}
	return resBody, nil
}

func Test_optimizeHandler_no_success(t *testing.T) {
	type args struct {
		body []byte
	}
	type expected struct {
		OptimizedPrice float64
		Status         string
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{"no_space",
			args{
				body: []byte(`{"id":"1234","price": 2.2,"floor_price": 0.2,
								"data_center": "us-east",
								"app_publisher_id": "app_id",
								"bundle_id": "bundle_id",
								"tag_id": "tag_id",
								"device_geo_country": "USA",
								"ext_ad_format": "banner"}`),
			},
			expected{
				OptimizedPrice: 2.2,
				Status:         "no space",
			},
		},
		{"unfeasible_price",
			args{
				body: []byte(`{"id":"1234","price": 1.2,"floor_price": 1.0,
								"data_center": "us-east4gcp",
      							"ext_ad_format": "banner",
      							"app_publisher_id": "1007950",
      							"bundle_id": "1207472156",
      							"tag_id": "BANNER",
      							"device_geo_country": "USA"}`),
			},
			expected{
				OptimizedPrice: 1.2,
				Status:         "unfeasible price 1.200000 [0.004100, 0.301500]",
			},
		},
		{"validation_error",
			args{
				body: []byte(`{"id":"1234","price": 1.2,"floor_price": 1.21,
								"data_center": "us-east4gcp",
      							"ext_ad_format": "banner",
      							"app_publisher_id": "1007950",
      							"bundle_id": "1207472156",
      							"tag_id": "BANNER",
      							"device_geo_country": "USA"}`),
			},
			expected{
				OptimizedPrice: 1.2,
				Status:         "validation error",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := sendToOptimize(tt.args.body)
			assert.Nil(t, err)
			assert.Equal(t, tt.expected.OptimizedPrice, res.OptimizedPrice)
			assert.Equal(t, tt.expected.Status, res.Status)
		})
	}
}

func Test_optimizeHandler_success(t *testing.T) {
	type args struct {
		body         []byte
		bodyFeedback []byte
	}
	type expected struct {
		FloorPrice float64
		Price      float64
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{"req_1",
			args{
				body: []byte(`{"id":"1001","price": 0.22,"floor_price": 0.01,
								"data_center": "us-east4gcp",
      							"ext_ad_format": "banner",
      							"app_publisher_id": "1007950",
      							"bundle_id": "1207472156",
      							"tag_id": "BANNER",
      							"device_geo_country": "USA"}`),
				bodyFeedback: []byte(`{"id":"1001","impression":true,"price":0.2}`),
			},
			expected{
				FloorPrice: 0.01,
				Price:      0.22,
			},
		},
		{"req_2",
			args{
				body: []byte(`{"id":"1002","price": 0.22,"floor_price": 0.01,
								"data_center": "us-east4gcp",
      							"ext_ad_format": "banner",
      							"app_publisher_id": "1007950",
      							"bundle_id": "1207472156",
      							"tag_id": "BANNER",
      							"device_geo_country": "USA"}`),
				bodyFeedback: []byte(`{"id":"1002","impression":true,"price":0.2}`),
			},
			expected{
				FloorPrice: 0.01,
				Price:      0.22,
			},
		},
		{"req_3",
			args{
				body: []byte(`{"id":"1003","price": 0.22,"floor_price": 0.01,
								"data_center": "us-east4gcp",
      							"ext_ad_format": "banner",
      							"app_publisher_id": "1007950",
      							"bundle_id": "1207472156",
      							"tag_id": "BANNER",
      							"device_geo_country": "USA"}`),
				bodyFeedback: []byte(`{"id":"1003","impression":true,"price":0.2}`),
			},
			expected{
				FloorPrice: 0.01,
				Price:      0.22,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := sendToOptimize(tt.args.body)
			time.Sleep(450 * time.Millisecond)
			assert.Nil(t, err)
			assert.GreaterOrEqual(t, res.OptimizedPrice, tt.expected.FloorPrice)
			assert.LessOrEqual(t, res.OptimizedPrice, tt.expected.Price)
			if res.Status == "explored" {
				r, err := sendFeedback(tt.args.bodyFeedback)
				assert.Nil(t, err)
				assert.True(t, r.Ack)
			} else {
				time.Sleep(450 * time.Millisecond)
			}
			if res.Status != "exploited" && res.Status != "explored" {
				t.Error("expected status is 'explored' or 'exploited'")
			}
		})
	}
}
