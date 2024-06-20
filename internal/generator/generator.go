package generator

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"os"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"

	"github.com/karamaru-alpha/mocka/internal/types"
)

//go:embed mock.go.tpl
var rawTemplate []byte
var parsedTemplate, _ = template.New("mock").Funcs(template.FuncMap{"add": func(n int) int { return n + 1 }}).Parse(string(rawTemplate))

type data struct {
	Package     string
	Interfaces  []*types.Interface
	ImportPaths []string
}

func Generate(pkg *packages.Package, fileMap map[string]*types.File) error {
	outputDir := strings.Join([]string{"mocks", pkg.PkgPath}, "/")
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return err
	}
	for fileName, fileInfo := range fileMap {
		filePath := strings.Join([]string{outputDir, fmt.Sprintf("mock_%s", fileName)}, "/")
		if err := generate(filePath, &data{
			Package:     pkg.Name,
			Interfaces:  fileInfo.Interfaces,
			ImportPaths: fileInfo.ImportPathSet.List(),
		}); err != nil {
			return err
		}
	}
	return nil
}

func generate(filePath string, data *data) error {
	var buf bytes.Buffer
	if err := parsedTemplate.Execute(&buf, data); err != nil {
		return err
	}
	source, err := imports.Process(filePath, buf.Bytes(), &imports.Options{
		Fragment:   true,
		AllErrors:  false,
		Comments:   true,
		TabIndent:  true,
		TabWidth:   8,
		FormatOnly: false,
	})
	if err != nil {
		return err
	}
	if source, err = format.Source(source); err != nil {
		return err
	}
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(source); err != nil {
		return err
	}
	return nil
}
