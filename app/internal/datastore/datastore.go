package datastore

import "time"

type Value struct {
	value      string
	insertedAt time.Time
	ttl        *time.Duration
}

type DataStore struct {
	hashTable map[string]Value
}

func New() *DataStore {
	return &DataStore{
		hashTable: make(map[string]Value),
	}
}

type SetArgs struct {
	Key   string
	Value string
	Ex    *time.Duration
}

func (ds *DataStore) Set(args SetArgs) {
	ds.hashTable[args.Key] = Value{
		value:      args.Value,
		insertedAt: time.Now(),
		ttl:        args.Ex,
	}
}

func (ds *DataStore) Get(key string) (string, bool) {
	value, ok := ds.hashTable[key]
	if !ok || (value.ttl != nil && time.Now().After(value.insertedAt.Add(*value.ttl))) {
		return "", false
	}
	return value.value, ok
}
