package vcfgo

import (
	"bytes"
	"testing"
)

// This is the example header from section 1.1 of the VCFv4.3 specification
// document at: https://samtools.github.io/hts-specs/VCFv4.2.pdf
var VCFv4_2eg = `##fileformat=VCFv4.2
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
##INFO=<ID=AA,Number=0,Type=Flag,Description="">
##INFO=<ID=LONG,Number=10,Type=Flag,Description="Large number of values">
##INFO=<ID=H2,Number=0,Type=Flag,Description="HapMap2 membership">
##FILTER=<ID=q10,Description="Quality below 10">
##FILTER=<ID=q0,Description="">
##FILTER=<ID=s50,Description="Less than 50% of samples have data">
##FORMAT=<ID=GT,Number=1,Type=String,Description="Genotype">
##FORMAT=<ID=GQ,Number=1,Type=Integer,Description="Genotype Quality">
##FORMAT=<ID=DP,Number=1,Type=Integer,Description="Read Depth">
##FORMAT=<ID=LONG,Number=10,Type=Integer,Description="Long number of values">
##FORMAT=<ID=HQ,Number=2,Type=Integer,Description="Haplotype Quality">
#CHROM	POS	ID	REF	ALT	QUAL	FILTER	INFO	FORMAT	NA00001	NA00002	NA00003
20	14370	rs6054257	G	A	29	PASS	NS=3;DP=14;AF=0.5;DB;H2	GT:GQ:DP:HQ	0|0:48:1:51,51	1|0:48:8:51,51	1/1:43:5:.,.
20	17330	.	T	A	3	q10	NS=3;DP=11;AF=0.017	GT:GQ:DP:HQ	0|0:49:3:58,50	0|1:3:5:65,3	0/0:41:3:.,.
20	1110696	rs6040355	A	G,T	67	PASS	NS=2;DP=10;AF=0.333,0.667;AA=T;DB	GT:GQ:DP:HQ	1|2:21:6:23,27	2|1:2:0:18,2	2/2:35:4:.,.
20	1230237	.	T	.	47	PASS	NS=3;DP=13;AA=T	GT:GQ:DP:HQ	0|0:54:7:56,60	0|0:48:4:51,51	0/0:61:2:.,.
20	1234567	microsat1	GTC	G,GTCT	50	PASS	NS=3;DP=9;AA=G	GT:GQ:DP	0/1:35:4	0/2:17:2	1/1:40:3
X	153171993	rs5201	A	.	.	.	.	GT	0	1	.
TRIPLOID	153171993	rs5201	A	.	.	.	.	GT	0|0|0	1/0/1	.`

var bedStr = `1	0	10000	0.061011
1	10000	10154	0.070013
1	10154	10200	0.126639
1	10400	10535	0.053691
1	10535	10625	0.078448
1	10625	11084	0.053691
1	11084	11159	0.078448
1	11159	11325	0.053691
1	11325	11400	0.078448
1	11400	11404	0.053691`

// no fileformat
var badVcfStr = `##source=myImputationProgramV3.1
##reference=file:///seq/references/1000GenomesPilot-NCBI36.fasta
##contig=<ID=20,length=62435964,assembly=B36,md5=f126cdf8a6e0c7f379d618ff66beb2da,species="Homo sapiens",taxonomy=x>
##phasing=partial
##INFO=<ID=NS,Number=1,Type=Integer,Description="Number of Samples With Data">
#CHROM	POS	ID	REF	ALT	QUAL	FILTER	INFO	FORMAT	NA00001	NA00002	NA00003
20	14370	rs6054257	G	A	29	PASS	NS=3;DP=14;AF=0.5;DB;H2	GT:GQ:DP:HQ	0|0:48:1:51,51	1|0:48:8:51,51	1/1:43:5:.,.`

func TestVCFv4_2(t *testing.T) {

	// Read the VCFv4.2 header string
	r := bytes.NewReader([]byte(VCFv4_2eg))
	rdr, err := NewReader(r, false)
	if err != nil {
		t.Errorf("Reading the VCFv4.2 example header threw an error: %v\n", err)
	}

	// Test the header

	var tests = []struct {
		label string
		obs   interface{}
		exp   interface{}
	}{
		{`FileFormat`, rdr.Header.FileFormat, `4.2`},
		{`SampleNames count`, len(rdr.Header.SampleNames), 3},
	}

	for _, v := range tests {
		if v.obs != v.exp {
			t.Errorf("%v is %v but expected %v\n", v.label, v.obs, v.exp)
		}
	}
}
