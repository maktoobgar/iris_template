package app

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/text/language"

	"service/build"
	iconfig "service/config"
	g "service/global"
	"service/pkg/config"
	db "service/pkg/database"
	"service/pkg/logging"
	"service/pkg/translator"
)

var (
	cfg       = &iconfig.Config{}
	languages = []language.Tag{language.English, language.Persian}
)

// Set Project PWD
func setPwd() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	for parent := pwd; true; parent = filepath.Dir(parent) {
		if _, err := os.Stat(filepath.Join(parent, "go.mod")); err == nil {
			cfg.PWD = parent
			break
		}
	}
	os.Chdir(cfg.PWD)
}

// Initialization for config files in configs folder
func initializeConfigs() {
	// Loads default config, you just have to hard code it
	if err := config.ParseYamlBytes(build.Config, cfg); err != nil {
		log.Fatalln(err)
	}

	if err1, err2 := config.Parse(cfg.PWD+"/env.yaml", cfg, false), config.Parse(cfg.PWD+"/env.yml", cfg, false); err1 != nil || err2 != nil {
		if err1 != nil {
			log.Fatalln(err1)
		} else if err2 != nil {
			log.Fatalln(err2)
		}
	}

	g.CFG = cfg
}

// Translator initialization
func initialTranslator() {
	t, err := translator.New(build.Translations, languages[0], languages[1:]...)
	if err != nil {
		log.Fatalln(err)
	}
	g.Translator = t
}

// Run dbs
func initialDBs() {
	var err error
	g.AllSQLCons, g.DB, err = db.New(cfg.Gateway.Databases, cfg.Debug)
	if err != nil {
		log.Fatalln(err)
	}

	var ok bool = false
	if !g.CFG.Debug {
		_, ok = g.AllSQLCons["main"]
		if !ok {
			log.Fatalln(errors.New("'main' db is not defined (required)"))
		}
	} else {
		_, ok = g.AllSQLCons["test"]
		if !ok {
			log.Fatalln(errors.New("'test' db is not defined"))
		}
	}
}

// Logger initialization
func initialLogger() {
	cfg.Logging.Path += "/" + g.Name
	k := cfg.Logging
	opt := logging.Option(k)
	l, err := logging.New(&opt, cfg.Debug)
	if err != nil {
		log.Fatalln(err)
	}
	g.Logger = l
}

// Server initialization
func init() {
	setPwd()
	initializeConfigs()
	initialDBs()
	initialTranslator()
	initialLogger()
}
