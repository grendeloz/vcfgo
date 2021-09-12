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
var structuredMetaRegexp = regexp.MustCompile(`^##(\w+)=<(.+)>`)

// Structured lines will also match this pattern so when parsing a new
// line, it should be checked against the structured pattern first and only
// if that fails should it be checked against the unstructured pattern.
var unstructuredMetaRegexp = regexp.MustCompile(`^##(\w+)=(.+)`)

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

	// 0-based index of where this Field appeared in the original
	// string. It is used to recreate meta lines as strings with the
	// key=value pairs in the same order as they were in the original.
	Index int

	// Quote character ([`'"]) if any that was used for the value of the
	// key-value pair. The spec does not state that double quotes must
	// be used for all quoting but it may be so. In any case, we can cope
	// with any of the 3 quoting characters shown above. Quote is empty
	// if the Value was not quoted.
	Quote rune
}

// MetaType - Create enum for header meta information line type.
type MetaType int

// Declare related constants for each MetaType starting with index 1
const (
	Unstructured MetaType = iota // EnumIndex = 0
	Structured                   // EnumIndex = 1
)

// String - Creating common behaviour - give the type a String function
func (m MetaType) String() string {
	return [...]string{"Unstructured", "Structured"}[m]
}

// EnumIndex - Creating common behaviour - give the type an EnumIndex function
func (m MetaType) EnumIndex() int {
	return int(m)
}

// MetaLine is designed to hold information from both structured and
// unstructured meta information lines from the VCF header. Fields and
// Order will only be set for structured lines and Value will only be set
// for unstructured lines.
// different fields set for the different MetaTypes.
type MetaLine struct {
	LineNumber int64

	// MetaType defaults to Unstructured. You can manually set this
	// value but it's best not to. Let the package do the work.
	MetaType MetaType

	// The basic XXX= value which is present in both STructured and
	// Unstructured MetaLines.
	LineKey string

	// Value is only used in Unstructured MetaLines - STructured
	// MetaLines use Fields and Order instead.
	Value string

	// Fields and Order contain the key=value items (as Fields) from a
	// Structured MetaLine plus the order in which they occured in the
	// OgString or the order in which they were added with AddField().
	// The Order is obeyed by String()
	Fields map[string]*Field
	Order  []string

	// OgString is only available if the MetaLine was created via
	// NewMetaLineFromString().
	OgString string
}

// NewMetaLine returns a pointer to a MetaLine. By default, the MetaType
// is Unstructured. If you use the AddField() function, MetaType will be
// automatically converted to Structured.
func NewMetaLine() *MetaLine {
	var m MetaLine
	m.Fields = make(map[string]*Field)
	m.Order = make([]string, 0)
	return &m
}

// NewMetaLineFromString matches the input string against the pattern
// for Structured and Unstructured MetaLines and returns a MetaLine. If
// neither pattern matches, it throws an error.
func NewMetaLineFromString(s string) (*MetaLine, error) {
	var m MetaLine

	if structuredMetaRegexp.Find([]byte(s)) != nil {
		res := structuredMetaRegexp.FindStringSubmatch(s)

		if len(res) != 3 {
			return &m, fmt.Errorf("%w - structured line: %s", ErrLinePattern, s)
		}

		fields, order, err := kvSplitter(res[2])
		if err != nil {
			return &m, err
		}
		m.MetaType = Structured
		m.LineKey = res[1]
		m.Fields = fields
		m.Order = order
		m.OgString = s

		return &m, nil
	} else if unstructuredMetaRegexp.Find([]byte(s)) != nil {

		res := unstructuredMetaRegexp.FindStringSubmatch(s)
		if len(res) != 3 {
			return &m, fmt.Errorf("%w - unstructured line: %s", ErrLinePattern, s)
		}

		m.MetaType = Unstructured
		m.LineKey = res[1]
		m.Value = res[2]
		m.OgString = s

		return &m, nil
	} else {
		return &m, fmt.Errorf("%w - %s", ErrLinePattern, s)
	}

}

