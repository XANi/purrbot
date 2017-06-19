package git

import (
)


func parseCfg(c map[string]interface{}) *Config {
	var cfg Config
	if searchPath, ok := c["search_path"].([]interface{}); ok {
		if len(searchPath) < 1 {
			log.Errorf("plugin git needs at least search_path array in config!")
			return &cfg
		}
		cfg.SearchPath = make([]string,len(searchPath))
		for i, val := range searchPath {
			var ok bool
			if cfg.SearchPath[i], ok = val.(string); !ok {
				log.Warningf("git path is not a string: %+v", val)
				return &cfg
			}
		}
		return &cfg
	} else {
		log.Errorf("V: %+v %#v", c["search_path"], c["search_path"].([]string))
	}

	return &cfg
}
