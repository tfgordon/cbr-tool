// Copyright © 2016 Thomas F. Gordon
// This Source Code Form is subject to the terms of the
// Mozilla Public License, v. 2.0. If a copy of the MPL
// was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.

package web

import (
	"encoding/json"
	// "fmt"
	"git.list.lu/eagle/argumentation-tool/internal/model"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func makeListDomainsHandler(cdb CouchDBConfig, tc TemplatesConfig) func(http.ResponseWriter, *http.Request) {
	couchdbURL := cdb.URL
	domainsDBName := cdb.Domains
	errorTemplate := tc.errorTemplate
	templatesDir := tc.templatesDir
	return func(w http.ResponseWriter, req *http.Request) {
		// Retrieve the list of domain models from CouchDB
		resp, err := http.Get(couchdbURL + "/" + domainsDBName + "/_all_docs")
		if err != nil {
			errorTemplate.Execute(w, err.Error())
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errorTemplate.Execute(w, err.Error())
			return
		}

		data := struct {
			Rows []struct {
				Id    string // domain model id in CouchB
				Value struct {
					Rev string
				}
			}
		}{}

		err = json.Unmarshal(body, &data)

		if err != nil {
			errorTemplate.Execute(w, err.Error())
			return
		}

		domains := []model.Domain{}
		for _, row := range data.Rows {
			// Retrieve the domain model from Couchd
			resp, err := http.Get(couchdbURL + "/" + domainsDBName + "/" + row.Id)
			if err != nil {
				errorTemplate.Execute(w, err.Error())
				return
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				errorTemplate.Execute(w, err.Error())
				return
			}
			domain := model.Domain{}
			err = json.Unmarshal(body, &domain)

			if err != nil {
				errorTemplate.Execute(w, err.Error())
				return
			}
			domain.Id = row.Id
			domains = append(domains, domain)
		}

		t := template.Must(template.ParseFiles(filepath.Join(templatesDir, "list-domains.html")))
		t.Execute(w, domains)
	}
}
