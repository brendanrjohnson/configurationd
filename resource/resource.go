package resource

import (
	"github.com/kelseyhightower/confd/backends"
)

type Config struct {
	ConfDir     string
	ConfigDir   string
	Prefix      string
	StoreClient backends.StoreClient
	ResourceDir string
}
