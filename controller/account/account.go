package account

import "q29"
import "q29/session"
import "q29/user"
import "blacklistme/controller/account/validate"
import "blacklistme/model/apikey"
import "blacklistme/util/mailgo"
import "fmt"
import "math/rand"

type TemplateVars struct {
	Vw q29.View
	Av validate.AccessVars
	RdoUsername string
	RdoEmail string
}

func verify_user_create_session_and_redirect(q *q29.ReqRsp, uname string, pword string, uri string) (error int) {
	var u *user.User
	u = user.FindByUname(q.M, uname)
	if u != nil {
		if u.Confirmed == false {
			return 2
		}
		encpw := session.EncryptPassword(pword, u.Passsalt)
		if encpw == u.Password {
			session.Create(q.M, q.W, q.Base, u.Username, u.Email)
			user.StampLogin(q.M, u)
			q29.Redirect(q, uri)
			return 0
		}
	}
	return 1
}

func BeforeFilter(q *q29.ReqRsp) bool {
	return true
}

func Login(q *q29.ReqRsp) {
	var templateVars TemplateVars

	if q.R.Method == "POST" {
		validate.Login(q, &templateVars.Av)
		if templateVars.Av.Error.Count == 0 {
			err := verify_user_create_session_and_redirect(q, templateVars.Av.Username, templateVars.Av.Password, "ulist/dashboard")
			if err == 0 {
				return
			}
			if err == 1 {
				templateVars.Av.Password = ""
				templateVars.Av.ErrorLabel.Password = "incorrect"			
				templateVars.Av.Error.Password = "Sorry, that password was incorrect"
				templateVars.Av.Error.Count++
			} else {
				templateVars.Av.Password = ""
				templateVars.Av.ErrorLabel.Username = "pending"			
				templateVars.Av.Error.Username = "This account is awaitng email confirmation"
				templateVars.Av.Error.Count++
			}
		}
	}
	if templateVars.Av.StateToken == "" {
		/* sometimes we get here from sign in form POST with no stateToken or GET */
		templateVars.Av.StateToken = session.AllocateClientStateToken(q.M, q29.RemoteIP(q))
	}
	q29.Render(q, &templateVars)	
}

func Logout(q *q29.ReqRsp) {
	session.Destroy(q.M, q.R, q.W)	
	q29.Redirect(q, "/")
}

func Register(q *q29.ReqRsp) {
	var templateVars TemplateVars
	var u user.User
	
	if q.R.Method == "GET" {
		templateVars.Av.StateToken = session.AllocateClientStateToken(q.M, q29.RemoteIP(q))
		q29.Render(q, &templateVars)
		return
	}

 /* request == POST */
	validate.Register(q, &templateVars.Av)
	if templateVars.Av.Error.Count != 0 {
		q29.Render(q, &templateVars)
		return
	}
	u.Username = templateVars.Av.Username
	u.Email    = templateVars.Av.Email
	u.Passsalt = session.ShakeSalt(u.Email)
	u.Password = session.EncryptPassword(templateVars.Av.Password, u.Passsalt)
	u.Firstname = templateVars.Av.Firstname
	u.Lastname  = templateVars.Av.Lastname
	u.Confirmed = false
	s := mailgo.Session {
		Fname: u.Firstname,
		Lname: u.Lastname,
		Email: u.Email,
		URL: "http://"+q.R.Host+q29.AssetURL(q, "account/confirm?vps=")+u.Passsalt,
	}
	mailgo.ConfirmRegistration(&s)
	user.Add(q.M, &u)
	q29.Redirect(q, "account/thanks")	
}

func Thanks(q *q29.ReqRsp) {
	var page struct {
		Vw q29.View
	}
	q29.Render(q, &page)
}

func Confirm(q *q29.ReqRsp) {
	var page struct {
		Vw q29.View
	}
	var u *user.User
	var a apikey.Apikey
	
	u = user.FindByPasssalt(q.M, q.R.URL.Query().Get("vps"))
	if u != nil {
		u.Confirmed = true
		user.Update(q.M, u)
		a.UserId = u.Id
		apikey.Upsert(q.M, &a)
	}
	q29.Render(q, &page)			
}

