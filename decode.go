package aproto

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"
)

//func IsAsciiPrintable(s string) bool {
//for _, r := range s {
//if r > unicode.MaxASCII || !unicode.IsPrint(r) {
//return false
//}
//}
//return true
//}

// 0	Varint	int32, int64, uint32, uint64, sint32, sint64, bool, enum
// 1	64-bit	fixed64, sfixed64, double
// 2	Length-delimited	string, bytes, embedded messages, packed repeated fields
// 3	Start group	groups (deprecated)
// 4	End group	groups (deprecated)
// 5	32-bit	fixed32, sfixed32, float
func decode_1_chunk(
	data []byte,
) (chk Chunk, chunk_len uint64, e error) {
	pos := 0

	id_type, id_type_len := proto.DecodeVarint(data)
	if id_type_len <= 0 || id_type_len > 16 {
		e = errors.New("malformed id_type_len")
		return
	}
	pos += id_type_len

	if pos >= len(data) {
		e = errors.New("not enough data for any furter wire type")
		return
	}

	id := int(id_type >> 3)

	if id > 536870911 { // max field: 2^29 - 1 == 536870911
		e = errors.New("field number > 2^29-1")
		return
	}
	_type := id_type & 7

	id_type_bytes := data[0:id_type_len]

	chunk_len += uint64(id_type_len)

	switch _type {
	case 0: // varint
		// overflow, not enough data
		if pos+1 > len(data) {
			e = errors.New("not enough data for wire type 0(varint)")
			return
		}

		u64, u64_len := proto.DecodeVarint(data[pos:])
		if u64 == 0 && u64_len == 0 {
			e = errors.New("fail DecodeVarint()")
			return
		}
		chk = &Varint{
			Value: u64,
			IdType: IdType{
				Id:   id,
				Type: _type,
				data: id_type_bytes,
			},
		}

		chunk_len += uint64(u64_len)

	case 1: // fixed 64 / double
		// overflow, not enough data
		if pos+8 > len(data) {
			// fmt.Println("pos: ", pos, ", len(data): ", len(data), ", s_len: ", int(s_len))
			e = errors.New("not enough data for wire type 1(fixed64)")
			return
		}
		u64 := binary.LittleEndian.Uint64(data[pos : pos+8])
		chk = &Fixed64{
			Value: u64,
			IdType: IdType{
				Id:   id,
				Type: _type,
				data: id_type_bytes,
			},
		}

		chunk_len += 8

	case 2: // struct / string
		s_len, s_len_len := proto.DecodeVarint(data[pos:])
		if s_len == 0 && s_len_len == 0 {
			e = errors.New("fail DecodeVarint()")
			return
		}
		chunk_len += uint64(s_len_len)
		pos += s_len_len

		// overflow, not enough data
		if uint64(pos)+s_len > uint64(len(data)) {
			e = errors.New("not enough data for wire type 2(string)")
			return
		}

		str := data[pos : pos+int(s_len)]

		_struct := &Struct{
			DataLen: len(str),
			IdType: IdType{
				Id:   id,
				Type: _type,
				data: id_type_bytes,
			},
		}

		// try to decode as inner struct first
		chunks, err2 := decode_all_chunks(str)

		// if decode success, treat as struct
		if err2 == nil {
			_struct.Children = chunks
			_struct.Str = str
			chk = _struct

			chunk_len += uint64(s_len)
			return
		} else { // decode fail, just treat as string
			_struct.Str = str
			chk = _struct

			chunk_len += uint64(s_len)
		}

	case 3:
		e = errors.New("[proto 3] not implemented")
		return
	case 4:
		e = errors.New("[proto 4] not implemented")
		return
	case 5: // fixed 32 / float
		if pos+4 > len(data) {
			e = errors.New("not enough data for wire type 5(fixed32)")
			return
		}
		u32 := binary.LittleEndian.Uint32(data[pos : pos+4])

		chk = &Fixed32{
			value: u32,
			IdType: IdType{
				Id:   id,
				Type: _type,
				data: id_type_bytes,
			},
		}

		chunk_len += 4
	default:
		e = errors.New(fmt.Sprintf("Unknown wire type %d of id_type %x", _type, id_type))
		return
	}
	return
}
func decode_all_chunks(data []byte) ([]Chunk, error) {
	var pos uint64 = 0
	var ret []Chunk

	for pos < uint64(len(data)) {
		chunk, chunk_len, e := decode_1_chunk(data[pos:])
		if e != nil {
			return ret, e
		}
		ret = append(ret, chunk)
		pos += chunk_len
	}

	return ret, nil
}

// dump with all kinds of Renderer
func TryDumpEx(data []byte, r Renderer) (string, error) {
	chunks, e := decode_all_chunks(data)
	if e != nil {
		return ``, e
	}
	ret := ``

	for _, ch := range chunks {
		ret += ch.Render(``, r)
		ret += r.NEWLINE()
	}
	return ret, nil
}

// dump to Console
func TryDump(data []byte) (string, error) {
	return TryDumpEx(data, &ConsoleRenderer{})
}

// dump to Console and ignore error
func Dump(data []byte) string {
	ret, e := TryDump(data)
	if e != nil {
		return ""
	}
	return ret
}
