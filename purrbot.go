package main

import (
	"github.com/op/go-logging"
	"os"
	//	"strings"
	"github.com/XANi/go-yamlcfg"
	"github.com/XANi/purrbot/config"
	"github.com/XANi/purrbot/plugins/git"
)

var version string
var log = logging.MustGetLogger("main")
var stdout_log_format = logging.MustStringFormatter("%{color:bold}%{time:2006-01-02T15:04:05.0000Z-07:00}%{color:reset}%{color} [%{level:.1s}] %{color:reset}%{shortpkg}[%{longfunc}] %{message}")

func main() {
	stderrBackend := logging.NewLogBackend(os.Stderr, "", 0)
	stderrFormatter := logging.NewBackendFormatter(stderrBackend, stdout_log_format)
	logging.SetBackend(stderrFormatter)
	logging.SetFormatter(stdout_log_format)

	log.Debugf("version: %s", version)
	cfgFiles := []string{
		"$HOME/.config/purrbot/config.yaml",
	}
	var cfg config.Config
	err := yamlcfg.LoadConfig(cfgFiles, &cfg)
	c := git.Config{
		SearchPath: []string{"$HOME/src/my/*", "$HOME/src/lib/go/src/github.com/XANi/"},
	}
	log.Noticef("Config: %+v", cfg)
	a, err := git.New(c)
	if err != nil {
		log.Errorf("can't start git plugin: %s", err)
	}
	a.Run()
}
