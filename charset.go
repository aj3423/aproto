package aproto

import (
	"errors"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
)

type Charset interface {
	Match([]byte) ([]byte, error)
	Name() string
}

// ---- utf8 ----
type Utf8 struct{}

func (x *Utf8) Match(data []byte) ([]byte, error) {
	if utf8.Valid(data) {
		return data, nil
	}
	return nil, errors.New(`not u8`)
}
func (x *Utf8) Name() string {
	return `utf8`
}

// ---- GBK ----
type GBK struct{}

func (x *GBK) Match(data []byte) ([]byte, error) {
	return simplifiedchinese.GBK.NewDecoder().Bytes(data)
}
func (x *GBK) Name() string {
	return `GBK`
}

/*
Used to detect string encoding, utf8 by default
Customize this array to decode other charset, eg: &GBK{}
Custom charset causes too much miss-decoding,
normal binary rendered as GBK, so disabled by default.
*/
var List []Charset = []Charset{
	&Utf8{},
	//&GBK{},
}

func detect_charset(data []byte) ([]byte, string, error) {
	for _, enc := range List {
		ret, e := enc.Match(data)
		if e == nil {
			return ret, enc.Name(), nil
		}
	}
	return nil, ``, errors.New(`unknown charset`)
}
