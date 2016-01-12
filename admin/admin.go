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
