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
	Mark string
	Comment string
}

type JudgeDashboard struct {
	Task string
	Homeworks []TeamHomeworks
}

func GetJudgeDashboard(username string, task string) JudgeDashboard {
	if task != "project" && task != "robot" {
		task = "team"
	}
	jd := JudgeDashboard{Task: task, Homeworks: make([]TeamHomeworks, 0)}
	tt := teams.GetTeams()
	tms := GetTeamMarks(username)
	
	for _, t := range tt {
		tm := tms[t.Id]
		m := tm.Marks[task]
		th := TeamHomeworks{Id: t.Id, Name: t.Name, Homeworks: make([]Homework, 0), Mark: m.Points, Comment: m.Comment}
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

type Mark struct {
	TaskId string
	Points string
	Comment string
}

type TeamMark struct {
	Marks map[string]Mark
}

func GetTeamMarks(username string) map[string]TeamMark {
	var m map[string]TeamMark
	file := ws.ReadFile(username, "points.json")
	err := json.Unmarshal(file, &m)
	if err != nil {
		m = make(map[string]TeamMark)
	}
	fmt.Println(m)
	return m
}

func GetTeamMark(username string, teamid string, taskname string) (string, string) { 
	tms := GetTeamMarks(username)
	t := tms[teamid]
	m := t.Marks[taskname]
	return m.Points, m.Comment
}

func writeTeamMarks(username string, m map[string]TeamMark) {
	file := ws.GetFilePath(username, "points.json")
	json, _ := json.Marshal(m)
	os.Create(file)
	ioutil.WriteFile(file, json, 0700)
}
	
func UpdatePoints(r *http.Request, username string) {
	r.ParseForm()
	teamId := r.Form["id"][0]
	task := r.Form["task"][0]
	value := r.Form["value"][0]
	changeType := r.Form["type"][0]
	if teams.GetTeamId(teamId) == "" {
		return
	}
	teamMarks := GetTeamMarks(username)
	tm, exist := teamMarks[teamId]
	if !exist {
		teamMarks[teamId] = TeamMark{Marks: make(map[string]Mark)}
	}
	tm, _ = teamMarks[teamId]
	m, exist := tm.Marks[task]
	if !exist {
		tm.Marks[task] = Mark{TaskId: task}
	}
	m, _ = tm.Marks[task]
	if changeType == "text" {
		m.Points = value
	} else {
		m.Comment = value
	}
	tm.Marks[task] = m
	writeTeamMarks(username, teamMarks)
}

type TeamResults struct {
	Id string
	Name string
	Stars []int
	NoStars []int
}

func GetResults() []TeamResults {
	tmrs := make([]TeamResults, 0)
	//tmrs = append(tmrs, TeamResults{Id: "Id", Name: "Otbor", Results: []string{"A", "B", "C"}})

	tt := teams.GetTeams()
	tms := GetTeamMarks("pesho")
	
	for _, t := range tt {
		tm := tms[t.Id]
		tr := TeamResults{Id: t.Id, Name: t.Name, Stars: make([]int, 0), NoStars: make([]int, 0)}

		v, _ := strconv.Atoi(tm.Marks["team"].Points)
		if v >= 3 {
			v = 3
		}
		tr.Stars = append(tr.Stars, v)
		tr.NoStars = append(tr.NoStars, 3-v)

		v, _ = strconv.Atoi(tm.Marks["project"].Points)
		if v >= 3 {
			v = 3
		}
		tr.Stars = append(tr.Stars, v)
		tr.NoStars = append(tr.NoStars, 3-v)

		v, _ = strconv.Atoi(tm.Marks["robot"].Points)
		if v >= 3 {
			v = 3
		}
		tr.Stars = append(tr.Stars, v)
		tr.NoStars = append(tr.NoStars, 3-v)

		tmrs = append(tmrs, tr)
	}
	fmt.Println("results", tmrs)
	return tmrs
}


