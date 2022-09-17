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

package raizumux

import (
	"net/http"
	"strings"
)

type Mux struct{ mux http.ServeMux }

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}

func (m *Mux) Mount(prefix string, h http.Handler) {
	if !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}

	if len(prefix) > 1 {
		h = http.StripPrefix(prefix[:len(prefix)-1], h)
	}
	m.mux.Handle(prefix, h)
}

func New() *Mux {
	return &Mux{}
}
