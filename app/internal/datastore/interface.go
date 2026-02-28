package datastore

type IDataStore interface {
	Set(key string, value string)
	Get(key string) (string, bool)
}
