// * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
// Copyright(c) 2022 individual contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// <https://www.apache.org/licenses/LICENSE-2.0>
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.
// * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *

package raizu

import (
	"bufio"
	"io"
	"net"
	"net/http"
)

type responseProxy struct {
	proxied     http.ResponseWriter
	wroteHeader bool
}

// Ensure responseProxy implements the required interfaces.
var (
	_ io.Writer       = &responseProxy{}
	_ io.StringWriter = &responseProxy{}
	_ io.ReaderFrom   = &responseProxy{}

	_ http.ResponseWriter = &responseProxy{}
	_ http.Hijacker       = &responseProxy{}
	_ http.Flusher        = &responseProxy{}
)

func (p *responseProxy) write(dataB []byte, dataS string) (n int, err error) {
	if !p.wroteHeader {
		p.WriteHeader(http.StatusOK)
	}

	if dataB != nil {
		return p.proxied.Write(dataB)
	} else {
		if sw, ok := p.proxied.(io.StringWriter); ok {
			return sw.WriteString(dataS)
		}
		return p.proxied.Write([]byte(dataS))
	}
}

func (p *responseProxy) Header() http.Header {
	return p.proxied.Header()
}

func (p *responseProxy) Write(data []byte) (n int, err error) {
	return p.write(data, "")
}

func (p *responseProxy) WriteString(data string) (n int, err error) {
	return p.write(nil, data)
}

func (p *responseProxy) WriteHeader(statusCode int) {
	p.wroteHeader = true
	p.proxied.WriteHeader(statusCode)
}

func (p *responseProxy) ReadFrom(r io.Reader) (n int64, err error) {
	rf, ok := p.proxied.(io.ReaderFrom)
	if !ok {
		return 0, &UnsupportedError{"proxied handler does not implement ReadFrom"}
	}
	return rf.ReadFrom(r)
}

func (p *responseProxy) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	rf, ok := p.proxied.(http.Hijacker)
	if !ok {
		return nil, nil, &UnsupportedError{"proxied handler does not implement http.Hijacker"}
	}
	return rf.Hijack()
}

func (p *responseProxy) Flush() {
	if fl, ok := p.proxied.(http.Flusher); ok {
		fl.Flush()
	}
}

func HasReadFrom(w http.ResponseWriter) bool {
	if p, ok := w.(*responseProxy); ok {
		_, ok = p.proxied.(io.ReaderFrom)
		return ok
	}
	_, ok := w.(io.ReaderFrom)
	return ok
}

func HasHijack(w http.ResponseWriter) bool {
	if p, ok := w.(*responseProxy); ok {
		_, ok = p.proxied.(http.Hijacker)
		return ok
	}
	_, ok := w.(http.Hijacker)
	return ok
}

func HasFlush(w http.ResponseWriter) bool {
	if p, ok := w.(*responseProxy); ok {
		_, ok = p.proxied.(http.Flusher)
		return ok
	}
	_, ok := w.(http.Flusher)
	return ok
}
