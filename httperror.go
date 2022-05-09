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
	ErrBadRequest                   = NewErrorOpts(http.StatusBadRequest)
	ErrUnauthorized                 = NewErrorOpts(http.StatusUnauthorized)
	ErrPaymentRequired              = NewErrorOpts(http.StatusPaymentRequired)
	ErrForbidden                    = NewErrorOpts(http.StatusForbidden)
	ErrNotFound                     = NewErrorOpts(http.StatusNotFound)
	ErrMethodNotAllowed             = NewErrorOpts(http.StatusMethodNotAllowed)
	ErrNotAcceptable                = NewErrorOpts(http.StatusNotAcceptable)
	ErrProxyAuthRequired            = NewErrorOpts(http.StatusProxyAuthRequired)
	ErrRequestTimeout               = NewErrorOpts(http.StatusRequestTimeout)
	ErrConflict                     = NewErrorOpts(http.StatusConflict)
	ErrGone                         = NewErrorOpts(http.StatusGone)
	ErrLengthRequired               = NewErrorOpts(http.StatusLengthRequired)
	ErrPreconditionFailed           = NewErrorOpts(http.StatusPreconditionFailed)
	ErrRequestEntityTooLarge        = NewErrorOpts(http.StatusRequestEntityTooLarge)
	ErrRequestURITooLong            = NewErrorOpts(http.StatusRequestURITooLong)
	ErrUnsupportedMediaType         = NewErrorOpts(http.StatusUnsupportedMediaType)
	ErrRequestedRangeNotSatisfiable = NewErrorOpts(http.StatusRequestedRangeNotSatisfiable)
	ErrExpectationFailed            = NewErrorOpts(http.StatusExpectationFailed)
	ErrTeapot                       = NewErrorOpts(http.StatusTeapot)
	ErrMisdirectedRequest           = NewErrorOpts(http.StatusMisdirectedRequest)
	ErrUnprocessableEntity          = NewErrorOpts(http.StatusUnprocessableEntity)
	ErrLocked                       = NewErrorOpts(http.StatusLocked)
	ErrFailedDependency             = NewErrorOpts(http.StatusFailedDependency)
	ErrTooEarly                     = NewErrorOpts(http.StatusTooEarly)
	ErrUpgradeRequired              = NewErrorOpts(http.StatusUpgradeRequired)
	ErrPreconditionRequired         = NewErrorOpts(http.StatusPreconditionRequired)
	ErrTooManyRequests              = NewErrorOpts(http.StatusTooManyRequests)
	ErrRequestHeaderFieldsTooLarge  = NewErrorOpts(http.StatusRequestHeaderFieldsTooLarge)
	ErrUnavailableForLegalReasons   = NewErrorOpts(http.StatusUnavailableForLegalReasons)

	ErrInternalServerError = NewErrorOpts(http.StatusInternalServerError)
)
