package dlist

//import "fmt"
//import "time"
import "q29"
import "blacklistme/model/domain"

func BeforeFilter(q *q29.ReqRsp) bool {
	if q.U != nil { return true }
	q29.Redirect(q, "account/login")
	return false
}

func Dlist(q *q29.ReqRsp) {
	var page struct {
		Vw q29.View
		Dlist []domain.Domain
		DlistCount int
		FlashMsg string
	}
	domain.ListByUid(q.M, q.U.Id, &page.Dlist)
	page.DlistCount = len(page.Dlist)
	page.Vw.Template = "dlist/dlist"
	q29.Render(q, &page)
}
