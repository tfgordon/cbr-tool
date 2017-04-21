// Copyright © 2015 Thomas F. Gordon
// This Source Code Form is subject to the terms of the
// Mozilla Public License, v. 2.0. If a copy of the MPL
// was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.

package web

import (
	"bytes"
	"github.com/carneades/carneades-4/src/engine/caes"
	"github.com/carneades/carneades-4/src/engine/caes/encoding/dot"
	"github.com/carneades/carneades-4/src/engine/caes/encoding/json"
	"io/ioutil"
	"net/http"
	"os/exec"
	"path"
)

func makeMapCaseHandler(cdb CouchDBConfig, tc TemplatesConfig) func(http.ResponseWriter, *http.Request) {
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

		var cag *caes.ArgGraph
		rd := bytes.NewReader(body)

		cag, err = json.Import(rd)
		if err != nil {
			errorTemplate.Execute(w, err.Error())
			return
		}

		cmd := exec.Command("dot", "-Tsvg")
		w2 := bytes.NewBuffer([]byte{})

		cmd.Stdin = w2
		cmd.Stdout = w

		err = dot.Export(w2, *cag)
		if err != nil {
			errorTemplate.Execute(w, err.Error())
			return
		}

		w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
		// w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		err = cmd.Run()
		if err != nil {
			errorTemplate.Execute(w, err.Error())
			return
		}
	}
}
