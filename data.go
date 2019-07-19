package main

import "errors"

type Contests map[string]*Contest

type Contest struct {
	Name    string    `json:"name"`
	Samples []*Sample `json:"sample"`
}

type Sample struct {
	Name   string `json:"name"`
	Input  string `json:"input"`
	Output string `json:"output"`
}

func getContest(name string) (*Contest, error) {
	if err := valifyName(name); err != nil {
		return nil, err
	}

	contests := make(Contests, 0)
	if err := readSamplesFromFile(name, contests); err != nil {
		return nil, err
	}

	var err error
	if _, ok := contests[name]; !ok {
		contests[name] = &Contest{
			Name:    name,
			Samples: make([]*Sample, 0),
		}
		err = readSamplesFromWeb(name, contests)
	}
	return contests[name], err
}

func valifyName(name string) error {
	if len(name) != 7 || (name[:3] != "abc" && name[:3] != "agc") {
		return errors.New("valify " + name + ": invalid contest name")
	}
	return nil
}
