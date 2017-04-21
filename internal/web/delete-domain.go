// Copyright © 2016 Thomas F. Gordon
// This Source Code Form is subject to the terms of the
// Mozilla Public License, v. 2.0. If a copy of the MPL
// was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.

package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
)

func makeDeleteDomainHandler(cdb CouchDBConfig) func(http.ResponseWriter, *http.Request) {
	couchdbURL := cdb.URL
	domainsDBName := cdb.Domains

	return func(w http.ResponseWriter, req *http.Request) {
		values := req.URL.Query()
		apikey := values.Get("apikey")
		user := values.Get("user")
		domainId := path.Base(req.URL.Path)

		// return a status and JSON encoded message to the client,
		// including the id of the newly created domain, if successful
		respond := func(status int, body string) {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(body)
		}

		if apikey != cdb.APIKey {
			respond(http.StatusUnauthorized, "Incorrect API Key")
			return
		}

		url := couchdbURL + domainsDBName + "/" + domainId
		resp, err := http.Get(url)
		if err != nil {
			respond(http.StatusBadRequest, "No domain model with this id could be retrieved from the Couchdb server.")
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			respond(http.StatusInternalServerError, "Unable to read the body of the response from the Couchdb server.")
			return
		}

		data := make(map[string]interface{})
		err = json.Unmarshal(body, &data)
		rev := data["_rev"]
		owner := data["Owner"]

		if user != owner {
			respond(http.StatusUnauthorized, "User is not the owner (creator) of the domain model.")
			return
		}

		req2, err := http.NewRequest("DELETE", couchdbURL+domainsDBName+"/"+domainId+"?rev="+rev.(string), nil)
		if err != nil {
			respond(http.StatusInternalServerError, "Failed to construct the delete request.")
			return
		}

		client := http.Client{}
		resp, err = client.Do(req2)
		if err != nil {
			respond(http.StatusInternalServerError, "Failed to delete the domain model from the Couchdb database.")
			return
		}

		defer resp.Body.Close()
		respond(http.StatusOK, fmt.Sprintf("Revision %s of the domain model with the id %s successully deleted.", rev, domainId))
	}
}
