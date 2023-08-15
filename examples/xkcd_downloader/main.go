package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/VisarBuza/taskflow/pkg/taskflow"
)

const urlTemplate = "https://xkcd.com/%s/info.0.json"

type comicTaskFactory struct{}

func (f *comicTaskFactory) Make(line string) taskflow.Task {
	return &ComicTask{url: fmt.Sprintf(urlTemplate, line)}
}

type ComicTask struct {
	url   string
	err   error
	comic *Comic
}

type Comic struct {
	Transcript string `json:"transcript"`
	Image      string `json:"img"`
	Title      string `json:"title"`
	SafeTitle  string `json:"safe_title"`
	Number     int    `json:"num"`
	Day        string `json:"day"`
	Month      string `json:"month"`
	Year       string `json:"year"`
}

func (t *ComicTask) Process() {
	resp, err := http.Get(t.url)
	if err != nil {
		fmt.Fprintln(os.Stderr, "err downloading resource")
		os.Exit(-1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.err = fmt.Errorf("urlTemplate: %s, status code: %d", t.url, resp.StatusCode)
		return
	}
	var cartoon Comic
	if err = json.NewDecoder(resp.Body).Decode(&cartoon); err != nil {
		t.err = err
		return
	}
	t.comic = &cartoon
}

func (t *ComicTask) Print() {
	if t.err != nil {
		fmt.Fprintln(os.Stderr, t.err.Error())
		return
	}
	if v, err := json.Marshal(t.comic); err == nil {
		fmt.Println(string(v))
	}
}

func main() {
	f := &comicTaskFactory{}
	taskflow.Run(f, 500)
}
