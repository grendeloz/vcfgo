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

	ii1 = &MetaLine{
		LineKey:  "INFO",
		MetaType: Structured,
		OgString: is1,
		KVs: map[string]*KV{
			`ID`:          &KV{`ID`, `NS`, 0, 0},
			`Number`:      &KV{`Number`, `1`, 1, 0},
			`Type`:        &KV{`Type`, `Integer`, 2, 0},
			`Description`: &KV{`Description`, `Number of Samples With Data`, 3, '"'}},
		Order: []string{`ID`, `Number`, `Type`, `Description`}}

	is2 string = `##INFO=<ID=DP,Number=1,Type=Integer,Description="Total Depth">`

	ii2 = &MetaLine{
		LineKey:  "INFO",
		MetaType: Structured,
		OgString: is2,
		KVs: map[string]*KV{
			`ID`:          &KV{`ID`, `DP`, 0, 0},
			`Number`:      &KV{`Number`, `1`, 1, 0},
			`Type`:        &KV{`Type`, `Integer`, 2, 0},
			`Description`: &KV{`Description`, `Total Depth`, 3, '"'}},
		Order: []string{`ID`, `Number`, `Type`, `Description`}}

	is3 string = `##INFO=<ID=AF,Number=A,Type=Float,Description="Allele Frequency">`

	ii3 = &MetaLine{
		LineKey:  "INFO",
		MetaType: Structured,
		OgString: is3,
		KVs: map[string]*KV{
			`ID`:          &KV{`ID`, `AF`, 0, 0},
			`Number`:      &KV{`Number`, `A`, 1, 0},
			`Type`:        &KV{`Type`, `Float`, 2, 0},
			`Description`: &KV{`Description`, `Allele Frequency`, 3, '"'}},
		Order: []string{`ID`, `Number`, `Type`, `Description`}}

	is4 string = `##INFO=<ID=AA,Number=1,Type=String,Description="Ancestral Allele">`

	ii4 = &MetaLine{
		LineKey:  "INFO",
		MetaType: Structured,
		OgString: is4,
		KVs: map[string]*KV{
			`ID`:          &KV{`ID`, `AA`, 0, 0},
			`Number`:      &KV{`Number`, `1`, 1, 0},
			`Type`:        &KV{`Type`, `String`, 2, 0},
			`Description`: &KV{`Description`, `Ancestral Allele`, 3, '"'}},
		Order: []string{`ID`, `Number`, `Type`, `Description`}}

	is5 string = `##INFO=<ID=DB,Number=0,Type=Flag,Description="dbSNP membership, build 129">`

	ii5 = &MetaLine{
		LineKey:  "INFO",
		MetaType: Structured,
		OgString: is5,
		KVs: map[string]*KV{
			`ID`:          &KV{`ID`, `DB`, 0, 0},
			`Number`:      &KV{`Number`, `0`, 1, 0},
			`Type`:        &KV{`Type`, `Flag`, 2, 0},
			`Description`: &KV{`Description`, `dbSNP membership, build 129`, 3, '"'}},
		Order: []string{`ID`, `Number`, `Type`, `Description`}}

	is6 string = `##INFO=<ID=H2,Number=2,Type=Flag,Description="HapMap2 membership">`

	ii6 = &MetaLine{
		LineKey:  "INFO",
		MetaType: Structured,
		OgString: is6,
		KVs: map[string]*KV{
			`ID`:          &KV{`ID`, `H2`, 0, 0},
			`Number`:      &KV{`Number`, `2`, 1, 0},
			`Type`:        &KV{`Type`, `Flag`, 2, 0},
			`Description`: &KV{`Description`, `HapMap2 membership`, 3, '"'}},
		Order: []string{`ID`, `Number`, `Type`, `Description`}}

	// Make sure that the original ordering can be recreated
	is7 string = `##INFO=<Type=Flag,ID=HX,Description="XapMap2 membership",Number=2>`

	ii7 = &MetaLine{
		LineKey:  "INFO",
		MetaType: Structured,
		OgString: is7,
		KVs: map[string]*KV{
			`Type`:        &KV{`Type`, `Flag`, 0, 0},
			`ID`:          &KV{`ID`, `HX`, 1, 0},
			`Description`: &KV{`Description`, `XapMap2 membership`, 2, '"'},
			`Number`:      &KV{`Number`, `2`, 3, 0}},
		Order: []string{`Type`, `ID`, `Description`, `Number`}}

	// Make sure that arbitrary fields are handled
	is8 string = `##INFO=<Type=Flag,Trick='1',ID=Hx,Description="XapMap2 membership",Number=2>`

	ii8 = &MetaLine{
		LineKey:  "INFO",
		MetaType: Structured,
		OgString: is8,
		KVs: map[string]*KV{
			`Type`:        &KV{`Type`, `Flag`, 0, 0},
			`Trick`:       &KV{`Trick`, `1`, 1, '\''},
			`ID`:          &KV{`ID`, `Hx`, 2, 0},
			`Description`: &KV{`Description`, `XapMap2 membership`, 3, '"'},
			`Number`:      &KV{`Number`, `2`, 4, 0}},
		Order: []string{`Type`, `Trick`, `ID`, `Description`, `Number`}}
)

