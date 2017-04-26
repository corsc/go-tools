package fiximports

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"sort"

	"strings"

	"github.com/corsc/go-tools/commons"
)

const lineBreak = '\n'

// ProcessFiles will process the supplied files and attempt to fix the imports
func ProcessFiles(files []string) {
	for _, filename := range files {
		source, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "skipping file '%s': failed to read with err err: %v\n", filename, err)
			continue
		}

		newCode, err := processFile(filename, source)
		if err != nil {
			fmt.Fprintf(os.Stderr, "skipping file '%s': failed to generate with err err: %v\n", filename, err)
			continue
		}

		err = ioutil.WriteFile(filename, newCode, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "skipping file '%s': failed to write with err err: %v\n", filename, err)
			continue
		}
	}
}

func processFile(filename string, source []byte) ([]byte, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, filename, source, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	visitor := &myVisitor{
		filename: filename,
		fileSet:  fileSet,
		source:   source,
	}
	ast.Walk(visitor, file)

	if visitor.err != nil {
		return nil, visitor.err
	}

	return visitor.output, nil
}

// implements ast.Visitor
type myVisitor struct {
	filename string
	fileSet  *token.FileSet
	source   []byte
	output   []byte
	err      error
}

// Visit implements ast.Visitor
func (v *myVisitor) Visit(node ast.Node) ast.Visitor {
	switch statement := node.(type) {
	case *ast.File:
		v.fixImports(statement)

	default:
		// intentionally do nothing
	}

	return v
}

func (v *myVisitor) fixImports(file *ast.File) {
	startPos := 0
	endPos := 0
	updatedImports := ""

	err := v.preClean(file)
	if err != nil {
		return
	}

	if len(file.Imports) > 0 {
		startPos, endPos = v.getImportBoundaries(file)
		v.orderImports(file)
		updatedImports = v.generateImportsFragment(file)
	}

	v.output = v.replaceImports(file, updatedImports, startPos, endPos)
	err = v.validate(v.output)
	if err != nil {
		v.err = fmt.Errorf("generated code was invalid, err: %s", err)
		return
	}
}

// make sure the imports are valid and clean first
func (v *myVisitor) preClean(file *ast.File) error {
	v.source, v.err = commons.GoFmt(v.source)
	return v.err
}

func (v *myVisitor) getImportBoundaries(file *ast.File) (startPos, endPos int) {
	startPos = math.MaxInt32
	endPos = -1

	for _, thisImport := range file.Imports {
		// Run to end of line in both directions if not at line start/end.
		thisStartPos, thisEndPos := int(thisImport.Pos()), int(thisImport.Pos())+1
		for thisStartPos > 0 && v.source[thisStartPos-1] != lineBreak {
			thisStartPos--
		}

		for thisEndPos < len(v.source) && v.source[thisEndPos-1] != lineBreak {
			thisEndPos++
		}

		if thisStartPos < startPos {
			startPos = thisStartPos
		}

		if thisEndPos > endPos {
			endPos = thisEndPos
		}
	}

	return
}

func (v *myVisitor) orderImports(file *ast.File) {
	sort.Sort(byImportPath(file.Imports))
}

func (v *myVisitor) generateImportsFragment(file *ast.File) string {
	stdLibFragment := ""
	customFragment := ""

	stdLibRegex := regexp.MustCompile(`(")[a-zA-Z/]+(")`)

	// special case: single import
	totalImports := len(file.Imports)

	for _, thisImport := range file.Imports {
		if stdLibRegex.MatchString(thisImport.Path.Value) {
			if totalImports == 1 {
				stdLibFragment += "import "
			} else {
				stdLibFragment += "\t"
			}
			stdLibFragment += v.buildImportLine(thisImport)
		} else {
			if totalImports == 1 {
				customFragment += "import "
			} else {
				customFragment += "\t"
			}
			customFragment += v.buildImportLine(thisImport)
		}
	}

	padding := ""
	if len(stdLibFragment) > 0 && len(customFragment) > 0 {
		padding = string(lineBreak)
	}
	return stdLibFragment + padding + customFragment
}

func (v *myVisitor) buildImportLine(thisImport *ast.ImportSpec) string {
	output := ""

	topComment := strings.TrimSpace(thisImport.Doc.Text())
	if len(topComment) > 0 {
		output += "// " + topComment + string(lineBreak) + "\t"
	}

	if thisImport.Name != nil {
		name := strings.TrimSpace(thisImport.Name.Name)
		if len(name) > 0 {
			output += name + " "
		}
	}

	output += thisImport.Path.Value

	commentAfter := strings.TrimSpace(thisImport.Comment.Text())
	if len(commentAfter) > 0 {
		output += " // " + commentAfter
	}

	return output + string(lineBreak)
}

func (v *myVisitor) replaceImports(file *ast.File, newImports string, startPos, endPos int) []byte {
	var output []byte

	// replace the imports section
	output = append(output, v.source[:startPos]...)
	output = append(output, newImports...)
	output = append(output, v.source[endPos:]...)

	return output
}

// validate the result by running it through GoFmt
func (v *myVisitor) validate(newCode []byte) error {
	// TODO: add "fast" mode that skips this check or remove this when we have handled all the weird cases
	_, err := commons.GoFmt(newCode)
	return err
}

// implements sort.Interface for []*ast.ImportSpec
type byImportPath []*ast.ImportSpec

func (a byImportPath) Len() int {
	return len(a)
}

func (a byImportPath) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a byImportPath) Less(i, j int) bool {
	return a[i].Path.Value < a[j].Path.Value
}
