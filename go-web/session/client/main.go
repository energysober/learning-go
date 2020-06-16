package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var cookies []*http.Cookie

func httpRequestFunc(method, url string) {
	client := http.Client{}
	var req *http.Request
	var err error
	req, err = http.NewRequest(method, url, nil)
	if err != nil {
		panic("new request error " + err.Error())
	}
	cookie := http.Cookie{}
	req.AddCookie(&cookie)
	resp, err := client.Do(req)
	if err != nil {
		panic("do request error " + err.Error())
	}
	defer resp.Body.Close()
	cookies = resp.Cookies()
	fmt.Println(cookies)
}

func Login() {
	httpRequestFunc("POST", "http://127.0.0.1:8080/login")
}

func GetSession() {
	client := http.Client{}
	var req *http.Request
	var err error
	req, err = http.NewRequest("GET", "http://127.0.0.1:8080/user/getSession", nil)
	if err != nil {
		panic("new request error " + err.Error())
	}
	req.AddCookie(cookies[0])
	resp, err := client.Do(req)
	if err != nil {
		panic("do request error " + err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func Logout() {
	httpRequestFunc("POST", "http://127.0.0.1:8080/logout")
}

func main() {
	Login()
	time.Sleep(time.Second * 10)
	GetSession()
	time.Sleep(time.Second * 10)
	Logout()
}
