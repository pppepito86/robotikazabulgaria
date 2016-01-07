package main

import (
	"fmt"
	"strings"
	"time"
)

import "net/http"
import "html/template"
import "io"

import "os"

import (
	"robotikazabulgaria/hw"
	"robotikazabulgaria/session"
	"robotikazabulgaria/user"
	"robotikazabulgaria/ws"
)

func main() {
	fmt.Println("working dir", ws.Getwd())
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	http.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir("docs"))))
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	loggedIn := false
	if r.URL.Path == "/login.html" && isLoggedIn(*r) {
		http.Redirect(w, r, "/home.html", http.StatusFound)
		return
	}
	if r.URL.Path == "/index.html" && !isLoggedIn(*r) {
		t, _ := template.ParseFiles("index.html")
		t.Execute(w, nil)
		return
	}
	if r.URL.Path == "/login.html" {
		t, _ := template.ParseFiles("login.html")
		t.Execute(w, nil)
		return
	} else if r.URL.Path == "/login" {
		if r.Method == "GET" {
			loggedIn = getLogin(w, *r)
		} else if r.Method == "POST" {
			loggedIn = postLogin(w, r)
		}
	} else if r.URL.Path == "/upload" {
		homework, err := upload(w, r)
		if err == nil {
			hw.AddHomework(getUser(*r), homework)
		}
		http.Redirect(w, r, "/tasks.html", http.StatusFound)
		return
	}
	if loggedIn || isLoggedIn(*r) {
		fmt.Println("******", r.URL.Path[1:])
		t, err := template.ParseFiles(r.URL.Path[1:])
		if err != nil {
			http.Redirect(w, r, "/home.html", http.StatusFound)
			return
		}
		if r.URL.Path == "/home.html" {
			// sss := []string{"aaa", "bbb", "ccc"}
			// pwd, _ := os.Getwd()
			// files, _ := filepath.Glob(pwd+"\\"+getUser(*r)+"\\*")
			t.Execute(w, hw.ReadHomeworks(getUser(*r)))
		} else {
			t.Execute(w, nil)
		}
	} else {
		http.Redirect(w, r, "/index.html", http.StatusFound)
	}
}

func upload(w http.ResponseWriter, r *http.Request) (hw.Homework, error) {
	file, header, err := r.FormFile("file")
	if err != nil {
		if strings.TrimSpace(r.Form["link"][0]) != "" {
			return hw.Homework{"", r.Form["link"][0], r.Form["description"][0], r.Form["task"][0], time.Now().UTC()}, nil
		} else {
			return hw.Homework{}, err
		}
	}
	defer file.Close()
	fp := ws.GetFilePath(getUser(*r), header.Filename)
	out, err := os.Create(fp)
	if err != nil {
		fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
		return hw.Homework{}, err
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Fprintln(w, err)
		return hw.Homework{}, err
	}
	return hw.Homework{header.Filename, r.Form["link"][0], r.Form["description"][0], r.Form["task"][0], time.Now().UTC()}, nil
}
func getUser(r http.Request) string {
	cookie := getSessionIdCookie(r)
	return session.GetAttribute(cookie.Value)
}
func isLoggedIn(r http.Request) bool {
	cookie := getSessionIdCookie(r)
	fmt.Println("session Cookie is:", cookie.Value)
	return session.ContainsKey(cookie.Value)
}
func getSessionIdCookie(r http.Request) *http.Cookie {
	for _, cookie := range r.Cookies() {
		if cookie.Name == "session.id" {
			return cookie
		}
	}
	return new(http.Cookie)
}
func getLogin(w http.ResponseWriter, r http.Request) bool {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")
	return login(w, r, username, password)
}
func postLogin(w http.ResponseWriter, r *http.Request) bool {
	r.ParseForm()
	user := r.Form["username"]
	pass := r.Form["password"]
	fmt.Println("user", user)
	fmt.Println("pass", pass)
	if len(user) != 1 || len(pass) != 1 {
		return false
	}
	return login(w, *r, user[0], pass[0])
}
func login(w http.ResponseWriter, r http.Request, username string, password string) bool {
	fmt.Println("username:", username, "password:", password)
	if !user.Authenticate(username, password) {
		return false
	}
	val := username + "-" + user.RandomString()
	cookie := http.Cookie{Name: "session.id", Value: val}
	http.SetCookie(w, &cookie)
	session.SetAttribute(val, username)
	fmt.Println("set session cookie", val)
	return true
}
