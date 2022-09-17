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
	"net/http"
)

type (
	HandlerAppConfig struct {
		Prefix  string
		Handler http.Handler
	}
	HandlerApp struct {
		config HandlerAppConfig
	}
)

func (a *HandlerApp) Prefix() string        { return a.config.Prefix }
func (a *HandlerApp) Handler() http.Handler { return a.config.Handler }

func NewHandlerApp(cfg HandlerAppConfig) (*HandlerApp, error) {
	return &HandlerApp{config: cfg}, nil
}

func HandlerAppBlueprint(cfg HandlerAppConfig) Blueprint {
	return NewBlueprint(NewHandlerApp, cfg)
}
