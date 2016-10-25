package ulist

import "fmt"
import "time"
import "q29"


func BeforeFilter(q *q29.ReqRsp) bool {
	return true
}

func Dashboard(q *q29.ReqRsp) {
	var page struct {
		Vw q29.View
		LastLoginTime string
	}
	t, _ := time.Parse("2006-01-02 15:04:05", q.U.LastLoginTime)
	page.LastLoginTime = fmt.Sprintf(t.Format("January 2, 2006 3:04PM"))
	page.Vw.Template = "ulist/dashboard"
	q29.Render(q, &page)
}

func Profile(q *q29.ReqRsp) {
	var page struct {
		Vw q29.View
		MemberSince string
	}
	t, _ := time.Parse("2006-01-02 15:04:05", q.U.Created)
	page.MemberSince = fmt.Sprintf(t.Format("January 2, 2006"))
	page.Vw.Template = "ulist/profile"
	q29.Render(q, &page)
}
