package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	server()
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	User        string `json:"user"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Id       string `json:"id"`
	Uuid     string `json:"uuid"`
}

func apiLogin(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Println(r.RequestURI, string(body))
	li := loginRequest{}
	_ = json.Unmarshal(body, &li)
	if li.Username == "admin" && li.Password == "password" {
		re := loginResponse{
			AccessToken: "AccessToken",
			User:        `{"username":"username"}`,
		}
		jse := json.NewEncoder(w)
		_ = jse.Encode(re)
	} else {
		_, _ = w.Write([]byte(`{"error":"账号密码错误!"}`))
	}
}

func apiCurrentUser(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Println(r.RequestURI, string(body))
	jse := json.NewEncoder(w)
	_ = jse.Encode(`{"username":"username"}`)
}

func apiLogout(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Println(r.RequestURI, string(body))
}

type response struct {
	Error string `json:"error,omitempty"`
	Data  string `json:"data,omitempty"`
}

//{
//    peers: [{id: "abcd", username: "", hostname: "", platform: "", alias: "", tags: ["", "", ...]}, ...],
//    tags: [],
//}

type peer struct {
	Id       string   `json:"id"`
	Username string   `json:"username"`
	Hostname string   `json:"hostname"`
	Platform string   `json:"platform"`
	Alias    string   `json:"alias"`
	Tags     []string `json:"tags"`
}
type addressBook struct {
	Peers []peer   `json:"peers,omitempty"`
	Tags  []string `json:"tags,omitempty"`
}

const addressFile = "address.json"

func apiAbGet(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Println(r.RequestURI, string(body))
	d, _ := ioutil.ReadFile(addressFile)
	re := response{
		Error: "",
		Data:  string(d),
	}
	jse := json.NewEncoder(w)
	_ = jse.Encode(re)
}

func apiAbPost(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	b := response{}
	e := json.Unmarshal(body, &b)
	if e != nil {
		w.Write([]byte(`{"error":"error"}`))
		return
	}
	fmt.Println(r.RequestURI, b.Data)
	ab := addressBook{}
	_ = json.Unmarshal([]byte(b.Data), &ab)
	ioutil.WriteFile(addressFile, []byte(b.Data), 0644)
}

func server() {
	http.HandleFunc("/api/login", apiLogin)
	http.HandleFunc("/api/currentUser", apiCurrentUser)
	http.HandleFunc("/api/logout", apiLogout)
	http.HandleFunc("/api/ab/get", apiAbGet)
	http.HandleFunc("/api/ab", apiAbPost)
	e := http.ListenAndServe(":21114", nil)
	if e != nil {
		panic(e)
	}
}
