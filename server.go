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
	"context"
	"net"
	"net/http"
)

// A HTTPServer represents an HTTP server.
type HTTPServer interface {
	// Close immediately closes all active [net.Listener] instances and
	// any connections in state [http.StateNew], [http.StateActive], or [http.StateIdle]. For a
	// graceful shutdown, use Shutdown.
	//
	// See [http.Server.Close] for more details.
	Close() error

	// Shutdown gracefully shuts down the server without interrupting any
	// active connections.
	//
	// See [http.Server.Shutdown] for more details.
	Shutdown(ctx context.Context) error

	// Serve starts accepting incoming connections on the [net.Listener] in l, creating a
	// new service goroutine for each. The service goroutines read requests and
	// then is expected to call srv.Handler of the corresponding NewBlueprint to reply to them.
	//
	// See [http.Server.Serve] for more details.
	Serve(l net.Listener) error
}

// NewServer creates a new [raizu.HTTPServer] based on the [raizu.ServerConfig] provided.
//
// All Blueprints provided in [raizu.ServerConfig.Blueprints] will be instantiated and
// mounted on the server. If any [raizu.Blueprint] instance fails, the
// whole function fails and returns an error instead.
func NewServer(config ServerConfig) (HTTPServer, error) {
	h, err := config.Blueprint.NewHandlerProvider()
	if err != nil {
		return nil, err
	}

	s := &http.Server{
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
		TLSConfig:    config.TLSConfig,
		Handler:      h.Handler(),
	}
	return s, nil
}
