var uname = document.querySelector('#username');
var output = document.querySelector('h1');

uname.addEventListener('input', function(){
  var xhr = new XMLHttpRequest();
  xhr.open('POST','/api/unique/');
  xhr.send(uname.value);
  xhr.addEventListener('readystatechange', function(){
       if (xhr.readyState === 4 && xhr.status === 200) {
           var taken = xhr.responseText;
           console.log('TAKEN:', taken, '\n\n');
           if (taken == 'true') {
               output.textContent = 'uname Taken!';
           } else {
               output.textContent = '';
           }
       }
   });
});
