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
	"time"
)

// Muxer is anything that can register HTTP handlers.
type Muxer interface {
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

// Mountable is anything that can be mounted onto a Server.
type Mountable interface {
	Mount(m Muxer)
}

// The MountableFunc type is an adapter to allow the use of
// ordinary functions as a Mountable. If f is a function
// with the appropriate signature, MountableFunc(f) is a
// Mountable that calls f.
type MountableFunc func(m Muxer)

func (f MountableFunc) Mount(m Muxer) { f(m) }

// Server represents a http server with a combined multiplexer.
type Server struct {
	srv *http.Server
	mux *http.ServeMux
}

func (s *Server) Close() error {
	return s.srv.Close()
}

// Mount mounts the given Mountable resource to the multiplexer.
func (s *Server) Mount(m Mountable) { m.Mount(s.mux) }

// Serve works just like net/http.Server.Serve and accepts incoming
// connections on the Listener l, creating a new service goroutine
// for each. The service goroutines read requests and then look for a
// handler depending on the referenced resource.
//
// Serve always returns a non-nil error and closes l.
//
// After Shutdown or Close, the returned error is ErrServerClosed.
//
// Serve will block until
func (s *Server) Serve(l net.Listener) error {
	return s.srv.Serve(l)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func NewServerOpt(s *http.Server) *Server {
	m := http.NewServeMux()

	s.Handler = m

	return &Server{srv: s, mux: m}
}

func NewServer() *Server {
	return NewServerOpt(&http.Server{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
	})
}
