package validate

import "strconv"
import "q29"
import "q29/user"
import "q29/session"
import "q29/validfield"

type AccessVars struct {
	Username   string
	Email      string
	Password   string
	Passagain  string
	Oldpassword string	
	Firstname  string
	Lastname   string
	StateToken string
	Error struct {
		Count     int
		Username  string
		Email     string
		Password  string
		Passagain string
		Oldpassword string			
		Firstname string
		Lastname  string
	}
	ErrorLabel struct {
		Username  string
		Email     string
		Password  string
		Passagain string
		Oldpassword string					
		Firstname string
		Lastname  string
	}
}

func vusername(av *AccessVars) {
	emsg := validfield.Username(validfield.F{"Username", av.Username, 0, 0, true})
	if emsg != "" {
		av.Error.Username = emsg
		av.ErrorLabel.Username = "invalid"
		av.Error.Count++
	}
}
func vemail(av *AccessVars) {
	emsg := validfield.Email(validfield.F{"Email address", av.Email, 0, 0, true})
	if emsg != "" {
		av.Error.Email = emsg
		av.ErrorLabel.Email = "invalid"
		av.Error.Count++
	}
}
func vpassword(av *AccessVars) {
	emsg := validfield.Password(validfield.F{"Password", av.Password, 0, 0, true})
	if emsg != "" {
		av.Error.Password = emsg
		av.ErrorLabel.Password = "invalid"		
		av.Error.Count++
	}
}
func voldpassword(av *AccessVars) {
	emsg := validfield.Password(validfield.F{"Password", av.Oldpassword, 0, 0, true})
	if emsg != "" {
		av.Error.Oldpassword = emsg
		av.ErrorLabel.Oldpassword = "invalid"		
		av.Error.Count++
	}
}

func vpassagain(av *AccessVars) {
	if av.Password != av.Passagain {
		av.Error.Passagain = "The 2 passwords do not match"
		av.ErrorLabel.Passagain = "invalid"
		av.Error.Count++
	}
}

func vernam_decipher(q string, x string) (r string) {
	var a, b int64
	r = ""
	for i := 0; i < len(q); i = i + 2 {
		a, _ = strconv.ParseInt(q[i:i+2], 16, 16)
		b, _ = strconv.ParseInt(x[i:i+2], 16, 16)
		r = r + string(a ^ b)
	}
	return r
}

func initAVs(q *q29.ReqRsp, av *AccessVars) {
	q.R.ParseForm()
	av.Username = q.R.FormValue("username")
	av.Email = q.R.FormValue("email")
	av.Password = q.R.FormValue("password")
	av.Passagain = q.R.FormValue("passagain")
	av.Oldpassword = q.R.FormValue("oldpassword")	
	av.Firstname = q.R.FormValue("firstname")
	av.Lastname = q.R.FormValue("lastname")
	_, stateTokenPresent := q.R.Form["stateToken"]
	if stateTokenPresent == true {
		av.StateToken = session.RetrieveClientStateToken(q.M, q29.RemoteIP(q))
		av.Password = vernam_decipher(av.Password, av.StateToken)
		av.Passagain = vernam_decipher(av.Passagain, av.StateToken)
		av.Oldpassword = vernam_decipher(av.Oldpassword, av.StateToken)
	}
}

func Register(q *q29.ReqRsp, av *AccessVars) {
	initAVs(q, av)
	vusername(av)
	vemail(av)
	vpassword(av)
	vpassagain(av)
	if av.Error.Count == 0 {
		u := user.FindByUname(q.M, av.Username)
		if u != nil {
			av.Error.Username = "Username is not available"
			av.ErrorLabel.Username = "in use"			
			av.Error.Count++
		}
		u = user.FindByEmail(q.M, av.Email)
		if u != nil {
			av.Error.Email = "Email address already has account"
			av.ErrorLabel.Email = "in use"			
			av.Error.Count++
		}
	}
}

func ChangePassword(q *q29.ReqRsp, av *AccessVars) {
	initAVs(q, av)
	vpassword(av)
	vpassagain(av)
	voldpassword(av)
}

func ChangeEmail(q *q29.ReqRsp, av *AccessVars) {
	initAVs(q, av)
	vpassword(av)	
	vemail(av)
}

func ChangeName(q *q29.ReqRsp, av *AccessVars) {
	initAVs(q, av)
	vpassword(av)	
}

func Login(q *q29.ReqRsp, av *AccessVars) {
	initAVs(q, av)
	vusername(av)
	vpassword(av)
	if av.Error.Count == 0 {
		u := user.FindByUname(q.M, av.Username)
		if u==nil {
			av.Error.Username = "Username is unknown"
			av.ErrorLabel.Username = "not found"			
			av.Error.Count++
		}
	}
}

