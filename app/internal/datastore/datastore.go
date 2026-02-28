package datastore

type DataStore struct {
	hashTable map[string]string
}

func New() *DataStore {
	return &DataStore{
		hashTable: make(map[string]string),
	}
}

func (ds *DataStore) Set(key string, value string) {
	ds.hashTable[key] = value
}

func (ds *DataStore) Get(key string) (string, bool) {
	value, ok := ds.hashTable[key]
	return value, ok
}
