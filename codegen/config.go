package codegen

// Prompt is a prompt for the user to enter.
type Prompt struct {
	Name    string `yaml:"name"`
	Message string `yaml:"message"`
}

// InOut is a mapping of input to output.
type InOut struct {
	Input  string `yaml:"input"`
	Output string `yaml:"output"`
}

// Config is the configuration for generating code.
type Config struct {
	Name    string   `yaml:"name"`
	Prompts []Prompt `yaml:"prompts"`
	InOuts  []InOut  `yaml:"mappings"`
}
