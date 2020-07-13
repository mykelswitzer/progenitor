package scaffolding

type scaffold interface {
	setupStructure() error
	cloneTemplates() error
}

type Scaffold struct {
	Structure []Structure
	TemplatePath string 
}

type Structure struct {
	Name string
	SubStructure []Structure
}