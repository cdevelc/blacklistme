package payment

import "q29"

func BeforeFilter(q *q29.ReqRsp) bool {
	if q.U != nil { return true }
	q29.Redirect(q, "account/login")
	return false
}

func Enroll(q *q29.ReqRsp) {
	q29.Redirect(q, "payment/confirm")
}

func Confirm(q *q29.ReqRsp) {
	var page struct {
		Vw q29.View
	}
	q29.Render(q, &page)
}
