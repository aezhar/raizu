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
	// Mounter is anything a [http.Handler] can be mounted on.
	Mounter interface {
		Mount(prefix string, h http.Handler)
	}

	// MounterHandler represents a [http.Handler] that can mount
	// other [http.Handler] instances.
	MounterHandler interface {
		http.Handler
		Mounter
	}

	AppMuxConfig[MT MounterHandler] struct {
		Prefix     string
		NewMuxFn   func() MT
		Blueprints []Blueprint
	}
	AppMux[MT MounterHandler] struct {
		cfg AppMuxConfig[MT]
		mux http.Handler
	}
)

func (ag *AppMux[MT]) Prefix() string        { return ag.cfg.Prefix }
func (ag *AppMux[MT]) Handler() http.Handler { return ag.mux }

func NewAppMux[MT MounterHandler](cfg AppMuxConfig[MT]) (*AppMux[MT], error) {
	mux := cfg.NewMuxFn()
	for _, af := range cfg.Blueprints {
		hp, err := af.NewHandlerProvider()
		if err != nil {
			return nil, err
		}
		mux.Mount(hp.Prefix(), hp.Handler())
	}

	h := http.Handler(mux)
	if cfg.Prefix != "/" {
		h = http.StripPrefix(cfg.Prefix, h)
	}

	return &AppMux[MT]{cfg: cfg, mux: h}, nil
}

func AppMuxBlueprint[MT MounterHandler](cfg AppMuxConfig[MT]) Blueprint {
	return NewBlueprint(NewAppMux[MT], cfg)
}
