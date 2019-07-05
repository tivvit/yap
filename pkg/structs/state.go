package structs

type State interface {
	Set(key string, value string)
	Get(key string) string
}