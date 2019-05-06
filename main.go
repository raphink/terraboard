package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/camptocamp/terraboard/api"
	"github.com/camptocamp/terraboard/auth"
	"github.com/camptocamp/terraboard/config"
	"github.com/camptocamp/terraboard/util"
	"github.com/camptocamp/terradb/pkg/client"
	log "github.com/sirupsen/logrus"
)

// idx serves index.html, always,
// so as to let AngularJS manage the app routing.
// The <base> HTML tag is edited on the fly
// to reflect the proper base URL
func idx(w http.ResponseWriter, r *http.Request) {
	idx, err := ioutil.ReadFile("static/index.html")
	if err != nil {
		log.Errorf("Failed to open index.html: %v", err)
		// TODO: Return error page
	}
	idxStr := string(idx)
	idxStr = util.ReplaceBase(idxStr, "base href=\"/\"", "base href=\"%s\"")
	io.WriteString(w, idxStr)
}

// Pass the DB to API handlers
// This takes a callback and returns a HandlerFunc
// which calls the callback with the DB
func handleWithDB(apiF func(w http.ResponseWriter, r *http.Request, d *client.Client), d *client.Client) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiF(w, r, d)
	})
}

var version = "undefined"

func getVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	j, err := json.Marshal(map[string]string{
		"version":   version,
		"copyright": "Copyright © 2017 Camptocamp",
	})
	if err != nil {
		api.JSONError(w, "Failed to marshal version", err)
		return
	}
	io.WriteString(w, string(j))
}

// Main
func main() {
	c := config.LoadConfig(version)

	log.Infof("Terraboard v%s is starting...", version)

	err := c.SetupLogging()
	if err != nil {
		log.Fatal(err)
	}

	terradb := client.NewClient(c.TerraDB.URL)

	// Set up auth
	auth.Setup(c)

	// Index is a wildcard for all paths
	http.HandleFunc(util.AddBase(""), idx)

	// Serve static files (CSS, JS, images) from dir
	staticFs := http.FileServer(http.Dir("static"))
	http.Handle(util.AddBase("static/"), http.StripPrefix(util.AddBase("static"), staticFs))

	// Handle API points
	http.HandleFunc(util.AddBase("api/version"), getVersion)
	http.HandleFunc(util.AddBase("api/user"), api.GetUser)
	http.HandleFunc(util.AddBase("api/states"), handleWithDB(api.ListStates, terradb))
	//	http.HandleFunc(util.AddBase("api/states/stats"), handleWithDB(api.ListStateStats, terradb))
	//http.HandleFunc(util.AddBase("api/states/tfversion/count"), handleWithDB(api.ListTerraformVersionsWithCount, terradb))
	http.HandleFunc(util.AddBase("api/state/"), handleWithDB(api.GetState, terradb))
	http.HandleFunc(util.AddBase("api/state/activity/"), handleWithDB(api.GetStateActivity, terradb))
	http.HandleFunc(util.AddBase("api/state/compare/"), handleWithDB(api.StateCompare, terradb))
	http.HandleFunc(util.AddBase("api/locks"), handleWithDB(api.GetLocks, terradb))
	//http.HandleFunc(util.AddBase("api/search/attribute"), handleWithDB(api.SearchAttribute, terradb))
	//http.HandleFunc(util.AddBase("api/resource/types"), handleWithDB(api.ListResourceTypes, terradb))
	//http.HandleFunc(util.AddBase("api/resource/types/count"), handleWithDB(api.ListResourceTypesWithCount, terradb))
	//http.HandleFunc(util.AddBase("api/resource/names"), handleWithDB(api.ListResourceNames, terradb))
	//http.HandleFunc(util.AddBase("api/attribute/keys"), handleWithDB(api.ListAttributeKeys, terradb))
	//http.HandleFunc(util.AddBase("api/tf_versions"), handleWithDB(api.ListTfVersions, terradb))

	// Start server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", c.Port), nil))
}
