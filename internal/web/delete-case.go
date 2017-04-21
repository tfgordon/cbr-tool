// Copyright © 2015 Thomas F. Gordon
// This Source Code Form is subject to the terms of the
// Mozilla Public License, v. 2.0. If a copy of the MPL
// was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.

package web

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
)

func makeDeleteCaseHandler(cdb CouchDBConfig, tc TemplatesConfig) func(http.ResponseWriter, *http.Request) {
	couchdbURL := cdb.URL
	dbName := cdb.Cases
	errorTemplate := tc.errorTemplate
	return func(w http.ResponseWriter, req *http.Request) {

		caseId := path.Base(req.URL.Path)

		url := couchdbURL + dbName + "/" + caseId
		resp, err := http.Get(url)
		if err != nil {
			errorTemplate.Execute(w, err.Error())
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errorTemplate.Execute(w, err.Error())
			return
		}

		data := make(map[string]interface{})
		err = json.Unmarshal(body, &data)

		rev := data["_rev"]
		req2, err := http.NewRequest("DELETE", couchdbURL+dbName+"/"+caseId+"?rev="+rev.(string), nil)
		if err != nil {
			errorTemplate.Execute(w, err.Error())
			return
		}

		client := http.Client{}
		resp, err = client.Do(req2)
		if err != nil {
			errorTemplate.Execute(w, err.Error())
			return
		}
		defer resp.Body.Close()

		http.Redirect(w, req, "/eagle-argumentation-tool/", http.StatusSeeOther)
	}
}
