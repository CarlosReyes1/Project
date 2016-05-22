package finalproject

import (
	"encoding/json"
	"net/http"

	"github.com/nu7hatch/gouuid"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

func makesess(res http.ResponseWriter, req *http.Request, myuser User) string {
	//make appengine context to start doing anything
	ctx := appengine.NewContext(req)
	id, _ := uuid.NewV4()

	//make cookie and place it
	cookie := &http.Cookie{
		Name:  "Session",
		Value: id.String(),
		Path:  "/",
		// Secure:   true,
		// HttpOnly: true,
	}
	http.SetCookie(res, cookie)

	jsonneduser, err := json.Marshal(myuser)
	if err != nil {
		//error marshalling json
		log.Errorf(ctx, "COULD NOT MARSHAL JSON IN MAKESESS", err)
		http.Error(res, err.Error(), 500)
		return "-1"
	}

	//puts it in KVPair
	memitem := memcache.Item{
		Key:   cookie.Value,
		Value: jsonneduser,
	}

	//stores in memcache
	memcache.Set(ctx, &memitem)

	//returns jsonned info
	return cookie.Value
}

func getsess(req *http.Request) (*memcache.Item, string, error) {
	//new context for goog app engine
	ctx := appengine.NewContext(req)

	cookie, err := req.Cookie("Session")

	if err != nil {
		key := "q"
		val := req.URL.Query().Get(key)

		if val != "" {
			item, err := memcache.Get(ctx, val)
			//not found
			if err != nil {
				return &memcache.Item{}, "", err
			}

			//returns item from memcache, session info in url, and no error!
			return item, val, nil
		}
		//error and no items
		return &memcache.Item{}, "", err
	}

	//if at this line, you had a cookie!
	//get the values from memcache from cookie value!
	item, err := memcache.Get(ctx, cookie.Value)

	if err != nil {
		return &memcache.Item{}, "", err
	}

	return item, cookie.Value, nil

}
