package storage

type Storage interface {
	Write([]byte)
	Read() []byte
}