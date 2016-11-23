package blist

import "net/http"
import "strings"
import "strconv"
import "q29"
import "q29/validfield"
import "blacklistme/model/emaddr"
import "blacklistme/util/mailgo"
import "bytes"
import "io"

func BeforeFilter(q *q29.ReqRsp) bool {
	return true
}

func Index(q *q29.ReqRsp) {
	var page struct {
		Vw q29.View
	}
	page.Vw.Template = "blist/index"
	q29.Render(q, &page)
}

func Inquire(q *q29.ReqRsp) {
	var page struct {
		Email       string
		Found       bool		
		EmAddr      emaddr.Emaddr
		Vw          q29.View
	}
	page.Email = q.R.URL.Query().Get("email")
	emsg := validfield.Email(validfield.F{ "Email address", page.Email, 0, 0, true})
	if emsg != "" {
		http.Error(q.W, q.R.URL.Path+" invalid email address", 404)
		return
	}
	page.Found = emaddr.Find(q.M, "blacklist", page.Email, &page.EmAddr)
	q29.Render(q, &page)	
}

func AddRem(q *q29.ReqRsp) {
	var page struct {
		Email   string
		AddRem  string
		Vw      q29.View
	}
	var found   bool		
	var emAddr  emaddr.Emaddr

	q.R.ParseForm()
	page.Email = q.R.FormValue("email")
	emsg := validfield.Email(validfield.F{ "Email address", page.Email, 0, 0, true})
	if emsg != "" {
		http.Error(q.W, q.R.URL.Path+" invalid email address", 404)
		return
	}
	found = emaddr.Find(q.M, "blacklist", page.Email, &emAddr)

	if found == true {
		page.AddRem = "remove"
		s := mailgo.Session {
			Email: page.Email,
			URL: "http://"+q.R.Host+q29.AssetURL(q, "blist/complete/remove?vps=")+emAddr.Sha256,
		}
		mailgo.ConfirmEmailAddressUnBlacklist(&s)
		
	} else {
		page.AddRem = "add"
		emAddr.Email = page.Email;		
		emaddr.Upsert(q.M, "blacklistwannabe", &emAddr);		
		s := mailgo.Session {
			Email: page.Email,
			URL: "http://"+q.R.Host+q29.AssetURL(q, "blist/complete?vps=")+emAddr.Sha256,
		}
		mailgo.ConfirmEmailAddressBlacklist(&s)		

	}
	page.Vw.Template = "blist/request"
	q29.Render(q, &page)
}

func Complete(q *q29.ReqRsp) {
	var page struct {
		Remove bool
		Email string
		Vw q29.View
	}
	var emAddr emaddr.Emaddr
	var found bool

	sig := q.R.URL.Query().Get("vps")
	chunks := strings.Split(q.R.URL.Path, "/")
	found = false

	if len(sig) > 0 {
		if chunks[len(chunks)-1] == "remove" {
			page.Remove = true
			found = emaddr.FindBySig(q.M, "blacklist", sig, &emAddr)
			if found == true {
				page.Email = emAddr.Email				
				emaddr.Delete(q.M, "blacklist", emAddr.Id)
			}
			
		} else { /* add request confirmed */
			page.Remove = false
			found = emaddr.FindBySig(q.M, "blacklistwannabe", sig, &emAddr)
			if found == true {
				page.Email = emAddr.Email
				emaddr.Delete(q.M, "blacklistwannabe", emAddr.Id)
				emAddr.Id = ""
				emaddr.Upsert(q.M, "blacklist", &emAddr)
			}
		}
	}
	if found == false {
		http.Error(q.W, q.R.URL.Path+" invalid confirmation", 404)
		return
	}
	q29.Render(q, &page)
}

func Upload(q *q29.ReqRsp) {
	var page struct {
		Vw q29.View
	}
	var emAddr  emaddr.Emaddr
	
	if q.R.Method == "POST" {
		q.R.ParseMultipartForm(256*1024)
		f, _, err := q.R.FormFile("emfile")
		if err != nil {
			page.Vw.FlashMsg = "emfile upload error: "+err.Error()
		} else {
			var b bytes.Buffer
			_, err = io.Copy(&b, f)
			if err != nil {
				page.Vw.FlashMsg = "emfile upload copy error: "+err.Error()
			} else {
				page.Vw.FlashMsg = "emfile upload complete"
				addrs := strings.Split(b.String(),"\n")
				page.Vw.FlashMsg = "emfile upload complete, "+strconv.Itoa(len(addrs))+" spammy email addresses processed"
				for i:=0; i<len(addrs); i = i + 1 {
					em := strings.TrimSpace(addrs[i])
					if len(em) > 0 {
						found := emaddr.Find(q.M, "blacklist", em, &emAddr)
						if found == false {
							emAddr.Id = ""
							emAddr.Email = em
							emaddr.Upsert(q.M, "blacklist", &emAddr);									
						}
					}					
				}				
			}
		}
	}
	q29.Render(q, &page)
}
