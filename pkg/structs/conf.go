package structs

type VisualizeConf struct {
	OutputFile        string
	OutputImage       string
	OutputConnections bool
	PipelineNodes     bool
	PipelineBoxes     bool
	RunDot            bool
	Legend            bool
	Check             bool
}