func Password(q *q29.ReqRsp) {
	var templateVars TemplateVars
	var u *user.User

	if q.R.Method == "POST" {
		validate.ChangePassword(q, &templateVars.Av)

		if templateVars.Av.Error.Count == 0 {
			u = user.FindByUname(q.M, q.U.Username)
			if u != nil {
				encpw := session.EncryptPassword(templateVars.Av.Oldpassword, u.Passsalt)
				if encpw == u.Password {
					u.Passsalt = session.ShakeSalt(u.Email)
					u.Password = session.EncryptPassword(templateVars.Av.Password, u.Passsalt)
					user.Update(q.M, u)
					session.Destroy(q.M, q.R, q.W)
					s := mailgo.Session { Fname: u.Firstname, Lname: u.Lastname, Email: u.Email,}
					mailgo.NotifyPasswordChange(&s)
					
					verify_user_create_session_and_redirect(q, q.U.Username, templateVars.Av.Password, "ulist/profile")
					return
				}
				templateVars.Av.Oldpassword = ""
				templateVars.Av.ErrorLabel.Oldpassword = "incorrect"			
				templateVars.Av.Error.Oldpassword = "Sorry, that password was incorrect"
				templateVars.Av.Error.Count++
			}
		}
	}
	if templateVars.Av.StateToken == "" {
		templateVars.Av.StateToken = session.AllocateClientStateToken(q.M, q29.RemoteIP(q))
	}
	q29.Render(q, &templateVars)		
}

func Email(q *q29.ReqRsp) {
	var templateVars TemplateVars
	var u *user.User
	var oldEmail string

	if q.R.Method == "POST" {
		validate.ChangeEmail(q, &templateVars.Av)

		if templateVars.Av.Error.Count == 0 {
			u = user.FindByUname(q.M, q.U.Username)
			if u != nil {
				encpw := session.EncryptPassword(templateVars.Av.Password, u.Passsalt)
				if encpw == u.Password {
					oldEmail = u.Email
					u.Email  = templateVars.Av.Email
					user.Update(q.M, u)
					session.Destroy(q.M, q.R, q.W)
					s := mailgo.Session { Fname: u.Firstname, Lname: u.Lastname, Email: u.Email,}
					mailgo.NotifyEmailAddressChange(&s, oldEmail)
					verify_user_create_session_and_redirect(q, q.U.Username, templateVars.Av.Password, "ulist/profile")
					return
				}
				templateVars.Av.Password = ""
				templateVars.Av.ErrorLabel.Password = "incorrect"			
				templateVars.Av.Error.Password = "Sorry, that password was incorrect"
				templateVars.Av.Error.Count++
			}
		}
	}
	if templateVars.Av.StateToken == "" {
		templateVars.Av.StateToken = session.AllocateClientStateToken(q.M, q29.RemoteIP(q))
	}
	q29.Render(q, &templateVars)		
}

func Rename(q *q29.ReqRsp) {
	var templateVars TemplateVars
	var u *user.User

	if q.R.Method == "POST" {
		validate.ChangeName(q, &templateVars.Av)
		if templateVars.Av.Error.Count == 0 {
			u = user.FindByUname(q.M, q.U.Username)
			if u != nil {
				encpw := session.EncryptPassword(templateVars.Av.Password, u.Passsalt)
				if encpw == u.Password {
					u.Firstname = templateVars.Av.Firstname
					u.Lastname  = templateVars.Av.Lastname
					user.Update(q.M, u)
					q29.Redirect(q, "ulist/profile")
					return
				}
				templateVars.Av.Password = ""
				templateVars.Av.ErrorLabel.Password = "incorrect"			
				templateVars.Av.Error.Password = "Sorry, that password was incorrect"
				templateVars.Av.Error.Count++
			}
		}
	}
	if templateVars.Av.StateToken == "" {
		templateVars.Av.StateToken = session.AllocateClientStateToken(q.M, q29.RemoteIP(q))
	}
	q29.Render(q, &templateVars)		
}

func randletter(length int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZa"
	var result string
	for i:=0; i<length; i++ {		
		result = result + fmt.Sprintf("%c", letterBytes[rand.Intn(52)])
	}
	return result
}

func Forgot(q *q29.ReqRsp) {
	var templateVars TemplateVars

	templateVars.RdoUsername = "checked"
	templateVars.RdoEmail = ""
	if q.R.Method == "POST" {
		validate.Forgot(q, &templateVars.Av)
		if templateVars.Av.Error.Count == 0 {

			if templateVars.Av.Radio1 == "email" {
				s := mailgo.Session {
					Email: templateVars.Av.Email,
					Username: templateVars.Av.Username,
				}
				mailgo.UsernameReminderEmail(&s)
				
			} else { //username argument provided, reset user password
				u := user.FindByUname(q.M, templateVars.Av.Username)
				if u != nil {

					pw := "Bm2017"+randletter(4)
					u.Password = session.EncryptPassword(pw, u.Passsalt)
					user.Update(q.M, u)			
					s := mailgo.Session {
						Email: u.Email,
						Username: templateVars.Av.Username,
						Password: pw,
					}
					mailgo.PasswordResetEmail(&s)
				}
			}			
			q29.Redirect(q, "account/recovery")
			return
		}
		if templateVars.Av.ErrorLabel.Email != "" {
			templateVars.RdoUsername = ""
			templateVars.RdoEmail = "checked"
		}
	}
	q29.Render(q, &templateVars)		
}
func Recovery(q *q29.ReqRsp) {
	var page struct {
		Vw q29.View
	}
	q29.Render(q, &page)
}
