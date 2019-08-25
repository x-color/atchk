package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type Config struct {
	Commands Commands `json:"commands"`
	System   System   `json:"system"`
}

type Submit struct {
	TestMode bool `json:"test_mode"`
}

type Test struct {
	Cmds     []string `json:"cmds"`
	Parallel bool     `json:"parallel"`
}

type Commands struct {
	Submit Submit `json:"submit"`
	Test   Test   `json:"test"`
}

type System struct {
	Contest    string            `json:"contest"`
	Cookies    map[string]string `json:"cookies"`
	Languageid int               `json:"languageid"`
}

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

func (conf *Config) Read() error {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// TODO: write error message
			return fmt.Errorf("%s\n%s", err, "sample message")
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

func (conf *Config) Format() (string, error) {
	b, err := ioutil.ReadFile(viper.ConfigFileUsed())
	if err != nil {
		return "", err
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "\t")
	return out.String(), err
}
