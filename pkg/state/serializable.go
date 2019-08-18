package state

type Serializable interface {
	Serialize() (string, error)
	Deserialize(s string) error
}
