package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

var contestsFile string

func init() {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	contestsFile = filepath.Join(home, ".atcoder_contests")
}

func readSamplesFromFile(name string, contests Contests) error {
	return loadContestsFile(contests)
}

func readSamplesFromWeb(name string, contests Contests) error {
	url, err := makeURLForTasks(name)
	if err != nil {
		return err
	}
	data, err := getHTML(url)
	if err != nil {
		return err
	}

	inputs, outputs := filterForInput(data), filterForOutput(data)
	if len(inputs) != len(outputs) {
		return errors.New("read contest from web: not equal number of input and output")
	}

	for i := range inputs {
		sample := &Sample{
			Name:   fmt.Sprintf("Sample %d", i+1),
			Input:  inputs[i],
			Output: outputs[i],
		}
		contests[name].Samples = append(contests[name].Samples, sample)
	}

	// Save new contest data to file in the background
	go saveContestsFile(contests)

	return nil
}

func getHTML(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", errors.New("invalid url: not found " + url)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func filterForURL(data, name string) string {
	raw := strings.Split(data, ">"+strings.ToUpper(name[6:7])+"</a>")[0]
	list := strings.Split(raw, "href=")
	return strings.Trim(list[len(list)-1], "'\"")
}

func filterForIO(data, word string) []string {
	list := strings.Split(data, "<pre")
	r := []string{}
	for i, s := range list {
		if strings.Contains(s, word) {
			raw := strings.Split(strings.SplitN(list[i+1], ">", 2)[1], "</pre>")[0]
			sample := strings.Replace(raw, "\r\n", "\n", -1)
			r = append(r, strings.TrimPrefix(sample, "\n"))
		}
	}
	return r
}

func filterForInput(data string) []string {
	return filterForIO(data, "入力例")
}

func filterForOutput(data string) []string {
	return filterForIO(data, "出力例")
}

func makeURLForContest(name string) string {
	return fmt.Sprintf("https://atcoder.jp/contests/%s/tasks", name[:6])
}

func makeURLForTasks(name string) (string, error) {
	data, err := getHTML(makeURLForContest(name))
	if err != nil {
		return "", err
	}
	return "https://atcoder.jp" + filterForURL(data, name), err
}

func loadContestsFile(contests Contests) error {
	data, err := ioutil.ReadFile(contestsFile)
	if os.IsNotExist(err) {
		// Ignore errors of parsing to json and writing file,
		// because it does not affect main process.
		go ioutil.WriteFile(contestsFile, []byte("{}"), 0644)
		return nil
	}
	return json.Unmarshal(data, &contests)
}

func saveContestsFile(contests Contests) error {
	data, err := json.Marshal(contests)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(contestsFile, data, 0644)
}
