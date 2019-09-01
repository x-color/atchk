package internal

import (
	"fmt"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"github.com/x-color/atchk/internal/atcoder/contest"
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

type Cache struct {
	Contest string            `mapstructure:"contest"`
	Cookies map[string]string `mapstructure:"cookies"`
	Samples contest.Samples   `mapstructure:"samples"`
	Task    string            `mapstructure:"task"`
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
	Command  string `mapstructure:"cmd"`
	LangID   int    `mapstructure:"lang_id"`
	TestMode bool   `mapstructure:"test_mode"`
}

func (conf *Config) String() string {
	return fmt.Sprintf("cmd = \"%s\"\nlang_id = %d\ntest_mode = %v",
		conf.Command, conf.LangID, conf.TestMode)
}

type ConfFile struct {
	Cache *Cache  `mapstructure:"cache"`
	Conf  *Config `mapstructure:"config"`
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

func NewConfAndCache() (*Config, *Cache, error) {
	cf := ConfFile{}
	if err := cf.Read(); err != nil {
		return nil, nil, err
	}
	return cf.Conf, cf.Cache, nil
}
