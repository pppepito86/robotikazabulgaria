package admin

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"robotikazabulgaria/hw"
	"robotikazabulgaria/teams"
	"robotikazabulgaria/ws"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	Name        string
	DisplayName string
	Category    string
	Time        time.Time
	Documents   []Document
}

type Document struct {
	Id      string
	Link    string
	DocType string
	Time    time.Time
}

type Challenge struct {
	Id                  string
	Name                string
	State               string
	EndTime             time.Time
	CreateTime          time.Time
	Tasks               []Task
	AdditionalDocuments []Task
}

type Challenges struct {
	ActiveChallenge string
	Challenges      []Challenge
}

type PageChallenges struct {
	CurrentIndex  int
	AllChallenges Challenges
}

func GetPageChallenges(id string) PageChallenges {
	challenges := GetChallenges()
	index := -1
	if len(challenges.Challenges) > 0 {
		index = 0
	}
	for i, element := range challenges.Challenges {
		if element.Id == id {
			index = i
			break
		}
	}
	return PageChallenges{
		CurrentIndex:  index,
		AllChallenges: *challenges,
	}
}

func UpdateChallenge(r *http.Request) {
	h := r.Header.Get("Content-Type")
	if !strings.HasPrefix(h, "multipart") {
		r.ParseForm()
		if len(r.Form["operation"]) != 1 {
			return
		}
		operation := r.Form["operation"][0]
		if operation == "new_challenge" {
			createChallenge(r)
		} else if operation == "challenge_task" {
			createTask1(r)
		} else if operation == "activate_challenge" {
			activateChallenge(r)
		} else if operation == "publish_results" {
			publishResults(r)
		} else if operation == "challenge_additional" {
			createAdditional(r)
		}
		return
	}

	file, header, _ := r.FormFile("file")
	if len(r.Form["operation"]) != 1 {
		return
	}
	operation := r.Form["operation"][0]
	if operation == "task_document" {
		uploadDocument(r, file, header)
	} else if operation == "additional_document" {
		additionalDocument(r, file, header)
	}
}

func activateChallenge(r *http.Request) {
	id := r.Form["challenge"][0]
	challenges := GetChallenges()
	challenges.ActiveChallenge = id
	writeChallenges(challenges)
	updateStatus(id, "")
}

func createAdditional(r *http.Request) {
	if len(r.Form["challenge"]) != 1 ||
		len(r.Form["category"]) != 1 ||
		len(r.Form["name"]) != 1 {
		return
	}
	cc := r.Form["challenge"][0]
	category := r.Form["category"][0]
	ttt := time.Now().UTC()
	id := "extra_" + strconv.FormatInt(ttt.UnixNano(), 16)
	name := r.Form["name"][0]

	t := Task{
		Name:        id,
		DisplayName: name,
		Category:    category,
		Time:        ttt,
		Documents:   make([]Document, 0),
	}
	challenges := GetChallenges()
	idx := -1
	for index, element := range challenges.Challenges {
		if element.Id == cc {
			idx = index
			break
		}
	}
	if idx == -1 {
		return
	}
	ch := &challenges.Challenges[idx]
	ch.AdditionalDocuments = append(ch.AdditionalDocuments, t)
	writeChallenges(challenges)
}

func publishResults(r *http.Request) {
	if len(r.Form["challenge"]) != 1 {
		return
	}
	cc := r.Form["challenge"][0]
	updateStatus(cc, "finished")
}

func updateStatus(cc string, status string) {
	challenges := GetChallenges()
	idx := -1
	for index, element := range challenges.Challenges {
		if element.Id == cc {
			idx = index
			break
		}
	}
	if idx == -1 {
		return
	}
	ch := &challenges.Challenges[idx]
	ch.State = status
	writeChallenges(challenges)
}

