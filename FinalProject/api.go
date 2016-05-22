package finalproject

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func unique(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	bs, err := ioutil.ReadAll(req.Body)
	sbs := string(bs)
	log.Infof(ctx, "REQUEST BODY: %v", sbs)
	var user User
	key := datastore.NewKey(ctx, "EMAILS", sbs, 0, nil)
	err = datastore.Get(ctx, key, &user)

	// if there is an err, there is NO user
	log.Infof(ctx, "ERR: %v", err)
	if err != nil {
		// there is an err, there is a NO user
		fmt.Fprint(res, "false")
		return
	}
	fmt.Fprint(res, "true")

}
