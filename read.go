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
	data, err := getHTML(makeURL(strings.ToLower(name)))
	if err != nil {
		return err
	}

	inputs, outputs := inputFilter(data), outputFilter(data)
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
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func filter(data, word string) []string {
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

func inputFilter(data string) []string {
	return filter(data, "入力例")
}

func outputFilter(data string) []string {
	return filter(data, "出力例")
}

func makeURL(name string) string {
	url := fmt.Sprintf("https://atcoder.jp/contests/%s/tasks/%s_", name[:6], name[:6])
	// abc019 ~ abc001 urls are old style. (e.g. abc001a => 'https://.../abc001_1')
	if name[:3] == "abc" && name[3:6] < "020" {
		url += fmt.Sprintf("%d", name[6]-byte('a')+1)
	} else {
		url += string(name[6])
	}
	return url
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
