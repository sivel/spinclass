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
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/GeertJohan/go.rice"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mholt/binding"
	"github.com/sivel/spinclass/common"
	"github.com/sivel/spinclass/forms"
	"github.com/sivel/spinclass/spin"
	"github.com/sivel/spinclass/utils"
	"gopkg.in/yaml.v2"
)

var (
	Roster common.RosterType
)

type SpinUp struct {
	Prefix    string   `json:"prefix"`
	ServerIDs []string `json:"server_ids"`
}

type Odometer struct {
	Instances map[string]interface{} `json:"instances"`
}

type Handler struct {
	Router    *mux.Router
	Config    common.Config
	Templates *rice.Box
	Class     spin.Class
}

func (h *Handler) Index(w http.ResponseWriter, req *http.Request) {
	templateString, err := h.Templates.String("index.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	tmplMessage, err := template.New("index").Parse(templateString)
	if err != nil {
		log.Fatal(err)
	}
	tmplMessage.Execute(w, nil)
}

func (h *Handler) SpinUp(w http.ResponseWriter, req *http.Request) {
	spinUp := SpinUp{
		Prefix: fmt.Sprintf("spinclass-%s", utils.RandStr(8)),
	}
	countForm := new(forms.CountForm)
	binding.Bind(req, countForm)
	serverIDs := h.Class.New(countForm.Count, spinUp.Prefix)
	spinUp.ServerIDs = serverIDs
	output, _ := json.Marshal(spinUp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

func (h *Handler) SpinDown(w http.ResponseWriter, req *http.Request) {
	prefixForm := new(forms.PrefixForm)
	binding.Bind(req, prefixForm)
	h.Class.Dismiss(prefixForm.Prefix)
}

func (h *Handler) Odometer(w http.ResponseWriter, req *http.Request) {
	prefixForm := new(forms.PrefixForm)
	binding.Bind(req, prefixForm)

	odometer := Odometer{
		Instances: Roster[prefixForm.Prefix],
	}
	output, _ := json.Marshal(odometer)
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

func (h *Handler) Final(w http.ResponseWriter, req *http.Request) {
	prefixForm := new(forms.PrefixForm)
	binding.Bind(req, prefixForm)
	delete(Roster, prefixForm.Prefix)
}

func ParseConfig() common.Config {
	var config common.Config
	configPath, err := filepath.Abs(".")
	configFile := filepath.Join(configPath, "spinclass.yaml")
	text, err := ioutil.ReadFile(configFile)
	if err == nil {
		yaml.Unmarshal(text, &config)
	}
	if config.Server.Port == "" {
		config.Server.Port = ":3000"
	}
	if config.OpenStack.IdentityEndpoint == "" {
		config.OpenStack.IdentityEndpoint = "https://identity.api.rackspacecloud.com/v2.0"
	}
	return config
}

func main() {
	config := ParseConfig()

	flag.StringVar(&config.Server.Port, "port", config.Server.Port, "HOST:PORT to listen on, HOST not required to listen on all addresses")
	flag.StringVar(&config.Server.Cert, "cert", config.Server.Cert, "SSL cert file path. This option with 'key' enables SSL communication")
	flag.StringVar(&config.Server.Key, "key", config.Server.Key, "SSL key file path. This option with 'cert' enables SSL communication")
	flag.StringVar(&config.OpenStack.IdentityEndpoint, "identity", config.OpenStack.IdentityEndpoint, "OpenStack Identity V2 Endpoint URL")
	flag.StringVar(&config.OpenStack.Username, "username", config.OpenStack.Username, "OpenStack username")
	flag.StringVar(&config.OpenStack.Password, "password", config.OpenStack.Password, "OpenStack password")
	flag.StringVar(&config.OpenStack.APIKey, "apikey", config.OpenStack.APIKey, "OpenStack API Key")
	flag.StringVar(&config.OpenStack.Region, "region", config.OpenStack.Region, "OpenStack Region")
	flag.StringVar(&config.OpenStack.ImageRef, "image", config.OpenStack.ImageRef, "OpenStack Image ID")
	flag.StringVar(&config.OpenStack.FlavorRef, "flavor", config.OpenStack.FlavorRef, "OpenStack Flavor ID")
	flag.Parse()

	if config.OpenStack.Username == "" || (config.OpenStack.Password == "" && config.OpenStack.APIKey == "") {
		log.Fatal("OpenStack Username and Password or APIKey are required")
	}
	if config.OpenStack.Region == "" || config.OpenStack.ImageRef == "" || config.OpenStack.FlavorRef == "" {
		log.Fatal("OpenStack Region, Image and Flavor are required")
	}

	Roster = make(common.RosterType)

	router := mux.NewRouter().StrictSlash(true)

	h := Handler{
		Router: router,
		Config: config,
		Class: spin.Class{
			Config: config,
			Roster: Roster,
		},
	}

	h.Class.Create()

	templateBox, err := rice.FindBox("templates")
	if err != nil {
		log.Fatal(err)
	}

	h.Templates = templateBox

	loggingHandler := handlers.CombinedLoggingHandler(os.Stdout, router)

	router.HandleFunc("/", h.Index)
	router.HandleFunc("/spin/up", h.SpinUp).Methods("POST")
	router.HandleFunc("/spin/down", h.SpinDown).Methods("POST")
	router.HandleFunc("/odometer", h.Odometer).Methods("POST")
	router.HandleFunc("/final", h.Final).Methods("POST")

	http.Handle("/", loggingHandler)

	server := &http.Server{
		Addr:    config.Server.Port,
		Handler: loggingHandler,
	}
	if len(config.Server.Cert) == 0 || len(config.Server.Key) == 0 {
		log.Fatal(server.ListenAndServe())
	} else {
		tlsConfig := &tls.Config{MinVersion: tls.VersionTLS10}
		server.TLSConfig = tlsConfig
		log.Fatal(server.ListenAndServeTLS(config.Server.Cert, config.Server.Key))
	}
}
