package git
import (
	"github.com/op/go-logging"
	"github.com/XANi/go-gitcli"
	"path/filepath"
	"os"
	"strings"
	"time"
)

var log = logging.MustGetLogger("main")

type Config struct {
	// Either repo itself or directory containing repos
	SearchPath []string `yaml:"search_path"`
	//update interval (in seconds). Key is substring of path, value is interval
	UpdateInterval map[string]int `yaml:"update_interval"`
	DefaultInterval int `yaml:"default_interval"`
	valid bool
}

type repo struct {
	updateInterval int
	repo *gitcli.Repo

}

type Plugin struct {
	repos map[string]*repo
}

func (g Plugin)add(path string) {
	log.Noticef("adding git repo in %s", path)
	r := gitcli.New(path,path)
	g.repos[path] = &repo{
		repo: &r,
	}
}



func New(rawCfg map[string]interface{}) (*Plugin, error) {
	g := Plugin{
		repos: make(map[string]*repo),
	}
	c := parseCfg(rawCfg)
	log.Errorf("%+v",c)
	for _, path :=  range c.SearchPath {
		fullPath,_ := filepath.Abs(os.ExpandEnv(path))
		// if main path is git dir, stop search there
		if looksLikeGitDir(fullPath) {
			g.add(fullPath)
			continue
		}
		// also try .git dir
		if looksLikeGitDir(fullPath + "/.git") {
			g.add(fullPath + "/.git")
			continue
		}
		if strings.HasSuffix(fullPath,"/") {
			fullPath = fullPath + "*"
		}
		files, err :=  filepath.Glob(fullPath)
		if err != nil {
			log.Warningf("Error in path [%s|%s]: %s", path, fullPath, err)
			continue
		}
		for _, file := range files {
//			log.Debugf("Trying %s",file)
			if looksLikeGitDir(file) {
				g.add(file)
			}
			if looksLikeGitDir(file + "/.git") {
				g.add(file + "/.git")
			}
		}
	}
	defaultInterval := 1200
	if c.DefaultInterval > 1 {
		defaultInterval = c.DefaultInterval
	}
	for k, _ := range g.repos {
		g.repos[k].updateInterval = defaultInterval
	}
	if c.UpdateInterval != nil {
		for pattern, interval := range c.UpdateInterval {
			for repoName,_ := range g.repos {
				if strings.Contains(repoName,pattern) {
					if g.repos[repoName].updateInterval == 0 || g.repos[repoName].updateInterval > interval {
						g.repos[repoName].updateInterval = interval
					}
				}
			}
		}
	}
	log.Noticef("Default git repo update interval: %d",defaultInterval)

	return &g, nil
}

func (p Plugin)Run() error {
	var nextRun time.Time
	for {
		nextRun = time.Now().Add(time.Duration(600) * time.Second)
		// TODO check if repo still exists
		for repoName, repo := range p.repos {
			log.Noticef("Fetching %s", repoName)
			err := repo.repo.Fetch()
			if err != nil {
				log.Errorf("Error when fetching [%s]: %s",repoName, err)
			}
		}
		t := time.Now()
		log.Noticef("next update in %s", nextRun.Sub(t))
		for nextRun.After(t) {
			diff := nextRun.Sub(t)
			// cap sleeping at 300s in case date changes between ticks
			if diff > time.Second * 300  {
				time.Sleep(time.Second * 10)
			} else {
				time.Sleep(diff)
			}
			t = time.Now()
		}
	}
	return nil
}

func looksLikeGitDir(path string) bool {
	if f, err := os.Stat(path + "/refs"); err != nil || !f.Mode().IsDir() {return false }
	if f, err := os.Stat(path + "/config"); err != nil || !f.Mode().IsRegular() {return false }
	if f, err := os.Stat(path + "/objects"); err != nil || !f.Mode().IsDir() {return false }
	return true
}
