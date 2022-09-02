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

type HandlerProvider interface {
	Prefix() string
	Handler() http.Handler
}

type BlueprintImpl[HPT HandlerProvider, CT any] struct {
	makeFn func(cf CT) (HPT, error)
	config CT
}

func (bp BlueprintImpl[HPT, CT]) NewHandlerProvider() (HandlerProvider, error) {
	return bp.makeFn(bp.config)
}

type Blueprint interface {
	NewHandlerProvider() (HandlerProvider, error)
}

func NewBlueprint[HPT HandlerProvider, CT any](makeFn func(cf CT) (HPT, error), cf CT) Blueprint {
	return BlueprintImpl[HPT, CT]{makeFn, cf}
}
