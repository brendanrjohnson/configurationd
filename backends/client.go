package backends

// The StoreClient interface is implemented by objects that can retrieve
// key/value pairs from a backend store.
type StoreClient interface {
	GetValues(key []string) (map[string]string, error)
	WatchPrefix(prefix string, waitIndex uint64, stopChan chan bool) (uint64, error)
}
