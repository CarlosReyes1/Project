var shape;
var comp;

function init() {
  // 0 = O, 1 = X
  switch (Math.floor(Math.random() * 2)) {
    case 0:
      shape = "images/piece_o.png";
      comp = "images/piece_x.png";
      document.getElementById("you_img").setAttribute("src", shape);
      break;
    case 1:
      shape = "images/piece_x.png";
      comp = "images/piece_o.png";
      document.getElementById("you_img").setAttribute("src", shape);
  }
  
  // 0 = human, 1 = comp
  switch (Math.floor(Math.random() * 2)) {
    case 0:
      break;
    case 1:
      place("");
  }
}

function place(id) {
  var xhttp = new XMLHttpRequest();
  
  xhttp.onreadystatechange = function() {
    if (xhttp.readyState == 4 && xhttp.status == 200) {
      var resp = xhttp.responseText;
      
      if (resp == "win") {
        var m = document.getElementById("msg");
        m.style.visibility = "visible";
        m.innerHTML = "Game Over! You won!";
        return;
      }
      else if (resp == "lose") {
        var m = document.getElementById("msg");
        m.style.visibility = "visible";
        m.innerHTML = "Game Over! You lose!";
        return;
      }
      else if (resp == "tie") {
        var m = document.getElementById("msg");
        m.style.visibility = "visible";
        m.innerHTML = "Game Over! It's a tie!";
        return;
      }
      
      /*if (resp == "-1") {
        document.getElementById("msg").style.visibility = "visible";
        return;
      }*/
      
      var s = document.getElementById("sq" + resp).style;
      s.backgroundImage = "url('" + comp + "')";
      s.backgroundRepeat = "no-repeat";
      s.backgroundPosition = "center center";
      
      unblock();
    }
  };
  
  if (id != "") {
    var s = document.getElementById(id).style;
    s.backgroundImage = "url('" + shape + "')";
    s.backgroundRepeat = "no-repeat";
    s.backgroundPosition = "center center";
  }
  
  block();
  var s = "/getmove?" + getBoardState();
  xhttp.open("GET", s, true);
  xhttp.send();
}

function block() {
  var d = document.getElementsByClassName("q");
  for (var i = 0; i < d.length; i++) {
    d[i].setAttribute("onclick", "");
  }
}

function unblock() {
  var d = document.getElementsByClassName("q");
  for (var i = 0; i < d.length; i++) {
    d[i].setAttribute("onclick", "place('sq" + i.toString() + "')");
  }
}

function getBoardState() {
  var d = document.getElementsByClassName("q");
  var s = "sq";
  var t = "";
  var r = 'url("' + shape + '")';
  
  for (var i = 0; i < d.length; i++) {
    var q = s + i.toString() + "=";
    
    if (d[i].style.backgroundImage == "") {
      q = q + "n";
    }
    else if (d[i].style.backgroundImage == r) {
      q = q + "u";
    }
    else {
      q = q + "c";
    }
    
    t = t + q;
    if (i != d.length - 1) {
      t = t + "&";
    }
  }
  
  return t;
}