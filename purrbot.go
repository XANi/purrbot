package main

import (
	"fmt"
	"github.com/op/go-logging"

	//	"strings"
	"github.com/XANi/go-yamlcfg"
	"github.com/XANi/purrbot/config"
	"github.com/XANi/purrbot/plugins/git"
	"github.com/XANi/purrbot/utils"
	"github.com/rivo/tview"
)

var version string
var log = logging.MustGetLogger("main")
var stdout_log_format = logging.MustStringFormatter("%{color:bold}%{time:2006-01-02T15:04:05.0000Z-07:00}%{color:reset}%{color} [%{level:.1s}] %{color:reset}%{shortpkg}[%{longfunc}] %{message}")

func main() {

	// TODO term detection ?
	utils.UpdateXtermTitle(fmt.Sprintf("Purrbot v%s", version))

	log.Debugf("version: %s", version)
	cfgFiles := []string{
		"$HOME/.config/purrbot/config.yaml",
	}
	var cfg config.Config
	err := yamlcfg.LoadConfig(cfgFiles, &cfg)
	if err != nil {
		log.Panicf("Config error: %s", err)
	}

	app := tview.NewApplication()
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(false).
		SetChangedFunc(func() {
			app.Draw()
		})
	textView.SetBorder(true).SetTitle("log")
	w := tview.ANSIWriter(textView)
	stderrBackend := logging.NewLogBackend(w, "", 0)
	stderrFormatter := logging.NewBackendFormatter(stderrBackend, stdout_log_format)
	logging.SetBackend(stderrFormatter)
	logging.SetFormatter(stdout_log_format)

	log.Noticef("Config: %+v", cfg)
	if pluginCfg, ok := cfg.Plugins["git"]; ok {
		gp, err := git.New(pluginCfg)
		if err != nil {
			log.Errorf("can't start git plugin: %s", err)
		}
		go gp.Run()
		_ = gp
	}

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("top"), 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("left"), 0, 1, false).
			//AddItem(tview.NewBox().SetBorder(true).SetTitle("Middle (3 x height of Top)"), 0, 3, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("right"), 0, 2, false), 0, 2, false).
		AddItem(textView, 10, 1, false)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}

}
