package influx

import "time"

// TODO: 기술적 요구에 맞게 수정
type OrderEvent struct {
	ProductID         string
	ProportionPercent int    // Proportion is the target percentage of the order
	Action            string // Action to be performed for the order Buy or Sell
	Task              string
	CreatedAt         time.Time
}
