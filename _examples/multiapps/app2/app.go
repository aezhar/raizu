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

package app2

import (
	"net/http"

	"github.com/aezhar/raizu"
)

type (
	Config struct {
		Prefix string
	}

	App struct {
		cfg Config
	}
)

func (a *App) Prefix() string {
	return a.cfg.Prefix
}

func (a *App) handleIndex(w http.ResponseWriter, r *http.Request) error {
	return raizu.WriteOKMessage(w, "Hello from App2")
}

func (a *App) Handler() http.Handler {
	m := http.NewServeMux()
	m.Handle("/", raizu.Wrap(a.handleIndex))
	return m
}

func New(cfg Config) (*App, error) {
	return &App{cfg: cfg}, nil
}

func Blueprint(cfg Config) raizu.Blueprint {
	return raizu.NewBlueprint(New, cfg)
}
