package finalproject

import (
	"encoding/json"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/appengine/urlfetch"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
	"google.golang.org/cloud/storage"
)

var temp *template.Template

func init() {
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public/"))))

	http.HandleFunc("/", serve)
	http.HandleFunc("/login/", loginpage)
	http.HandleFunc("/register/", registerpage)
	http.HandleFunc("/state/", statepage)
	http.HandleFunc("/home/", homepage)
	http.HandleFunc("/settings/", settingspage)

	//API CALLS
	http.HandleFunc("/api/unique/", unique)

	temp = template.Must(template.ParseGlob("public/*.html"))

}

func homepage(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Errorf(ctx, "error making new client")
	}
	defer client.Close()

	var mysess Session
	var myuser User
	_, session, err := getsess(req)

	//if not logged in
	if err != nil {
		http.Redirect(res, req, `/`, http.StatusSeeOther)
		return
	}

	mysess.id = session

	info, err := memcache.Get(ctx, session)

	json.Unmarshal(info.Value, &myuser)
	mysess.User = myuser

	if req.Method == "POST" {

		mpf, hdr, err := req.FormFile("uploader")
		if err != nil {
			log.Errorf(ctx, "ERROR handler req.FormFile: ", err)
			http.Error(res, "We were unable to upload your file\n", http.StatusInternalServerError)
			return
		}
		defer mpf.Close()

		err = uploadFile(req, mpf, hdr, myuser)
		if err != nil {
			log.Errorf(ctx, "ERROR handler uploadFile: ", err)
			http.Error(res, "We were unable to accept your file\n"+err.Error(), http.StatusUnsupportedMediaType)
			return
		}

		// err = uploadSmallerFile(req, mpf, hdr, myuser)
		// if err != nil {
		// 	log.Errorf(ctx, "ERROR handler uploadFile: ", err)
		// 	http.Error(res, "We were unable to accept your file\n"+err.Error(), http.StatusUnsupportedMediaType)
		// 	return
		// }

	}

	query := &storage.Query{
		Prefix: myuser.Email + "/",
	}

	objs, err := client.Bucket(gcsBucket).List(ctx, query)
	if err == nil {

		var bucketphotos []string

		for _, obj := range objs.Results {
			bucketphotos = append(bucketphotos, obj.Name)
			//	html += `<img src=https://storage.googleapis.com/meme-1299.appspot.com/` + obj.Name + `>` + "\n"
		}

		myuser.Photos = bucketphotos

	} else {
		log.Errorf(ctx, "ERROR in query_delimiter")

	}

	mysess.User = myuser

	temp.ExecuteTemplate(res, "homepage.html", mysess)
}

func serve(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	// Fill the record with the data from the JSON
	var record SynApiData
	var mysess Session
	_, session, err := getsess(req)
	mysess.id = session

	if err == nil {
		http.Redirect(res, req, `/home?q=`+mysess.id, http.StatusSeeOther)
	}

	if req.Method == "POST" {
		//get data from form
		word := req.FormValue("new-word")
		apikey := "AIzaSyAeBfv6lcL1FQEWkNxWujXhE0x5LxKqFmY"
		url := "http://words.bighugelabs.com/api/2/" + apikey + "/" + word + "/json"
		//	http://words.bighugelabs.com/api/2/b5a921e2e1ca6ec782682b1b4654f3f1/word/json

		// Build the request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Infof(ctx, "NewRequest: ", err)
			return
		}

		// For control over HTTP client headers,
		// redirect policy, and other settings,
		// create a Client
		// A Client is an HTTP client
		client := urlfetch.Client(ctx)

		// Send the request via a client
		// Do sends an HTTP request and
		// returns an HTTP response
		resp, err := client.Do(req)
		if err != nil {
			log.Infof(ctx, "Do: ", err)
			return
		}

		// Callers should close resp.Body
		// when done reading from it
		// Defer the closing of the body
		defer resp.Body.Close()

		// Use json.Decode for reading streams of JSON data
		if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
			log.Infof(ctx, "Decodin probs:", err)
		}

	}

	temp.ExecuteTemplate(res, "index.html", record)
}

func loginpage(res http.ResponseWriter, req *http.Request) {
	var mysess Session
	var myuser User
	var retrieveduser User
	ctx := appengine.NewContext(req)

	if req.Method == "POST" {
		//get data from form
		useremail := req.FormValue("username")
		userpass := req.FormValue("password")

		//validate email
		if validateEmail(useremail) == false {
			http.Redirect(res, req, "/login/?m=BadEmail", http.StatusSeeOther)
			return
		}

		//check the datastore for that info
		key := datastore.NewKey(ctx, "EMAILS", useremail, 0, nil)
		err := datastore.Get(ctx, key, &retrieveduser)

		log.Infof(ctx, "UNAME:", useremail)
		log.Infof(ctx, "PASS:", userpass)
		log.Infof(ctx, "myuser.email:", myuser.Email)
		log.Infof(ctx, "myuser.PASS:", myuser.Password)
		log.Infof(ctx, "ret.email:", retrieveduser.Email)
		log.Infof(ctx, "ret.PASS:", retrieveduser.Password)

		hiddenpass, err := bcrypt.GenerateFromPassword([]byte(userpass), bcrypt.DefaultCost)
		log.Infof(ctx, "hiddenpass:", string(hiddenpass))
		log.Infof(ctx, "myuser.pass:", myuser.Password)

		if err == nil && bcrypt.CompareHashAndPassword([]byte(retrieveduser.Password), []byte(userpass)) == nil {
			mysess.id = makesess(res, req, retrieveduser)
			http.Redirect(res, req, `/home?q=`+mysess.id, http.StatusSeeOther)
		} else {
			log.Infof(ctx, "User information was not found in datastore, Not Logged in!")
		}

	}

	temp.ExecuteTemplate(res, "login.html", mysess)
}

