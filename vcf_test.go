package vcfgo

import (
	"os"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type VCFSuite struct{}

var _ = Suite(&VCFSuite{})

var (
	kv1 = map[string]string{`ID`: `NS`, `Number`: `1`, `Type`: `Integer`,
		`Description`: `Number of Samples With Data`}
	kv2 = map[string]string{`ID`: `DP`, `Number`: `1`, `Type`: `Integer`,
		`Description`: `Total Depth`}
	kv3 = map[string]string{`ID`: `AF`, `Number`: `A`, `Type`: `Float`,
		`Description`: `Allele Frequency`}
)

var kvtests = []struct {
	input string
	exp   map[string]string
}{
	{`ID=NS,Number=1,Type=Integer,Description="Number of Samples With Data"`, kv1},
	{`ID=DP,Number=1,Type=Integer,Description="Total Depth"`, kv2},
	{`ID=AF,Number=A,Type=Float,Description="Allele Frequency"`, kv3},
}

// This var block holds pairs of strings and expected data structures
// created by parsing the strings. The data structures are complicated to
// construct. Apologies if you must make more of these - grendeloz.

var (
	is1 string = `##INFO=<ID=NS,Number=1,Type=Integer,Description="Number of Samples With Data">`

	ii1 = &Info{
		Id:          "NS",
		Number:      "1",
		Type:        "Integer",
		Description: "Number of Samples With Data",
		fields: map[string]*Field{
			`ID`:          &Field{`ID`, `NS`, 0, 0},
			`Number`:      &Field{`Number`, `1`, 1, 0},
			`Type`:        &Field{`Type`, `Integer`, 2, 0},
			`Description`: &Field{`Description`, `Number of Samples With Data`, 3, '"'}},
		order: []string{`ID`, `Number`, `Type`, `Description`}}

	is2 string = `##INFO=<ID=DP,Number=1,Type=Integer,Description="Total Depth">`

	ii2 = &Info{
		Id:          "DP",
		Number:      "1",
		Type:        "Integer",
		Description: "Total Depth",
		fields: map[string]*Field{
			`ID`:          &Field{`ID`, `DP`, 0, 0},
			`Number`:      &Field{`Number`, `1`, 1, 0},
			`Type`:        &Field{`Type`, `Integer`, 2, 0},
			`Description`: &Field{`Description`, `Total Depth`, 3, '"'}},
		order: []string{`ID`, `Number`, `Type`, `Description`}}

	is3 string = `##INFO=<ID=AF,Number=A,Type=Float,Description="Allele Frequency">`

	ii3 = &Info{
		Id:          "AF",
		Number:      "A",
		Type:        "Float",
		Description: "Allele Frequency",
		fields: map[string]*Field{
			`ID`:          &Field{`ID`, `AF`, 0, 0},
			`Number`:      &Field{`Number`, `A`, 1, 0},
			`Type`:        &Field{`Type`, `Float`, 2, 0},
			`Description`: &Field{`Description`, `Allele Frequency`, 3, '"'}},
		order: []string{`ID`, `Number`, `Type`, `Description`}}

	is4 string = `##INFO=<ID=AA,Number=1,Type=String,Description="Ancestral Allele">`

	ii4 = &Info{
		Id:          "AA",
		Number:      "1",
		Type:        "String",
		Description: "Ancestral Allele",
		fields: map[string]*Field{
			`ID`:          &Field{`ID`, `AA`, 0, 0},
			`Number`:      &Field{`Number`, `1`, 1, 0},
			`Type`:        &Field{`Type`, `String`, 2, 0},
			`Description`: &Field{`Description`, `Ancestral Allele`, 3, '"'}},
		order: []string{`ID`, `Number`, `Type`, `Description`}}

	is5 string = `##INFO=<ID=DB,Number=0,Type=Flag,Description="dbSNP membership, build 129">`

	ii5 = &Info{
		Id:          "DB",
		Number:      "0",
		Type:        "Flag",
		Description: "dbSNP membership, build 129",
		fields: map[string]*Field{
			`ID`:          &Field{`ID`, `DB`, 0, 0},
			`Number`:      &Field{`Number`, `0`, 1, 0},
			`Type`:        &Field{`Type`, `Flag`, 2, 0},
			`Description`: &Field{`Description`, `dbSNP membership, build 129`, 3, '"'}},
		order: []string{`ID`, `Number`, `Type`, `Description`}}

	is6 string = `##INFO=<ID=H2,Number=2,Type=Flag,Description="HapMap2 membership">`

	ii6 = &Info{
		Id:          "H2",
		Number:      "2",
		Type:        "Flag",
		Description: "HapMap2 membership",
		fields: map[string]*Field{
			`ID`:          &Field{`ID`, `H2`, 0, 0},
			`Number`:      &Field{`Number`, `2`, 1, 0},
			`Type`:        &Field{`Type`, `Flag`, 2, 0},
			`Description`: &Field{`Description`, `HapMap2 membership`, 3, '"'}},
		order: []string{`ID`, `Number`, `Type`, `Description`}}

	// Make sure that the original ordering can be recreated
	is7 string = `##INFO=<Type=Flag,ID=HX,Description="XapMap2 membership",Number=2>`

	ii7 = &Info{
		Id:          "HX",
		Number:      "2",
		Type:        "Flag",
		Description: "XapMap2 membership",
		fields: map[string]*Field{
			`Type`:        &Field{`Type`, `Flag`, 0, 0},
			`ID`:          &Field{`ID`, `HX`, 1, 0},
			`Description`: &Field{`Description`, `XapMap2 membership`, 2, '"'},
			`Number`:      &Field{`Number`, `2`, 3, 0}},
		order: []string{`Type`, `ID`, `Description`, `Number`}}

	// Make sure that arbitrary fields are handled
	is8 string = `##INFO=<Type=Flag,Trick='1',ID=Hx,Description="XapMap2 membership",Number=2>`

	ii8 = &Info{
		Id:          "Hx",
		Number:      "2",
		Type:        "Flag",
		Description: "XapMap2 membership",
		fields: map[string]*Field{
			`Type`:        &Field{`Type`, `Flag`, 0, 0},
			`Trick`:       &Field{`Trick`, `1`, 1, '\''},
			`ID`:          &Field{`ID`, `Hx`, 2, 0},
			`Description`: &Field{`Description`, `XapMap2 membership`, 3, '"'},
			`Number`:      &Field{`Number`, `2`, 4, 0}},
		order: []string{`Type`, `Trick`, `ID`, `Description`, `Number`}}
)

var infotests = []struct {
	input string
	exp   *Info
}{
	{is1, ii1},
	{is2, ii2},
	{is3, ii3},
	{is4, ii4},
	{is5, ii5},
	{is6, ii6},
	{is7, ii7},
	{is8, ii8},
}

var (
	f1      *Field = &Field{`ID`, `GT`, 0, 0}
	f2      *Field = &Field{`Number`, `1`, 1, 0}
	f3      *Field = &Field{`Type`, `String`, 2, 0}
	f4      *Field = &Field{`Description`, `Genotype`, 3, '"'}
	fkv1           = map[string]*Field{`ID`: f1, `Number`: f2, `Type`: f3, `Description`: f4}
	forder1        = []string{`ID`, `Number`, `Type`, `Description`}

	f5      *Field = &Field{`ID`, `GQ`, 0, 0}
	f6      *Field = &Field{`Number`, `1`, 1, 0}
	f7      *Field = &Field{`Type`, `Integer`, 2, 0}
	f8      *Field = &Field{`Description`, `Genotype Quality`, 3, '"'}
	fkv2           = map[string]*Field{`ID`: f5, `Number`: f6, `Type`: f7, `Description`: f8}
	forder2        = []string{`ID`, `Number`, `Type`, `Description`}

	f10     *Field = &Field{`ID`, `HQ`, 0, 0}
	f11     *Field = &Field{`Number`, `2`, 1, 0}
	f12     *Field = &Field{`Type`, `Integer`, 2, 0}
	f13     *Field = &Field{`Description`, `Haplotype Quality`, 3, '"'}
	fkv3           = map[string]*Field{`ID`: f10, `Number`: f11, `Type`: f12, `Description`: f13}
	forder3        = []string{`ID`, `Number`, `Type`, `Description`}

	f14     *Field = &Field{`ID`, `DP`, 0, 0}
	f15     *Field = &Field{`Number`, `1`, 1, 0}
	f16     *Field = &Field{`Type`, `Integer`, 2, 0}
	f17     *Field = &Field{`Description`, `Read Depth`, 3, '"'}
	fkv4           = map[string]*Field{`ID`: f14, `Number`: f15, `Type`: f16, `Description`: f17}
	forder4        = []string{`ID`, `Number`, `Type`, `Description`}

	// FORMAT 5 is the same as 4 except for the order
	f18     *Field = &Field{`ID`, `DP`, 2, 0}
	f19     *Field = &Field{`Number`, `1`, 1, 0}
	f20     *Field = &Field{`Type`, `Integer`, 3, 0}
	f21     *Field = &Field{`Description`, `Read Depth`, 0, '"'}
	fkv5           = map[string]*Field{`ID`: f18, `Number`: f19, `Type`: f20, `Description`: f21}
	forder5        = []string{`Description`, `Number`, `ID`, `Type`}
)

// grendeloz: The DeepEqual testing requires us to recreate the expected
// structure of the parsed *Info so the code below looks bollocks. Also
// note that 'zero' runes take the value 0.
var formattests = []struct {
	input string
	exp   *SampleFormat
}{
	{`##FORMAT=<ID=GT,Number=1,Type=String,Description="Genotype">`,
		&SampleFormat{Id: "GT", Number: "1", Type: "String", Description: "Genotype",
			fields: fkv1, order: forder1}},
	{`##FORMAT=<ID=GQ,Number=1,Type=Integer,Description="Genotype Quality">`,
		&SampleFormat{Id: "GQ", Number: "1", Type: "Integer", Description: "Genotype Quality",
			fields: fkv2, order: forder2}},
	{`##FORMAT=<ID=HQ,Number=2,Type=Integer,Description="Haplotype Quality">`,
		&SampleFormat{Id: "HQ", Number: "2", Type: "Integer", Description: "Haplotype Quality",
			fields: fkv3, order: forder3}},
	{`##FORMAT=<ID=DP,Number=1,Type=Integer,Description="Read Depth">`,
		&SampleFormat{Id: "DP", Number: "1", Type: "Integer", Description: "Read Depth",
			fields: fkv4, order: forder4}},
	{`##FORMAT=<Description="Read Depth",Number=1,ID=DP,Type=Integer>`,
		&SampleFormat{Id: "DP", Number: "1", Type: "Integer", Description: "Read Depth",
			fields: fkv5, order: forder5}},
}

var filtertests = []struct {
	filter string
	exp    []string
}{
	{`##FILTER=<ID=q10,Description="Quality below 10">`,
		[]string{"q10", "Quality below 10"}},
	{`##FILTER=<ID=s50,Description="Less than 50% of samples have data">`,
		[]string{"s50", "Less than 50% of samples have data"}},
}

var samplelinetests = []struct {
	line string
	exp  []string
}{
	{`#CHROM	POS	ID	REF	ALT	QUAL	FILTER	INFO	FORMAT	NA00001	NA00002	NA00003`, []string{"NA00001", "NA00002", "NA00003"}},
	{`#CHROM	POS	ID	REF	ALT	QUAL	FILTER	INFO	FORMAT`, []string{}},
}

func (s *VCFSuite) TestKvSplitter(c *C) {
	for _, v := range kvtests {
		obs, _, err := kvSplitter(v.input)
		c.Assert(err, IsNil)
		fields := make(map[string]string)
		for _, f := range obs {
			fields[f.Key] = f.Value
		}
		c.Assert(fields, DeepEquals, v.exp)
	}
}

func (s *VCFSuite) TestHeaderInfoParse(c *C) {
	for _, v := range infotests {
		obs, err := parseHeaderInfo(v.input)
		c.Assert(err, IsNil)
		c.Assert(obs, DeepEquals, v.exp)
		c.Assert(obs.String(), Equals, v.input)
	}
}

func (s *VCFSuite) TestHeaderFormatParse(c *C) {
	for _, v := range formattests {
		obs, err := parseHeaderFormat(v.input)
		c.Assert(err, IsNil)
		c.Assert(obs, DeepEquals, v.exp)
		c.Assert(obs.String(), Equals, v.input)

	}
}

func (s *VCFSuite) TestHeaderFilterParse(c *C) {

	for _, v := range filtertests {
		obs, err := parseHeaderFilter(v.filter)
		c.Assert(err, IsNil)
		c.Assert(obs, DeepEquals, v.exp)
	}
}

func (s *VCFSuite) TestHeaderVersionParse(c *C) {
	obs, err := parseHeaderFileVersion(`##fileformat=VCFv4.2`)
	c.Assert(err, IsNil)
	c.Assert(obs, Equals, "4.2")
}

func (s *VCFSuite) TestHeaderBadVersionParse(c *C) {
	_, err := parseHeaderFileVersion(`##fileformat=VFv4.2`)
	c.Assert(err, ErrorMatches, "file format error.*")
}

func (s *VCFSuite) TestHeaderContigParse(c *C) {
	m, err := parseHeaderContig(`##contig=<ID=20,length=62435964,assembly=B36,md5=f126cdf8a6e0c7f379d618ff66beb2da,species="Homo sapiens",taxonomy=x>`)
	c.Assert(err, IsNil)
	c.Assert(m, DeepEquals, map[string]string{"assembly": "B36", "md5": "f126cdf8a6e0c7f379d618ff66beb2da", "species": "\"Homo sapiens\"", "taxonomy": "x", "ID": "20", "length": "62435964"})
}

func (s *VCFSuite) TestHeaderExtra(c *C) {
	obs, err := parseHeaderExtraKV("##key=value")
	c.Assert(err, IsNil)
	c.Assert(obs[0], Equals, "key")
	c.Assert(obs[1], Equals, "value")
}

func (s *VCFSuite) TestHeaderSampleLine(c *C) {

	for _, v := range samplelinetests {
		r, err := parseSampleLine(v.line)
		c.Assert(err, IsNil)
		c.Assert(r, DeepEquals, v.exp)
	}
}

func (s *VCFSuite) TestIssue5(c *C) {
	rdr, err := os.Open("test-multi-allelic.vcf")
	c.Assert(err, IsNil)
	vcf, err := NewReader(rdr, false)
	c.Assert(err, IsNil)

	variant := vcf.Read()
	samples := variant.Samples

	c.Assert(samples[0].GT, DeepEquals, []int{2, 2})
	c.Assert(samples[1].GT, DeepEquals, []int{2, 2})
	c.Assert(samples[2].GT, DeepEquals, []int{2, 2})
}
