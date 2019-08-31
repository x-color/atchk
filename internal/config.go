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

type System struct {
	Contest    string            `mapstructure:"contest"`
	Cookies    map[string]string `mapstructure:"cookies"`
	Language   string            `mapstructure:"language"`
	Languageid int               `mapstructure:"languageid"`
}

func (sys *System) String() string {
	return fmt.Sprintf("contest = %s\nlanguage = %s",
		sys.Contest, sys.Language)
}

type Config struct {
	Commands []string `mapstructure:"cmds"`
	System   *System  `mapstructure:"system"`
	TestMode bool     `mapstructure:"test_mode"`
}

func (conf *Config) Read() error {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return fmt.Errorf("%s\n%s", err, "Please execute `atchk init`.")
		}
		return err
	}
	if err := viper.Unmarshal(conf); err != nil {
		return err
	}
	if conf.System == nil {
		conf.System = new(System)
	}
	return nil
}

func (conf *Config) Update() error {
	viper.Set("cmds", conf.Commands)
	viper.Set("system", conf.System)
	viper.Set("test_mode", conf.TestMode)
	return viper.WriteConfig()
}

func (conf *Config) Set(key, value string) error {
	viper.Set(key, value)
	return viper.Unmarshal(conf)
}

func (conf *Config) String() string {
	return fmt.Sprintf("test_mode = %v\ncmds = [%s]\n%s",
		conf.TestMode, strings.Join(conf.Commands, ", "), conf.System)
}
