package vcfgo

import (
	"reflect"
	"testing"
)

func TestKVCreation(t *testing.T) {
	var tests = []struct {
		KV *KV
		Key   string
		Value string
		Index int
		Quote rune
	}{
		{&KV{`ID`, `NS`, 0, 0}, `ID`, `NS`, 0, 0},
		{&KV{`Number`, `1`, 1, 0}, `Number`, `1`, 1, 0},
		{&KV{`Type`, `Integer`, 2, 0}, `Type`, `Integer`, 2, 0},
		{&KV{`Description`, `Number of Samples With Data`, 3, '"'},
			`Description`, `Number of Samples With Data`, 3, '"'},
	}

	for _, v := range tests {
		if v.KV.Key != v.Key {
			t.Errorf("Expected KV.Key %v but got %v\n", v.KV.Key, v.Key)
		}
		if v.KV.Value != v.Value {
			t.Errorf("Expected KV.Value %v but got %v\n", v.KV.Value, v.Value)
		}
		if v.KV.Index != v.Index {
			t.Errorf("Expected KV.Index %v but got %v\n", v.KV.Index, v.Index)
		}
		if v.KV.Quote != v.Quote {
			t.Errorf("Expected KV.Quote %v but got %v\n", v.KV.Quote, v.Quote)
		}
	}
}

// This var block holds pairs of strings and expected data structures
// created by parsing the strings. The data structures are complicated to
// construct. Apologies if you must make more of these - grendeloz.

var (
	mi2 = &MetaLine{
		LineKey:  `PICKLE`,
		MetaType: Structured,
		KVs: map[string]*KV{
			`ID`:          &KV{`ID`, `NS`, 0, 0},
			`Number`:      &KV{`Number`, `1`, 1, 0},
			`Type`:        &KV{`Type`, `Integer`, 2, 0},
			`Description`: &KV{`Description`, `Number of Samples With Data`, 3, '"'}},
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
