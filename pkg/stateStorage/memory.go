package stateStorage

type MemoryStorage map[string]string

func (ms *MemoryStorage) Set(key string, value string) {
	(*ms)[key] = value
}

func (ms *MemoryStorage) Get(key string) string {
	return (*ms)[key]
}

func (ms *MemoryStorage) Delete(key string) {
	delete(*ms, key)
}
