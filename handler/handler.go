package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
	"github.com/mholt/binding"
	"github.com/sivel/spinclass/common"
	"github.com/sivel/spinclass/forms"
	"github.com/sivel/spinclass/spin"
	"github.com/sivel/spinclass/utils"
)

type Handler struct {
	Router    *mux.Router
	Config    common.Config
	Templates *rice.Box
	Class     spin.Class
	Roster    common.RosterType
}

func (h *Handler) SpinUp(w http.ResponseWriter, req *http.Request) {
	spinUp := common.SpinUp{
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

	odometer := common.Odometer{
		Instances: h.Roster[prefixForm.Prefix],
	}
	output, _ := json.Marshal(odometer)
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

func (h *Handler) Final(w http.ResponseWriter, req *http.Request) {
	prefixForm := new(forms.PrefixForm)
	binding.Bind(req, prefixForm)
	delete(h.Roster, prefixForm.Prefix)
}
