package mongo


type Trade struct {
	Price     float64 `bson:"price"`
	Size      int64   `bson:"size"`
	Timestamp int64   `bson:"timestamp"`
}


type Aggregate struct {
	Open      float64 `bson:"open"`
	Closed    float64 `bson:"closed"`
	Min       float64 `bson:"min"`
	Max       float64 `bson:"max"`
	Volume    int64   `bson:"volume"`
	Timestamp int64   `bson:"timestamp"`
}