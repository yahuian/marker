package config

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"regexp"

	_ "embed"

	"github.com/spf13/viper"
	"github.com/yahuian/gox/filex"
)

//go:embed marker.yaml
var defaultFile []byte

// Val global config value
var Val struct {
	SkipFiles  []string `mapstructure:"skip_files,omitempty"`
	ImageTypes []string `mapstructure:"image_types,omitempty"`
}

// TODO write comments

// Init config file
// when there is not config file, auto create a default config in $HOME/.marker.yaml
// otherwise merge user config file and default file
func Init() {
	// check default config file path
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}
	cfgPath := path.Join(home, ".marker.yaml")
	exist, err := filex.Exist(cfgPath)
	if err != nil {
		log.Fatalln(err)
	}

	// create default config file
	if !exist {
		if err := os.WriteFile(cfgPath, defaultFile, 0600); err != nil {
			log.Fatalln(err)
		}
	}

	userViper := viper.New()
	userViper.SetConfigFile(cfgPath)
	if err := userViper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}

	sysViper := viper.New()
	sysViper.SetConfigType("yaml")
	if err := sysViper.ReadConfig(bytes.NewReader(defaultFile)); err != nil {
		log.Fatalln(err)
	}

	// merge config
	newViper := viper.New()
	if exist {
		for key, sysVal := range sysViper.AllSettings() {
			if userVal, ok := userViper.AllSettings()[key]; ok {
				newViper.Set(key, userVal)
			} else { // there will be new config key when software upgrade
				newViper.Set(key, sysVal)
			}
		}
	} else {
		newViper = sysViper
	}

	if err := newViper.Unmarshal(&Val); err != nil {
		log.Fatalln(err)
	}

	if err := validate(); err != nil {
		log.Fatalln(err)
	}

	// write new config
	newViper.SetConfigFile(cfgPath)
	if err := newViper.WriteConfig(); err != nil {
		log.Fatalln(err)
	}
}

func validate() error {
	for _, v := range Val.SkipFiles {
		if _, err := regexp.Compile(v); err != nil {
			return fmt.Errorf("%s is invalid regex: %w", v, err)
		}
	}
	return nil
}

func SkipFiles(d fs.DirEntry) bool {
	for _, v := range Val.SkipFiles {
		if regexp.MustCompile(v).MatchString(d.Name()) {
			return true
		}
	}
	return false
}
