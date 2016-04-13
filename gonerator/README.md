This tool provide a simple go code generator using "text/template" as a template engine

## Sample Usage:

* `$ gonerator -i myType -o mytype_gonerated.go -t template.tmpl`

## Sample 

### Template
```
func populate{{firstUpper .TypeName}}(row *sql.Row, in *{{.TypeName}}) {
	row.Scan({{$len := len .Fields }}{{range $index, $value := .Fields}}&in.{{$value.Name}}{{isNotLast $len $index ", "}}{{end}})
}
```

### Type
```
type myType struct {
	ID      int64
	Name    string
	Balance float64
}
```

### Result
```
func populateMyType(row *sql.Row, in *myType) {
	row.Scan(&in.ID, &in.Name, &in.Balance)
}
```

## Sample Using go generate syntax
```
//go:generate gonerator -i=myType -t=template.tmpl -o=mytype_gonerated.go
type myType struct {
	ID      int64
	Name    string
	Balance float64
}
```