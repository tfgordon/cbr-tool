// Copyright © 2016 Thomas F. Gordon
// This Source Code Form is subject to the terms of the
// Mozilla Public License, v. 2.0. If a copy of the MPL
// was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.

package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"git.list.lu/eagle/argumentation-tool/internal/model"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

type updateDomainResponse struct {
	Revision string `json:"rev"` // Revision id of the updated model
	Message  string `json:"message"`
}

func makeUpdateDomainHandler(cdb CouchDBConfig) func(http.ResponseWriter, *http.Request) {
	couchdbURL := cdb.URL
	domainsDBName := cdb.Domains

	return func(w http.ResponseWriter, req *http.Request) {
		values := req.URL.Query()
		apikey := values.Get("apikey")
		user := values.Get("user")
		domainId := path.Base(req.URL.Path)

		// return a status and JSON encoded message to the client,
		// including the id of the newly created domain, if successful
		respond := func(status int, id string, msg string) {
			response := createDomainResponse{
				Id:      id,
				Message: msg,
			}
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(response)
		}

		if apikey != cdb.APIKey {
			respond(http.StatusUnauthorized, domainId, "Incorrect API Key")
			return
		}

		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			respond(http.StatusBadRequest, domainId, "Body of request could not be read")
			return
		}

		defer req.Body.Close()

		// To do: Validation of the body, to check that it is well-formed JSON
		// and correctly represents an domain model

		// Get the latest revision of the domain model

		domain, err := model.GetDomain(couchdbURL, domainsDBName, domainId)
		if err != nil {
			respond(http.StatusInternalServerError, domainId, "Failed to retrieve a domain model with the given Id from the Couchdb server.")
			return
		}

		fmt.Printf("Domain: Description=%v; Owner=%v\n", domain.Description, domain.Owner)

		// Check that the user owns the domain model
		if user != domain.Owner {
			respond(http.StatusUnauthorized, domainId, "User is not the owner (creator) of the domain model.")
			fmt.Printf("user=%v; owner=%v\n", user, domain.Owner)
			return
		}

		// add a field for the revision number to the JSON representation of
		// the argument graph

		json1 := strings.Replace(string(body), "{", "{\"_rev\":\""+domain.Revision+"\",", 1)
		// Add a property recording the user modifying the domain model
		// to the JSON file
		json1 = strings.Replace(json1, "{", "{\"Owner\":\""+user+"\",", 1)

		// Construct the PUT request
		req2, err := http.NewRequest("PUT", couchdbURL+domainsDBName+"/"+domainId, bytes.NewBufferString(json1))
		if err != nil {
			respond(http.StatusInternalServerError, domainId, "Failed to create the PUT request to the Couchdb server.")
			return
		}

		// upload the revised JSON document to Couchdb

		client := http.Client{}
		resp2, err := client.Do(req2)
		if err != nil {
			respond(http.StatusInternalServerError, "", "The Couchdb database failed to handle the PUT request.")
			return
		}

		defer resp2.Body.Close()

		// check that uploading was successful

		m := make(map[string]interface{})
		v, err := ioutil.ReadAll(resp2.Body)
		if err != nil {
			respond(http.StatusInternalServerError, "", "Unable to read the response from the Couchdb database.")
			return
		}

		err = json.Unmarshal(v, &m)
		if err != nil {
			respond(http.StatusInternalServerError, "", "Failed to unmarshall the response from the Couchdb database.")
			return
		}

		// m should have the form:
		// {"ok":true,"id":"6e1295ed6c29495e54cc05947f18c8af","rev":"1-2902191555"}

		if m["ok"] != true {
			respond(http.StatusInternalServerError, "", "Failed to upload the revised domain model to the Couchdb database.")
			return
		}

		respond(http.StatusOK, domainId, "Successfully uploaded the revised new domain model to the Couchdb database")
	}
}
