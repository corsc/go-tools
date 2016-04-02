package test

//go:generate gonerator -i=myType -t=template.tmpl -o=mytype_gonerated.go
type myType struct {
	ID      int64
	Name    string
	Balance float64
}
