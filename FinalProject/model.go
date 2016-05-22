package finalproject

type User struct {
	Email    string //email is the prim key, has to be unique
	Password string
	Photos   []string
}

type Session struct {
	User User
	id   string
}

type SynApiData struct {
	Noun struct {
		Syn []string `json:"syn"`
	} `json:"noun"`
	Verb struct {
		Syn []string `json:"syn"`
	} `json:"verb"`
}
