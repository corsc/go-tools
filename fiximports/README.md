This tool aims to properly format Go imports

## Examples

* `fiximports` - Fix all the Go files in the current directory
* `fiximports ./...` - Fix all the Go files in the current directory and all sub directories
* `fiximports ./dir1/...` - Fix all the Go files in the specific directory and all sub directories
* `fiximports file1.go file2.go` - Fix only the supplied Go file(s)

### Notes:
* This tool does not call `goimports`  (could be added; ping me if this would be useful)
* This tool does not attempt to resolve missing or unnecessary deps
