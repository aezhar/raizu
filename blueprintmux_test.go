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
	"net/http"
	"path"
	"testing"

	"github.com/appleboy/gofight/v2"
	assertPkg "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	requirePkg "github.com/stretchr/testify/require"

	"github.com/aezhar/raizu"
	"github.com/aezhar/raizu/raizumux"
)

type HandlerMock struct{ mock.Mock }

func (m *HandlerMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Called(r.URL.String())
}

func TestAppMux2(t *testing.T) {

}

func Test(t *testing.T) {
	tt := []struct {
		name   string
		prefix string
	}{
		// TODO: tc cases
		{"NoPrefix", "/"},
		{"Prefixed", "/foo"},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			hm1 := &HandlerMock{}
			hm1.On("ServeHTTP", "/").Return()
			bp1 := raizu.HandlerAppBlueprint(raizu.HandlerAppConfig{
				Prefix:  "/app1",
				Handler: hm1,
			})

			hm2 := &HandlerMock{}
			bp2 := raizu.HandlerAppBlueprint(raizu.HandlerAppConfig{
				Prefix:  "/app2",
				Handler: hm2,
			})

			bpm := raizu.AppMuxBlueprint(raizu.AppMuxConfig[*raizumux.Mux]{
				Prefix:     tc.prefix,
				NewMuxFn:   raizumux.New,
				Blueprints: []raizu.Blueprint{bp1, bp2},
			})

			hp, err := bpm.NewHandlerProvider()
			requirePkg.NoError(t, err)

			h := hp.Handler()

			r := gofight.New()
			r.GET(path.Clean(tc.prefix+"/app1/")+"/").
				Run(h, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
					assertPkg.Equal(t, http.StatusOK, r.Code)
				})

			hm1.AssertExpectations(t)
			hm2.AssertNumberOfCalls(t, "ServeHTTP", 0)
		})
	}
}
