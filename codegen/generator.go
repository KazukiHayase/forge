package codegen

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const RootDir = ".forge"

type Data map[string]string

type Generator struct {
	Config Config
	data   Data
}

// NewGenerator returns a new Generator.
// It finds the codegen configuration file by name and marshals it into the Config struct.
func NewGenerator(name string) (Generator, error) {
	var g Generator
	if err := filepath.Walk(
		RootDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			ext := filepath.Ext(path)
			if info.IsDir() || (ext != ".yml" && ext != ".yaml") {
				return nil
			}

			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			var c Config
			if err := yaml.Unmarshal(data, &c); err != nil {
				return err
			}

			if c.Name == name {
				g = Generator{Config: c}
				return filepath.SkipAll
			}

			return nil
		},
	); err != nil {
		return g, err
	}

	if g.Config.Name == "" {
		return g, fmt.Errorf("generator '%s' not found", name)
	}

	return g, nil
}

// BuildData builds the data for the template execution.
func (g *Generator) BuildData(data Data) {
	for _, p := range g.Config.Prompts {
		if _, ok := data[p.Name]; ok {
			continue
		}

		message := p.Message
		if message == "" {
			message = fmt.Sprintf("input %s", p.Name)
		}

		var input string
		fmt.Printf("%s: ", message)
		fmt.Scanln(&input)
		data[p.Name] = input
	}

	g.data = data
}

// ParseInOuts parses the input and output paths by the template execution.
func (g *Generator) ParseInOuts() error {
	data, err := yaml.Marshal(g.Config.InOuts)
	if err != nil {
		return err
	}

	t, err := template.New("").Parse(string(data))
	if err != nil {
		return err
	}

	w := new(bytes.Buffer)
	if err := t.Execute(w, g.data); err != nil {
		return err
	}

	if err := yaml.Unmarshal(w.Bytes(), &g.Config.InOuts); err != nil {
		return err
	}

	return nil
}

// Generate generates the files by the template execution.
func (g *Generator) Generate() error {
	fmt.Println("🔨 Start generating...")
	defer fmt.Println("✨ Done.")

	var inputPaths []string
	for _, m := range g.Config.InOuts {
		inputPaths = append(inputPaths, fmt.Sprintf("%s/%s", RootDir, m.Input))
	}

	tmpls, err := template.ParseFiles(inputPaths...)
	if err != nil {
		return err
	}

	var files []string
	for _, inOut := range g.Config.InOuts {
		operation := "added"
		if _, err := os.Stat(inOut.Output); err == nil {
			var input string
			fmt.Printf("overwrite %s ? [y/N]: ", inOut.Output)
			fmt.Scanln(&input)
			if input != "y" {
				continue
			}
			operation = "updated"
		}

		dir := filepath.Dir(inOut.Output)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
		outputFile, err := os.Create(inOut.Output)
		if err != nil {
			return err
		}
		defer outputFile.Close()

		tmplName := filepath.Base(inOut.Input)
		if err := tmpls.ExecuteTemplate(outputFile, tmplName, g.data); err != nil {
			return err
		}

		files = append(files, fmt.Sprintf("      %s: %s", operation, inOut.Output))
	}

	fmt.Println("📝 Files:")
	fmt.Println(strings.Join(files, "\n"))

	return nil
}
