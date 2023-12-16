package model


type TickerDetail struct {
	Ticker      string
	Name        string
	Description string
	Exchange    string
}


type TickerDetailResult struct {
	TickerDetail
	Status  string
	Message string
}