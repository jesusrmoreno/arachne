package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Port string `toml:"port"`
	Root string `toml:"root"`
}


var domain Graph
var GraphState = NewGraphSubject(domain)


func main() {
	graphChan := make(chan Graph) // uwu senpai~

	conf := Config{}

	if _, err := toml.DecodeFile("./config.toml", &conf); err != nil {
		log.Fatal("Could not decode config file, make sure it exists and that it defines a port and root directory.")
	}


	go StartGraphStructureService(conf.Root, graphChan)
	go StartServer(conf.Port, conf.Root)
	GraphState.Subscribe(func(g Graph) {

	})
	for {
		GraphState.Next(<-graphChan)
	}
}
