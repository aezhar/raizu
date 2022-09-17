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

package raizumux_test

import (
	"net/http"
	"testing"

	"github.com/appleboy/gofight/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aezhar/raizu/raizumux"
)

type AppMock struct{ mock.Mock }

func (m *AppMock) handleIndex(w http.ResponseWriter, r *http.Request) {
	m.Called(r.URL.String())
}

func (m *AppMock) handleFoo(w http.ResponseWriter, r *http.Request) {
	m.Called(r.URL.String())
}

func (m *AppMock) handleBaa(w http.ResponseWriter, r *http.Request) {
	m.Called(r.URL.String())
}

func TestMux_Mount(t *testing.T) {
	am := &AppMock{}

	am.On("handleIndex", "/").Return().Once()
	am.On("handleFoo", "/one/").Return().Once()
	am.On("handleBaa", "/two/").Return().Once()

	m := raizumux.New()

	m.Mount("/", http.HandlerFunc(am.handleIndex))
	m.Mount("/foo", http.HandlerFunc(am.handleFoo))
	m.Mount("/baa", http.HandlerFunc(am.handleBaa))

	r := gofight.New()

	for _, url := range []string{"/", "/foo/one/", "/baa/two/"} {
		r.GET(url).
			Run(m, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusOK, r.Code, "url %s", url)
			})
	}

	am.AssertExpectations(t)
}
