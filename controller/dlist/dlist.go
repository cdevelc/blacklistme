package dlist

//import "log"
//import "time"
import "q29"
import "q29/validfield"
import "blacklistme/model/domain"

func BeforeFilter(q *q29.ReqRsp) bool {
	if q.U != nil { return true }
	q29.Redirect(q, "account/login")
	return false
}

func Index(q *q29.ReqRsp) {
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

func Add(q *q29.ReqRsp) {
	var dm domain.Domain

	q.R.ParseForm()
	domname := q.R.FormValue("domain")
	emsg := validfield.Domain(validfield.F{"Domain Name", domname, 0, 0, true})
	if emsg != "" {
		q29.SetFlash(q, "The domain name "+domname+" is invalid.")
		q29.Redirect(q, "dlist/index")
		return
	}
	found := domain.Find(q.M, domname, &dm)
	if found == true && dm.UserId != q.U.Id {
		q29.SetFlash(q, "The domain "+domname+" is already under BlackList control.")				
		q29.Redirect(q, "dlist/index")
		return
	}
	if found == false {
		dm.Id = ""
		dm.Domain = domname
		dm.UserId = q.U.Id
		domain.Upsert(q.M, &dm)
	}
	q29.SetFlash(q, "The domain "+domname+" has been added.")	
	q29.Redirect(q, "dlist/index")
}

func Del(q *q29.ReqRsp) {
	var dm domain.Domain

	q.R.ParseForm()
	domname := q.R.FormValue("domain")
	emsg := validfield.Domain(validfield.F{"Domain Name", domname, 0, 0, true})
	if emsg == "" {
		found := domain.Find(q.M, domname, &dm)
		if found == true && dm.UserId == q.U.Id {
			domain.Delete(q.M, dm.Id)
		}
	}
	q29.SetFlash(q, "The domain "+domname+" has been removed.")
	q29.Redirect(q, "dlist/index")
}
