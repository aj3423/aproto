package main

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aj3423/aproto"

	"github.com/fatih/color"
)

func ReadFile(fn string) []byte {
	s, err := ioutil.ReadFile(fn)
	if err != nil {
		return []byte{}
	}
	return s
}
func UnZlib(b []byte) ([]byte, error) {
	z, err := zlib.NewReader(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(z)
}

func has_flag(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func main() {

	file := flag.String("file", "", "file contain proto like 080d....")
	flag.Bool("bin", false, "content is binary")
	is_b64 := flag.Bool("b64", false, "content is base64")
	is_zlib := flag.Bool("zlib", false, "content is zlib")
	is_xxd := flag.Bool("xxd", false, "content is hexdump (without offset and text)")

	flag.Parse()
	if len(os.Args) == 1 {
		flag.Usage()
		return
	}

	var proto []byte // binary

	// 1. file ?
	if has_flag("file") {
		proto = ReadFile(*file)
	} else {
		proto = []byte(flag.Arg(0)) // last param
	}

	// base64 ?
	if *is_b64 {
		var e error
		proto, e = base64.StdEncoding.DecodeString(string(proto))
		if e != nil {
			color.HiRed(e.Error())
			return
		}
	} else {
		// binary ?
		if !has_flag("bin") {
			proto, _ = hex.DecodeString(string(proto))
		}
	}

	// 4. zlib ?
	if *is_zlib {
		var e error
		proto, e = UnZlib(proto)
		if e != nil {
			color.HiRed(e.Error())
			return
		}
	}

	if strings.ContainsAny(flag.Arg(0), " \r\n") {
		*is_xxd = true
	}

	if *is_xxd {
		s := flag.Arg(0)
		s = strings.ReplaceAll(s, " ", "")
		s = strings.ReplaceAll(s, "\n", "")
		s = strings.ReplaceAll(s, "\n", "")
		s = strings.ReplaceAll(s, "\r", "")
		s = strings.ReplaceAll(s, "\t", "")
		proto, _ = hex.DecodeString(s)
	}

	p, err := aproto.TryDump(proto)
	if err != nil {
		panic(err)
	}
	fmt.Println(p)
}
