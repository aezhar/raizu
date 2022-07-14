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
	"errors"
	"net/http"
)

type responseWriterProxy struct {
	proxied     http.ResponseWriter
	wroteHeader bool
}

func (p *responseWriterProxy) Header() http.Header { return p.proxied.Header() }

func (p *responseWriterProxy) Write(bytes []byte) (int, error) {
	p.wroteHeader = true
	return p.proxied.Write(bytes)
}

func (p *responseWriterProxy) WriteHeader(statusCode int) {
	p.wroteHeader = true
	p.proxied.WriteHeader(statusCode)
}

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request) error
}

type Wrapper struct{ Handler }

func (w Wrapper) ServeHTTP(wr http.ResponseWriter, rq *http.Request) {
	p := &responseWriterProxy{proxied: wr}
	if err := w.Handler.ServeHTTP(p, rq); err != nil {
		var raizuError *Error
		if !errors.As(err, &raizuError) {
			// TODO: log error here before dropping it.
			raizuError = ErrInternalServerError
		}

		if !p.wroteHeader {
			WriteError(p, raizuError)
		}

		return
	}
}

func Wrap(h Handler) http.Handler { return &Wrapper{h} }

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers. If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler that calls f.
type HandlerFunc func(http.ResponseWriter, *http.Request) error

// ServeHTTP calls f(w, r).
func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}

// WrapFunc wraps an ordinary function as a HTTP handler.
func WrapFunc(f HandlerFunc) http.Handler { return Wrap(f) }
