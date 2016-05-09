function submitHandler() {
  var u = document.getElementById("username").value;
  var p = document.getElementById("password").value;
  
  if (u == "" || p == "") {
    return false;
  }
  else {
    return true;
  }
}