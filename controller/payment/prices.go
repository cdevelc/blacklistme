package payment

type Price struct {
	UserCount int
	DisplayCharge string
	StripeCharge int
}

var Prices = []Price {
	{    10,     "$24",    2400},
	{   100,    "$240",   24000},
	{  1000,  "$2,400",  240000},
	{ 10000, "$24,000", 2400000},
}
