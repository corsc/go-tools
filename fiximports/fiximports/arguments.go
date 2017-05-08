package fiximports

import (
	"flag"
	"fmt"
	"strings"

	"github.com/corsc/go-tools/commons"
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
	return commons.GetGoFilesFromCurrentDir()
}

type singleArg struct{}

func (f *singleArg) FileNames() ([]string, error) {
	arg := flag.Arg(0)
	if strings.HasSuffix(arg, "/...") && commons.IsDir(arg[:len(arg)-4]) {
		return commons.GetGoFilesFromDirectoryRecursive(arg)
	}

	if commons.IsDir(arg) {
		return commons.GetGoFilesFromDir(arg)
	}

	if commons.FileExists(arg) {
		return commons.GetGoFiles(arg)
	}

	return nil, fmt.Errorf("'%s' did not resolve to a directory or file", arg)
}

type unknownArgs struct{}

func (f *unknownArgs) FileNames() ([]string, error) {
	return commons.GetGoFiles(flag.Args()...)
}
