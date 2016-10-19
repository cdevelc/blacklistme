package ulist

//import "net/http"
//import "strings"
import "q29"
//import "q29/user"

func BeforeFilter(q *q29.ReqRsp) bool {
	return true
}

func Dashboard(q *q29.ReqRsp) {
	var page struct {
		Vw q29.View
	}
	page.Vw.Template = "ulist/dashboard"
	q29.Render(q, &page)
}

func Profile(q *q29.ReqRsp) {
	var page struct {
		Vw q29.View
	}
	page.Vw.Template = "ulist/profile"
	q29.Render(q, &page)
}
