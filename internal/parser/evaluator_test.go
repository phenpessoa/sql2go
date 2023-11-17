package parser

import (
	"reflect"
	"testing"

	"github.com/phenpessoa/sql2go/internal/testdata"
)

func TestParser(t *testing.T) {
	type queries struct {
		Foo            string
		Bar            string
		Baz            string
		Qux            string
		Quux           string
		Corge          string
		Grault         string
		HardToLex      string
		Empty          string
		Garply         string
		Waldo          string
		Fred           string
		Whatif         string
		WhatAboutThis  string
		WhatAboutThis2 string
		Plugh          string
		Xyzzy          string
		Thud           string
	}

	want := queries{
		Foo:            "SELECT * FROM foo;",
		Bar:            "SELECT * FROM bar\nWHERE id = 123;",
		Baz:            "SELECT\n*\nFROM\nbaz\nWHERE\nbaz = 123 AND\nbaz = baz;",
		Qux:            "SELECT * FROM qux;",
		Quux:           "SELECT * FROM quux\nWHERE quux = 123;",
		Corge:          "SELECT '--' FROM corge;",
		Grault:         "SELECT '\n-- name: Grault\n' FROM grault;",
		HardToLex:      "SELECT;",
		Empty:          "",
		Garply:         "SELECT 'garply-string-literal' FROM garply;",
		Waldo:          "SELECT \"waldo_identifier_1\" FROM waldo;",
		Fred:           "SELECT `fred_identifier_2` FROM fred;",
		Whatif:         "SELECT * FROM whatif;",
		WhatAboutThis:  "SELECT 'foo--hard--string--literal' FROM whatAboutThis; `foo\"_identifier_3'`",
		WhatAboutThis2: "SELECT\n`foo\"_identifier_4`",
		Plugh:          "SELECT * FROM plugh",
		Xyzzy:          "SELECT * FROM xyzzy",
		Thud:           "SELECT * FROM thud;\nSELECT * FROM thud2;",
	}

	f, err := testdata.TestFS.Open("files/initial.sql")
	if err != nil {
		t.Errorf("failed to open inital.sql: %s", err)
		t.FailNow()
		return
	}

	var got queries
	if err := Parse(&got, f); err != nil {
		t.Errorf("failed to parse initial.sql: %s", err)
		t.FailNow()
		return
	}

	if got != want {
		t.Error("initial.sql not parsed properly\n")
		typ := reflect.TypeOf(got)
		for i := 0; i < typ.NumField(); i++ {
			f := typ.Field(i)
			fv1 := reflect.ValueOf(got).Field(i).Interface()
			fv2 := reflect.ValueOf(want).Field(i).Interface()

			if fv1 == fv2 {
				continue
			}

			t.Errorf(
				"field: %s\nwanted: %#+v\ngot: %#+v\n",
				f.Name, fv2, fv1,
			)
		}

		t.FailNow()
		return
	}
}
