// Gow: A Lightweight GO Web Framework.
// Copyright 2016 The Gow Author. All Rights Reserved.
//
// This file is provided to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file
// except in compliance with the License.  You may obtain
// a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
// -------------------------------------------------------------------

package gow

import (
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

type Config struct {
	BaseUrl      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type GowServer struct {
	config      *Config
	handlers    map[string]http.Handler
	controllers map[string]interface{}
}

func New() *GowServer {

	return &GowServer{
		handlers:    make(map[string]http.Handler),
		controllers: make(map[string]interface{}),
	}
}

func (gw *GowServer) SetBaseUrl(u string) {
	gw.config.BaseUrl = "/" + strings.Trim(u, "/")
}

func (gw *GowServer) SetConfig(cfg *Config) {
	gw.config = cfg
}

func (gw *GowServer) Run(addr string) error {

	serveMux := http.NewServeMux()

	for hs, hd := range gw.handlers {
		serveMux.Handle(hs, hd)
	}

	serveMux.HandleFunc("/", gw.handler)

	if len(gw.config.BaseUrl) < 1 {
		gw.config.BaseUrl = "/"
	}

	if gw.config.ReadTimeout < 1 {
		gw.config.ReadTimeout = 3 * time.Second
	}

	if gw.config.WriteTimeout < 1 {
		gw.config.WriteTimeout = 10 * time.Second
	}

	s := &http.Server{
		Addr:           addr,
		Handler:        serveMux,
		ReadTimeout:    gw.config.ReadTimeout,
		WriteTimeout:   gw.config.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	return s.ListenAndServe()
}

func (gw *GowServer) RegisterHandler(r string, h http.Handler) {
	gw.handlers[r] = h
}

func (gw *GowServer) RegisterController(r string, c interface{}) {
	r = "/" + strings.Trim(r, "/")
	gw.controllers[r] = c
}

func (gw *GowServer) RegisterStaticRoute(r, path string) {

	if r[len(r)-1] != '/' {
		r = r + "/"
	}

	gw.RegisterHandler(r, http.StripPrefix(r, http.FileServer(http.Dir(path))))
}

func (gw *GowServer) handler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/favicon.ico" {
		return
	}

	var (
		path = r.URL.Path[len(gw.config.BaseUrl):]
		spi  = strings.LastIndex(path, "/")
	)

	if controller, ok := gw.controllers[path[:spi]]; ok {

		value_c := reflect.ValueOf(controller)

		value_c.Elem().FieldByName("Request").Set(reflect.ValueOf(r))
		value_c.Elem().FieldByName("Params").Set(reflect.ValueOf(r.URL.Query()))
		value_c.Elem().FieldByName("Response").FieldByName("Response").Set(reflect.ValueOf(w))

		if exec_method := value_c.MethodByName(strings.Replace(strings.Title(path[spi+1:]), "-", "", -1) + "Action"); exec_method.Kind() == reflect.Func {
			exec_method.Call([]reflect.Value{})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Not Found"))
}

type Controller struct {
	Request  *http.Request
	Response ResponseWriter
	Params   url.Values
}

type ResponseWriter struct {
	Response http.ResponseWriter
}

func (c ResponseWriter) SetHeader(key, value string) {
	c.Response.Header().Set(key, value)
}

func (c ResponseWriter) WriteHeader(h int) {
	c.Response.WriteHeader(h)
}

func (c ResponseWriter) RenderBytes(b []byte) {
	c.Response.Write(b)
}

func (c ResponseWriter) RenderString(s string) {
	c.Response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Response.Write([]byte(s))
}

func (c ResponseWriter) RenderJson(v interface{}) {
	c.Response.Header().Set("Content-Type", "application/json;charset=utf-8")

	if jrsp, err := json.Marshal(&v); err == nil {
		c.Response.Write([]byte(jrsp))
	}
}
