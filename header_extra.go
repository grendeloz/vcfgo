package vcfgo

import (
	"errors"
	"fmt"
	"regexp"
	//"reflect"
	"sort"
)

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
)

// Runes (like bytes) use single quotes
var fieldSeparator rune = ','
var kvSeparator rune = '='

var infoRegexp = regexp.MustCompile(`##INFO=<(.+)>`)
var formatRegexp = regexp.MustCompile(`##FORMAT=<(.+)>`)

// Field holds key=value fields from the INFO and FORMAT Meta-information
// lines. See section 1.2 from the v4.2 VCF specification (retrieved
// 2021-08-24) at https://samtools.github.io/hts-specs/VCFv4.2.pdf
type Field struct {
	Key   string
	Value string
	Index int  // 0-based index of this Field in the parsed string
	Quote rune // empty if the value was not quoted
}

// GetValue returns the value for a given key. If the key does not exist,
// an ErrKeyNotFound error is returned.
func (i *Info) GetValue(k string) (string, error) {
	if f, found := i.kvs[k]; found {
		return f.Value, nil
	} else {
		return ``, ErrKeyNotFound
	}
}

// GetField returns the Field for a given key. If the key does not exist,
// an ErrKeyNotFound error is returned.
func (i *Info) GetField(k string) (*Field, error) {
	if f, found := i.kvs[k]; found {
		return f, nil
	} else {
		return nil, ErrKeyNotFound
	}
}

// GetValue returns the value for a given key. If the key does not exist,
// an ErrKeyNotFound error is returned.
func (s *SampleFormat) GetValue(k string) (string, error) {
	if f, found := s.kvs[k]; found {
		return f.Value, nil
	} else {
		return ``, ErrKeyNotFound
	}
}

// GetField returns the Field for a given key. If the key does not exist,
// an ErrKeyNotFound error is returned.
func (s *SampleFormat) GetField(k string) (*Field, error) {
	if f, found := s.kvs[k]; found {
		return f, nil
	} else {
		return nil, ErrKeyNotFound
	}
}

func infoSplitter(s string) (*Info, error) {
	var i Info
	kvs, err := kvSplitter(s)
	if err != nil {
		return &i, err
	}

	// Work out original order of keys
	positions := make([]int, 0, len(kvs))
	ogorder := make(map[int]string)

	// Establish new map based on position and use that to get the keys
	// in the order they appeared in the og line.
	keys := make([]string, 0)
	for _, v := range kvs {
		ogorder[v.Index] = v.Key
		positions = append(positions, v.Index)
	}
	sort.Ints(positions)
	for _, k := range positions {
		keys = append(keys, ogorder[k])
	}

	i.kvs = kvs
	i.order = keys

	return &i, nil
}

func formatSplitter(s string) (*SampleFormat, error) {
	var f SampleFormat
	kvs, err := kvSplitter(s)
	if err != nil {
		return &f, err
	}

	// Work out original order of keys
	positions := make([]int, 0, len(kvs))
	ogorder := make(map[int]string)

	// Establish new map based on position and use that to get the keys
	// in the order they appeared in the og line.
	keys := make([]string, 0)
	for _, v := range kvs {
		ogorder[v.Index] = v.Key
		positions = append(positions, v.Index)
	}
	sort.Ints(positions)
	for _, k := range positions {
		keys = append(keys, ogorder[k])
	}

	f.kvs = kvs
	f.order = keys

	return &f, nil
}

// kvSplitter returns building blocks for Info and SampleFormat structs.
func kvSplitter(s string) (map[string]*Field, error) {
	//info := NewInfo()

	//fmt.Printf("string: %s\n", s)
	kvs := make(map[string]*Field)

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
				kvs[k] = &f
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
				kvs[k] = &f
				//fmt.Printf("> i:%d r:%c k:%s v:%s ctr:%d state:%v\n", i, r, k, v, ctr, state)
				ctr++
				k = ``
				v = ``
				state = inQuote
			} else {
				v = v + string(r)
			}
		default:
			return kvs, fmt.Errorf("kvSplitter: problem parsing Key=value string i:%d r:%c k:%s v:%s state:%v\n", i, r, k, v, state)
		}
		//fmt.Printf("  i:%d r:%c k:%s v:%s state:%v\n", i, r, k, v, state)
	}

	// Quoted values have an explicit signal that the value in finished
	// - the quote character.  For unquoted values, the loop will exit when
	// the last rune has been read so we need to explicitly capture the last
	// key=value pair.
	if state == inValue {
		//fmt.Printf("> k:%s v:%s ctr:%d state:%v\n", k, v, ctr, state)
		f := Field{Key: k, Value: v, Index: ctr}
		//fmt.Printf("field: %v\n", f)
		kvs[k] = &f
	}

	return kvs, nil
}
