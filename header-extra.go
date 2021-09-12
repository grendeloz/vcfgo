package vcfgo

import (
	"errors"
	"regexp"
)

// GrendelOz 2021
// This was supposed to be a small redesign to allow for more permissive
// parsing of headers and then at some point it turned into a full
// (possibly poorly considered) redesign of the header logic to allow
// for easier mods going forward.
// There is a natural tension between the generic and the specific in
// VCF headers because while there is ONLY ONE mandatory meta-information
// line (##fileformat=, at least in VCFv4.3) there are a great many
// common and expected meta-info lines. And while there are only two
// underlying types of meta-info lines (structured and unstructured), each
// different type (as defined by the KEY) of meta-line has different
// The original vcfgo went down the route of specific parsing and
// validation code and golang structs/types for each different type of
// meta-line.
// In the redesign, I'm trying to make the parsing and use of meta-lines
// generic, the validation will still have to be specific.

type parserState int

const (
	inKey         parserState = iota + 1 // EnumIndex = 1
	inValue                              // EnumIndex = 2
	inQuote                              // EnumIndex = 3
	inQuotedValue                        // EnumIndex = 4
	inKvSeparator                        // EnumIndex = 5
)

var (
	ErrKeyNotFound  = errors.New("vcfgo: key not found")
	ErrDuplicateKey = errors.New("vcfgo: key cannot be added multiple times")
	ErrLinePattern  = errors.New("vcfgo: unexpected header line")
)

var infoRegexp = regexp.MustCompile(`##INFO=<(.+)>`)
var formatRegexp = regexp.MustCompile(`##FORMAT=<(.+)>`)

// GetValue returns the value for a given key. If the key does not exist,
// an ErrKeyNotFound error is returned.
func (i *Info) GetValue(k string) (string, error) {
	if f, found := i.fields[k]; found {
		return f.Value, nil
	} else {
		return ``, ErrKeyNotFound
	}
}

// GetField returns the Field for a given key. If the key does not exist,
// an ErrKeyNotFound error is returned.
func (i *Info) GetField(k string) (*Field, error) {
	if f, found := i.fields[k]; found {
		return f, nil
	} else {
		return nil, ErrKeyNotFound
	}
}

// GetValue returns the value for a given key. If the key does not exist,
// an ErrKeyNotFound error is returned.
func (s *SampleFormat) GetValue(k string) (string, error) {
	if f, found := s.fields[k]; found {
		return f.Value, nil
	} else {
		return ``, ErrKeyNotFound
	}
}

// GetField returns the Field for a given key. If the key does not exist,
// an ErrKeyNotFound error is returned.
func (s *SampleFormat) GetField(k string) (*Field, error) {
	if f, found := s.fields[k]; found {
		return f, nil
	} else {
		return nil, ErrKeyNotFound
	}
}

// NewInfo allocates the internals and returns a *Info
func NewInfo() *Info {
	var i Info
	i.fields = make(map[string]*Field)
	i.order = make([]string, 0)
	return &i
}

// NewInfoFromString parses a key=value string and returns a *Info.
// Note that the string is not a full INFO line from the header but just
// that portion between < and > in the INFO line. For example
// `ID=DP,Number=1,Type=Integer,Description="Total Depth"` not
// `##INFO=<ID=DP,Number=1,Type=Integer,Description="Total Depth">`
func NewInfoFromString(s string) (*Info, error) {
	return infoSplitter(s)
}

func infoSplitter(s string) (*Info, error) {
	var i Info
	fields, order, err := kvSplitter(s)
	if err != nil {
		return &i, err
	}

	i.fields = fields
	i.order = order

	return &i, nil
}

func formatSplitter(s string) (*SampleFormat, error) {
	var f SampleFormat
	fields, order, err := kvSplitter(s)
	if err != nil {
		return &f, err
	}

	f.fields = fields
	f.order = order

	return &f, nil
}
