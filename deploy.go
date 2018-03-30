package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"bytes"
	"flag"
	"encoding/json"
	"io/ioutil"
)

func main() {
	app := LoadConf()
	log.Fatal(http.ListenAndServe("0.0.0.0:"+app.Port, app))
}

type App struct {
	Workdir string `json:"workdir"`
	Shell   string `json:"shell"`
	Port    string `json:"port"`
}

func IsDir(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func LoadConf() *App {
	conf := flag.String("conf", "conf.json", "config file")
	flag.Parse()

	content, err := ioutil.ReadFile(*conf)
	if err != nil {
		log.Fatal(err)
	}

	app := &App{}
	err = json.Unmarshal(content, &app)
	if err != nil {
		log.Fatal(err)
	}
	return app
}

func Pull(app App, project string, branch string) []byte {
	cmd := exec.Command(app.Shell, project, branch)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return out.Bytes()
}

func (app App) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	project_arr := req.Form["project"]
	if len(project_arr) < 1 {
		w.Write([]byte("invalid parameter!"))
		return
	}

	branch_arr := req.Form["branch"]
	var branch string
	if len(branch_arr) < 1 {
		branch = ""
	} else {
		branch = branch_arr[0]
	}

	project := app.Workdir + project_arr[0]

	if exist, _ := IsDir(project); !exist {
		w.Write([]byte("invalid parameter!"))
		return
	}

	w.Write(Pull(app, project, branch))
}
