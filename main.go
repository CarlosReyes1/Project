package main

import (
	"html/template"
	"log"
	"net/http"
	"github.com/nu7hatch/gouuid"
	//"fmt"
)

func serve(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("session-info")
	//if cookie not foud retrun error
	if err != nil {
		id,_ := uuid.NewV4()

		cookie = &http.Cookie{//stores cookie to make sure it works
			Name: "session-info",//we dont need a UUID we used a string
			Value: id.String(),
			HttpOnly: true,//dont know what it does
		}
		http.SetCookie(res, cookie)//if no cookie found it will set a new cookie
	}


	//create template from index.html file
	temp, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatalln(err, "Failed Parsing File")
	}
	temp.Execute(res, nil)
	//fmt.Fprintf(res, "UUID:"+cookie.Value)
}

func init() {
	http.HandleFunc("/", serve)
	http.Handle("/css/",http.StripPrefix("/css", http.FileServer(http.Dir("css"))))
	http.Handle("/img/",http.StripPrefix("/img", http.FileServer(http.Dir("img"))))

}

