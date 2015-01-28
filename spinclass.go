// Copyright 2015 Matt Martz <matt@sivel.net>
// All Rights Reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"github.com/GeertJohan/go.rice"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sivel/spinclass/common"
	"github.com/sivel/spinclass/config"
	"github.com/sivel/spinclass/handler"
	"github.com/sivel/spinclass/spin"
)

var (
	Roster common.RosterType
)

func main() {
	conf := config.Config()

	Roster = make(common.RosterType)

	router := mux.NewRouter().StrictSlash(true)

	h := handler.Handler{
		Router: router,
		Config: conf,
		Roster: Roster,
		Class: spin.Class{
			Config: conf,
			Roster: Roster,
		},
	}

	h.Class.Create()

	loggingHandler := handlers.CombinedLoggingHandler(os.Stdout, router)

	box := rice.MustFindBox("static")
	router.HandleFunc("/spin/up", h.SpinUp).Methods("POST")
	router.HandleFunc("/spin/down", h.SpinDown).Methods("POST")
	router.HandleFunc("/odometer", h.Odometer).Methods("POST")
	router.HandleFunc("/final", h.Final).Methods("POST")
	router.PathPrefix("/").Handler(http.FileServer(box.HTTPBox()))

	http.Handle("/", loggingHandler)

	server := &http.Server{
		Addr:    conf.Server.Port,
		Handler: loggingHandler,
	}
	if len(conf.Server.Cert) == 0 || len(conf.Server.Key) == 0 {
		log.Fatal(server.ListenAndServe())
	} else {
		tlsConfig := &tls.Config{MinVersion: tls.VersionTLS10}
		server.TLSConfig = tlsConfig
		log.Fatal(server.ListenAndServeTLS(conf.Server.Cert, conf.Server.Key))
	}
}
