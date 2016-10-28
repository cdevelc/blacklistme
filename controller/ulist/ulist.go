package ulist

import "fmt"
import "time"
import "q29"
import "blacklistme/model/apikey"

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

func Apikey(q *q29.ReqRsp) {
	var page struct {
		Vw q29.View
		Apikey string
		Created string
	}
	var apk apikey.Apikey
	apikey.FindByUserId(q.M, q.U.Id, &apk)

	t, _ := time.Parse("2006-01-02 15:04:05", apk.Created)	
	page.Created = fmt.Sprintf(t.Format("January 2, 2006"))
	page.Apikey = apk.APIkey
	
	q29.Render(q, &page)
}

func ApikeyRegen(q *q29.ReqRsp) {
	var apk apikey.Apikey
	if q.U != nil {
		apikey.FindByUserId(q.M, q.U.Id, &apk)
		apk.UserId = q.U.Id
		apikey.Upsert(q.M, &apk)
	}
	q29.Redirect(q, "ulist/apikey")
}
