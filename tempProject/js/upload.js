function submitHandler() {
  var u = document.getElementById("uploader").value;

  if (!u.files[0].name.match(/\.(JPG|jpg|JPEG|jpeg|PNG|png|GIF|gif)$/)) {
    alert("We were unable to accept your file")
    return false;
  }
  else {
    alert("We were unable to accept your file2")
    return true;
  }
}

function successUpload() {
  alert("successfully uploaded")
  return true
}