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

type ReporterConf struct {
	Storages []ReporterStorageConf
}

type ReporterStorageConf interface{}

type ReporterStorageConfGeneric struct {
	Type string
}

type ReporterStorageConfStdout struct{}

type ReporterStorageConfJson struct {
	FileName string
}
