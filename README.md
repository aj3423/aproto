# Description
Golang module/tool for decoding proto buffer without message definition.

# Try it online
https://168.138.55.177/

# Screenshot
![pb](https://user-images.githubusercontent.com/4710875/122819817-7ff91d00-d30d-11eb-9c0a-c8d46ee2b821.png)

# Usage of library aproto
The library provides two Renderers: **Console** and **Html**, which are used for Console program and web site.

- ConsoleRenderer
```
  out, e := aproto.TryDump([]byte{...})
  fmt.Println(out)
  
  // or:
  out, e := aproto.TryDumpEx([]byte{...}, &aproto.ConsoleRenderer{})

```
- HtmlRenderer (used on demo site)
```
  out, e := aproto.TryDumpEx([]byte{...}, &aproto.HtmlRenderer{})
  // transfer the output to client browser, render it with
  $('#div').text(out)
```
- Or create other custom Renders, just follow the `Renderer` interface

# Use the prebuilt tool "proto_decoder"
supported: text/file with hex-string/binary/base64/zlib encoding

eg:
- hex string: `pro 120123`
- space/tab/newline will be trimmed: 
```
pro "08 01 12 03   04 05
  06 07 08
  09 10 111213"
```
- base64 + zlib string: `pro -b64 EgEj`
- binary file: `pro -bin -file a.bin`
- zlib+base64 file: `pro -zlib -b64 -file a.bin`
- ...

# Or build it yourself
1. Install Golang
2. clone this repo: `git clone https://github.com/aj3423/aproto`
3. go to binary dir: `cd aproto/proto_decoder`
4. `go build .`

# License
MIT
