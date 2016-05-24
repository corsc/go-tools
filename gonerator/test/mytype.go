package test

// MyType ...
//go:generate gonerator -i=MyType -t=template.tmpl -o=mytype_gonerated.go
type MyType struct {
	ID      int64
	Name    string
	Balance float64
}
