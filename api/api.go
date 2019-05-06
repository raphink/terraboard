package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/camptocamp/terraboard/auth"
	"github.com/camptocamp/terraboard/compare"
	"github.com/camptocamp/terraboard/types"
	"github.com/camptocamp/terraboard/util"
	"github.com/camptocamp/terradb/pkg/client"
)

var states []string

// JSONError is a wrapper function for errors
// which prints them to the http.ResponseWriter as a JSON response
func JSONError(w http.ResponseWriter, message string, err error) {
	errObj := make(map[string]string)
	errObj["error"] = message
	errObj["details"] = fmt.Sprintf("%v", err)
	j, _ := json.Marshal(errObj)
	io.WriteString(w, string(j))
}

// ListStates lists States
func ListStates(w http.ResponseWriter, r *http.Request, d *client.Client) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	states, err := d.ListStates()
	if err != nil {
		JSONError(w, "Failed to list states", err)
		return
	}

	j, err := json.Marshal(states)
	if err != nil {
		JSONError(w, "Failed to marshal states", err)
		return
	}
	io.WriteString(w, string(j))
}

// GetState provides information on a State
func GetState(w http.ResponseWriter, r *http.Request, d *client.Client) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	st := util.TrimBase(r, "api/state/")
	versionID := r.URL.Query().Get("versionid")
	var err error
	var serial int
	if versionID == "" {
		serial = 0
	} else {
		serial, err = strconv.Atoi(versionID)
		if err != nil {
			JSONError(w, "Failed to parse versionID", err)
			return
		}
	}
	state, err := d.GetState(st, serial)
	if err != nil {
		JSONError(w, "Failed to get state", err)
		return
	}

	jState, err := json.Marshal(state)
	if err != nil {
		JSONError(w, "Failed to marshal state", err)
		return
	}
	io.WriteString(w, string(jState))
}

// GetStateActivity returns the activity (version history) of a State
func GetStateActivity(w http.ResponseWriter, r *http.Request, d *client.Client) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	st := util.TrimBase(r, "api/state/activity/")
	activity, err := d.ListStateSerials(st)
	if err != nil {
		JSONError(w, "Failed to get state activity", err)
		return
	}

	jActivity, err := json.Marshal(activity)
	if err != nil {
		JSONError(w, "Failed to marshal state activity", err)
		return
	}
	io.WriteString(w, string(jActivity))
}

// StateCompare compares two versions ('from' and 'to') of a State
func StateCompare(w http.ResponseWriter, r *http.Request, d *client.Client) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	st := util.TrimBase(r, "api/state/compare/")
	query := r.URL.Query()
	fromVersion := query.Get("from")
	fromSerial, err := strconv.Atoi(fromVersion)
	if err != nil {
		JSONError(w, "Failed to parse from serial", err)
		return
	}
	toVersion := query.Get("to")
	toSerial, err := strconv.Atoi(toVersion)
	if err != nil {
		JSONError(w, "Failed to parse to serial", err)
		return
	}

	from, err := d.GetState(st, fromSerial)
	if err != nil {
		JSONError(w, "Failed to get from state", err)
		return
	}

	fromInt := types.State{
		Path: from.Name,
		Version: types.Version{
			VersionID:    fmt.Sprintf("%v", from.Serial),
			LastModified: from.LastModified,
		},
		TFVersion: from.TFVersion,
		Serial:    from.Serial,
	}

	to, err := d.GetState(st, toSerial)
	if err != nil {
		JSONError(w, "Failed to get to state", err)
		return
	}

	toInt := types.State{
		Path: to.Name,
		Version: types.Version{
			VersionID:    fmt.Sprintf("%v", to.Serial),
			LastModified: to.LastModified,
		},
		TFVersion: to.TFVersion,
		Serial:    to.Serial,
	}

	compare, err := compare.Compare(fromInt, toInt)
	if err != nil {
		JSONError(w, "Failed to compare state versions", err)
		return
	}

	jCompare, err := json.Marshal(compare)
	if err != nil {
		JSONError(w, "Failed to marshal state compare", err)
		return
	}
	io.WriteString(w, string(jCompare))
}

// GetLocks returns information on locked States
func GetLocks(w http.ResponseWriter, r *http.Request, d *client.Client) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	/*
			locks, err := s3.GetLocks()
			if err != nil {
				JSONError(w, "Failed to get locks", err)
				return
			}

			j, err := json.Marshal(locks)
			if err != nil {
				JSONError(w, "Failed to marshal locks", err)
				return
			}
		io.WriteString(w, string(j))
	*/
}

// GetUser returns information about the logged user
func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	name := r.Header.Get("X-Forwarded-User")
	email := r.Header.Get("X-Forwarded-Email")

	user := auth.UserInfo(name, email)

	j, err := json.Marshal(user)
	if err != nil {
		JSONError(w, "Failed to marshal user information", err)
		return
	}
	io.WriteString(w, string(j))
}