func createTask1(r *http.Request) {
	if len(r.Form["challenge"]) != 1 ||
		len(r.Form["category"]) != 1 ||
		len(r.Form["name"]) != 1 {
		return
	}
	cc := r.Form["challenge"][0]
	category := r.Form["category"][0]
	ttt := time.Now().UTC()
	id := "task_" + strconv.FormatInt(ttt.UnixNano(), 16)
	name := r.Form["name"][0]

	t := Task{
		Name:        id,
		DisplayName: name,
		Category:    category,
		Time:        ttt,
		Documents:   make([]Document, 0),
	}
	challenges := GetChallenges()
	idx := -1
	for index, element := range challenges.Challenges {
		if element.Id == cc {
			idx = index
			break
		}
	}
	if idx == -1 {
		return
	}
	ch := &challenges.Challenges[idx]
	ch.Tasks = append(ch.Tasks, t)
	writeChallenges(challenges)
}

func createChallenge(r *http.Request) {
	if len(r.Form["name"]) != 1 ||
		len(r.Form["end_time"]) != 1 {
		return
	}
	ttt := time.Now().UTC()
	id := "challenge_" + strconv.FormatInt(ttt.UnixNano(), 16)
	name := r.Form["name"][0]
	//	startTime := r.Form["start_time"][0]
	endTime := r.Form["end_time"][0]
	timeSplit := strings.Split(endTime, " ")
	dd := strings.Split(timeSplit[0], "-")
	hh := strings.Split(timeSplit[1], ":")
	y, _ := strconv.Atoi(dd[0])
	M, _ := strconv.Atoi(dd[1])
	d, _ := strconv.Atoi(dd[2])
	h, _ := strconv.Atoi(hh[0])
	m, _ := strconv.Atoi(hh[1])
	location, _ := time.LoadLocation("Europe/Sofia")
	deadline := time.Date(y, time.Month(M), d, h, m, 0, 0, location)

	c := Challenge{
		Id:                  id,
		Name:                name,
		CreateTime:          ttt,
		EndTime:             deadline,
		Tasks:               make([]Task, 0),
		AdditionalDocuments: make([]Task, 0),
	}
	challenges := GetChallenges()
	challenges.Challenges = append(challenges.Challenges, c)
	writeChallenges(challenges)
}

func writeChallenges(challenges *Challenges) {
	file := ws.GetFilePath("challenges.json")
	json, _ := json.Marshal(challenges)
	os.Create(file)
	ioutil.WriteFile(file, json, 0700)
}

func GetChallenges() *Challenges {
	var c Challenges
	file := ws.ReadFile("challenges.json")
	err := json.Unmarshal(file, &c)
	if err != nil {
		c = Challenges{
			ActiveChallenge: "",
			Challenges:      make([]Challenge, 0),
		}
	}
	return &c
}

func UploadTask(w http.ResponseWriter, r *http.Request) error {
	task, err := createTask(w, r)
	if err != nil {
		return err
	}
	writeTasks(task)
	return nil
}

func additionalDocument(r *http.Request, file multipart.File, header *multipart.FileHeader) {
	//r.Body = http.MaxBytesReader(w, r.Body, 20*1024*1024)
	challenges := GetChallenges()
	challengeStr := r.Form["challenge"][0]
	challenge := &Challenge{}
	for idx, cc := range challenges.Challenges {
		if cc.Id == challengeStr {
			challenge = &challenges.Challenges[idx]
			break
		}
	}
	task := &Task{}
	for idx, tt := range challenge.AdditionalDocuments {
		if tt.Name == r.Form["task"][0] {
			task = &challenge.AdditionalDocuments[idx]
			break
		}
	}
	link := r.Form["link"][0]
	ttt := time.Now().UTC()
	if len(link) == 0 {
		defer file.Close()
		fn := task.Category + "_" + strconv.FormatInt(ttt.UnixNano(), 16) + filepath.Ext(header.Filename)
		fp := ws.GetFilePath("docs", fn)
		out, err := os.Create(fp)
		if err != nil {
			return
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			return
		}
		link = "/docs/" + fn
	}
	document := Document{
		Id:      "document_" + strconv.FormatInt(ttt.UnixNano(), 16),
		Link:    link,
		DocType: r.Form["type"][0],
		Time:    ttt,
	}
	task.Documents = append(task.Documents, document)
	writeChallenges(challenges)
}

