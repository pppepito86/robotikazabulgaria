package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"robotikazabulgaria/ws"
	"robotikazabulgaria/teams"
	"robotikazabulgaria/hw"
	"strconv"
	"time"
)

type Task struct {
	Name      string
	Type      string
	Time      time.Time
	Documents []Document
}

type Document struct {
	Link string
	Text string
}

func UploadTask(w http.ResponseWriter, r *http.Request) error {
	task, err := createTask(w, r)
	if err != nil {
		return err
	}
	writeTasks([]Task{task})
	return nil
}

func createTask(w http.ResponseWriter, r *http.Request) (Task, error) {
	r.Body = http.MaxBytesReader(w, r.Body, 20*1024*1024)
	file, header, _ := r.FormFile("file")
	task := Task{}
	task.Name = r.Form["name"][0]
	task.Type = r.Form["type"][0]
	task.Time = time.Now().UTC()
	link := r.Form["link"][0]
	if link != "" {
		defer file.Close()
		fn := task.Type + "_" + strconv.FormatInt(task.Time.UnixNano(), 16) + filepath.Ext(header.Filename)
		fp := ws.GetFilePath("tasks", fn)
		fmt.Println("Path***", fp)
		out, err := os.Create(fp)
		if err != nil {
			fmt.Println(err)
			return task, errors.New("Problems writing the file. Contact system admins for help")
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			fmt.Println(err)
			return task, errors.New("Problems writing the file. Contact system admins for help")
		}
		link = "/docs/" + fn
	}
	document := Document{
		Link: link,
		Text: r.Form["text"][0]}
	task.Documents = []Document{document}
	return task, nil
}

func writeTasks(tasks []Task) {
	file := ws.GetFilePath("tasks.json")
	json, _ := json.Marshal(tasks)
	os.Create(file)
	ioutil.WriteFile(file, json, 0700)
}

type Homework struct {
	Link string
	Description string
}

type TeamHomeworks struct {
	Id string
	Name string
	Homeworks []Homework
}

type JudgeDashboard struct {
	Task string
	Homeworks []TeamHomeworks
}

func GetJudgeDashboard(task string) JudgeDashboard {
	if task != "project" && task != "robot" {
		task = "team"
	}
	jd := JudgeDashboard{Task: task, Homeworks: make([]TeamHomeworks, 0)}
	tt := teams.GetTeams()
	for _, t := range tt {
		th := TeamHomeworks{Id: t.Id, Name: t.Name, Homeworks: make([]Homework, 0)}
		hws := hw.ReadHomeworks(t.Id)
		for _, hw := range hws {
			if hw.Task == task {
				h := Homework{Link:hw.Link, Description:hw.Description}
				th.Homeworks = append(th.Homeworks, h)
			}
		}
		jd.Homeworks = append(jd.Homeworks, th)
	}
	fmt.Println("judge dashboar", jd)
	return jd
}
	
