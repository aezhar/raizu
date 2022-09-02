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

package raizu_test

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"

	assertPkg "github.com/stretchr/testify/assert"

	"github.com/aezhar/raizu"
)

func TestWrapFn(t *testing.T) {
	h := raizu.Wrap(func(w http.ResponseWriter, r *http.Request) error {
		return fs.ErrNotExist
	})

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	assertPkg.Equal(t, rr.Code, http.StatusInternalServerError)
}
