function checkUsername() {
  var xhttp;
  var u = document.getElementById("username");
  var i = document.getElementById("username_img");
  
  if (u.value == "") {
    u.style.border = "solid 5px #cc0000";
    i.setAttribute("src", "images/wrong.png");
    return;
  }
  
  xhttp = new XMLHttpRequest();
  
  xhttp.onreadystatechange = function() {
    if (xhttp.readyState == 4 && xhttp.status == 200) {
      if (xhttp.responseText == "good") {
        u.style.border = "solid 5px #330099";
        i.setAttribute("src", "images/check.png");
      }
      else {
        u.style.border = "solid 5px #cc0000";
        i.setAttribute("src", "images/wrong.png");
      }
    }
  };
  
  xhttp.open("POST", "/checkUsername", true);
  xhttp.send(new FormData(document.querySelector("form")));
}

function checkPassword() {
  var f = document.getElementById("firstpass");
  var s = document.getElementById("secondpass");
  var fi = document.getElementById("firstpass_img");
  var si = document.getElementById("secondpass_img");
  
  if (s.value != "" && f.value == s.value) {
    f.style.border = "solid 5px #330099";
    s.style.border = "solid 5px #330099";
    fi.setAttribute("src", "images/check.png");
    si.setAttribute("src", "images/check.png");
  }
  else {
    f.style.border = "solid 5px #cc0000";
    s.style.border = "solid 5px #cc0000";
    fi.setAttribute("src", "images/wrong.png");
    si.setAttribute("src", "images/wrong.png");
  }
}

function checkEmail() {
  var e = document.getElementById("email");
  var ei = document.getElementById("email_img");
  var r = /.+@.+\..+/;
  
  if (r.test(e.value)) {
    e.style.border = "solid 5px #330099";
    ei.setAttribute("src", "images/check.png");
  }
  else {
    e.style.border = "solid 5px #cc0000";
    ei.setAttribute("src", "images/wrong.png");
  }
}

function submitHandler() {
  var ui = document.getElementById("username_img").getAttribute("src");
  var fi = document.getElementById("firstpass_img").getAttribute("src");
  var si = document.getElementById("secondpass_img").getAttribute("src");
  var ei = document.getElementById("email_img").getAttribute("src");
  var ch = "images/check.png";
  
  if (ui != ch || fi != ch || si != ch || ei != ch) {
    return false;
  }
  else {
    return true;
  }
}

