package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/kelseyhightower/confd/log"
)

var (
	configFile        = ""
	defaultConfigFile = "etc/loadconfd/loadconfd.toml"
	backend           string
	clientCaKeys      string
	clientCert        string
	clientKey         string
	confdir           string
	config            Config
	debug             bool
	interval          int
	nodes             Nodes
	prefix            string
	printVersion      bool
	quiet             bool
	scheme            string
	verbose           bool
)

//A config struct is used to configure loadconfd
type Config struct {
	Backend      string   `toml:"backends"`
	BackendNodes []string `toml:"nodes"`
	ClientCaKeys string   `toml:"client_cakeys"`
	ClientCert   string   `toml:"client_cert"`
	ClientKey    string   `toml:"client_key"`
	ConfDir      string   `toml:"confdir"`
	Debug        bool     `toml:"debug"`
	Interval     int      `toml:"interval"`
	Prefix       string   `toml:"prefix"`
	Quiet        bool     `toml:"quiet"`
	Scheme       string   `toml:"scheme"`
	Verbose      bool     `toml:"verbose"`
}

func init() {
	flag.StringVar(&backend, "backend", "etcd", "backend to use")
	flag.StringVar(&clientCaKeys, "client-ca-keys", "", "client ca keys")
	flag.StringVar(&clientCert, "client-cert", "", "the client cert")
	flag.StringVar(&clientKey, "client-key", "", "the client key")
	flag.StringVar(&confdir, "confdir", "/etc/loadconfd", "loadconfd conf directory")
	flag.StringVar(&configFile, "config-file", "", "the loadconfd config file")
	flag.BoolVar(&debug, "debug", false, "enable debug logging")
	flag.IntVar(&interval, "interval", 600, "backend polling interval")
	flag.Var(&nodes, "node", "list of backend nodes")
	flag.StringVar(&prefix, "prefix", "/", "key path prefix")
	flag.BoolVar(&printVersion, "version", false, "print version and exit")
	flag.BoolVar(&quiet, "quiet", false, "enable quiet logging")
	flag.StringVar(&scheme, "scheme", "http", "the backend URI scheme (http or https)")
	flag.BoolVar(&verbose, "verbose", false, "enable verbose logging")
}

// initConfig initializes the loadconfd configuration by first setting defaults,
// then overriding setting from the loadconfd config file, and finally overriding
// settings from flags set on the command line.
// It returns an error if any.
func initConfig() error {
	if configFile == "" {
		if _, err := os.Stat(defaultConfigFile); !os.IsNotExist(err) {
			configFile = defaultConfigFile
		}
	}
	// Set defaults.
	config = Config{
		Backend:  "etcd",
		ConfDir:  "/etc/loadconfd",
		Interval: 600,
		Prefix:   "/",
		Scheme:   "http",
	}
	// Update config from the TOML configuration file.
	if configFile == "" {
		log.Warning("Skipping loadconfd config file")
	} else {
		log.Debug("Loading " + configFile)
		configBytes, err := ioutil.ReadFile(configFile)
		if err != nil {
			return err
		}
		_, err = toml.Decode(string(configBytes), &config)
		if err != nil {
			return err
		}
	}
	// Update config from commandline flags.
	processFlags()

	// Configure logging.
	log.SetQuiet(config.Quiet)
	log.SetVerbose(config.Verbose)
	log.SetDebug(config.Debug)

	if len(config.BackendNodes) == 0 {
		switch config.Backend {
		case "consul":
			config.BackendNodes = []string{"127.0.0.1:8500"}
		case "etcd":
			peerstr := os.Getenv("ETCDCTL_PEERS")
			if len(peerstr) > 0 {
				config.BackendNodes = strings.Split(peerstr, ",")
			} else {
				config.BackendNodes = []string{"127.0.0.1:4100"}
			}
		}
	}
	// Initialize the storage client

	fmt.Println(config.Backend)
	fmt.Println(config.BackendNodes)
	fmt.Println(config.ConfDir)
	fmt.Println(config.Debug)

	return nil
}

func processFlags() {
	flag.Visit(setConfigFromFlag)
}
func setConfigFromFlag(f *flag.Flag) {
	switch f.Name {
	case "backend":
		config.Backend = backend
	case "debug":
		config.Debug = debug
	case "client-cert":
		config.ClientCert = clientCert
	case "client-key":
		config.ClientKey = clientKey
	case "client-cakeys":
		config.ClientCaKeys = clientCaKeys
	case "confdir":
		config.ConfDir = confdir
	case "node":
		config.BackendNodes = nodes
	case "interval":
		config.Interval = interval
	case "prefix":
		config.Prefix = prefix
	case "quiet":
		config.Quiet = quiet
	case "scheme":
		config.Scheme = scheme
	case "verbose":
		config.Verbose = verbose
	}
}
