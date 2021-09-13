package vcfgo

import (
	"io"
	"os"
	"strings"

	"github.com/brentp/irelate/interfaces"

	. "gopkg.in/check.v1"
)

type HeaderSuite struct {
	reader io.Reader
	vcfStr string
}

type BadVcfSuite HeaderSuite

var a = Suite(&HeaderSuite{vcfStr: VCFv4_2eg})
var b = Suite(&BadVcfSuite{vcfStr: bedStr})

func (s *HeaderSuite) SetUpTest(c *C) {
	s.reader = strings.NewReader(s.vcfStr)
}

func (s *BadVcfSuite) SetUpTest(c *C) {
	s.reader = strings.NewReader(s.vcfStr)
}

func (s *HeaderSuite) TestReaderHeaderParseSample(c *C) {
	r, err := NewReader(s.reader, false)
	c.Assert(err, IsNil)
	v := r.Read()
	c.Assert(r.Error(), IsNil)

	fmt := v.Format
	c.Assert(fmt, DeepEquals, []string{"GT", "GQ", "DP", "HQ"})
}

func (b *BadVcfSuite) TestReaderHeaderParseSample(c *C) {
	r, err := NewReader(b.reader, false)
	c.Assert(r, IsNil)
	c.Assert(err, NotNil)
}

func (s *HeaderSuite) TestSamples(c *C) {
	r, err := NewReader(s.reader, false)
	c.Assert(err, IsNil)

	v := r.Read()
	samp := v.Samples[0]

	c.Assert(samp.DP, Equals, 1)
	c.Assert(samp.GQ, Equals, 48)

	f, err := v.GetGenotypeField(samp, "HQ", -1)
	c.Assert(err, IsNil)
	c.Assert(f, DeepEquals, []int{51, 51})

	samp2 := v.Samples[2]
	f, err = v.GetGenotypeField(samp2, "HQ", -1)
	c.Assert(err, IsNil)
	c.Assert(f, DeepEquals, []int{-1, -1})

	variants := make([]*Variant, 0)
	chromCount := make(map[string]int)
	var vv interfaces.IVariant
	for vv = r.Read(); vv != nil; vv = r.Read() {
		v := vv.(*Variant)
		if v == nil {
			break
		}

		variants = append(variants, v)
		chromCount[v.Chromosome]++
	}

	c.Assert(chromCount["20"], Equals, 4)
	c.Assert(chromCount["X"], Equals, 1)

	c.Assert(int(variants[len(variants)-1].Pos), Equals, int(153171993))
	c.Assert(variants[3].Filter, Equals, "PASS")
}

func (s *HeaderSuite) TestWeirdHeader(c *C) {
	rr, err := os.Open("test-weird-header.vcf")
	c.Assert(err, IsNil)
	_, err = NewReader(rr, false)
	c.Assert(err, IsNil)
}

func (s *HeaderSuite) TestSampleGenotypes(c *C) {
	r, err := NewReader(s.reader, false)
	c.Assert(err, IsNil)

	variants := make([]*Variant, 0)
	for {
		v := r.Read()
		if v == nil {
			break
		}

		variants = append(variants, v)
	}

	// validate diploid parsing works
	firstVariant := variants[0]
	c.Assert(firstVariant.Samples[0].GT, DeepEquals, []int{0, 0})
	c.Assert(firstVariant.Samples[0].Phased, Equals, true)

	c.Assert(firstVariant.Samples[1].GT, DeepEquals, []int{1, 0})
	c.Assert(firstVariant.Samples[1].Phased, Equals, true)

	c.Assert(firstVariant.Samples[2].GT, DeepEquals, []int{1, 1})
	c.Assert(firstVariant.Samples[2].Phased, Equals, false)

	// validate haploid parsing works
	hapVariant := variants[5]
	c.Assert(hapVariant.Samples[0].GT, DeepEquals, []int{0})
	c.Assert(hapVariant.Samples[0].Phased, Equals, false)

	c.Assert(hapVariant.Samples[1].GT, DeepEquals, []int{1})
	c.Assert(hapVariant.Samples[1].Phased, Equals, false)

	c.Assert(hapVariant.Samples[2].GT, DeepEquals, []int{-1})
	c.Assert(hapVariant.Samples[2].Phased, Equals, false)

	// validate triploid parsing works
	tripVariant := variants[6]
	c.Assert(tripVariant.Samples[0].GT, DeepEquals, []int{0, 0, 0})
	c.Assert(tripVariant.Samples[0].Phased, Equals, true)

	c.Assert(tripVariant.Samples[1].GT, DeepEquals, []int{1, 0, 1})
	c.Assert(tripVariant.Samples[1].Phased, Equals, false)

	c.Assert(tripVariant.Samples[2].GT, DeepEquals, []int{-1})
	c.Assert(tripVariant.Samples[2].Phased, Equals, false)
}

/*
func (s *VariantSuite) TestParseOne(c *C) {

	v, err := parseOne("key", "123", "Integer")
	c.Assert(err, IsNil)
	c.Assert(v, Equals, 123)

	v1, err := parseOne("key", "a123", "String")
	c.Assert(err, IsNil)
	c.Assert(v1, Equals, "a123")

	v2, err := parseOne("key", "asdf", "Flag")
	c.Assert(err, ErrorMatches, ".*flag field .* had value")
	c.Assert(v2, Equals, true)

	v3, err := parseOne("key", "", "Flag")
	c.Assert(err, IsNil)
	c.Assert(v3, Equals, true)

}
*/
