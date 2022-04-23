package main

import "text/template"

type FloatCompareCasesGenerator struct{}

func (g FloatCompareCasesGenerator) Template() *template.Template {
	return floatCmpCasesTmpl
}

func (g FloatCompareCasesGenerator) Data() any {
	return struct {
		Bits          []string
		Pkgs          []string
		VarSets       [][]any
		InvalidChecks []Check
		ValidChecks   []Check
	}{
		Bits: []string{"32", "64"},
		Pkgs: []string{"assert", "require"},
		VarSets: [][]any{
			{"a", "1.01"},
			{"1.01", "a"},
			{"b", "bb"},
			{"bb", "b"},
			{"s.c", "h.Calculate()"},
			{"d", "e"},
			{"e", "d"},
			{"(*f).c", "*g"},
			{"*g", "f.c"},
			{"h.Calculate()", "floatOp()"},
			{"floatOp()", "s.c"},
		},
		InvalidChecks: []Check{
			{Fn: "Equal", ArgsTmpl: "t, %s, %s", ReportedMsg: "use %s.InDelta"},
			{Fn: "NotEqual", ArgsTmpl: "t, %s, %s", ReportedMsg: "use %s.InDelta"},
			{Fn: "Greater", ArgsTmpl: "t, %s, %s", ReportedMsg: "use %s.InDelta"},
			{Fn: "GreaterOrEqual", ArgsTmpl: "t, %s, %s", ReportedMsg: "use %s.InDelta"},
			{Fn: "Less", ArgsTmpl: "t, %s, %s", ReportedMsg: "use %s.InDelta"},
			{Fn: "LessOrEqual", ArgsTmpl: "t, %s, %s", ReportedMsg: "use %s.InDelta"},

			{Fn: "True", ArgsTmpl: "t, %s == %s", ReportedMsg: "use %s.InDelta"},
			{Fn: "True", ArgsTmpl: "t, %s != %s", ReportedMsg: "use %s.InDelta"},
			{Fn: "True", ArgsTmpl: "t, %s > %s", ReportedMsg: "use %s.InDelta"},
			{Fn: "True", ArgsTmpl: "t, %s >= %s", ReportedMsg: "use %s.InDelta"},
			{Fn: "True", ArgsTmpl: "t, %s < %s", ReportedMsg: "use %s.InDelta"},
			{Fn: "True", ArgsTmpl: "t, %s <= %s", ReportedMsg: "use %s.InDelta"},

			{Fn: "False", ArgsTmpl: "t, %s == %s", ReportedMsg: "use %s.InDelta"},
			{Fn: "False", ArgsTmpl: "t, %s != %s", ReportedMsg: "use %s.InDelta"},
			{Fn: "False", ArgsTmpl: "t, %s > %s", ReportedMsg: "use %s.InDelta"},
			{Fn: "False", ArgsTmpl: "t, %s >= %s", ReportedMsg: "use %s.InDelta"},
			{Fn: "False", ArgsTmpl: "t, %s < %s", ReportedMsg: "use %s.InDelta"},
			{Fn: "False", ArgsTmpl: "t, %s <= %s", ReportedMsg: "use %s.InDelta"},
		},
		ValidChecks: []Check{
			{Fn: "InDelta", ArgsTmpl: "t, %s, %s, 0.001"},
		},
	}
}

var floatCmpCasesTmpl = template.Must(template.New("floatCmpCasesTmpl").
	Funcs(template.FuncMap{
		"ExpandCheck": ExpandCheck,
	}).
	Parse(`// Code generated by testifylint/internal/cmd/testgen. DO NOT EDIT.

package basic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)
{{ range $bi, $bits := .Bits }}
func TestFloat{{ $bits }}Compare(t *testing.T) {
	type number float{{ $bits }}
	type withFloat{{ $bits }} struct{ c float{{ $bits }} }
	floatOp := func() float{{ $bits }} { return 0. }

	var a float{{ $bits }}
	var b, bb number
	var s withFloat{{ $bits }}
	d := float{{ $bits }}(1.01)
	const e = float{{ $bits }}(2.02)
	f := new(withFloat{{ $bits }})
	var g *float{{ $bits }}
	var h withFloat{{ $bits }}Method

	{{ range $pi, $pkg := $.Pkgs }}
	t.Run("{{ $pkg }}", func(t *testing.T) {
		{{- range $vi, $vars := $.VarSets }}
		{
			{{- range $ci, $check := $.InvalidChecks }}
			{{ ExpandCheck $check $pkg $vars }}
			{{ end }}}
		{{ end }}
		// Valid {{ $pkg }}s.
		{{ range $vi, $vars := $.VarSets }}
		{
			{{- range $ci, $check := $.ValidChecks }}
			{{ ExpandCheck $check $pkg $vars }}
			{{ end }}}
		{{ end -}}
	})
	{{ end }}
}

type withFloat{{ $bits }}Method struct{}

func (withFloat{{ $bits }}Method) Calculate() float{{ $bits }} { return 0. }
{{ end }}
`))
