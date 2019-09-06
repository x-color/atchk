package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"github.com/x-color/atchk/internal/atcoder/contest"
)

var dir, path string

func init() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	dir = filepath.Join(home, ".atchk")
	path = filepath.Join(home, ".atchk", "config.json")
	viper.SetConfigType("json")
	viper.AddConfigPath(dir)
	viper.SetConfigName("config")
}

type Cache struct {
	Contest string            `json:"contest" mapstructure:"contest"`
	Cookies map[string]string `json:"cookies" mapstructure:"cookies"`
	Samples contest.Samples   `json:"samples" mapstructure:"samples"`
	Task    string            `json:"task" mapstructure:"task"`
}

func (cache *Cache) Update() error {
	viper.Set("cache", cache)
	return viper.WriteConfig()
}

func (cache *Cache) Set(key, value string) error {
	viper.Set(key, value)
	return viper.Unmarshal(cache)
}

type Config struct {
	Command  string `json:"cmd" mapstructure:"cmd"`
	LangID   string `json:"lang_id" mapstructure:"lang_id"`
	TestMode bool   `json:"test_mode" mapstructure:"test_mode"`
}

func (conf *Config) String() string {
	return fmt.Sprintf("cmd = \"%s\"\nlang_id = %s\ntest_mode = %v",
		conf.Command, conf.LangID, conf.TestMode)
}

type ConfFile struct {
	Cache *Cache  `json:"cache" mapstructure:"cache"`
	Conf  *Config `json:"config" mapstructure:"config"`
}

func (cf *ConfFile) Read() error {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return fmt.Errorf("%s\n%s", err, "Please execute `atchk init`.")
		}
		return err
	}
	if err := viper.Unmarshal(cf); err != nil {
		return err
	}
	if cf.Cache == nil {
		cf.Cache = new(Cache)
	}
	if cf.Conf == nil {
		cf.Conf = new(Config)
	}
	return nil
}

func CreateNewConfFile() error {
	if err := os.Mkdir(dir, 0644); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(&ConfFile{Conf: &Config{}}, "", "\t")
	if _, err := f.Write(b); err != nil {
		return err
	}
	return nil
}

func NewConfAndCache() (*Config, *Cache, error) {
	cf := ConfFile{}
	if err := cf.Read(); err != nil {
		return nil, nil, err
	}
	return cf.Conf, cf.Cache, nil
}
