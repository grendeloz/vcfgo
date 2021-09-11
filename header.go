package vcfgo

import (
	"encoding/csv"
	"errors"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

//var infoRegexp = regexp.MustCompile(fmt.Sprintf(`##INFO=<ID=(.+),Number=([\dAGR\.]*),Type=(%s),Description="(.*)">`, typeRe))
//var formatRegexp = regexp.MustCompile(fmt.Sprintf(`##FORMAT=<ID=(.+),Number=([\dAGR\.]*),Type=(%s),Description="(.*)">`, typeRe))

var typeRe = `String|Integer|Float|Flag|Character|Unknown`
var filterRegexp = regexp.MustCompile(`##FILTER=<ID=(.+),Description="(.*)">`)
var contigRegexp = regexp.MustCompile(`contig=<.*((\w+)=([^,>]+))`)
var sampleRegexp = regexp.MustCompile(`SAMPLE=<ID=([^,>]+)`)
var pedRegexp = regexp.MustCompile(`PEDIGREE=<=([^,>]+)`)

//var headerIdRegexp = regexp.MustCompile(`##([^=]+)=<ID=([^,]+)`)
var fileVersionRegexp = regexp.MustCompile(`##fileformat=VCFv(.+)`)

// Info holds the Info and Format fields
type Info struct {
	Id          string
	Description string
	Number      string            // A G R . ''
	Type        string            // STRING INTEGER FLOAT FLAG CHARACTER UNKNOWN
	fields      map[string]*Field // grendeloz
	order       []string          // grendeloz
}

// SampleFormat holds the type info for Format fields.
type SampleFormat Info

// Header holds all the type and format information for the variants.
type Header struct {
	sync.RWMutex

	SampleNames   []string
	Infos         map[string]*Info
	SampleFormats map[string]*SampleFormat
	Filters       map[string]string
	Extras        []string
	FileFormat    string
	// Contigs is a list of maps of length, URL, etc.
	Contigs []map[string]string
	// ##SAMPLE
	Samples   map[string]string
	Pedigrees []string
}

// String returns a string representation.
func (i *Info) String() string {
	// Work out original order of fields
	positions := make([]int, 0)
	ogorder := make(map[int]*Field)

	// New position-based map of fields
	for _, f := range i.fields {
		ogorder[f.Index] = f
		positions = append(positions, f.Index)
	}
	sort.Ints(positions)

	// Create field strings in original order
	fieldStrings := make([]string, 0)
	for _, k := range positions {
		f := ogorder[k]
		thisStr := f.Key + `=`
		if f.Quote != 0 {
			thisStr += string(f.Quote) + f.Value + string(f.Quote)
		} else {
			thisStr += f.Value
		}
		fieldStrings = append(fieldStrings, thisStr)
	}

	// Assemble final string
	newStr := `##INFO=<` + strings.Join(fieldStrings, `,`) + `>`
	//fmt.Println(newStr)
	return newStr
}

// String returns a string representation.
func (s *SampleFormat) String() string {
	// Work out original order of fields
	positions := make([]int, 0)
	ogorder := make(map[int]*Field)

	// New position-based map of fields
	for _, f := range s.fields {
		ogorder[f.Index] = f
		positions = append(positions, f.Index)
	}
	sort.Ints(positions)

	// Create field strings in original order
	fieldStrings := make([]string, 0)
	for _, k := range positions {
		f := ogorder[k]
		thisStr := f.Key + `=`
		if f.Quote != 0 {
			thisStr += string(f.Quote) + f.Value + string(f.Quote)
		} else {
			thisStr += f.Value
		}
		fieldStrings = append(fieldStrings, thisStr)
	}

	// Assemble final string
	newStr := `##FORMAT=<` + strings.Join(fieldStrings, `,`) + `>`
	//fmt.Println(newStr)
	return newStr
}

/*
func (h *Header) Validate(verr *VCFError) []error {
	var errs []error
	return errs
}*/

func (h *Header) parseSample(format []string, s string) (*SampleGenotype, []error) {
	values := strings.Split(s, ":")
	if len(format) != len(values) {
		return NewSampleGenotype(), []error{fmt.Errorf("bad sample string: %s", s)}
	}
	//if geno == nil {
	var value string
	var geno = NewSampleGenotype()
	var errs []error
	//}
	var e error

	for i, field := range format {
		value = values[i]
		switch field {
		case "GT":
			e = h.setSampleGT(geno, value)
		case "DP":
			e = h.setSampleDP(geno, value)
		case "GL":
			e = h.setSampleGL(geno, value, false)
		case "PL":
			e = h.setSampleGL(geno, value, true)
		case "GQ":
			if format, ok := h.SampleFormats[field]; ok {
				e = h.setSampleGQ(geno, value, format.Type)
			}
		}
		geno.Fields[field] = value
		if e != nil {
			errs = append(errs, e)
		}
	}
	return geno, errs
}

func (h *Header) setSampleDP(geno *SampleGenotype, value string) error {
	var err error
	geno.DP, err = strconv.Atoi(value)
	if err != nil && (value == "" || value == ".") {
		return nil
	}
	return err
}

func (h *Header) setSampleGQ(geno *SampleGenotype, value string, Type string) error {
	var err error
	if Type == "Integer" {
		geno.GQ, err = strconv.Atoi(value)
	} else if Type == "Float" {
		var v float64
		v, err = strconv.ParseFloat(value, 32)
		if err == nil {
			err = errors.New("setSampleGQ: GQ reported as float. rounding to int")
			geno.GQ = int(math.Floor(v + 0.5))
		}
	}

	if err != nil && (value == "" || value == ".") {
		return nil
	}
	return err
}

func (h *Header) setSampleGL(geno *SampleGenotype, value string, isPL bool) error {
	var err error
	if len(geno.GL) != 0 {
		geno.GL = geno.GL[:0]
	}
	if value == "." || value == "" {
		return nil
	}
	vals := strings.Split(value, ",")
	var v float64
	for _, val := range vals {
		v, err = strconv.ParseFloat(val, 64)
		if isPL {
			v /= -10.0
		}
		geno.GL = append(geno.GL, v)
	}
	return err
}

func (h *Header) setSampleGT(geno *SampleGenotype, value string) error {
	geno.Phased = strings.Contains(value, "|")
	splitString := "/"
	if geno.Phased {
		splitString = "|"
	}

	alleles := strings.Split(value, splitString)
	for _, allele := range alleles {
		switch allele {
		case ".":
			geno.GT = append(geno.GT, -1)
		case "1":
			geno.GT = append(geno.GT, 1)
		case "0":
			geno.GT = append(geno.GT, 0)
		default:
			v, err := strconv.Atoi(allele)
			if err != nil {
				return err
			}
			geno.GT = append(geno.GT, v)
		}
	}

	return nil
}

// NewHeader returns a Header with the requisite allocations.
func NewHeader() *Header {
	var h Header
	h.Filters = make(map[string]string)
	h.Infos = make(map[string]*Info)
	h.SampleFormats = make(map[string]*SampleFormat)
	h.SampleNames = make([]string, 0)
	h.Pedigrees = make([]string, 0)
	h.Samples = make(map[string]string)
	h.Extras = make([]string, 0)
	h.Contigs = make([]map[string]string, 0, 64)
	return &h
}

func parseHeaderInfo(info string) (*Info, error) {
	res := infoRegexp.FindStringSubmatch(info)
	i, err := infoSplitter(res[1])
	if err != nil {
		return i, err
	}

	// Set these for backwards compatibility. Note also that we are
	// ignoring errors here which is not a good idea. TO DO - fix this.
	// Also note that we need to GetValue(`ID`) to set Id.
	v, err := i.GetValue(`ID`)
	if err == nil {
		i.Id = v
	}
	v, err = i.GetValue(`Number`)
	if err == nil {
		i.Number = v
	}
	v, err = i.GetValue(`Type`)
	if err == nil {
		i.Type = v
	}
	v, err = i.GetValue(`Description`)
	if err == nil {
		i.Description = v
	}

	return i, nil
}

// Old version replaced by grendeloz
//func parseHeaderInfo(info string) (*Info, error) {
//	res := infoRegexp.FindStringSubmatch(info)
//	if len(res) != 5 {
//		return nil, fmt.Errorf("INFO error: %s, %v", info, res)
//	}
//	var i Info
//	i.Id = res[1]
//	i.Number = res[2]
//	i.Type = res[3]
//	i.Description = res[4]
//	return &i, nil
//}

func parseHeaderContig(contig string) (map[string]string, error) {
	vmap := make(map[string]string)
	contig = strings.TrimSuffix(strings.TrimPrefix(contig, "##contig=<"), ">")
	rdr := csv.NewReader(strings.NewReader(contig))
	rdr.LazyQuotes = true
	rdr.TrimLeadingSpace = true
	contigs, err := rdr.Read()

	for _, pair := range contigs {
		kv := strings.SplitN(pair, "=", 2)
		vmap[kv[0]] = kv[1]
	}
	return vmap, err
}

func parseHeaderExtraKV(kv string) ([]string, error) {
	kv = strings.TrimLeft(kv, "##")
	kv = strings.TrimLeft(kv, " ")
	kvpair := strings.SplitN(kv, "=", 2)

	if len(kvpair) != 2 {
		return nil, fmt.Errorf("header error in extra field: %s", kv)
	}
	return kvpair, nil
}

func parseHeaderFormat(info string) (*SampleFormat, error) {
	res := formatRegexp.FindStringSubmatch(info)
	f, err := formatSplitter(res[1])
	if err != nil {
		return f, err
	}

	// Set these for backwards compatibility. Note also that we are
	// ignoring errors here which is not a good idea. TO DO - fix this.
	// Also note that we need to GetValue(`ID`) to set Id.
	v, err := f.GetValue(`ID`)
	if err == nil {
		f.Id = v
	}
	v, err = f.GetValue(`Number`)
	if err == nil {
		f.Number = v
	}
	v, err = f.GetValue(`Type`)
	if err == nil {
		f.Type = v
	}
	v, err = f.GetValue(`Description`)
	if err == nil {
		f.Description = v
	}

	return f, nil
}

// Old version replaced by grendeloz
//func parseHeaderFormat(info string) (*SampleFormat, error) {
//	res := formatRegexp.FindStringSubmatch(info)
//	if len(res) != 5 {
//		return nil, fmt.Errorf("FORMAT error: %s", info)
//	}
//	var i SampleFormat
//	i.Id = res[1]
//	i.Number = res[2]
//	i.Type = res[3]
//	i.Description = res[4]
//	return &i, nil
//}

func parseHeaderFilter(info string) ([]string, error) {
	res := filterRegexp.FindStringSubmatch(info)
	if len(res) != 3 {
		return nil, fmt.Errorf("FILTER error: %s", info)
	}
	return res[1:3], nil
}

// return just the sample id.
func parseHeaderSample(line string) (string, error) {
	res := sampleRegexp.FindStringSubmatch(line)
	if len(res) != 2 {
		return "", fmt.Errorf("error parsing ##SAMPLE")
	}
	return res[1], nil
}

func parseHeaderFileVersion(format string) (string, error) {
	res := fileVersionRegexp.FindStringSubmatch(format)
	if len(res) != 2 {
		return "-1", fmt.Errorf("file format error: %s", format)
	}

	return res[1], nil
}

func parseSampleLine(line string) ([]string, error) {
	fields := strings.Split(line, "\t")
	var samples []string
	if len(fields) > 9 {
		samples = fields[9:]
	} else {
		samples = []string{}
	}
	return samples, nil
}

func parseOne(key, val, itype string) (interface{}, error) {
	var v interface{}
	var err error
	switch itype {
	case "Integer", "INTEGER":
		v, err = strconv.Atoi(val)
	case "Float", "FLOAT":
		v, err = strconv.ParseFloat(val, 32)
	case "Flag", "FLAG":
		if val != "" {
			err = fmt.Errorf("Info Error: flag field (%s) had value", key)
		}
		v = true
	default:
		v = val
	}
	return v, err
}
