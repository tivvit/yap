package structs

type PipelineBlock interface {
	Run(state State)
	Checkable
}

type Checkable interface {
	Changed(state State) bool
}

type PipelineBlockBase struct {
	Deps     []string  `yaml:"deps"`
	Parent   *Pipeline `yaml:"-"`
	FullName string    `yaml:"-"`
	DepsFull []string  `yaml:"-"`
}
