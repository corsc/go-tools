package gonerator

import (
	"go/ast"
)

// File holds a single parsed file and associated data.
type File struct {
	pkg  *Package  // Package to which this file belongs.
	file *ast.File // Parsed AST.
}

// Package ...
type Package struct {
	dir      string
	name     string
	files    []*File
	astFiles []*ast.File
}
