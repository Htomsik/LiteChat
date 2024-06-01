package main

import (
	"Chat/internal/app/Server"
	"flag"
	"github.com/BurntSushi/toml"
	"log"
)

var (
	cfgPath string
)

func init() {
	flag.StringVar(&cfgPath, "config-path", "configs/apiServer.toml", "path to cfg file")
}

// @ title 			Simple site go
// @version         1.0

// @host      localhost:8080
// @BasePath  /
func main() {
	flag.Parse()

	cfg := Server.NewConfig()

	if _, err := toml.DecodeFile(cfgPath, cfg); err != nil {
		log.Fatal(err)
	}

	if err := Server.Start(cfg); err != nil {
		log.Fatal(err)
	}
}
