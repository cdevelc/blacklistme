package services

import "q29"
import "blacklistme/controller/payment"


func Index(q *q29.ReqRsp) {
	var page struct {
		Vw q29.View
		Tab1 bool
		Tab2 bool
		Tab3 bool
		Tab4 bool
		Tab5 bool
		Prices *[]payment.Price
		Stripe_api_key string
	}
	switch (q.Action) {
	default: fallthrough
	case "emailblacklist":   page.Tab1 = true
	case "domainblacklist":  page.Tab2 = true
	case "privateblacklist": page.Tab3 = true
	case "api":              page.Tab4 = true
	case "pricing":          page.Tab5 = true		
	}
	page.Prices = &payment.Prices
	page.Stripe_api_key = payment.Stripe_api_key
	page.Vw.Template = "services/index"
	q29.Render(q, &page)
}

