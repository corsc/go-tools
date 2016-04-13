package gonerator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/corsc/go-tools/gonerator/tmpl"
)

// Gonerator co-ordinates the generation of code
type Gonerator struct {
	buf  bytes.Buffer
	pkg  *Package
	data tmpl.TemplateData
}

// Build generates the code based on supplied values (request call to preceeding ParsePackageDir()
func (g *Gonerator) Build(dir string, typeName string, templateFile string, outputFile string) {
	g.buildHeader()

	templateContent, err := ioutil.ReadFile(dir + templateFile)
	if err != nil {
		panic(err)
	}

	g.findTypeFields(typeName)
	g.generate(string(templateContent))
	g.writeFile(dir, outputFile, typeName)
}

func (g *Gonerator) findTypeFields(typeName string) {
	for _, file := range g.pkg.astFiles {
		fields := tmpl.GetFields(file, typeName)
		if len(fields) > 0 {
			g.data = tmpl.TemplateData{
				TypeName: typeName,
				Fields:   fields,
			}
			return
		}
	}
	log.Panicf("failed to locate definition for %s type", typeName)
}

// generate produces the code for the named type.
func (g *Gonerator) generate(templateContent string) {
	buffer := &bytes.Buffer{}
	tmpl.Generate(buffer, g.data, templateContent)
	g.printf(buffer.String())
}

func (g *Gonerator) printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

// ParsePackageDir parses the package residing in the directory.
func (g *Gonerator) ParsePackageDir(directory string) {
	pkg, err := build.Default.ImportDir(directory, 0)
	if err != nil {
		log.Fatalf("cannot process directory %s: %s", directory, err)
	}
	var names []string
	names = append(names, pkg.GoFiles...)
	names = prefixDirectory(directory, names)
	g.parsePackage(directory, names)
}

// prefixDirectory places the directory name on the beginning of each name in the list.
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

// parsePackage analyzes the single package constructed from the named files.
// If text is non-nil, it is a string to be used instead of the content of the file,
// to be used for testing. parsePackage exits if there is an error.
func (g *Gonerator) parsePackage(directory string, names []string) {
	var files []*File
	var astFiles []*ast.File
	g.pkg = new(Package)
	fs := token.NewFileSet()
	for _, name := range names {
		if !strings.HasSuffix(name, ".go") {
			continue
		}
		parsedFile, err := parser.ParseFile(fs, name, nil, 0)
		if err != nil {
			log.Fatalf("parsing package: %s: %s", name, err)
		}
		astFiles = append(astFiles, parsedFile)
		files = append(files, &File{
			file: parsedFile,
			pkg:  g.pkg,
		})
	}
	if len(astFiles) == 0 {
		log.Fatalf("%s: no buildable Go files", directory)
	}
	g.pkg.name = astFiles[0].Name.Name
	g.pkg.files = files
	g.pkg.dir = directory
	g.pkg.astFiles = astFiles
}

// format returns the gofmt-ed contents of the Gonerator's buffer.
func (g *Gonerator) format() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go gonerated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return g.buf.Bytes()
	}
	return src
}

func (g *Gonerator) buildHeader() {
	g.printf("// Code gonerated by \"github.com/corsc/go-tools/gonerator\"\n// DO NOT EDIT\n")
	g.printf("\n")
	g.printf("package %s", g.pkg.name)
	g.printf("\n")
}

func (g *Gonerator) writeFile(dir string, filename string, typeName string) {
	// Format the output.
	src := g.format()

	// Write to file.
	outputName := filepath.Join(dir, strings.ToLower(filename))

	err := ioutil.WriteFile(outputName, src, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}
