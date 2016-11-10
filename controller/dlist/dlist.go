package dlist

//import "log"
//import "time"
import "q29"
import "q29/validfield"
import "blacklistme/model/domain"
import "blacklistme/model/emaddr"

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
			//need a lot of code here to delete email addresses in the list before domain is deleted.
			domain.Delete(q.M, dm.Id)
		}
	}
	q29.SetFlash(q, "The domain "+domname+" has been removed.")
	q29.Redirect(q, "dlist/index")
}

func Elist(q *q29.ReqRsp) {
	var page struct {
		Vw q29.View
		Dname string
		Domain domain.Domain
		Elist []emaddr.Emaddr
		ElistCount int
		FlashMsg string
	}
	page.Dname = q.R.URL.Query().Get("dname")
	domain.Find(q.M, page.Dname, &page.Domain)
	emaddr.ListByDid(q.M, "domainaddrs", page.Domain.Id, &page.Elist)
	page.ElistCount = len(page.Elist)
	q29.Render(q, &page)
}

func ElistAdd(q *q29.ReqRsp) {
	var emAddr emaddr.Emaddr
	var dm domain.Domain
	q.R.ParseForm()
	dname := q.R.FormValue("dname")
  email := q.R.FormValue("email")
	emsg := validfield.Email(validfield.F{"Email address", email, 0, 0, true})
	if emsg != "" {
		q29.SetFlash(q, "That email address was invalid.")
	} else {
		domain.Find(q.M, dname, &dm)
		found := emaddr.FindByDid(q.M, "domainaddrs", dm.Id, email, &emAddr)
		if found == false {
			emAddr.Id = ""
			emAddr.Email = email
			emAddr.UserId = q.U.Id
			emAddr.DomainId = dm.Id
			emaddr.Upsert(q.M, "domainaddrs", &emAddr)
		}
		q29.SetFlash(q, "The email address "+email+" has been added to this Domain BlackList.")
	}
	q29.Redirect(q, "dlist/elist?dname="+dname)
}

func ElistDel(q *q29.ReqRsp) {
	var emAddr emaddr.Emaddr
	var dm domain.Domain	
	q.R.ParseForm()
	dname := q.R.FormValue("dname")
  email := q.R.FormValue("email")
	emsg := validfield.Email(validfield.F{"Email address", email, 0, 0, true})
	if emsg != "" {
		q29.SetFlash(q, "That email address was invalid.")		
	} else {
		domain.Find(q.M, dname, &dm)		
		found := emaddr.FindByDid(q.M, "domainaddrs", dm.Id, email, &emAddr)
		if found == false {
			q29.SetFlash(q, "That email address was not found.");
		} else {
			emaddr.Delete(q.M, "domainaddrs", emAddr.Id)
			q29.SetFlash(q, "The email address "+email+" has been removed from this Domain BlackList.")		
		}		
	}
	q29.Redirect(q, "dlist/elist?dname="+dname)
}