func settingspage(res http.ResponseWriter, req *http.Request) {
	var mysess Session
	var myuser User
	var retrieveduser User
	ctx := appengine.NewContext(req)
	_, session, _ := getsess(req)

	mysess.id = session

	info, _ := memcache.Get(ctx, session)

	json.Unmarshal(info.Value, &myuser)
	mysess.User = myuser

	if req.Method == "POST" {
		//get data from form
		newpass := req.FormValue("newpass")
		mypass1 := req.FormValue("password1")
		mypass2 := req.FormValue("password2")

		//check the datastore for that info
		key := datastore.NewKey(ctx, "EMAILS", myuser.Email, 0, nil)
		err := datastore.Get(ctx, key, &retrieveduser)

		//passwords dont match
		if mypass1 != mypass2 {
			log.Infof(ctx, "Password does not match confirmed pass")
			temp.ExecuteTemplate(res, "register.html", mysess)
		}

		//hash pass
		hiddenpass, err := bcrypt.GenerateFromPassword([]byte(mypass1), bcrypt.DefaultCost)
		log.Infof(ctx, "hiddenpass:", string(hiddenpass))
		if err != nil {
			log.Errorf(ctx, "Could not bcrypt the pass", err)
			http.Error(res, err.Error(), 500)
		}

		if err == nil && bcrypt.CompareHashAndPassword([]byte(retrieveduser.Password), []byte(mypass1)) == nil {

			hiddenpass, err := bcrypt.GenerateFromPassword([]byte(newpass), bcrypt.DefaultCost)
			log.Infof(ctx, "hiddenpass:", string(hiddenpass))
			if err != nil {
				log.Errorf(ctx, "Could not bcrypt the pass", err)
				http.Error(res, err.Error(), 500)
			}

			//create the registered user
			reggeduser := User{
				Email:    myuser.Email,
				Password: string(hiddenpass),
			}

			//make key for hash table
			userkey := datastore.NewKey(ctx, "EMAILS", reggeduser.Email, 0, nil)
			//save to datastore
			userkey, err = datastore.Put(ctx, userkey, &reggeduser)

			mysess.id = makesess(res, req, retrieveduser)

			http.Redirect(res, req, `/home?q=`+mysess.id, http.StatusSeeOther)

		} else {
			log.Infof(ctx, "User information was not found in datastore, Not Logged in!")
		}

	}

	temp.ExecuteTemplate(res, "settings.html", mysess)
}

func statepage(res http.ResponseWriter, req *http.Request) {

	key := "q"
	Val := req.URL.Query().Get(key)

	temp.ExecuteTemplate(res, "state.html", Val)
}

func registerpage(res http.ResponseWriter, req *http.Request) {

	var mysess Session
	var myuser User

	//string slice photos
	var myphotos []string

	ctx := appengine.NewContext(req)

	if req.Method == "POST" {
		myemail := req.FormValue("username")
		mypass1 := req.FormValue("password1")
		mypass2 := req.FormValue("password2")
		myuser.Email = myemail
		myuser.Password = mypass1
		myuser.Photos = myphotos

		//validate email
		if validateEmail(myemail) == false {
			http.Redirect(res, req, "/register/?m=BadEmail", http.StatusSeeOther)
			return
		}

		log.Infof(ctx, "UNAME:", myemail)
		log.Infof(ctx, "PASS:", mypass1)

		userkey := datastore.NewKey(ctx, "EMAILS", myuser.Email, 0, nil)
		err := datastore.Get(ctx, userkey, &myuser)
		log.Infof(ctx, "myuser.email:", myuser.Email)
		log.Infof(ctx, "myuser.pass:", myuser.Password)

		//Username not unique
		if err == nil {
			log.Infof(ctx, "!!problem in register, User already created!!")

			temp.ExecuteTemplate(res, "register.html", mysess)
			return
		}

		//passwords dont match
		if mypass1 != mypass2 {
			log.Infof(ctx, "Password does not match confirmed pass")
			temp.ExecuteTemplate(res, "register.html", mysess)
		}

		//hash pass
		hiddenpass, err := bcrypt.GenerateFromPassword([]byte(mypass1), bcrypt.DefaultCost)
		log.Infof(ctx, "hiddenpass:", string(hiddenpass))
		if err != nil {
			log.Errorf(ctx, "Could not bcrypt the pass", err)
			http.Error(res, err.Error(), 500)
		}

		//create the registered user
		reggeduser := User{
			Email:    myemail,
			Password: string(hiddenpass),
		}

		//make key for hash table
		userkey = datastore.NewKey(ctx, "EMAILS", reggeduser.Email, 0, nil)
		//save to datastore
		userkey, err = datastore.Put(ctx, userkey, &reggeduser)

		if err != nil {
			log.Errorf(ctx, "While Registering, could not store user to datastore", err)
			http.Error(res, err.Error(), 500)
			return
		}

		newsessionid := makesess(res, req, reggeduser)
		log.Errorf(ctx, "myuser", reggeduser)

		http.Redirect(res, req, "/home/?q="+newsessionid, http.StatusSeeOther)

	}

	temp.ExecuteTemplate(res, "register.html", mysess)
}
