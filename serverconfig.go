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
	"crypto/tls"
	"time"

	"github.com/go-logr/logr"
)

// ServerConfig describes the Server runtime configuration.
type ServerConfig struct {
	// AccessLog will receive http access requests logs as Info
	// level logging and any error encountered during the handling
	// of http access requests.
	Logger logr.Logger

	// TLSConfig
	TLSConfig *tls.Config

	// ReadTimeout see [http.Server.ReadTimeout]
	ReadTimeout time.Duration

	// ReadTimeout see [http.Server.ReadHeaderTimeout]
	ReadHeaderTimeout time.Duration

	// ReadTimeout see [http.Server.WriteTimeout]
	WriteTimeout time.Duration

	// IdleTimeout see [http.Server.IdleTimeout]
	IdleTimeout time.Duration

	// Blueprint is a factory that produces an app which will
	// be served via the http server.
	Blueprint Blueprint
}

func (c ServerConfig) modify(fn func(c *ServerConfig)) ServerConfig {
	fn(&c)
	return c
}

func (c ServerConfig) WithLogger(logger logr.Logger) ServerConfig {
	return c.modify(func(c *ServerConfig) { c.Logger = logger })
}

func (c ServerConfig) WithTLSConfig(cfg *tls.Config) ServerConfig {
	return c.modify(func(c *ServerConfig) { c.TLSConfig = cfg })
}

func (c ServerConfig) WithReadTimeout(v time.Duration) ServerConfig {
	return c.modify(func(c *ServerConfig) { c.ReadTimeout = v })
}

func (c ServerConfig) WithReadHeaderTimeout(v time.Duration) ServerConfig {
	return c.modify(func(c *ServerConfig) { c.ReadHeaderTimeout = v })
}

func (c ServerConfig) WithWriteTimeout(v time.Duration) ServerConfig {
	return c.modify(func(c *ServerConfig) { c.WriteTimeout = v })
}

func (c ServerConfig) WithIdleTimeout(v time.Duration) ServerConfig {
	return c.modify(func(c *ServerConfig) { c.IdleTimeout = v })
}

func (c ServerConfig) WithBlueprint(af Blueprint) ServerConfig {
	return c.modify(func(c *ServerConfig) { c.Blueprint = af })
}

func NewServerConfigDefault() ServerConfig {
	return ServerConfig{
		Logger: logr.Discard(),

		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
}
