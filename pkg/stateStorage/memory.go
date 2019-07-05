package stateStorage

type MemoryStorage map[string]string

func (ms *MemoryStorage) Set(key string, value string) {
	(*ms)[key] = value
}

func (ms *MemoryStorage) Get(key string) string {
	return (*ms)[key]
}