var metaLineInfoTests = []struct {
	input string
	exp   *MetaLine
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
	f1      *KV = &KV{`ID`, `GT`, 0, 0}
	f2      *KV = &KV{`Number`, `1`, 1, 0}
	f3      *KV = &KV{`Type`, `String`, 2, 0}
	f4      *KV = &KV{`Description`, `Genotype`, 3, '"'}
	fkv1        = map[string]*KV{`ID`: f1, `Number`: f2, `Type`: f3, `Description`: f4}
	forder1     = []string{`ID`, `Number`, `Type`, `Description`}

	f5      *KV = &KV{`ID`, `GQ`, 0, 0}
	f6      *KV = &KV{`Number`, `1`, 1, 0}
	f7      *KV = &KV{`Type`, `Integer`, 2, 0}
	f8      *KV = &KV{`Description`, `Genotype Quality`, 3, '"'}
	fkv2        = map[string]*KV{`ID`: f5, `Number`: f6, `Type`: f7, `Description`: f8}
	forder2     = []string{`ID`, `Number`, `Type`, `Description`}

	f10     *KV = &KV{`ID`, `HQ`, 0, 0}
	f11     *KV = &KV{`Number`, `2`, 1, 0}
	f12     *KV = &KV{`Type`, `Integer`, 2, 0}
	f13     *KV = &KV{`Description`, `Haplotype Quality`, 3, '"'}
	fkv3        = map[string]*KV{`ID`: f10, `Number`: f11, `Type`: f12, `Description`: f13}
	forder3     = []string{`ID`, `Number`, `Type`, `Description`}

	f14     *KV = &KV{`ID`, `DP`, 0, 0}
	f15     *KV = &KV{`Number`, `1`, 1, 0}
	f16     *KV = &KV{`Type`, `Integer`, 2, 0}
	f17     *KV = &KV{`Description`, `Read Depth`, 3, '"'}
	fkv4        = map[string]*KV{`ID`: f14, `Number`: f15, `Type`: f16, `Description`: f17}
	forder4     = []string{`ID`, `Number`, `Type`, `Description`}

	// FORMAT 5 is the same as 4 except for the order
	f18     *KV = &KV{`ID`, `DP`, 2, 0}
	f19     *KV = &KV{`Number`, `1`, 1, 0}
	f20     *KV = &KV{`Type`, `Integer`, 3, 0}
	f21     *KV = &KV{`Description`, `Read Depth`, 0, '"'}
	fkv5        = map[string]*KV{`ID`: f18, `Number`: f19, `Type`: f20, `Description`: f21}
	forder5     = []string{`Description`, `Number`, `ID`, `Type`}
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
	for _, v := range metaLineInfoTests {
		obs, err := NewMetaLineFromString(v.input)
		c.Assert(err, IsNil)
		c.Assert(obs, DeepEquals, v.exp)
		s, err := obs.String()
		c.Assert(err, IsNil)
		c.Assert(s, Equals, v.input)
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
