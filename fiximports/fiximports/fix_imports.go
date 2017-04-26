package fiximports

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"math"
	"regexp"
	"sort"
)

const lineBreak = "\n"

// ProcessFiles will process the supplied files and attempt to fix the imports
func ProcessFiles(files []string) error {
	for _, filename := range files {
		source, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		newCode, err := processFile(filename, source)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(filename, newCode, 0644)
		if err != nil {
			log.Fatalf("error writing output to file '%s'. err: %v", filename, err)
		}
	}

	return nil
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
	}
	ast.Walk(visitor, file)

	if visitor.err != nil {
		return nil, err
	}

	return visitor.output, nil
}

// implements ast.Visitor
type myVisitor struct {
	filename string
	fileSet  *token.FileSet
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

	if len(file.Imports) > 0 {
		startPos, endPos = v.getImportBoundaries(file)
		v.orderImports(file)
		updatedImports = v.generateImportsFragment(file)
	}

	v.output = v.replaceImports(file, updatedImports, startPos, endPos)
}

func (v *myVisitor) getImportBoundaries(file *ast.File) (startPos, endPos int) {
	startPos = math.MaxInt32
	endPos = -1

	for _, thisImport := range file.Imports {
		thisStartPos := int(thisImport.Pos())
		if thisStartPos < startPos {
			startPos = thisStartPos
		}

		thisEndPos := int(thisImport.Path.End())
		if thisEndPos > endPos {
			endPos = thisEndPos
		}
	}

	// -1 accounts for strange indexing on the imports
	startPos -= 1

	return
}

func (v *myVisitor) orderImports(file *ast.File) {
	sort.Sort(byImportPath(file.Imports))
}

func (v *myVisitor) generateImportsFragment(file *ast.File) string {
	stdLibFragment := ""
	customFragment := ""

	stdLibRegex := regexp.MustCompile(`(")[a-zA-Z\/]+(")`)

	for index, thisImport := range file.Imports {
		if stdLibRegex.MatchString(thisImport.Path.Value) {
			if index > 0 {
				stdLibFragment += "\t"
			}
			stdLibFragment += v.buildImportLine(thisImport)
		} else {
			if index > 0 {
				customFragment += "\t"
			}
			customFragment += v.buildImportLine(thisImport)
		}
	}

	padding := ""
	if len(stdLibFragment) > 0 && len(customFragment) > 0 {
		padding = lineBreak
	}
	return stdLibFragment + padding + customFragment
}

func (v *myVisitor) buildImportLine(thisImport *ast.ImportSpec) string {
	output := ""

	if len(thisImport.Comment.Text()) > 0 {
		output += thisImport.Comment.Text() + lineBreak
	}

	if thisImport.Name != nil {
		output += thisImport.Name.Name + " "
	}

	output += thisImport.Path.Value + lineBreak

	return output
}

func (v *myVisitor) replaceImports(file *ast.File, newImports string, startPos, endPos int) []byte {
	var buf bytes.Buffer
	var output []byte

	// convert AST back to string
	if v.err = format.Node(&buf, v.fileSet, file); v.err != nil {
		return nil
	}

	// replace the imports section
	orig := buf.Bytes()
	output = append(output, orig[:startPos]...)
	output = append(output, newImports...)
	output = append(output, orig[endPos:]...)

	return output
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
