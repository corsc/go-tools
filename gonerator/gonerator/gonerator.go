// Copyright 2017 Corey Scott http://www.sage42.org/
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gonerator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/corsc/go-tools/gonerator/internal/gotools"
	"github.com/corsc/go-tools/gonerator/tmpl"
)

// Gonerator co-ordinates the generation of code
type Gonerator struct {
	buf  bytes.Buffer
	pkg  *Package
	data tmpl.TemplateData
}

// Build is the main method of this struct/program
func (g *Gonerator) Build(dir string, typeName string, templateFile string, outputFile string, extras string, dryRun, noop bool) {
	g.data = tmpl.TemplateData{
		TypeName:     typeName,
		TemplateFile: templateFile,
		OutputFile:   outputFile,

		Fields:  []tmpl.Field{},
		Extras:  []string{},
		Methods: []tmpl.Method{},
	}

	g.buildHeader()

	outputName := g.buildOutputName(dir, outputFile)

	var templateContent string
	var err error

	if noop {
		fmt.Fprintf(os.Stdout, "Gonerating NOOP for %s in file %s with extras [%v]\n", typeName, outputName, extras)
		templateContent = tmpl.NoopTemplate
	} else {
		fmt.Fprintf(os.Stdout, "Gonerating for %s with template %s in file %s with extras [%v]\n", typeName, templateFile, outputName, extras)

		var path string
		if !strings.HasPrefix(templateFile, "/") {
			path = dir + templateFile
		} else {
			path = templateFile
		}

		templateContentRaw, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}

		templateContent = string(templateContentRaw)
	}

	g.buildTemplateData(extras)

	g.generate(templateContent)

	contents, err := g.gonerate(outputFile)
	if dryRun {
		fmt.Fprintf(os.Stdout, "\n%s", string(contents))
	} else {
		if err != nil {
			os.Exit(-1)
		}

		err = g.writeFile(outputName, contents)
		if err != nil {
			os.Exit(-1)
		}
	}
}

func (g *Gonerator) buildTemplateData(extras string) {
	g.findTypeFields()

	g.data.Extras = strings.Split(extras, ",")
	g.data.PackageName = g.pkg.name
}

func (g *Gonerator) findTypeFields() {
	for _, file := range g.pkg.astFiles {
		fields := tmpl.GetFields(file, g.data.TypeName)
		if len(fields) > 0 {
			g.data.Fields = append(g.data.Fields, fields...)
		}
		methods := tmpl.GetMethods(file, g.data.TypeName)
		if len(methods) > 0 {
			g.data.Methods = append(g.data.Methods, methods...)
		}
	}
}

func (g *Gonerator) generate(templateContent string) {
	buffer := &bytes.Buffer{}
	tmpl.Generate(buffer, g.data, templateContent)
	g.printf(buffer.String())
}

func (g *Gonerator) printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

// ParsePackageDir finds all go files in the supplied directory
func (g *Gonerator) ParsePackageDir(directory string) {
	pkg, err := build.Default.ImportDir(directory, 0)
	if err != nil {
		log.Fatalf("[%s] cannot process directory %s: %s", g.data.OutputFile, directory, err)
	}
	var names []string
	names = append(names, pkg.GoFiles...)
	names = prefixDirectory(directory, names)
	g.parsePackage(directory, names)
}

func prefixDirectory(directory string, names []string) []string {
	if directory == "." {
		return names
	}
	ret := make([]string, len(names))
	for i, name := range names {
		ret[i] = filepath.Join(directory, name)
	}
	return ret
}

func (g *Gonerator) parsePackage(directory string, names []string) {
	var files []*File
	var astFiles []*ast.File
	g.pkg = &Package{}
	fs := token.NewFileSet()
	for _, name := range names {
		if !strings.HasSuffix(name, ".go") {
			continue
		}
		parsedFile, err := parser.ParseFile(fs, name, nil, 0)
		if err != nil {
			log.Fatalf("[%s] parsing package: %s: %s", g.data.OutputFile, name, err)
		}
		astFiles = append(astFiles, parsedFile)
		files = append(files, &File{
			file: parsedFile,
			pkg:  g.pkg,
		})
	}
	if len(astFiles) == 0 {
		log.Fatalf("[%s] %s: no buildable Go files", g.data.OutputFile, directory)
	}
	g.pkg.name = astFiles[0].Name.Name
	g.pkg.files = files
	g.pkg.dir = directory
	g.pkg.astFiles = astFiles
}

// format returns the contents of the Gonerator's buffer after processing by gofmt and goimports
func (g *Gonerator) format() ([]byte, error) {
	original := g.buf.Bytes()

	result, err := gotools.GoFmt(original)
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("[%s] warning: internal error: invalid Go gonerated: %s", g.data.OutputFile, err)
		log.Printf("[%s] warning: compile the package to analyze the error", g.data.OutputFile)
		return original, err
	}

	result, err = gotools.GoImports(g.data.OutputFile, result)
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("[%s] warning: internal error: invalid Go gonerated: %s", g.data.OutputFile, err)
		log.Printf("[%s] warning: compile the package to analyze the error", g.data.OutputFile)
		return original, err
	}

	return result, err
}

func (g *Gonerator) toBytes() []byte {
	return g.buf.Bytes()
}

func (g *Gonerator) buildHeader() {
	if isGo(g.data.OutputFile) {
		g.printf("// Code gonerated by \"github.com/corsc/go-tools/gonerator\"\n// DO NOT EDIT\n")
		g.printf("// @" + "generated \n")
		g.printf("//\n")
		g.printf("// Args:\n")
		g.printf("// TypeName: %s\n", g.data.TypeName)
		g.printf("// Template: %s\n", g.data.TemplateFile)
		g.printf("// Destination: %s\n", g.data.OutputFile)
		g.printf("\n")
	}
}

func (g *Gonerator) gonerate(filename string) ([]byte, error) {
	if isGo(filename) {
		return g.format()
	}
	return g.toBytes(), nil
}

func (g *Gonerator) writeFile(outputName string, contents []byte) error {
	directory := filepath.Dir(outputName)
	err := os.MkdirAll(directory, 0700)
	if err != nil {
		log.Fatalf("[%s] error creating destination directory: %s", g.data.OutputFile, err)
	}

	err = ioutil.WriteFile(outputName, contents, 0600)
	if err != nil {
		log.Fatalf("[%s] error writing output: %s", g.data.OutputFile, err)
	}

	return err
}

func (g *Gonerator) buildOutputName(dir, filename string) string {
	return filepath.Join(dir, strings.ToLower(filename))
}

func isGo(filename string) bool {
	return strings.Contains(filename, ".go")
}
