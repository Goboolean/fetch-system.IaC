package model



type TradeJson struct {
	Price     float64 `json:"price,omitempty"`
	Size      int64   `json:"size,omitempty"`
	Timestamp int64   `json:"timestamp,omitempty"`
}

type AggregateJson struct {
	Open      float64 `json:"open,omitempty"`
	Close     float64 `json:"close,omitempty"`
	High      float64 `json:"high,omitempty"`
	Low       float64 `json:"low,omitempty"`
	Volume    int64   `json:"volume,omitempty"`
	Timestamp int64   `json:"timestamp"`
}