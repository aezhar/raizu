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
	"fmt"
	"io"
	"net/http"

	"github.com/golang/gddo/httputil/header"
)

var (
	ErrUnsupportedMediaType = errors.New("unsupported media type")
	ErrMalformedContent     = errors.New("content has syntax error")
	ErrTooLarge             = errors.New("request is too large")
)

type decoderFn func(r io.Reader, out any) error

func getDecoder(contentType string) decoderFn {
	if contentType == "application/json" {
		return decodeJsonBody
	}
	return nil
}

func DecodeBody(res http.ResponseWriter, req *http.Request, out any) error {
	var decoder decoderFn = decodeJsonBody
	if req.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(req.Header, "Content-Type")

		decoder = getDecoder(value)
		if decoder == nil {
			msg := fmt.Sprintf("unsupported Content-Type %q", value)
			return NewError(http.StatusUnsupportedMediaType, msg, ErrUnsupportedMediaType)
		}
	}

	return decoder(http.MaxBytesReader(res, req.Body, 64*1024), out)
}
