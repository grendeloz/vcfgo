package vcfgo

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// The VCFv4.3 spec appears to be silent on what characters are
// permissible in the "key" within a structured meta line, i.e. the XXX
// in XXX=<(.*)>.  For simplicity, I've gone with \w as a sensible
// default but that could be too restrictive.
var metaStructuredRegexp = regexp.MustCompile(`^##(\w+)=<(.+)>`)

// Structured lines will also match this pattern so every line must be
// checked against the Structured pattern first.
var metaUnstructuredRegexp = regexp.MustCompile(`^##(\w+)=(.+)`)

// Runes (like bytes) use single quotes
var fieldSeparator rune = ','
var kvSeparator rune = '='

// Field holds key=value fields from structured meta-information lines
// such as INFO, FORMAT and FILTER. See section 1.4 from the VCFv4.3
// specification (version 27 Jul 2021; retrieved 2021-09-05) at:
// https://samtools.github.io/hts-specs/VCFv4.3.pdf
type Field struct {
	Key   string
	Value string
	Index int  // 0-based index of this Field in the parsed string
	Quote rune // empty if the value was not quoted
}

// metaStructured can hold any of the structured meta-information
// lines, i.e. those that have the pattern '##KEY=<(key=value)+>'.
type metaStructured struct {
	lineNumber int64
	metaType   string
	fields     map[string]*Field
	order      []string
	ogString   string // only available if created via NewMetaStructuredFromString()
}

// NewMetaStructured allocates the internals and returns a *Info
func NewMetaStructured() *metaStructured {
	var m metaStructured
	m.fields = make(map[string]*Field)
	m.order = make([]string, 0)
	return &m
}

func NewMetaStructuredFromString(s string) (*metaStructured, error) {
	var m metaStructured
	res := metaStructuredRegexp.FindStringSubmatch(s)
	if len(res) != 3 {
		return &m, fmt.Errorf("vcfgo: line did not match unstructured line pattern [%s]", s)
	}

	fields, order, err := kvSplitter(res[2])
	if err != nil {
		return &m, err
	}
	m.metaType = res[1]
	m.fields = fields
	m.order = order
	m.ogString = s

	return &m, nil
}

// String returns a string representation.
func (m *metaStructured) String() string {
	// Work out original order of fields
	positions := make([]int, 0)
	ogorder := make(map[int]*Field)

	// New position-based map of fields
	for _, f := range m.fields {
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
	newStr := `##` + m.metaType + `=<` + strings.Join(fieldStrings, `,`) + `>`
	//fmt.Println(newStr)
	return newStr
}

// A meta Header line type using metaStructured by composition
type Pickle metaStructured

func NewPickle() *Pickle {
	var p Pickle
	p.metaType = `PICKLE`
	p.fields = make(map[string]*Field)
	p.order = make([]string, 0)
	return &p
}

// GetValue returns the value for a given key. If the key does not exist,
// an ErrKeyNotFound error is returned.
func (m *metaStructured) GetValue(k string) (string, error) {
	if f, found := m.fields[k]; found {
		return f.Value, nil
	} else {
		return ``, ErrKeyNotFound
	}
}

// GetField returns the Field for a given key. If the key does not exist,
// an ErrKeyNotFound error is returned.
func (m *metaStructured) GetField(k string) (*Field, error) {
	if f, found := m.fields[k]; found {
		return f, nil
	} else {
		return nil, ErrKeyNotFound
	}
}

// kvSplitter parses a structured meta-information line into a map of
// Fields where each Field is a key=value pair from the string.  This map
// can be used as the building block for structs such as Info, Format etc.
func kvSplitter(s string) (map[string]*Field, []string, error) {
	//info := NewInfo()

	//fmt.Printf("string: %s\n", s)
	fields := make(map[string]*Field)
	var order []string

	runes := []rune(s)
	state := inKey

	var k, v string
	var ctr int
	var quote rune

	for i, r := range runes {
		//fmt.Printf("> i:%d r:%c k:%s v:%s state:%v\n", i, r, k, v, state)
		switch state {
		case inKey:
			if r == kvSeparator {
				state = inKvSeparator
			} else {
				k = k + string(r)
			}
		case inKvSeparator:
			if r == '\'' || r == '"' || r == '`' {
				state = inQuotedValue
				quote = r
			} else {
				state = inValue
				v = v + string(r)
			}
		case inValue:
			if r == fieldSeparator {
				//fmt.Printf("> i:%d r:%c k:%s v:%s ctr:%d state:%v\n", i, r, k, v, ctr, state)
				f := Field{Key: k, Value: v, Index: ctr}
				//fmt.Printf("field: %v\n", f)
				fields[k] = &f
				order = append(order, k)
				ctr++
				k = ``
				v = ``
				state = inKey
			} else {
				v = v + string(r)
			}
		case inQuote:
			// The next rune *must* be fieldSeparator
			if r != fieldSeparator {
				//fmt.Printf("inQuote wtf! k:%s v:%s state:%v\n", k, v, state)
			}
			state = inKey
		case inQuotedValue:
			if r == quote {
				f := Field{Key: k, Value: v, Index: ctr, Quote: quote}
				//fmt.Printf("field: %v\n", f)
				fields[k] = &f
				order = append(order, k)
				//fmt.Printf("> i:%d r:%c k:%s v:%s ctr:%d state:%v\n", i, r, k, v, ctr, state)
				ctr++
				k = ``
				v = ``
				state = inQuote
			} else {
				v = v + string(r)
			}
		default:
			return fields, order, fmt.Errorf("kvSplitter: problem parsing Key=value string i:%d r:%c k:%s v:%s state:%v\n", i, r, k, v, state)
		}
		//fmt.Printf("  i:%d r:%c k:%s v:%s state:%v\n", i, r, k, v, state)
	}

	// Quoted values have an explicit signal that the value is finished
	// - the quote character.  For unquoted values, the loop will exit when
	// the last rune has been read so we need to explicitly capture the last
	// key=value pair.
	if state == inValue {
		//fmt.Printf("> k:%s v:%s ctr:%d state:%v\n", k, v, ctr, state)
		f := Field{Key: k, Value: v, Index: ctr}
		//fmt.Printf("field: %v\n", f)
		fields[k] = &f
		order = append(order, k)
	}

	return fields, order, nil
}
