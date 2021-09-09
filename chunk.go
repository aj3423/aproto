package aproto

import (
	"fmt"
	"math"
	"strings"
)

type ChunkType int

// https://developers.google.com/protocol-buffers/docs/encoding
const (
	CT_Varint ChunkType = iota
	CT_Fixed64
	CT_Struct
	CT_Deprecated_3
	CT_Deprecated_4
	CT_Fixed32
)

func type_str(t ChunkType) string {
	switch t {
	case CT_Varint:
		return `varint`
	case CT_Fixed64:
		return `fixed64/double`
	case CT_Struct:
		return `string`
	case CT_Fixed32:
		return `fixed32/float`
	}
	return `TODO chunk type`
}

// ---- interface ----
type Chunk interface {
	Render(indent string, r Renderer) string
	Type() ChunkType
}

// ---- IdType ----
type IdType struct {
	Id   int
	Type uint64

	data []byte
}

func (x *IdType) Render(indent string, r Renderer) string {
	return fmt.Sprintf("%s[%s] %s %s: ",
		indent,
		r.IDTYPE(fmt.Sprintf("% x", x.data)),
		r.ID(fmt.Sprintf("%d", x.Id)),
		r.TYPE(type_str(ChunkType(x.Type))),
	)
}

// ---- varint ----
type Varint struct {
	IdType
	/*
		The value is either of:
			int32, int64, uint32, uint64, sint32, sint64, bool, enum
		So an 'uint64' should be enough for 8 bytes
	*/
	Value uint64
}

func (x *Varint) Render(indent string, r Renderer) string {
	return x.IdType.Render(indent, r) +
		fmt.Sprintf("%s (%s)",
			r.NUM(fmt.Sprintf("%d", int64(x.Value))),
			r.NUM(fmt.Sprintf("0x%x", x.Value)),
		)
}
func (x *Varint) Type() ChunkType {
	return CT_Varint
}

// ---- fixed64 ----
type Fixed64 struct {
	Value uint64
	IdType
}

func (x *Fixed64) Render(indent string, r Renderer) string {
	f64 := math.Float64frombits(x.Value)

	return x.IdType.Render(indent, r) +
		fmt.Sprintf("%s (%s) (%s)",
			fmt.Sprintf("%d", uint64(x.Value)),
			fmt.Sprintf("0x%x", x.Value),
			fmt.Sprintf("%f", float64(f64)),
		)
}
func (x *Fixed64) Type() ChunkType {
	return CT_Fixed64
}

// ---- struct ----
type Struct struct {
	IdType

	DataLen int

	Str      []byte  // is string, not struct
	Children []Chunk // is struct, not string
}

func (x *Struct) Render(indent string, r Renderer) string {
	ret := x.IdType.Render(indent, r)

	ret += fmt.Sprintf("(%d): ", x.DataLen) // show length first

	if x.Children != nil { // is struct

		/*
			Sometimes a binary string may be miss parsed to struct, eg:
				"31303030303633303535313238384846"
			is parsed to:
				[31] 6 fixed64/double: 3832619590722007088 (0x3530333630303030) (0.000000)
				[35] 6 fixed32/float: 943206961 (0x38383231) (0.000044)
				[48] 9 varint: 70 (0x46)
			The longer the data is, the less likely it could be miss parsed.
			So print the hex dump first, if it's not very long(<=32 bytes)
		*/
		//if len(x.Str) <= 32 {
		//ret += fmt.Sprintf("% x", x.Str)
		//}

		ret += r.NEWLINE()

		lines := []string{}
		for _, ch := range x.Children {
			lines = append(lines, ch.Render(indent+r.INDENT(), r))
		}
		ret += strings.Join(lines, r.NEWLINE())

	} else { // is string

		// detect charset like GBK...
		bs, charset, e := detect_charset(x.Str)

		// find a charset, it's printable
		if e == nil { // empty string also goes here
			// show name of charset if not utf8
			if charset != `utf8` {
				ret += `[` + charset + `] `
			}
			ret += r.STR(string(bs))

			// in practice, string may contain special character, also print hex dump if string is short
			if len(x.Str) <= 8 && len(x.Str) > 0 {
				ret += " (" + r.STR(fmt.Sprintf("% x", x.Str)) + ")"
			}
		} else {
			if len(x.Str) > 32 { // too long, only show first 32 bytes, and "..."
				ret += r.STR(fmt.Sprintf("% x", x.Str[0:32]))
				ret += " ..."
			} else {
				ret += r.STR(fmt.Sprintf("% x", x.Str))
			}
		}
	}
	return ret
}
func (x *Struct) Type() ChunkType {
	return CT_Struct
}

// ---- fixed32 ----
type Fixed32 struct {
	IdType
	value uint32
}

func (x *Fixed32) Render(indent string, r Renderer) string {

	f32 := math.Float32frombits(x.value)
	return x.IdType.Render(indent, r) +
		fmt.Sprintf("%s (%s) (%s)",
			fmt.Sprintf("%d", int32(x.value)),
			fmt.Sprintf("0x%x", x.value),
			fmt.Sprintf("%f", f32),
		)
}
func (x *Fixed32) Type() ChunkType {
	return CT_Fixed32
}
