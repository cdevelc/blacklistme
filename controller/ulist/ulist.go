package ulist

import "fmt"
import "time"
import "q29"
import "q29/validfield"
import "blacklistme/model/apikey"
import "blacklistme/model/emaddr"
import "blacklistme/model/domain"

func BeforeFilter(q *q29.ReqRsp) bool {
	return true
}

func Dashboard(q *q29.ReqRsp) {
	var page struct {
		Vw q29.View
		LastLoginTime string
		Apikey string
		Plistinfo string
		Dlistinfo string
	}
	t, _ := time.Parse("2006-01-02 15:04:05", q.U.LastLoginTime)
	page.LastLoginTime = fmt.Sprintf(t.Format("January 2, 2006 3:04PM"))

	var apk apikey.Apikey
	apikey.FindByUserId(q.M, q.U.Id, &apk)
	page.Apikey = apk.APIkey
	plistcount := emaddr.ListByUidCount(q.M, "blacklistprivate", q.U.Id)
	if plistcount == 0 {
		page.Plistinfo = "none"
	} else {
		page.Plistinfo = fmt.Sprintf("%d email address entries", plistcount)
	}
	dlistcount := domain.ListByUidCount(q.M, q.U.Id)	
	if dlistcount == 0 {
		page.Dlistinfo = "none"
	} else {
		page.Dlistinfo = fmt.Sprintf("%d domain lists", dlistcount)
	}
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

func Plist(q *q29.ReqRsp) {
	var page struct {
		Vw q29.View
		Plist []emaddr.Emaddr
		PlistCount int
		FlashMsg string
	}
	emaddr.ListByUid(q.M, "blacklistprivate", q.U.Id, &page.Plist)
	page.PlistCount = len(page.Plist)
	q29.Render(q, &page)
}

func PlistAdd(q *q29.ReqRsp) {
	var emAddr emaddr.Emaddr
	q.R.ParseForm()
  email := q.R.FormValue("email")
	emsg := validfield.Email(validfield.F{"Email address", email, 0, 0, true})
	if emsg != "" {
		q29.SetFlash(q, "That email address was invalid.")
	} else {
		found := emaddr.FindByUid(q.M, "blacklistprivate", q.U.Id, email, &emAddr)
		if found == false {
			emAddr.Id = ""
			emAddr.Email = email
			emAddr.UserId = q.U.Id
			emaddr.Upsert(q.M, "blacklistprivate", &emAddr)
		}
		q29.SetFlash(q, "The email address "+email+" has been added to your Private blacklist.")
	}
	q29.Redirect(q, "ulist/plist")
}

func PlistDel(q *q29.ReqRsp) {
	var emAddr emaddr.Emaddr	
	q.R.ParseForm()
  email := q.R.FormValue("email")
	emsg := validfield.Email(validfield.F{"Email address", email, 0, 0, true})
	if emsg != "" {
		q29.SetFlash(q, "That email address was invalid.")		
	} else {
		found := emaddr.FindByUid(q.M, "blacklistprivate", q.U.Id, email, &emAddr)
		if found == false {
			q29.SetFlash(q, "That email address was not found.");
		} else {
			emaddr.Delete(q.M, "blacklistprivate", emAddr.Id)
			q29.SetFlash(q, "The email address "+email+" has been removed from your Private blacklist.")		
		}		
	}
	q29.Redirect(q, "ulist/plist")
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
	q29.Render(q, &page)
}
