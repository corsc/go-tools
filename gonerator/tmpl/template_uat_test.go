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
		Methods:  GetMethods(getASTFromSrc(src), typeName),
	}

	masterTmpl, err := getTemplate().Parse(tmpl)
	if err != nil {
		log.Fatal(err)
	}
	buffer := &bytes.Buffer{}
	_ = masterTmpl.Execute(buffer, vars)

	assert.Equal(t, "row.Scan(&in.ID, &in.Name, &in.Balance)", buffer.String())
}

func TestUATInterface(t *testing.T) {
	tmpl := `{{ $typeName := .TypeName }}{{ $len := len .Methods }}{{ range $index, $value := .Methods }}func (impl {{ $typeName }}Impl) {{ $value.Name }}({{ $plen := len $value.Params }}{{ range $pindex, $pvalue := $value.Params }}{{ $pplen := len $pvalue.Names }}{{ range $ppindex, $ppname := $pvalue.Names }}{{ $ppname }}{{ isNotLast $pplen $ppindex ", " }}{{ end }} {{ $pvalue.Type }}{{ isNotLast $plen $pindex ", " }}{{ end }}) {{ $rlen := len .Results }}({{ range $rindex, $rvalue := .Results }}{{ $rvalue.Type }}{{ isNotLast $rlen $rindex ", " }}{{ end }}) {}{{ end }}`

	src := `package test

type myType interface {
	LoadByID(ctx context.Context, id int64) tType
}
`
	typeName := "myType"
	vars := TemplateData{
		TypeName: typeName,
		Fields:   GetFields(getASTFromSrc(src), typeName),
		Methods:  GetMethods(getASTFromSrc(src), typeName),
	}

	masterTmpl, err := getTemplate().Parse(tmpl)
	if err != nil {
		log.Fatal(err)
	}
	buffer := &bytes.Buffer{}
	execErr := masterTmpl.Execute(buffer, vars)
	if execErr != nil {
		assert.Fail(t, execErr.Error())
	}

	assert.Equal(t, "func (impl myTypeImpl) LoadByID(ctx context.Context, id int64) (tType) {}", buffer.String())
}
