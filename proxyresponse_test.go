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
	"testing"

	assertPkg "github.com/stretchr/testify/assert"
)

type minimalResponse struct{}

func (r minimalResponse) Header() http.Header {
	return make(http.Header)
}

func (r minimalResponse) Write(bytes []byte) (int, error) {
	return len(bytes), nil
}

func (r minimalResponse) WriteHeader(statusCode int) {
}

type readerResponse struct{ minimalResponse }

func (r readerResponse) ReadFrom(rd io.Reader) (n int64, err error) {
	return 0, nil
}

type hijackResponse struct{ minimalResponse }

func (h hijackResponse) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return nil, nil, nil
}

type flushResponse struct{ minimalResponse }

func (h flushResponse) Flush() {
	return
}

func TestHasFunctionsFamily(t *testing.T) {
	type args struct {
		w http.ResponseWriter
	}
	tt := []struct {
		name string
		fn   func(w http.ResponseWriter) bool
		args args
		want bool
	}{
		// HasReadFrom
		{"HasReadFrom/proxied-wo", HasReadFrom, args{w: &responseProxy{proxied: minimalResponse{}}}, false},
		{"HasReadFrom/proxied-w", HasReadFrom, args{w: &responseProxy{proxied: readerResponse{}}}, true},
		{"HasReadFrom/direct-wo", HasReadFrom, args{w: minimalResponse{}}, false},
		{"HasReadFrom/direct-w", HasReadFrom, args{w: readerResponse{}}, true},

		// HasHijack
		{"HasHijack/proxied-wo", HasHijack, args{w: &responseProxy{proxied: minimalResponse{}}}, false},
		{"HasHijack/proxied-w", HasHijack, args{w: &responseProxy{proxied: hijackResponse{}}}, true},
		{"HasHijack/direct-wo", HasHijack, args{w: minimalResponse{}}, false},
		{"HasHijack/direct-w", HasHijack, args{w: hijackResponse{}}, true},

		// HasFlush
		{"HasFlush/proxied-wo", HasFlush, args{w: &responseProxy{proxied: minimalResponse{}}}, false},
		{"HasFlush/proxied-w", HasFlush, args{w: &responseProxy{proxied: flushResponse{}}}, true},
		{"HasFlush/direct-wo", HasFlush, args{w: minimalResponse{}}, false},
		{"HasFlush/direct-w", HasFlush, args{w: flushResponse{}}, true},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assertPkg.Equalf(t, tc.want, tc.fn(tc.args.w), "HasReadFrom(%v)", tc.args.w)
		})
	}
}
