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

var (
	i1      *Field = &Field{`ID`, `NS`, 0, 0}
	i2      *Field = &Field{`Number`, `1`, 1, 0}
	i3      *Field = &Field{`Type`, `Integer`, 2, 0}
	i4      *Field = &Field{`Description`, `Number of Samples With Data`, 3, '"'}
	ikv1           = map[string]*Field{`ID`: i1, `Number`: i2, `Type`: i3, `Description`: i4}
	iorder1        = []string{`ID`, `Number`, `Type`, `Description`}

	i5      *Field = &Field{`ID`, `DP`, 0, 0}
	i6      *Field = &Field{`Number`, `1`, 1, 0}
	i7      *Field = &Field{`Type`, `Integer`, 2, 0}
	i8      *Field = &Field{`Description`, `Total Depth`, 3, '"'}
	ikv2           = map[string]*Field{`ID`: i5, `Number`: i6, `Type`: i7, `Description`: i8}
	iorder2        = []string{`ID`, `Number`, `Type`, `Description`}

	i10     *Field = &Field{`ID`, `AF`, 0, 0}
	i11     *Field = &Field{`Number`, `A`, 1, 0}
	i12     *Field = &Field{`Type`, `Float`, 2, 0}
	i13     *Field = &Field{`Description`, `Allele Frequency`, 3, '"'}
	ikv3           = map[string]*Field{`ID`: i10, `Number`: i11, `Type`: i12, `Description`: i13}
	iorder3        = []string{`ID`, `Number`, `Type`, `Description`}

	i14     *Field = &Field{`ID`, `AA`, 0, 0}
	i15     *Field = &Field{`Number`, `1`, 1, 0}
	i16     *Field = &Field{`Type`, `String`, 2, 0}
	i17     *Field = &Field{`Description`, `Ancestral Allele`, 3, '"'}
	ikv4           = map[string]*Field{`ID`: i14, `Number`: i15, `Type`: i16, `Description`: i17}
	iorder4        = []string{`ID`, `Number`, `Type`, `Description`}

	i18     *Field = &Field{`ID`, `DB`, 0, 0}
	i19     *Field = &Field{`Number`, `0`, 1, 0}
	i20     *Field = &Field{`Type`, `Flag`, 2, 0}
	i21     *Field = &Field{`Description`, `dbSNP membership, build 129`, 3, '"'}
	ikv5           = map[string]*Field{`ID`: i18, `Number`: i19, `Type`: i20, `Description`: i21}
	iorder5        = []string{`ID`, `Number`, `Type`, `Description`}

	i22     *Field = &Field{`ID`, `H2`, 0, 0}
	i23     *Field = &Field{`Number`, `2`, 1, 0}
	i24     *Field = &Field{`Type`, `Flag`, 2, 0}
	i25     *Field = &Field{`Description`, `HapMap2 membership`, 3, '"'}
	ikv6           = map[string]*Field{`ID`: i22, `Number`: i23, `Type`: i24, `Description`: i25}
	iorder6        = []string{`ID`, `Number`, `Type`, `Description`}

	// INFO 7 is the same as 6 except for the order
	i26     *Field = &Field{`ID`, `HX`, 0, 0}
	i27     *Field = &Field{`Number`, `2`, 1, 0}
	i28     *Field = &Field{`Type`, `Flag`, 2, 0}
	i29     *Field = &Field{`Description`, `HapMap2 membership`, 3, '"'}
	ikv7           = map[string]*Field{`Type`: i28, `ID`: i26, `Description`: i29, `Number`: i27}
	//ikv7           = map[string]*Field{`ID`: i26, `Number`: i27, `Type`: i28, `Description`: i29}
	//ikv7           = map[string]*Field{`Type`: i24, `ID`: i22, `Description`: i25, `Number`: i23}
	iorder7 = []string{`Type`, `ID`, `Description`, `Number`}
)

var infotests = []struct {
	input string
	exp   *Info
}{
	{`##INFO=<ID=NS,Number=1,Type=Integer,Description="Number of Samples With Data">`,
		&Info{Id: "NS", Number: "1", Type: "Integer", Description: "Number of Samples With Data",
			kvs: ikv1, order: iorder1}},
	{`##INFO=<ID=DP,Number=1,Type=Integer,Description="Total Depth">`,
		&Info{Id: "DP", Number: "1", Type: "Integer", Description: "Total Depth",
			kvs: ikv2, order: iorder2}},
	{`##INFO=<ID=AF,Number=A,Type=Float,Description="Allele Frequency">`,
		&Info{Id: "AF", Number: "A", Type: "Float", Description: "Allele Frequency",
			kvs: ikv3, order: iorder3}},
	{`##INFO=<ID=AA,Number=1,Type=String,Description="Ancestral Allele">`,
		&Info{Id: "AA", Number: "1", Type: "String", Description: "Ancestral Allele",
			kvs: ikv4, order: iorder4}},
	{`##INFO=<ID=DB,Number=0,Type=Flag,Description="dbSNP membership, build 129">`,
		&Info{Id: "DB", Number: "0", Type: "Flag", Description: "dbSNP membership, build 129",
			kvs: ikv5, order: iorder5}},
	{`##INFO=<ID=H2,Number=2,Type=Flag,Description="HapMap2 membership">`,
		&Info{Id: "H2", Number: "2", Type: "Flag", Description: "HapMap2 membership",
			kvs: ikv6, order: iorder6}},
	{`##INFO=<Type=Flag,ID=H2,Description="HapMap2 membership",Number=2>`,
		&Info{Id: "HX", Number: "2", Type: "Flag", Description: "HapMap2 membership",
			kvs: ikv7, order: iorder7}},
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
	forder5 = []string{`Description`, `Number`, `ID`, `Type`}
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
			kvs: fkv1, order: forder1}},
	{`##FORMAT=<ID=GQ,Number=1,Type=Integer,Description="Genotype Quality">`,
		&SampleFormat{Id: "GQ", Number: "1", Type: "Integer", Description: "Genotype Quality",
			kvs: fkv2, order: forder2}},
	{`##FORMAT=<ID=HQ,Number=2,Type=Integer,Description="Haplotype Quality">`,
		&SampleFormat{Id: "HQ", Number: "2", Type: "Integer", Description: "Haplotype Quality",
			kvs: fkv3, order: forder3}},
	{`##FORMAT=<ID=DP,Number=1,Type=Integer,Description="Read Depth">`,
		&SampleFormat{Id: "DP", Number: "1", Type: "Integer", Description: "Read Depth",
			kvs: fkv4, order: forder4}},
	//{`##FORMAT=<Description="Read Depth",Number=1,ID=DP,Type=Integer>`,
	//	&SampleFormat{Id: "DP", Number: "1", Type: "Integer", Description: "Read Depth",
	//		kvs: fkv4, order: forder5}},
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

func (s *VCFSuite) TestAKvSplitter(c *C) {
	for _, v := range kvtests {
		obs, err := kvSplitter(v.input)
		c.Assert(err, IsNil)
		kvs := make(map[string]string)
		for _, f := range obs {
			kvs[f.Key] = f.Value
		}
		c.Assert(kvs, DeepEquals, v.exp)
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
