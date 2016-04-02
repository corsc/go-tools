package tmpl

import (
	"bytes"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUATSimple(t *testing.T) {
	tmpl := `row.Scan({{$len := len .Fields }}{{range $index, $value := .Fields}}&in.{{$value.Name}}{{isNotLast $len $index ", "}}{{end}})`

	src := `package test

type myType struct {
	ID      int64
	Name    string
	Balance float64
}
`
	typeName := "myType"
	vars := TemplateData{
		TypeName: typeName,
		Fields:   GetFields(getASTFromSrc(src), typeName),
	}

	masterTmpl, err := getTemplate().Parse(tmpl)
	if err != nil {
		log.Fatal(err)
	}
	buffer := &bytes.Buffer{}
	_ = masterTmpl.Execute(buffer, vars)

	assert.Equal(t, "row.Scan(&in.ID, &in.Name, &in.Balance)", buffer.String())
}
