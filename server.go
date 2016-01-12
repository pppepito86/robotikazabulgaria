package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

import "net/http"
import "html/template"
import "io"

import "os"

import (
	"robotikazabulgaria/admin"
	"robotikazabulgaria/dashboard"
	"robotikazabulgaria/hw"
	"robotikazabulgaria/session"
	"robotikazabulgaria/teams"
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

	if !isLoggedIn(*r) {
		handleGuest(w, r)
	} else if isAdmin(*r) {
		handleAdmin(w, r)
	} else {
		handleTeam(w, r)
	}
}

func handleGuest(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	if r.URL.Path == "/login.html" {
		if r.Method == "POST" {
			if postLogin(w, r) {
				fmt.Println("login successful")
				http.Redirect(w, r, "/home.html", http.StatusFound)
			} else {
				fmt.Println("login failed")
				sendError(w, r, "Грешни данни за вход", "/login.html")
			}
		} else {
			t, _ := template.ParseFiles("login.html")
			t.Execute(w, nil)
		}
	} else if r.URL.Path == "/register.html" {
		if r.Method == "POST" {
			err := register(r)
			if err == nil {
				http.Redirect(w, r, "/login.html", http.StatusFound)
			} else {
				sendError(w, r, err.Error(), "/register.html")
			}
		} else {
			t, _ := template.ParseFiles("register.html")
			t.Execute(w, nil)
		}
	} else if r.URL.Path == "/index.html" {
		t, _ := template.ParseFiles("index.html")
		location, _ := time.LoadLocation("Europe/Sofia")
		deadline := time.Date(2016, 1, 24, 18, 0, 0, 0, location)
		t.Execute(w, deadline.UnixNano()/1000000)
	} else {
		http.Redirect(w, r, "/index.html", http.StatusFound)
	}
}

func handleAdmin(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/admin.html" {
		http.Redirect(w, r, "/admin.html", http.StatusFound)
		return
	}
	if r.URL.Query().Get("page") == "registered_teams" {
		t, _ := template.ParseFiles("admin_registered_teams.html")
		t.Execute(w, teams.GetTeams())
	} else if r.URL.Query().Get("page") == "tasks" {
		if r.Method == "POST" {
			admin.UploadTask(w, r)
		}
		t, _ := template.ParseFiles("admin_tasks.html")
		t.Execute(w, teams.GetTeams())
	} else {
		if r.Method == "POST" {
			r.ParseForm()
			fmt.Println("id", r.Form["id"])
			fmt.Println("city", r.Form["city"])
			fmt.Println("school", r.Form["school"])
			teams.AddTeamId(r.Form["id"][0], r.Form["city"][0], r.Form["school"][0])
		}
		t, _ := template.ParseFiles("admin.html")
		t.Execute(w, teams.GetTeamsIdInfo())
	}
}

func handleTeam(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/login.html" ||
		r.URL.Path == "/register.html" ||
		r.URL.Path == "/index.html" ||
		r.URL.Path == "/admin.html" {
		http.Redirect(w, r, "/home.html", http.StatusFound)
		return
	}
	if r.URL.Path == "/tasks.html" && r.Method == "POST" {
		r.ParseForm()
		fmt.Println("len", r.Form["operation"])
		if len(r.Form["operation"]) == 0 {
			homework, err := upload(w, r)
			fmt.Println("Upload error", err)
			if err == nil {
				hw.AddHomework(getUser(*r), homework)
				http.Redirect(w, r, "/tasks.html", http.StatusFound)
				return
			} else {
				sendError(w, r, err.Error(), "/tasks.html")
				return
			}
		}
	}
	if r.URL.Path == "/download" {
		download(w, r)
		return
	}
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
		location, _ := time.LoadLocation("Europe/Sofia")
		tt := time.Date(2016, 1, 24, 18, 0, 0, 0, location)
		t.Execute(w, tt.UnixNano()/1000000)
	} else if r.URL.Path == "/tasks.html" {
		r.ParseForm()
		fmt.Println(r.Form["operation"])
		if len(r.Form["operation"]) != 0 {
			fmt.Println("delete")
			if r.Form["operation"][0] == "delete" {
				hw.DeleteHomework(getUser(*r), r.Form["id"][0])
			}
		}
		t.Execute(w, dashboard.GetDashboard(getUser(*r)))
	} else {
		t.Execute(w, nil)
	}

}

func register(r *http.Request) error {
	r.ParseForm()
	return teams.RegisterTeam(
		strings.TrimSpace(r.Form["username"][0]),
		r.Form["password1"][0],
		r.Form["password2"][0],
		strings.TrimSpace(r.Form["city"][0]),
		strings.TrimSpace(r.Form["school"][0]),
		strings.TrimSpace(r.Form["identification_number"][0]))
}

func sendError(w http.ResponseWriter, r *http.Request, msg string, page string) {
	t, _ := template.ParseFiles("error.html")
	t.Execute(w,
		struct {
			Message string
			Page    string
		}{msg, page})
}

func download(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")
	file := r.URL.Query().Get("file")
	if user != getUser(*r) {
		return
	}

	http.ServeFile(w, r, ws.GetFilePath(user, file))
}

func upload(w http.ResponseWriter, r *http.Request) (hw.Homework, error) {
	fmt.Println("ContentLength", r.Header.Get("Content-Length"))
	len, _ := strconv.ParseInt(r.Header.Get("Content-Length"), 10, 64)
	fmt.Println("Len", len)
	if len > 20*1024*1024 {
		return hw.Homework{}, errors.New("File should not exceed 20MB")
	}
	r.Body = http.MaxBytesReader(w, r.Body, 20*1024*1024)
	file, header, err := r.FormFile("file")
	if err != nil {
		if strings.Contains(err.Error(), "body too large") {
			return hw.Homework{}, errors.New("File should not exceed 20MB")
		}
		if strings.TrimSpace(r.Form["link"][0]) != "" {
			return hw.Homework{"", r.Form["link"][0], r.Form["description"][0], r.Form["task"][0], time.Now().UTC()}, nil
		} else {
			return hw.Homework{}, errors.New("You should either select a file or enter link")
		}
	} else {
		if strings.TrimSpace(r.Form["link"][0]) != "" {
			return hw.Homework{}, errors.New("You cannot both select a file and enter link")
		}
	}
	t := time.Now().UTC()
	defer file.Close()
	fn := strconv.FormatInt(t.UnixNano(), 16) + filepath.Ext(header.Filename)
	fp := ws.GetFilePath(getUser(*r), fn)
	out, err := os.Create(fp)
	if err != nil {
		fmt.Println(err)
		return hw.Homework{}, errors.New("Problems writing the file. Contact system admins for help")
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Println(err)
		return hw.Homework{}, errors.New("Problems writing the file. Contact system admins for help")

	}
	return hw.Homework{header.Filename, "/download?user=" + getUser(*r) + "&file=" + fn, r.Form["description"][0], r.Form["task"][0], t}, nil
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

func isAdmin(r http.Request) bool {
	cookie := getSessionIdCookie(r)
	fmt.Println("session Cookie is:", cookie.Value)
	name := session.GetAttribute(cookie.Value)
	return user.ContainsUser(name)
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
	if !user.Authenticate(username, password) &&
		!teams.Authenticate(username, password) {
		return false
	}
	if teams.Authenticate(username, password) {
		username = teams.GetTeamId(username)
	}
	val := username + "-" + user.RandomString()
	cookie := http.Cookie{Name: "session.id", Value: val}
	http.SetCookie(w, &cookie)
	session.SetAttribute(val, username)
	fmt.Println("set session cookie", val)
	return true
}
