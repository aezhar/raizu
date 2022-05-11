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
	"fmt"
	"net/http"
)

func StartJSONResponse(wr http.ResponseWriter, statusCode int) {
	wr.Header().Set("X-Content-Type-Options", "nosniff")
	wr.Header().Set("Content-Type", "text/plain; charset=utf-8")
	wr.WriteHeader(statusCode)
}

func WriteJSON(wr http.ResponseWriter, statusCode int, body any) error {
	StartJSONResponse(wr, statusCode)
	return json.NewEncoder(wr).Encode(body)
}

func WriteError(wr http.ResponseWriter, err *Error) {
	http.Error(wr, err.Error(), err.status)
}

func WriteOKMessage(wr http.ResponseWriter, msg string) error {
	wr.Header().Set("Content-Type", "text/plain; charset=utf-8")
	wr.Header().Set("X-Content-Type-Options", "nosniff")
	wr.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(wr, "%s\n", msg)
	return nil
}

func WriteOK(wr http.ResponseWriter) error {
	return WriteOKMessage(wr, "ok")
}
