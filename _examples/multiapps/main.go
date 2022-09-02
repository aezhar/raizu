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

package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/aezhar/raizu"
	"github.com/aezhar/raizu/_examples/multiapps/app1"
	"github.com/aezhar/raizu/_examples/multiapps/app2"
	"github.com/aezhar/raizu/raizumux"
)

func mainE() error {
	bp := raizu.AppMuxBlueprint(raizu.AppMuxConfig{
		Prefix:   "/",
		NewMuxFn: raizumux.New,
		Blueprints: []raizu.Blueprint{
			app1.Blueprint(app1.Config{Prefix: "/app1"}),
			app2.Blueprint(app2.Config{Prefix: "/app2"}),
		},
	})

	cfg := raizu.NewServerConfigDefault().
		WithBlueprint(bp)

	s, err := raizu.NewServer(cfg)
	if err != nil {
		return err
	}

	err = raizu.ListenAndServe(s, ":8000")
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func main() {
	if err := mainE(); err != nil {
		log.Fatal(err)
	}
}
