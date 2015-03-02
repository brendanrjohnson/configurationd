package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kelseyhightower/confd/log"
)

type Conf struct {
	Conffile string
	Groups   []Group `toml:"group"`
}

type Group struct {
	Identifier  string
	OptionPairs []optionPair `toml:"pair"`
}

type optionPair struct {
	Key   string `toml:"key"`
	Value string `toml:"value"`
}

func main() {
	flag.Parse()
	if printVersion {
		fmt.Printf("loadconfd %s\n", Version)
		os.Exit(0)
	}
	if err := initConfig(); err != nil {
		log.Fatal(err.Error())
	}
}
