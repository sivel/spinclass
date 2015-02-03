package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sivel/spinclass/common"
	"github.com/sivel/spinclass/spin"
	"github.com/sivel/spinclass/utils"
)

type Handler struct {
	Router *mux.Router
	Config common.Config
	Class  spin.Class
	Roster common.RosterType
}

type CountStruct struct {
	Count int `json:"count"`
}

type PrefixStruct struct {
	Prefix string `json:"prefix"`
}

func (h *Handler) SpinUp(w http.ResponseWriter, req *http.Request) {
	spinUp := common.SpinUp{
		Prefix: fmt.Sprintf("spinclass-%s", utils.RandStr(8)),
	}

	decoder := json.NewDecoder(req.Body)
	decoder.UseNumber()
	var count CountStruct
	decoder.Decode(&count)

	serverIDs := h.Class.New(count.Count, spinUp.Prefix)
	spinUp.ServerIDs = serverIDs
	output, _ := json.Marshal(spinUp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

func (h *Handler) SpinDown(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var prefix PrefixStruct
	decoder.Decode(&prefix)
	h.Class.Dismiss(prefix.Prefix)
}

func (h *Handler) Odometer(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var prefix PrefixStruct
	decoder.Decode(&prefix)

	odometer := common.Odometer{
		Instances: h.Roster[prefix.Prefix],
	}
	output, _ := json.Marshal(odometer)
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

func (h *Handler) Final(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var prefix PrefixStruct
	decoder.Decode(&prefix)
	delete(h.Roster, prefix.Prefix)
}
