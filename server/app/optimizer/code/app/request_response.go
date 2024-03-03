package app

type Response struct {
	OptimizedPrice float64 `json:"optimized_price"`
	Status         string  `json:"status"`
}

type Request struct {
	ID          string  `json:"id"`
	Price       float64 `json:"price"`
	FloorPrice  float64 `json:"floor_price"`
	DC          string  `json:"data_center"`
	PublisherID string  `json:"app_publisher_id"`
	BundleID    string  `json:"bundle_id"`
	TagID       string  `json:"tag_id"`
	GeoCountry  string  `json:"device_geo_country"`
	AdFormat    string  `json:"ext_ad_format"`
}

type FeedBackRequest struct {
	ID         string  `json:"id"`
	Impression bool    `json:"impression"`
	Price      float64 `json:"price"`
}

type FeedBackResponse struct {
	Ack bool `json:"ack"`
}