func uploadDocument(r *http.Request, file multipart.File, header *multipart.FileHeader) {
	//r.Body = http.MaxBytesReader(w, r.Body, 20*1024*1024)
	challenges := GetChallenges()
	challengeStr := r.Form["challenge"][0]
	challenge := &Challenge{}
	for idx, cc := range challenges.Challenges {
		if cc.Id == challengeStr {
			challenge = &challenges.Challenges[idx]
			break
		}
	}
	task := &Task{}
	for idx, tt := range challenge.Tasks {
		if tt.Name == r.Form["task"][0] {
			task = &challenge.Tasks[idx]
			break
		}
	}
	link := r.Form["link"][0]
	ttt := time.Now().UTC()
	if len(link) == 0 {
		defer file.Close()
		fn := task.Category + "_" + strconv.FormatInt(ttt.UnixNano(), 16) + filepath.Ext(header.Filename)
		fp := ws.GetFilePath("docs", fn)
		out, err := os.Create(fp)
		if err != nil {
			return
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			return
		}
		link = "/docs/" + fn
	}
	document := Document{
		Id:      "document_" + strconv.FormatInt(ttt.UnixNano(), 16),
		Link:    link,
		DocType: r.Form["type"][0],
		Time:    ttt,
	}
	task.Documents = append(task.Documents, document)
	writeChallenges(challenges)
}

func GetActiveChallenge() Challenge {
	challenges := GetChallenges()
	for _, element := range challenges.Challenges {
		if challenges.ActiveChallenge == element.Id {
			return element
		}
	}
	return Challenge{}
}

func createTask(w http.ResponseWriter, r *http.Request) ([]Task, error) {
	r.Body = http.MaxBytesReader(w, r.Body, 20*1024*1024)
	file, header, _ := r.FormFile("file")
	ttt := GetTasks()
	task := &Task{}
	for idx, tt := range ttt {
		if tt.Name == r.Form["name"][0] {
			task = &ttt[idx]
			break
		}
	}
	task.DisplayName = "pesho"
	if len(task.Name) == 0 {
		task.Name = r.Form["name"][0]
		task.Time = time.Now().UTC()
		ttt = append(ttt, *task)
	}
	if len(r.Form["display_name"][0]) != 0 {
		task.DisplayName = r.Form["display_name"][0]
	}
	if len(r.Form["category"][0]) != 0 {
		task.Category = r.Form["category"][0]
	}
	link := r.Form["link"][0]
	if len(link) == 0 {
		defer file.Close()
		fn := task.Category + "_" + strconv.FormatInt(task.Time.UnixNano(), 16) + filepath.Ext(header.Filename)
		fp := ws.GetFilePath("tasks", fn)
		out, err := os.Create(fp)
		if err != nil {
			return ttt, errors.New("Problems writing the file. Contact system admins for help")
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			return ttt, errors.New("Problems writing the file. Contact system admins for help")
		}
		link = "/docs/" + fn
	}
	document := Document{
		Link:    link,
		DocType: r.Form["doc_type"][0],
		Time:    time.Now().UTC(),
	}
	if task.Documents == nil {
		task.Documents = make([]Document, 0)
	}
	task.Documents = append(task.Documents, document)
	return ttt, nil
}

func writeTasks(tasks []Task) {
	file := ws.GetFilePath("tasks.json")
	json, _ := json.Marshal(tasks)
	os.Create(file)
	ioutil.WriteFile(file, json, 0700)
}

func GetTasks() []Task {
	var t []Task
	file := ws.ReadFile("tasks.json")
	err := json.Unmarshal(file, &t)
	if err != nil {
		t = make([]Task, 0)
	}
	return t
}

type Homework struct {
	Link        string
	Description string
	Extension   string
}

