package main

import (
	"fmt"
	"net/http"

	"github.com/learning-go/go-web/session/session"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		sess := session.GlobalSessionManager.SessionStart(w, r)
		if err := sess.Set("userName", "test_user"); err != nil {
			fmt.Println("session: set session error: " + err.Error())
		}
	} else {
		http.Error(w, "please use post method", 405)
	}

}

func GetSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		sess := session.GlobalSessionManager.SessionStart(w, r)
		if userName, err := sess.Get("userName"); err != nil {
			fmt.Println("session: get userName error")
		} else {
			fmt.Println(userName)
			_, _ = w.Write([]byte(fmt.Sprintf("%s", userName)))
		}
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		session.GlobalSessionManager.SessionDestroy(w, r)
	} else {
		http.Error(w, "please use post method", 405)
	}

}

func main() {
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/user/getSession", GetSessionHandler)
	http.HandleFunc("/logout", LogoutHandler)
	_ = http.ListenAndServe("127.0.0.1:8080", nil)
}
