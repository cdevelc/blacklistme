

function emTableInit() { 
   var tb = document.getElementById("emtab");
   for (var i=0; i<plist_len; i++) {
      var row = tb.insertRow(i+1);
      row.insertCell(0);
      row.insertCell(1);
      row.onclick = function() {
         document.forms["emaildel"].email.value = this.cells[0].innerHTML;
         document.getElementById("modaldel").style.display = "block";
      }
   }
   $(".gomodaladd").click(function(ev) {
      ev.preventDefault();
      document.getElementById("modaladd").style.display = "block";
   });
   $(".unmodaladd").click(function(ev) {
      ev.preventDefault();
      document.getElementById("modaladd").style.display = "none";
   });
   $(".unmodaldel").click(function(ev) {
      ev.preventDefault();
      document.getElementById("modaldel").style.display = "none";
   });
   $(".embtn").click(function(ev) {
      ev.preventDefault();
      if (escape(document.getElementById("c1").innerHTML) == '%u2193')
         emZ2A();
      else
         emA2Z();
   });
   $(".dtbtn").click(function(ev) {
      ev.preventDefault();
      if (escape(document.getElementById("c2").innerHTML) == '%u2193')
         dtZ2A();
      else
         dtA2Z();
   });
   
   emA2Z();
}
 
function emA2Z() {
   document.getElementById("c1").innerHTML = "&darr;";
   document.getElementById("c2").innerHTML = "&nbsp;";
   var tb = document.getElementById("emtab");
   for (var i=0; i<plist_len; i++) {
     tb.rows[i+1].cells[0].innerHTML = plistByEmail[i].em;
     tb.rows[i+1].cells[1].innerHTML = plistByEmail[i].dt.split(" ")[0];
   }
}
function emZ2A() {
   document.getElementById("c1").innerHTML = "&uarr;";
   document.getElementById("c2").innerHTML = "&nbsp;";
   var tb = document.getElementById("emtab");
   for (var i=plist_len-1; i>= 0; i--) {
     tb.rows[(plist_len-i)].cells[0].innerHTML = plistByEmail[i].em;
     tb.rows[(plist_len-i)].cells[1].innerHTML = plistByEmail[i].dt.split(" ")[0];
   }
}
function dtA2Z() {
   document.getElementById("c1").innerHTML = "&nbsp;";
   document.getElementById("c2").innerHTML = "&darr;";
   var tb = document.getElementById("emtab");
   for (var i=0; i<plist_len; i++) {
     tb.rows[i+1].cells[0].innerHTML = plistByDate[i].em;
     tb.rows[i+1].cells[1].innerHTML = plistByDate[i].dt.split(" ")[0];
   }
}
function dtZ2A() {
   document.getElementById("c1").innerHTML = "&nbsp;";
   document.getElementById("c2").innerHTML = "&uarr;";
   var tb = document.getElementById("emtab");
   for (var i=plist_len-1; i>= 0; i--) {
     tb.rows[(plist_len-i)].cells[0].innerHTML = plistByDate[i].em;
     tb.rows[(plist_len-i)].cells[1].innerHTML = plistByDate[i].dt.split(" ")[0];
   }
}
