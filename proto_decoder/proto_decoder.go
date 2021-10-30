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
)

func UnZlib(b []byte) ([]byte, error) {
	z, err := zlib.NewReader(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(z)
}

func error_assert(e error) {
	if e != nil {
		fmt.Println("Error:")
		fmt.Println(e.Error())
		os.Exit(1)
	}
}

func main() {

	file := flag.String("file", "", "file contain proto like 080d....")
	is_bin := flag.Bool("bin", false, "content is binary")
	is_b64 := flag.Bool("b64", false, "content is base64")
	is_zlib := flag.Bool("zlib", false, "content is zlib")
	is_xxd := flag.Bool("xxd", false, "content is hexdump (without offset and text)")

	flag.Parse()
	if len(os.Args) == 1 {
		flag.Usage()
		return
	}

	var e error

	var proto []byte // binary

	// 1. file ?
	if *file != "" {
		proto, e = ioutil.ReadFile(*file)
		error_assert(e)
	} else {
		proto = []byte(flag.Arg(0)) // last param
	}

	// base64 ?
	if *is_b64 {
		proto, e = base64.StdEncoding.DecodeString(string(proto))
		error_assert(e)
	} else {
		// binary ?
		if !*is_bin {
			proto, e = hex.DecodeString(string(proto))
			error_assert(e)
		}
	}

	// 4. zlib ?
	if *is_zlib {
		var e error
		proto, e = UnZlib(proto)
		error_assert(e)
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
		proto, e = hex.DecodeString(s)
		error_assert(e)
	}

	p, e := aproto.TryDump(proto)
	error_assert(e)

	fmt.Println(p)
}
