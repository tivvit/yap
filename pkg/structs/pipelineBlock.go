package structs

type PipelineBlock interface {
	Run(state State)
}

type PipelineBlockBase struct {
	Deps     []string  `yaml:"deps"`
	Parent   *Pipeline `yaml:"-"`
	FullName string    `yaml:"-"`
	DepsFull []string  `yaml:"-"`
}
