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

// A Handler responds to an HTTP request.
//
// ServeHTTP should write reply headers and data to the ResponseWriter
// and then return nil to indicate that the request was successfully
// processed or return an error if handling the request failed for
// whatever reason.
//
// All other notes from [http.Handler] apply here as well.
type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request) error
}

// HTTPHandler wraps a [raizu.Handler] as [http.Handler].
//
// Use the WrapHandler constructor function as a shorthand for wrapping a
// [raizu.Handler] as an [http.Handler].
type HTTPHandler struct{ Handler }

func (h HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := &responseProxy{proxied: w}
	if err := h.Handler.ServeHTTP(p, r); err != nil {
		var raizuError *Error
		if !errors.As(err, &raizuError) {
			raizuError = ErrInternalServerError
		}

		if !p.wroteHeader {
			WriteError(p, raizuError)
		}
	}
}

// WrapHandler wraps a [raizu.Handler] as [http.Handler].
func WrapHandler(h Handler) http.Handler { return &HTTPHandler{h} }

// HandlerFn is an adapter to allow the use of ordinary functions as
// [raizu.Handler] handlers. If f is a function with the appropriate
// signature, HandlerFn(f) is a [raizu.Handler] that calls f.
type HandlerFn func(http.ResponseWriter, *http.Request) error

// ServeHTTP calls f(w, r).
func (f HandlerFn) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}

// WrapFn wraps a function as an [http.HandlerFunc].
func WrapFn(f HandlerFn) http.HandlerFunc { return WrapHandler(f).ServeHTTP }

// Wrap wraps a function as an [http.Handler].
func Wrap(f HandlerFn) http.Handler { return WrapHandler(f) }
