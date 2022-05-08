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
	"net/http"
)

type Error struct {
	status int
	msg    *string
	inner  error
}

func (err *Error) Unwrap() error { return err.inner }

func (err *Error) Error() (s string) {
	if err.msg != nil {
		return *err.msg
	}
	return http.StatusText(err.status)
}

func NewError(status int, msg string, inner error) *Error {
	return &Error{status: status, msg: &msg, inner: inner}
}

type ErrorOption func(e *Error)

func WithMessage(msg string) ErrorOption { return func(e *Error) { e.msg = &msg } }
func WithInner(err error) ErrorOption    { return func(e *Error) { e.inner = err } }

func NewErrorOpts(status int, opts ...ErrorOption) (e *Error) {
	e = &Error{status: status}
	for _, opt := range opts {
		opt(e)
	}
	return
}

var (
	ErrInternalServerError = NewErrorOpts(http.StatusInternalServerError)
	ErrNotFound            = NewErrorOpts(http.StatusNotFound)
	ErrBadRequest          = NewErrorOpts(http.StatusBadRequest)
)
