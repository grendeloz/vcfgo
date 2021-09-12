package vcfgo

import (
	"io"
	"os"
	"strings"

	. "gopkg.in/check.v1"
)

type ReaderSuite struct {
	reader io.Reader
}

var _ = Suite(&ReaderSuite{})

func (s *ReaderSuite) SetUpTest(c *C) {
	s.reader = strings.NewReader(VCFv4_2eg)

}

func (s *ReaderSuite) TestReaderHeaderSamples(c *C) {
	v, err := NewReader(s.reader, false)
	c.Assert(err, IsNil)
	c.Assert(v.Header.SampleNames, DeepEquals, []string{"NA00001", "NA00002", "NA00003"})

}

func (s *ReaderSuite) TestLazyReader(c *C) {
	rdr, err := NewReader(s.reader, true)
	c.Assert(err, IsNil)
	rec := rdr.Read() //.(*Variant)
	c.Assert(rec.String(), Equals, "20\t14370\trs6054257\tG\tA\t29.0\tPASS\tNS=3;DP=14;AF=0.5;DB;H2\tGT:GQ:DP:HQ\t0|0:48:1:51,51\t1|0:48:8:51,51\t1/1:43:5:.,.")
}

func (s *ReaderSuite) TestReaderHeaderInfos(c *C) {

	parsedNS := &Info{Id: "NS",
		Number:      "1",
		Type:        "Integer",
		Description: "Number of Samples With Data",
		fields: map[string]*KV{
			`ID`:          &KV{"ID", "NS", 0, 0},
			`Number`:      &KV{"Number", "1", 1, 0},
			`Type`:        &KV{"Type", "Integer", 2, 0},
			`Description`: &KV{"Description", "Number of Samples With Data", 3, '"'}},
		order: []string{`ID`, `Number`, `Type`, `Description`}}

	parsedGT := &SampleFormat{Id: "GT",
		Number:      "1",
		Type:        "String",
		Description: "Genotype",
		fields: map[string]*KV{
			`ID`:          &KV{"ID", "GT", 0, 0},
			`Number`:      &KV{"Number", "1", 1, 0},
			`Type`:        &KV{"Type", "String", 2, 0},
			`Description`: &KV{"Description", "Genotype", 3, '"'}},
		order: []string{`ID`, `Number`, `Type`, `Description`}}

	v, err := NewReader(s.reader, false)
	c.Assert(err, IsNil)
	c.Assert(v.Header.Infos["NS"], DeepEquals, parsedNS)
	c.Assert(v.Header.Filters["q10"], Equals, "Quality below 10")
	c.Assert(v.Header.SampleFormats["GT"], DeepEquals, parsedGT)
}

func (s *ReaderSuite) TestReaderHeaderExtras(c *C) {
	v, err := NewReader(s.reader, true)
	c.Assert(err, IsNil)
	c.Assert(len(v.Header.Lines), Equals, 4)
	c.Assert(v.Header.Lines[0].OgString, Equals, `##fileDate=20090805`)
	c.Assert(v.Header.Lines[0].LineKey, Equals, `fileDate`)
	c.Assert(v.Header.Lines[0].Value, Equals, `20090805`)
	c.Assert(v.Header.Lines[3].LineKey, Equals, `phasing`)
	c.Assert(v.Header.Lines[3].Value, Equals, `partial`)
}

func (s *ReaderSuite) TestReaderRead(c *C) {
	rdr, err := NewReader(s.reader, false)
	c.Assert(err, IsNil)

	rec := rdr.Read() //.(*Variant)
	c.Assert(rec.Chromosome, Equals, "20")
	c.Assert(rec.Pos, Equals, uint64(14370))
	c.Assert(rec.Id(), Equals, "rs6054257")
	c.Assert(rec.Ref(), Equals, "G")
	c.Assert(rec.Alt()[0], Equals, "A")
	c.Assert(rec.Quality, Equals, float32(29.0))
	c.Assert(rec.Filter, Equals, "PASS")

	//20	17330	.	T	A	3	q10	NS=3;DP=11;AF=0.017	GT:GQ:DP:HQ	0|0:49:3:58,50	0|1:3:5:65,3	0/0:41:3
	rec0 := rdr.Read() //.(*Variant)
	c.Assert(rec0.Chromosome, Equals, "20")
	c.Assert(rec0.Pos, Equals, uint64(17330))
	c.Assert(rec0.Id(), Equals, ".")
	c.Assert(rec0.Ref(), Equals, "T")
	c.Assert(rec0.Alt()[0], Equals, "A")
	c.Assert(rec0.Quality, Equals, float32(3))
	c.Assert(rec0.Filter, Equals, "q10")

	//20	1110696	rs6040355	A	G,T	67	PASS	NS=2;DP=10;AF=0.333,0.667;AA=T;DB	GT:GQ:DP:HQ	1|2:21:6:23,27	2|1:2:0:18,2	2/2:35:4
	rec = rdr.Read() //.(*Variant)
	c.Assert(rec.Chromosome, Equals, "20")
	c.Assert(int(rec.Pos), Equals, 1110696)
	c.Assert(rec.Id(), Equals, "rs6040355")
	c.Assert(rec.Ref(), Equals, "A")
	c.Assert(rec.Alt(), DeepEquals, []string{"G", "T"})
	c.Assert(rec.Quality, Equals, float32(67))
	c.Assert(rec.Filter, Equals, "PASS")

	c.Assert(rec0.Chromosome, Equals, "20")
	c.Assert(rec0.Pos, Equals, uint64(17330))
	c.Assert(rec0.Id(), Equals, ".")
	c.Assert(rec0.Ref(), Equals, "T")
	c.Assert(rec0.Alt()[0], Equals, "A")
	c.Assert(rec0.Quality, Equals, float32(3))
	c.Assert(rec0.Filter, Equals, "q10")
}

//func (s *ReaderSuite) TestReaderRVGBug(c *C) {
//	v, err := os.Open("test-h.vcf")
//	if err != nil {
//		c.Fatalf("error opening test-h.vcf")
//	}
//	rdr, err := NewReader(v, false)
//	c.Assert(err, IsNil)
//	rec := rdr.Read()
//	_ = rec
//}

func (s *ReaderSuite) TestSampleParsingErrors(c *C) {
	sr := strings.NewReader(`##fileformat=VCFv4.0
##FORMAT=<ID=GT,Number=1,Type=String,Description="Genotype">
#CHROM	POS	ID	REF	ALT	QUAL	FILTER	INFO	FORMAT	S-1	S-2	S-3
1	100000	.	C	G	.	.	.	GT	0|M	1/.	1/0
2	200000	.	C	G	.	.	.	GT	0|0	0|1	1|E`)

	rdr, err := NewReader(sr, false)

	c.Assert(err, IsNil)

	c.Assert(rdr.Read(), NotNil)
	c.Assert(rdr.Error(), ErrorMatches, ".*M.* invalid syntax.*")

	rdr.Clear()

	c.Assert(rdr.Read(), NotNil)
	c.Assert(rdr.Error(), ErrorMatches, ".*E.* invalid syntax.*")
}

func (s *ReaderSuite) TestSampleParsingErrors2(c *C) {
	f, err := os.Open("test-dp.vcf")
	c.Assert(err, IsNil)

	rdr, err := NewReader(f, false)
	c.Assert(err, IsNil)

	variant := rdr.Read()
	c.Assert(variant, NotNil)

	c.Assert(rdr.Error(), IsNil)

}
