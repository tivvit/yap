package stateStorage

type State interface {
	Set(key string, value string)
	Get(key string) string
	Delete(key string)
}