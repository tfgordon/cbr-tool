// Copyright © 2016 Thomas F. Gordon
// This Source Code Form is subject to the terms of the
// Mozilla Public License, v. 2.0. If a copy of the MPL
// was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.

package web

import (
	"bytes"
	"encoding/json"
	//	"fmt"
	"github.com/pborman/uuid"
	"io/ioutil"
	"net/http"
	"strings"
)

type createDomainResponse struct {
	Id      string `json:"id"` // UUID of the newly created domain model
	Message string `json:"message"`
}

func makeCreateDomainHandler(cdb CouchDBConfig) func(http.ResponseWriter, *http.Request) {
	couchdbURL := cdb.URL
	domainsDBName := cdb.Domains

	return func(w http.ResponseWriter, req *http.Request) {
		apikey := req.URL.Query().Get("apikey")
		user := req.URL.Query().Get("user")

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
			respond(http.StatusUnauthorized, "", "Incorrect API Key")
			return
		}

		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			respond(http.StatusBadRequest, "", "Body could not be read")
			return
		}

		defer req.Body.Close()

		// To do: Validation of the body, to check that it is well-formed JSON
		// and correctly represents an domain model

		// Add a property recording the user creating the domain model
		// to the JSON file
		domainModel := strings.Replace(string(body), "{", "{\"Owner\":\""+user+"\",", 1)

		// upload the JSON document to Couchdb

		id := uuid.NewUUID().String()
		req2, err := http.NewRequest("PUT", couchdbURL+domainsDBName+"/"+id, bytes.NewBufferString(domainModel))
		if err != nil {
			respond(http.StatusInternalServerError, "", "Failed to create the PUT request to the Couchdb server.")
			return
		}

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
			respond(http.StatusInternalServerError, "", "Failed to upload the domain model to the Couchdb database.")
			return
		}

		respond(http.StatusOK, id, "Successfully uploaded the new domain model to the Couchdb database")
	}
}
