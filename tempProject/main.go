package main

import (
  "html/template"
  "net/http"
  "encoding/json"
  "time"

  
  "google.golang.org/appengine"
  "google.golang.org/appengine/datastore"
  "google.golang.org/appengine/memcache"
  "google.golang.org/appengine/log"
  "google.golang.org/cloud/storage"
  
  "golang.org/x/crypto/bcrypt"
  "github.com/nu7hatch/gouuid"
)

const gcsBucket = "meme-1299.appspot.com"
var tpl *template.Template
var fuq bool

type User struct {
  Username  string
  Password  string
  Email     string
  Photos   []string
}

func init() {
  http.HandleFunc("/", root)
  http.HandleFunc("/signin", signin)
  http.HandleFunc("/register", register)
  http.HandleFunc("/main", mainz)
  http.HandleFunc("/upload", upload)
  http.HandleFunc("/uploadMore", uploadMore)
  http.HandleFunc("/uploadFail", uploadFail)
  http.HandleFunc("/profile", profile)
  http.HandleFunc("/signout", signout)
  
  http.HandleFunc("/validate", validate)
  http.HandleFunc("/create", create)
  http.HandleFunc("/checkUsername", checkUsername)
  http.HandleFunc("/handler", handler)
  
  tpl = template.Must(template.ParseGlob("html/*.html"))
  fuq = false
}

func root(w http.ResponseWriter, r *http.Request) {
  fuq = false
  tpl.ExecuteTemplate(w, "index.html", nil)
}

func signin(w http.ResponseWriter, r *http.Request) {
  tpl.ExecuteTemplate(w, "signin.html", fuq)
}

func register(w http.ResponseWriter, r *http.Request) {
  tpl.ExecuteTemplate(w, "register.html", nil)
}

func mainz(w http.ResponseWriter, r *http.Request) {
  tpl.ExecuteTemplate(w, "main.html", nil)
}

func upload(w http.ResponseWriter, r *http.Request) {
  tpl.ExecuteTemplate(w, "upload.html", nil)
}

func uploadMore(w http.ResponseWriter, r *http.Request) {
  tpl.ExecuteTemplate(w, "uploadMore.html", nil)
}

func uploadFail(w http.ResponseWriter, r *http.Request) {
  tpl.ExecuteTemplate(w, "uploadFail.html", nil)
}

func profile(w http.ResponseWriter, r *http.Request) {
  i := getSession(r)
  if len(i.Value) > 0 {
    var u User
    json.Unmarshal(i.Value, &u)
    tpl.ExecuteTemplate(w, "profile.html", u)
  } else {
    tpl.ExecuteTemplate(w, "profile.html", nil)
  }
}

func signout(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  o, _ := r.Cookie("session")
  sd := memcache.Item {
    Key:        o.Value,
    Value:      []byte(""),
    Expiration: time.Duration(1 * time.Microsecond),
  }
  memcache.Set(c, &sd)
  o.MaxAge = -1
  http.SetCookie(w, o)
  http.Redirect(w, r, "/", 302)
}

func validate(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  k := datastore.NewKey(c, "Users", r.FormValue("username"), 0, nil)
  var u User
  err := datastore.Get(c, k, &u)
  if err != nil ||
     bcrypt.CompareHashAndPassword([]byte(u.Password),
     []byte(r.FormValue("password"))) != nil {
    fuq = true
    http.Redirect(w, r, "/signin", 302)
  } else {
    fuq = false
    createSession(w, r, u)
    http.Redirect(w, r, "/main", 302)
  }
}

func create(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  h, _ := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
  u := User {
    Username: r.FormValue("username"),
    Password: string(h),
    Email:    r.FormValue("email"),
  }
  k := datastore.NewKey(c, "Users", u.Username, 0, nil)
  _, err := datastore.Put(c, k, &u)
  if err != nil {
    http.Redirect(w, r, "/register", 302)
  }
  createSession(w, r, u)
  http.Redirect(w, r, "/main", 302)
}

func checkUsername(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  var u User
  k := datastore.NewKey(c, "Users", r.FormValue("username"), 0, nil)
  err := datastore.Get(c, k, &u)
  if err != nil {
    w.Write([]byte("good"))
  } else {
    w.Write([]byte("bad"))
  }
}

func createSession(w http.ResponseWriter, r *http.Request, u User) {
  c := appengine.NewContext(r)
  
  id, _ := uuid.NewV4()
  o := &http.Cookie {
    Name:   "session",
    Value:  id.String(),
    Path:   "/",
    //Secure: true,
    //HttpOnly: true,
  }
  http.SetCookie(w, o)
  
  j, _ := json.Marshal(u)
  sd := memcache.Item {
    Key:        id.String(),
    Value:      j,
    Expiration: time.Duration(20 * time.Second),
  }
  memcache.Set(c, &sd)
}

func getSession(r *http.Request) *memcache.Item {
  o, err := r.Cookie("session")
  if err != nil {
    return &memcache.Item{}
  }
  c := appengine.NewContext(r)
  i, err := memcache.Get(c, o.Value)
  if err != nil {
    return &memcache.Item{}
  }
  return i
}

func handler(res http.ResponseWriter, req *http.Request) {
  ctx := appengine.NewContext(req)
  client, err := storage.NewClient(ctx)

  if err != nil {
    log.Errorf(ctx, "error making new client")
  }
  defer client.Close()

  if req.Method == "POST" {

    mpf, hdr, err := req.FormFile("uploader")
    if err != nil {
      http.Redirect(res, req, "/uploadFail", 302)
      return
    }


    defer mpf.Close()

    _, err = uploadFile(req, mpf, hdr)
    if err != nil {
      http.Redirect(res, req, "/uploadFail", 302)
      return
    }
  }
  
  http.Redirect(res, req, "/uploadMore", 302)
}
