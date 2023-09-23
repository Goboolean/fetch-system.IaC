package mongo



type Aggregate struct {
	Timestamp int64   `bson:"timestamp"`
	Open      float64 `bson:"open"`
	Closed    float64 `bson:"closed"`
	Min       float64 `bson:"min"`
	Max       float64 `bson:"max"`
	Volume    int64   `bson:"volume"`
}