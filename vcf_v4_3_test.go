package vcfgo

import (
	"bytes"
	"testing"
)

// This is the example header from section 1.1 of the VCFv4.3 specification
// document at: https://samtools.github.io/hts-specs/VCFv4.3.pdf
var VCFv4_3eg = `##fileformat=VCFv4.3
##fileDate=20090805
##source=myImputationProgramV3.1
##reference=file:///seq/references/1000GenomesPilot-NCBI36.fasta
##contig=<ID=20,length=62435964,assembly=B36,md5=f126cdf8a6e0c7f379d618ff66beb2da,species="Homo sapiens",taxonomy=x>
##phasing=partial
##INFO=<ID=NS,Number=1,Type=Integer,Description="Number of Samples With Data">
##INFO=<ID=DP,Number=1,Type=Integer,Description="Total Depth">
##INFO=<ID=AF,Number=A,Type=Float,Description="Allele Frequency">
##INFO=<ID=AA,Number=1,Type=String,Description="Ancestral Allele">
##INFO=<ID=DB,Number=0,Type=Flag,Description="dbSNP membership, build 129">
##INFO=<ID=H2,Number=0,Type=Flag,Description="HapMap2 membership">
##FILTER=<ID=q10,Description="Quality below 10">
##FILTER=<ID=s50,Description="Less than 50% of samples have data">
##FORMAT=<ID=GT,Number=1,Type=String,Description="Genotype">
##FORMAT=<ID=GQ,Number=1,Type=Integer,Description="Genotype Quality">
##FORMAT=<ID=DP,Number=1,Type=Integer,Description="Read Depth">
##FORMAT=<ID=HQ,Number=2,Type=Integer,Description="Haplotype Quality">
#CHROM	POS	ID	REF	ALT	QUAL	FILTER	INFO	FORMAT	NA00001	NA00002	NA00003
20	14370	rs6054257	G	A	29	PASS	NS=3;DP=14;AF=0.5;DB;H2	GT:GQ:DP:HQ	0|0:48:1:51,51	1|0:48:8:51,51	1/1:43:5:.,.
20	17330	.	T	A	3	q10	NS=3;DP=11;AF=0.017	GT:GQ:DP:HQ	0|0:49:3:58,50	0|1:3:5:65,3	0/0:41:3
20	1110696	rs6040355	A	G,T	67	PASS	NS=2;DP=10;AF=0.333,0.667;AA=T;DB	GT:GQ:DP:HQ	1|2:21:6:23,27	2|1:2:0:18,2	2/2:35:4
20	1230237	.	T	.	47	PASS	NS=3;DP=13;AA=T	GT:GQ:DP:HQ 0|0:54:7:56,60	0|0:48:4:51,51	0/0:61:2
20	1234567	microsat1	GTC	G,GTCT	50	PASS	NS=3;DP=9;AA=G	GT:GQ:DP	0/1:35:4	0/2:17:2	1/1:40:3`

func TestVCFv4_3_Header(t *testing.T) {

	// Read the VCFv4.3 header string
	r := bytes.NewReader([]byte(VCFv4_3eg))
	rdr, err := NewReader(r, false)
	if err != nil {
		t.Errorf("Reading the VCFv4.3 example header threw an error: %v\n", err)
	}

	// Test the header

	var tests = []struct {
		label string
		obs   interface{}
		exp   interface{}
	}{
		{`FileFormat`, rdr.Header.FileFormat, `4.3`},
		{`SampleNames count`, len(rdr.Header.SampleNames), 3},
		{`MetaLines count`, len(rdr.Header.Lines), 17},
	}

	for _, v := range tests {
		if v.obs != v.exp {
			t.Errorf("%v is %v but expected %v\n", v.label, v.obs, v.exp)
		}
	}
}

func TestVCFv4_3_Variants(t *testing.T) {

	// Read the VCFv4.3 header string
	r := bytes.NewReader([]byte(VCFv4_3eg))
	rdr, err := NewReader(r, false)
	if err != nil {
		t.Errorf("Reading the VCFv4.3 example header threw an error: %v\n", err)
	}

    _ = rdr

	// Read and test the variants
	//	for {
	//		variant := rdr.Read()
	//		if variant == nil {
	//			break
	//		}
	//	}
}

func TestVCFv4_3_INFO(t *testing.T) {

	// Read the VCFv4.3 header string
	r := bytes.NewReader([]byte(VCFv4_3eg))
	rdr, err := NewReader(r, false)
	if err != nil {
		t.Errorf("Reading the VCFv4.3 example header threw an error: %v\n", err)
	}

	// Test INFO lines
    // ##INFO=<ID=NS,Number=1,Type=Integer,Description="Number of Samples With Data">
    // ##INFO=<ID=DP,Number=1,Type=Integer,Description="Total Depth">
    // ##INFO=<ID=AF,Number=A,Type=Float,Description="Allele Frequency">
    // ##INFO=<ID=AA,Number=1,Type=String,Description="Ancestral Allele">
    // ##INFO=<ID=DB,Number=0,Type=Flag,Description="dbSNP membership, build 129">
    // ##INFO=<ID=H2,Number=0,Type=Flag,Description="HapMap2 membership">

    lines := rdr.Header.GetLinesByType(`INFO`)

	var tests = []struct {
		label string
		obs   interface{}
		exp   interface{}
	}{
		{`INFO lines count`, len(lines), 6},
		{`Line 0 ID`, lines[0].GetValue(`ID`), `NS`},
		{`Line 0 Number`, lines[0].GetValue(`Number`), `1`},
		{`Line 0 Type`, lines[0].GetValue(`Type`), `Integer`},
		{`Line 0 Description`, lines[0].GetValue(`Description`), `Number of Samples With Data`},
		{`Line 1 ID`, lines[1].GetValue(`ID`), `DP`},
		{`Line 2 ID`, lines[2].GetValue(`ID`), `AF`},
		{`Line 2 Number`, lines[2].GetValue(`Number`), `A`},
		{`Line 2 Type`, lines[2].GetValue(`Type`), `Float`},
		{`Line 3 ID`, lines[3].GetValue(`ID`), `AA`},
		{`Line 3 Type`, lines[3].GetValue(`Type`), `String`},
		{`Line 4 ID`, lines[4].GetValue(`ID`), `DB`},
		{`Line 5 ID`, lines[5].GetValue(`ID`), `H2`},
		{`Line 5 Number`, lines[5].GetValue(`Number`), `0`},
		{`Line 5 Type`, lines[5].GetValue(`Type`), `Flag`},
	}

	for _, v := range tests {
		if v.obs != v.exp {
			t.Errorf("%v is %v but expected %v\n", v.label, v.obs, v.exp)
		}
	}
}