type TeamHomeworks struct {
	Id        string
	Name      string
	Homeworks []Homework
	Mark      string
	Comment   string
}

type JudgeDashboard struct {
	Task             string
	Homeworks        []TeamHomeworks
	CurrentChallenge Challenge
}

func GetJudgeDashboard(username string, task string) JudgeDashboard {
	challenge := GetActiveChallenge()
	if len(challenge.Tasks) == 0 {
		return JudgeDashboard{Homeworks: make([]TeamHomeworks, 0)}
	}
	if task == "" {
		task = challenge.Tasks[0].Name
	}
	jd := JudgeDashboard{Task: task, CurrentChallenge: challenge, Homeworks: make([]TeamHomeworks, 0)}
	tt := teams.GetTeams("0")
	tms := GetTeamMarks(username)

	for _, t := range tt {
		tm := tms[t.Id]
		m := tm.Marks[task]
		th := TeamHomeworks{Id: t.Id, Name: t.Name, Homeworks: make([]Homework, 0), Mark: m.Points, Comment: m.Comment}
		hws := hw.ReadHomeworks(t.Id)
		for _, hw := range hws {
			if hw.Task == task {
				h := Homework{Link: hw.Link, Description: hw.Description}
				if hw.Filename != "" && strings.Contains(hw.Filename, ".") {
					lastIndex := strings.LastIndex(hw.Filename, ".")
					ext := hw.Filename[lastIndex:]
					h.Extension = ext
				}
				th.Homeworks = append(th.Homeworks, h)
			}
		}
		jd.Homeworks = append(jd.Homeworks, th)
	}
	return jd
}

type Mark struct {
	TaskId  string
	Points  string
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
	Id      string
	Name    string
	Stars   []int
	NoStars []int
}

type Results []TeamResults

type DisplayResults struct {
	CurrentChallenge Challenge
	AllResults       Results
}

func (r Results) Len() int {
	return len(r)
}

func (r Results) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r Results) Less(i, j int) bool {
	countI := 0
	for _, element := range r[i].Stars {
		countI += element
	}
	countJ := 0
	for _, element := range r[j].Stars {
		countJ += element
	}
	return countI > countJ
}

func GetCurrentResults() DisplayResults {
	challenge := GetActiveChallenge()
	return GetResults(challenge, "0")
}

func GetFinishedResults(division string) DisplayResults {
	challenges := GetChallenges()
	for i := len(challenges.Challenges) - 1; i >= 0; i-- {
		c := challenges.Challenges[i]
		if c.State == "finished" {
			return GetResults(c, division)
		}
	}
	return DisplayResults{}
}

func GetLastFinishedChallenge() time.Time {
	challenges := GetChallenges()
	for i := len(challenges.Challenges) - 1; i >= 0; i-- {
		c := challenges.Challenges[i]
		if c.State == "finished" {
			return c.EndTime
		}
	}
	return time.Now()
}

func GetResults(challenge Challenge, division string) DisplayResults {
	tmrs := make(Results, 0)
	//tmrs = append(tmrs, TeamResults{Id: "Id", Name: "Otbor", Results: []string{"A", "B", "C"}})

	tt := teams.GetTeams(division)
	tms := GetTeamMarks("pesho")

	for _, t := range tt {
		tm := tms[t.Id]
		tr := TeamResults{Id: t.Id, Name: t.Name, Stars: make([]int, 0), NoStars: make([]int, 0)}

		for _, ttttt := range challenge.Tasks {
			v, _ := strconv.Atoi(tm.Marks[ttttt.Name].Points)
			if v >= 5 {
				v = 5
			}
			tr.Stars = append(tr.Stars, v)
			tr.NoStars = append(tr.NoStars, 5-v)
		}

		tmrs = append(tmrs, tr)
	}
	sort.Sort(tmrs)
	return DisplayResults{
		CurrentChallenge: challenge,
		AllResults:       tmrs,
	}
}
