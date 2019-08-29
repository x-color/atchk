package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

func init() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	viper.SetConfigType("json")
	viper.AddConfigPath(filepath.Join(home, ".atchk"))
	viper.SetConfigName("config")
}

type Submit struct {
	TestMode bool `json:"test_mode"`
}

func (s *Submit) String() string {
	return fmt.Sprintf("commands.submit.test_mode = %v", s.TestMode)
}

type Test struct {
	Cmds     []string `json:"cmds"`
	Parallel bool     `json:"parallel"`
}

func (t *Test) String() string {
	return fmt.Sprintf("commands.test.cmds = %s\ncommands.test.parallel = %v",
		strings.Join(t.Cmds, ", "), t.Parallel)
}

type Commands struct {
	Submit *Submit `json:"submit"`
	Test   *Test   `json:"test"`
}

func (c *Commands) String() string {
	return fmt.Sprintf("%s\n%s", c.Submit, c.Test)
}

type System struct {
	Contest    string            `json:"contest"`
	Cookies    map[string]string `json:"cookies"`
	Language   string            `json:"language"`
	Languageid int               `json:"languageid"`
}

func (sys *System) String() string {
	return fmt.Sprintf("system.contest = %s\nsystem.language = %s",
		sys.Contest, sys.Language)
}

type Config struct {
	Commands *Commands `json:"commands"`
	System   *System   `json:"system"`
}

func (conf *Config) Read() error {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return fmt.Errorf("%s\n%s", err, "Please execute `atchk init`.")
		}
		return err
	}
	return viper.Unmarshal(conf)
}

func (conf *Config) Update() error {
	viper.Set("commands", conf.Commands)
	viper.Set("system", conf.System)
	return viper.WriteConfig()
}

func (conf *Config) Set(key, value string) error {
	viper.Set(key, value)
	return viper.Unmarshal(conf)
}

func (conf *Config) String() string {
	return fmt.Sprintf("%s\n%s", conf.Commands, conf.System)
}
