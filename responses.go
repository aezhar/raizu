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

func StartJSONResponse(w http.ResponseWriter, statusCode int) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
}

func WriteJSON(w http.ResponseWriter, statusCode int, body any) error {
	StartJSONResponse(w, statusCode)
	return json.NewEncoder(w).Encode(body)
}

func WriteError(w http.ResponseWriter, err *Error) {
	http.Error(w, err.Error(), err.status)
}

func WriteOKMessage(w http.ResponseWriter, msg string) error {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "%s\n", msg)
	return nil
}

func WriteOK(w http.ResponseWriter) error {
	return WriteOKMessage(w, "ok\n")
}
