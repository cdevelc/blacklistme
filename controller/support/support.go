package support

import "q29"

func BeforeFilter(q *q29.ReqRsp) bool {
	return true
}

func FAQ(q *q29.ReqRsp) {
	var page struct {
		Vw q29.View
	}
	q29.Render(q, &page)		
}
