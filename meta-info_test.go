package vcfgo

import (
	. "gopkg.in/check.v1"
)

type MetaInfoSuite struct{}

var _ = Suite(&MetaInfoSuite{})

// This var block holds pairs of strings and expected data structures
// created by parsing the strings. The data structures are complicated to
// construct. Apologies if you must make more of these - grendeloz.

var (
	ms1 string = `##PICKLE=<ID=NS,Number=1,Type=Integer,Description="Number of Samples With Data">`

	mi1 = &metaStructured{
		metaType: `PICKLE`,
		fields: map[string]*Field{
			`ID`:          &Field{`ID`, `NS`, 0, 0},
			`Number`:      &Field{`Number`, `1`, 1, 0},
			`Type`:        &Field{`Type`, `Integer`, 2, 0},
			`Description`: &Field{`Description`, `Number of Samples With Data`, 3, '"'}},
		order:    []string{`ID`, `Number`, `Type`, `Description`},
		ogString: `##PICKLE=<ID=NS,Number=1,Type=Integer,Description="Number of Samples With Data">`}
)

var mstests = []struct {
	input string
	exp   *metaStructured
}{
	{ms1, mi1},
}

//func (s *VCFSuite) TestKvSplitter(c *C) {
//	for _, v := range kvtests {
//		obs, _, err := kvSplitter(v.input)
//		c.Assert(err, IsNil)
//		fields := make(map[string]string)
//		for _, f := range obs {
//			fields[f.Key] = f.Value
//		}
//		c.Assert(fields, DeepEquals, v.exp)
//	}
//}

func (s *MetaInfoSuite) TestMetaStructuredParse(c *C) {
	for _, v := range mstests {
		obs, err := NewMetaStructuredFromString(v.input)
		c.Assert(err, IsNil)
		c.Assert(obs, DeepEquals, v.exp)
		c.Assert(obs.String(), Equals, v.input)

	}
}
