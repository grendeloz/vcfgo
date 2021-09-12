package vcfgo

import (
	"reflect"
	"testing"
)

func TestFieldCreation(t *testing.T) {
	var tests = []struct {
		Field *Field
		Key   string
		Value string
		Index int
		Quote rune
	}{
		{&Field{`ID`, `NS`, 0, 0}, `ID`, `NS`, 0, 0},
		{&Field{`Number`, `1`, 1, 0}, `Number`, `1`, 1, 0},
		{&Field{`Type`, `Integer`, 2, 0}, `Type`, `Integer`, 2, 0},
		{&Field{`Description`, `Number of Samples With Data`, 3, '"'},
			`Description`, `Number of Samples With Data`, 3, '"'},
	}

	for _, v := range tests {
		if v.Field.Key != v.Key {
			t.Errorf("Expected Field.Key %v but got %v\n", v.Field.Key, v.Key)
		}
		if v.Field.Value != v.Value {
			t.Errorf("Expected Field.Value %v but got %v\n", v.Field.Value, v.Value)
		}
		if v.Field.Index != v.Index {
			t.Errorf("Expected Field.Index %v but got %v\n", v.Field.Index, v.Index)
		}
		if v.Field.Quote != v.Quote {
			t.Errorf("Expected Field.Quote %v but got %v\n", v.Field.Quote, v.Quote)
		}
	}
}

// This var block holds pairs of strings and expected data structures
// created by parsing the strings. The data structures are complicated to
// construct. Apologies if you must make more of these - grendeloz.

var (
	ms1 string = `##PICKLE=<ID=NS,Number=1,Type=Integer,Description="Number of Samples With Data">`

	mi1 = &StructuredMeta{
		LineKey: `PICKLE`,
		Fields: map[string]*Field{
			`ID`:          &Field{`ID`, `NS`, 0, 0},
			`Number`:      &Field{`Number`, `1`, 1, 0},
			`Type`:        &Field{`Type`, `Integer`, 2, 0},
			`Description`: &Field{`Description`, `Number of Samples With Data`, 3, '"'}},
		Order:    []string{`ID`, `Number`, `Type`, `Description`},
		OgString: `##PICKLE=<ID=NS,Number=1,Type=Integer,Description="Number of Samples With Data">`}
)

func TestStructuredMetaParse(t *testing.T) {
	var mstests = []struct {
		input string
		exp   *StructuredMeta
	}{
		{ms1, mi1},
	}

	for _, v := range mstests {
		obs, err := NewStructuredMetaFromString(v.input)

		if err != nil {
			t.Errorf("got an error when we should not have: %v", err)
		}
		if eq := reflect.DeepEqual(obs, v.exp); !eq {
			t.Errorf("incorrect structure\n  wanted: %+v\n  got: %+v\n", obs, v.exp)
		}
		if obs.String() != v.input {
			t.Errorf("String() gave string %v but wanted %v", obs.String(), v.input)
		}
	}
}

var (
	mi2 = &MetaLine{
		LineKey:  `PICKLE`,
		MetaType: Structured,
		Fields: map[string]*Field{
			`ID`:          &Field{`ID`, `NS`, 0, 0},
			`Number`:      &Field{`Number`, `1`, 1, 0},
			`Type`:        &Field{`Type`, `Integer`, 2, 0},
			`Description`: &Field{`Description`, `Number of Samples With Data`, 3, '"'}},
		Order:    []string{`ID`, `Number`, `Type`, `Description`},
		OgString: `##PICKLE=<ID=NS,Number=1,Type=Integer,Description="Number of Samples With Data">`}

	ms3 string = `##STOOGES=this is another fine kettle of fish`

	mi3 = &MetaLine{
		LineKey:  `STOOGES`,
		MetaType: Unstructured,
		Value:    `##STOOGES=this is another fine kettle of fish`,
		OgString: `##STOOGES=this is another fine kettle of fish`}
)

func TestMetaLine(t *testing.T) {
	var mstests = []struct {
		input string
		exp   *MetaLine
	}{
		{ms1, mi2},
		{ms3, mi3},
	}

	for _, v := range mstests {
		obs, err := NewMetaLineFromString(v.input)
		if err != nil {
			t.Errorf("NewMetaLineFromString() returned an error: %v", err)
		}

		if eq := reflect.DeepEqual(obs, v.exp); !eq {
			t.Errorf("incorrect structure\n  wanted: %+v\n  got: %+v\n", obs, v.exp)
		}

		s, err := obs.String()
		if err != nil {
			t.Errorf("String() returned an error: %v", err)
		}
		if s != v.input {
			t.Errorf("String() gave string %v but wanted %v", s, v.input)
		}
	}
}