// String returns a string representation.
func (m *MetaLine) String() (string, error) {
	if m.MetaType == Structured {
		// Work out original order of fields
		positions := make([]int, 0)
		ogorder := make(map[int]*Field)

		// New position-based map of fields
		for _, f := range m.Fields {
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
		newStr := `##` + m.LineKey + `=<` + strings.Join(fieldStrings, `,`) + `>`
		return newStr, nil
	} else if m.MetaType == Unstructured {
		// This is a trivial case
		newStr := `##` + m.LineKey + `=` + m.Value
		return newStr, nil
	}
	// If we get to here then m.MetaType is borked.
	return ``, fmt.Errorf("MetaType has an unexpected value: %v", m.MetaType)
}

// StructuredMeta can hold any of the structured meta-information
// lines, i.e. those that have the pattern '##KEY=<(key=value)+>'.
type StructuredMeta struct {
	LineNumber int64
	LineKey    string
	Fields     map[string]*Field
	Order      []string
	OgString   string // only available if created via NewMetaStructuredFromString()
}

// metaUnstructured can hold any of the structured meta-information
// lines, i.e. those that have the pattern '##KEY=<(key=value)+>'.
type UnstructuredMeta struct {
	LineNumber int64
	Key        string
	Value      string
	OgString   string // only available if created via NewMetaUnstructuredFromString()
}

// NewStructuredMeta allocates the internals and returns a *Info
func NewStructuredMeta() *StructuredMeta {
	var m StructuredMeta
	m.Fields = make(map[string]*Field)
	m.Order = make([]string, 0)
	return &m
}

// NewUnstructuredMeta allocates the internals and returns a *Info
func NewUnstructuredMeta() *UnstructuredMeta {
	var m UnstructuredMeta
	return &m
}

func NewStructuredMetaFromString(s string) (*StructuredMeta, error) {
	var m StructuredMeta
	res := structuredMetaRegexp.FindStringSubmatch(s)
	if len(res) != 3 {
		return &m, fmt.Errorf("vcfgo: line did not match structured line pattern [%s]", s)
	}

	fields, order, err := kvSplitter(res[2])
	if err != nil {
		return &m, err
	}
	m.LineKey = res[1]
	m.Fields = fields
	m.Order = order
	m.OgString = s

	return &m, nil
}

func NewUnstructuredMetaFromString(s string) (*UnstructuredMeta, error) {
	var m UnstructuredMeta
	res := unstructuredMetaRegexp.FindStringSubmatch(s)
	if len(res) != 3 {
		return &m, fmt.Errorf("vcfgo: line did not match unstructured line pattern [%s]", s)
	}

	m.Key = res[1]
	m.Value = res[2]
	m.OgString = s

	return &m, nil
}

// String returns a string representation.
func (m *StructuredMeta) String() string {
	// Work out original order of fields
	positions := make([]int, 0)
	ogorder := make(map[int]*Field)

	// New position-based map of fields
	for _, f := range m.Fields {
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
	newStr := `##` + m.LineKey + `=<` + strings.Join(fieldStrings, `,`) + `>`
	//fmt.Println(newStr)
	return newStr
}

// A meta Header line type using StructuredMeta by composition
type Pickle StructuredMeta

func NewPickle() *Pickle {
	var p Pickle
	p.LineKey = `PICKLE`
	p.Fields = make(map[string]*Field)
	p.Order = make([]string, 0)
	return &p
}

// GetValue returns the value for a given key. If the key does not exist,
// an ErrKeyNotFound error is returned.
func (m *StructuredMeta) GetValue(k string) (string, error) {
	if f, found := m.Fields[k]; found {
		return f.Value, nil
	} else {
		return ``, ErrKeyNotFound
	}
}

// GetField returns the Field for a given key. If the key does not exist,
// an ErrKeyNotFound error is returned.
func (m *StructuredMeta) GetField(k string) (*Field, error) {
	if f, found := m.Fields[k]; found {
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
	var quote, lastrune rune

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
			// quotes only count if they are not backspaced
			if r == quote && lastrune != '\\' {
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

		lastrune = r
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
