package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/goproxy/goproxy"
)

type Config struct {
	Renames map[string]string
}

func main() {
	proxy := goproxy.New()
	/*
		var (
			rename string
			config string
		)
		flag.StringVar(&rename, "rename", "", "--rename=github.com/xxx:github.com/abc,github.com/123:github.com/def")
		flag.StringVar(&config, "config", "", "--config=config.json")
		flag.Parse()
		if rename != "" {
			items := strings.Split(rename, ",")
			for _, item := range items {
				if pair := strings.Split(item, ":"); len(pair) == 2 {
					proxy.Renames[pair[0]] = pair[1]
				} else {
					panic("invailed args of rename:" + item)
				}
			}
		} else if config != "" {
			f, err := ioutil.ReadFile(config)
			if err != nil {
				panic(err)
			}
			cfg := Config{}
			err = json.Unmarshal(f, &cfg)
			if err != nil {
				panic(err)
			}
			for k, v := range cfg.Renames {
				proxy.Renames[k] = v
			}
		}*/
	os.Setenv("GOPRIVATE", "github.com/rio-cdvr")
	os.Setenv("GOMODCACHE", "/tmp/goproxy/pkg/mod/")
	os.Setenv("GOSUMDB", "sum.golang.google.cn")
	os.Setenv("GOPROXY", "https://goproxy.cn")
	proxy.Renames["github.comcast.com/viper-cog"] = "github.com/rio-cdvr"
	fmt.Println(proxy.GoBinEnv)
	http.ListenAndServe("localhost:8080", proxy)
}
