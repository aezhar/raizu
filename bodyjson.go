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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func decodeJsonBody(r io.Reader, out any) error {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&out); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return NewError(http.StatusBadRequest, msg, ErrMalformedContent)

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return NewError(http.StatusBadRequest, msg, ErrMalformedContent)

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return NewError(http.StatusBadRequest, msg, ErrMalformedContent)

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return NewError(http.StatusBadRequest, msg, ErrMalformedContent)

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return NewError(http.StatusBadRequest, msg, ErrMalformedContent)

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 64KiB"
			return NewError(http.StatusRequestEntityTooLarge, msg, ErrTooLarge)

		default:
			return err
		}
	}

	// Call decode again to make sure the request body only
	// contains a single JSON object. Decode will return io.EOF
	// if there was only one JSON object. So if we get anything else,
	// we know that there is additional data in the request body.
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return NewError(http.StatusBadRequest, msg, ErrMalformedContent)
	}

	return nil
}
