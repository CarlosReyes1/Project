package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/log"

	"github.com/nu7hatch/gouuid"
	"golang.org/x/crypto/bcrypt"
)

var tpl *template.Template
var fuq bool

type User struct {
	Username string
	Password string
	Email    string
}

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/signin", signin)
	http.HandleFunc("/register", register)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/main", mainz)
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

func upload(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "upload.html", nil)
}

func mainz(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "main.html", nil)
}

func profile(w http.ResponseWriter, r *http.Request) {
	i, _ := getSession(r)
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
	sd := memcache.Item{
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
	u := User{
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
	o := &http.Cookie{
		Name:  "session",
		Value: id.String(),
		Path:  "/",
		//Secure: true,
		//HttpOnly: true,
	}
	http.SetCookie(w, o)

	j, _ := json.Marshal(u)
	sd := memcache.Item{
		Key:        id.String(),
		Value:      j,
		Expiration: time.Duration(20 * time.Second),
	}
	memcache.Set(c, &sd)
}

func getSession(r *http.Request) (*memcache.Item, error) {
	o, err := r.Cookie("session")
	if err != nil {
		return &memcache.Item{}, err
	}
	c := appengine.NewContext(r)
	i, err := memcache.Get(c, o.Value)
	if err != nil {
		return &memcache.Item{}, err
	}
	return i, nil
}

func handler(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)

	if req.Method == "POST" {

		mpf, hdr, err := req.FormFile("uploader")
		if err != nil {
			fuq = true
			log.Errorf(ctx, "ERROR handler req.FormFile: ", err)
			http.Error(res, "We were unable to upload your file\n", http.StatusInternalServerError)
			http.Redirect(res, req, "/upload", 302)
		} else{
			http.Error(res, "Success\n", http.StatusInternalServerError)
			http.Redirect(res, req, "/upload", 302)

		}

		defer mpf.Close()

		_, err = uploadFile(req, mpf, hdr)
		if err != nil {
			fuq = false
			log.Errorf(ctx, "ERROR handler uploadFile: ", err)
			http.Error(res, "We were unable to accept your file\n" + err.Error(), http.StatusUnsupportedMediaType)
			http.Redirect(res, req, "/main", 302)
		} else{
			http.Redirect(res, req, "/upload", 302)

		}
	} else {
		http.Redirect(res, req, "/profile", 302)
	}


}