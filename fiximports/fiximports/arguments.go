package fiximports

import (
	"flag"
	"fmt"
	"github.com/corsc/go-tools/fiximports/fiximports/internal/filetools"
	"strings"
)

// FilesFromArgs translates supplied inputs into a list of filenames
type FilesFromArgs interface {
	FileNames() ([]string, error)
}

func FilesFromArgsFactory(numArgs int) FilesFromArgs {
	switch numArgs {
	case 0:
		return &noArgs{}

	case 1:
		return &singleArg{}

	default:
		return &unknownArgs{}
	}
}

type noArgs struct{}

func (f *noArgs) FileNames() ([]string, error) {
	return filetools.GetGoFilesFromCurrentDir()
}

type singleArg struct{}

func (f *singleArg) FileNames() ([]string, error) {
	arg := flag.Arg(0)
	if strings.HasSuffix(arg, "/...") && filetools.IsDir(arg[:len(arg)-4]) {
		return filetools.GetGoFilesFromDirectoryRecursive(arg)
	}

	if filetools.IsDir(arg) {
		return filetools.GetGoFilesFromDir(arg)
	}

	if filetools.FileExists(arg) {
		return filetools.GetGoFiles(arg)
	}

	return nil, fmt.Errorf("'%s' did not resolve to a directory or file", arg)
}

type unknownArgs struct{}

func (f *unknownArgs) FileNames() ([]string, error) {
	return filetools.GetGoFiles(flag.Args()...)
}
