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

package fiximports

import (
	"bytes"
	"fmt"
	"github.com/corsc/go-tools/fiximports/fiximports/internal/gotools"
	"github.com/corsc/go-tools/fiximports/fiximports/internal/log"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"sort"
	"strings"
)

const lineBreak = '\n'

// ProcessFiles will process the supplied files and attempt to fix the imports
func ProcessFiles(files []string, outputWriter io.Writer) {
	for _, filename := range files {
		originalCode, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Error("skipping file '%s': failed to read with err err: %v\n", filename, err)
			continue
		}

		newCode, err := processFile(filename, originalCode)
		if err != nil {
			log.Error("skipping file '%s': failed to generate with err err: %v\n", filename, err)
			continue
		}

		if outputWriter == nil {
			updateSourceFile(filename, originalCode, newCode)
		} else {
			outputNewCode(outputWriter, newCode)
		}
	}
}

func updateSourceFile(filename string, originalCode, newCode []byte) {
	if !bytes.Equal(originalCode, newCode) {
		fmt.Fprintf(os.Stdout, "%s\n", filename)
		err := ioutil.WriteFile(filename, newCode, 0)
		if err != nil {
			log.Error("skipping file '%s': failed to write with err err: %v\n", filename, err)
		}
	}
}

func outputNewCode(writer io.Writer, newCode []byte) {
	fmt.Fprintf(writer, "%s", string(newCode))
}

func processFile(filename string, source []byte) ([]byte, error) {
	fileSet := token.NewFileSet()
	// TODO: change to 	parser.ParseDir() and use filter for single file mode
	file, err := parser.ParseFile(fileSet, filename, source, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// special case: no imports
	if len(file.Imports) == 0 {
		return source, nil
	}

	visitor := newVisitor(filename, fileSet, source)
	ast.Walk(visitor, file)

	if visitor.err != nil {
		return nil, visitor.err
	}

	visitor.updateFile(file)

	return visitor.output, nil
}

func newVisitor(filename string, fileSet *token.FileSet, source []byte) *myVisitor {
	return &myVisitor{
		filename: filename,
		fileSet:  fileSet,
		source:   source,
		startPos: math.MaxInt32,
		endPos:   -1,
	}
}

// implements ast.Visitor
type myVisitor struct {
	filename string
	fileSet  *token.FileSet
	source   []byte
	output   []byte
	err      error
	startPos int
	endPos   int
}

// Visit implements ast.Visitor
func (v *myVisitor) Visit(node ast.Node) ast.Visitor {
	switch statement := node.(type) {
	case *ast.GenDecl:
		v.detectImportDecl(statement)

	default:
		// intentionally do nothing
	}

	return v
}

func (v *myVisitor) updateFile(file *ast.File) {
	updatedImports := ""

	if len(file.Imports) > 0 {
		v.orderImports(file)
		updatedImports = v.generateImportsFragment(file)
	}
	v.output = v.replaceImports(updatedImports)

	err := v.validate(v.output)
	if err != nil {
		v.err = fmt.Errorf("generated code was invalid, err: %s", err)
		return
	}
}

func (v *myVisitor) orderImports(file *ast.File) {
	sort.Sort(byImportPath(file.Imports))
}

func (v *myVisitor) generateImportsFragment(file *ast.File) string {
	stdLibFragment := ""
	customFragment := ""

	stdLibRegex := regexp.MustCompile(`(")[a-zA-Z0-9/]+(")`)

	for _, thisImport := range file.Imports {
		if stdLibRegex.MatchString(thisImport.Path.Value) {
			stdLibFragment += "\t" + v.buildImportLine(thisImport)
		} else {
			customFragment += "\t" + v.buildImportLine(thisImport)
		}
	}

	padding := ""
	if len(stdLibFragment) > 0 && len(customFragment) > 0 {
		padding = string(lineBreak)
	}

	output := "import (" + string(lineBreak)
	output += stdLibFragment + padding + customFragment
	output += ")" + string(lineBreak)

	return output
}

func (v *myVisitor) buildImportLine(thisImport *ast.ImportSpec) string {
	output := ""

	topComment := strings.TrimSpace(thisImport.Doc.Text())
	if len(topComment) > 0 {
		output += "// " + topComment + string(lineBreak) + "\t"
	}

	if thisImport.Name != nil {
		name := strings.TrimSpace(thisImport.Name.Name)

		// remove redundant names
		if !v.nameIsRedundant(name, thisImport.Path.Value) {
			if len(name) > 0 {
				output += name + " "
			}
		}
	}

	output += thisImport.Path.Value

	commentAfter := strings.TrimSpace(thisImport.Comment.Text())
	if len(commentAfter) > 0 {
		output += " // " + commentAfter
	}

	return output + string(lineBreak)
}

func (v *myVisitor) nameIsRedundant(name string, path string) bool {
	// case: `import io "io"`
	// compare name and path after trimming the quotes
	if name == path[1:len(path)-1] {
		return true
	}

	// case: `import proto "github.com/golang/protobuf/proto"`
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash > -1 {
		// trim the slash and the quotes
		pkgDir := path[lastSlash+1 : len(path)-1]
		if name == pkgDir {
			return true
		}
	}

	return false
}

func (v *myVisitor) replaceImports(newImports string) []byte {
	var output []byte

	// replace the imports section
	output = append(output, v.source[:v.startPos]...)
	output = append(output, newImports...)
	output = append(output, v.source[v.endPos:]...)

	return output
}

// validate the result by running it through GoFmt
func (v *myVisitor) validate(newCode []byte) error {
	// TODO: add "fast" mode that skips this check or remove this when we have handled all the weird cases
	_, err := gotools.GoFmt(newCode)
	return err
}

func (v *myVisitor) detectImportDecl(decl *ast.GenDecl) {
	if decl.Tok != token.IMPORT {
		return
	}

	thisStartPos, thisEndPos := gotools.GetLineBoundary(v.source, decl.Pos())
	if thisStartPos < v.startPos {
		v.startPos = thisStartPos
	}

	if decl.Rparen.IsValid() {
		// override with `)` if exists
		// NOTE: add 1 for the line break
		thisEndPos = int(decl.Rparen) + 1
	}

	if thisEndPos > v.endPos {
		v.endPos = thisEndPos
	}
}
