/*
 */
var erro_count = 0;
var erro = {
 username:    { str: '', lab: '' },
 password:    { str: '', lab: '' },
 passagain:   { str: '', lab: '' },
 oldpassword: { str: '', lab: '' },    
 email:       { str: '', lab: '' },
 firstname:   { str: '', lab: '' },
 lastname:    { str: '', lab: '' },
 domain:      { str: '', lab: '' }
};
function erroSet(k, s, l) {
 erro_count++;    
 erro[k].str = s;
 erro[k].lab = l;
}

var formfields = {
 register:   ['username', 'email', 'password', 'passagain', 'firstname', 'lastname'],
 cpassword:  ['password', 'passagain', 'oldpassword'],
 cemail:     ['email', 'password'],
 cname:      ['firstname', 'lastname', 'password'],
 login:      ['username', 'password'],
 email:      ['email'],
 domainform: ['domain']
};

function validateForm(fname) {
 erro_count = 0;
 for (var i in formfields[fname]) {
  var key = formfields[fname][i];
  var val = document.forms[fname][key].value;
  erro[key].str = '';
  erro[key].lab = '';
  switch (key) {
  case 'username': validate_username( key, val); break;
  case 'password': validate_password( key, val); break;
  case 'passagain': validate_passagain( key, document.forms[fname]['password'].value, val); break;
  case 'oldpassword': validate_password( key, val); break;      
  case 'email': validate_email( key, val); break;
  case 'firstname': validate_humanname( key, val, "First"); break;
  case 'lastname': validate_humanname( key, val, "Last"); break;
  case 'domain': validate_domainname( key, val); break;
  default: return false; break;
  }
  document.getElementById(key+".error").innerHTML = erro[key].str;
  document.getElementById(key+".errlabel").innerHTML = erro[key].lab;
 }
 if (erro_count) { return false; }
 if (document.forms[fname]['password'].value != undefined) {
   var st = document.forms[fname]['stateToken'].value;        
   var pw = document.forms[fname]['password'].value;
   var ve = vernam_encipher(pw, st);
   document.forms[fname]['password'].value = ve;
   if (document.forms[fname]['passagain'].value != undefined) {
     document.forms[fname]['passagain'].value = ve;
   }
   if (document.forms[fname]['oldpassword'].value != undefined) {
     var op = document.forms[fname]['oldpassword'].value;
     ve = vernam_encipher(op, st);         
     document.forms[fname]['oldpassword'].value = ve;
   }
   document.forms[fname]['stateToken'].value = "";
 }
 return true;    
}

function validate_username(k, u) {
 if (!u || u.length == 0)
  erroSet(k, "", "required");
 else if (u.length < 3) 
  erroSet(k, "Username is too short","invalid");
 else if (u.length > 20) 
  erroSet(k, "Username is too long", "invalid");
 else if (!u.match(/^[a-zA-Z0-9_.]+$/)) 
  erroSet(k, "Username contains invalid characters", "invalid");
}

function validate_humanname(k, h, labl) {
 if (!h || h.length == 0)
  erroSet(k, "", "required");
 else if (h.length < 1) 
  erroSet(k, labl+" name is too short","invalid");
 else if (h.length > 32) 
  erroSet(k, labl+" name is too long", "invalid");
 else if (!h.match(/^[a-zA-Z][0-9a-zA-Z .,'-]*$/)) 
  erroSet(k, labl+" name contains invalid characters", "invalid");
}

function validate_email(k, e) {
 if (!e || e.length == 0)
  erroSet(k, "", "required");
 else if (e.length < 5)
  erroSet(k, "Email address is too short", "invalid");
 else if (e.length > 64)
  erroSet(k, "Email address is too long", "invalid");
 else if (!e.match(/^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$/))
  erroSet(k, "Email address contains invalid characters", "invalid");
}

function validate_domainname(k, e) {
 if (!e || e.length == 0)
  erroSet(k, "", "required");
 else if (!e.match(/^((?!-))(xn--)?[a-z0-9][a-z0-9-_]{0,61}[a-z0-9]{0,1}\.(xn--)?([a-z0-9\-]{1,61}|[a-z0-9-]{1,30}\.[a-z]{2,})$/))
  erroSet(k, "Domain contains unexpected characters or invalid formatting", "invalid");
}
    
function validate_password(k, p) {
 if (!p || p.length == 0)
  erroSet(k, "", "required");
 else if (p.length < 8) 
  erroSet(k, "Password is too short", "invalid");
 else if (p.length > 128) 
  erroSet(k, "Password is too long", "invalid");
 else if (!p.match(/^[a-zA-Z0-9!"#$%&'()*+,.\/:;<=>?@\[\] ^_{|}~-]+$/))
  erroSet(k, "Password contains invalid characters", "invalid");
}

function validate_passagain(k, pw, pa) {
 if (!pa || pa.length == 0)
  erroSet(k, "", "required");
 else if (pw != pa)
  erroSet(k, "The 2 passwords do not match", "mismatch");
}

function vernam_encipher(p,k) {
 var r="", ax, kx;
 for (var i=0; i<p.length; i++) { 
   kx = (i*2) % (k.length);
   ax = (parseInt(p.charCodeAt(i)) ^ parseInt(k.substr(kx,2),16)) & 0xff;
   r += (ax < 16)? '0'+ax.toString(16): ax.toString(16);
 }
 return r;
}

function vernam_decipher(q,x) {
 var r="";
   for (var i=0; i<q.length; i+=2) { 
     r += String.fromCharCode(parseInt(q.substr(i,2), 16) ^ parseInt(x.substr(i,2),16));
   }
 return r;
}


$(document).ready(function() {
  $('input').focus(function() {
    var key=$(this).attr("name");
    document.getElementById(key+".error").innerHTML = "";
    document.getElementById(key+".errlabel").innerHTML = "";        
  });
});

