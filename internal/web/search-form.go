// Copyright © 2016 Thomas F. Gordon
// This Source Code Form is subject to the terms of the
// Mozilla Public License, v. 2.0. If a copy of the MPL
// was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.

package web

import (
	// "fmt"
	"git.list.lu/eagle/argumentation-tool/internal/model"
	"html/template"
	"net/http"
	"path"
	"path/filepath"
)

func makeSearchFormHandler(cdb CouchDBConfig, tc TemplatesConfig) func(http.ResponseWriter, *http.Request) {
	couchdbURL := cdb.URL
	domainsDBName := cdb.Domains
	errorTemplate := tc.errorTemplate
	templatesDir := tc.templatesDir
	return func(w http.ResponseWriter, req *http.Request) {

		domainId := path.Base(req.URL.Path)

		domain, err := model.GetDomain(couchdbURL, domainsDBName, domainId)
		if err != nil {
			errorTemplate.Execute(w, err.Error())
			return
		}

		t := template.Must(template.ParseFiles(filepath.Join(templatesDir, "search-form.html")))
		t.Execute(w, domain)
	}
}
