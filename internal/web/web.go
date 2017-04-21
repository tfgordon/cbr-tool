// Copyright © 2015 Thomas F. Gordon
// This Source Code Form is subject to the terms of the
// Mozilla Public License, v. 2.0. If a copy of the MPL
// was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.

package web

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

type templateHandler struct {
	once         sync.Once
	filename     string
	templatesDir string
	templ        *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join(t.templatesDir, t.filename)))
	})
	t.templ.Execute(w, nil)
}

type CouchDBConfig struct {
	URL     string
	Domains string
	Cases   string
	APIKey  string
}

type TemplatesConfig struct {
	templatesDir  string
	errorTemplate *template.Template
}

func ArgumentationToolServer(httpPort string, cdb CouchDBConfig, templatesDir string) {
	errorTemplate := template.Must(template.ParseFiles(filepath.Join(templatesDir, "error.html")))
	tc := TemplatesConfig{templatesDir, errorTemplate}

	root := "/eagle-argumentation-tool"

	// Domain Model Commands
	http.HandleFunc(root+"/", makeListDomainsHandler(cdb, tc))
	http.HandleFunc(root+"/view-domain/", makeViewDomainHandler(cdb, tc))
	http.HandleFunc(root+"/search-form/", makeSearchFormHandler(cdb, tc))
	http.HandleFunc(root+"/search", makeSearchHandler(cdb, tc))

	// Domain Model CRUD Web Service API
	http.HandleFunc(root+"/create-domain", makeCreateDomainHandler(cdb))
	http.HandleFunc(root+"/read-domain/", makeReadDomainHandler(cdb))
	http.HandleFunc(root+"/update-domain/", makeUpdateDomainHandler(cdb))
	http.HandleFunc(root+"/delete-domain/", makeDeleteDomainHandler(cdb))

	// Case Commands
	http.HandleFunc(root+"/view-case/", makeViewCaseHandler(cdb, tc))
	http.HandleFunc(root+"/new-case-form/", makeNewCaseFormHandler(cdb, tc))
	http.HandleFunc(root+"/new-case", makeNewCaseHandler(cdb, tc))
	http.HandleFunc(root+"/map-case/", makeMapCaseHandler(cdb, tc))
	http.HandleFunc(root+"/delete-case/", makeDeleteCaseHandler(cdb, tc))
	http.HandleFunc(root+"/edit-case/", makeEditCaseFormHandler(cdb, tc))
	http.HandleFunc(root+"/save-case", makeSaveCaseHandler(cdb, tc))

	// Other Commands
	http.Handle(root+"/help", &templateHandler{filename: "help.html", templatesDir: templatesDir})
	http.Handle(root+"/imprint", &templateHandler{filename: "imprint.html", templatesDir: templatesDir})

	// Start the web server
	if err := http.ListenAndServe(":"+httpPort, nil); err != nil {
		log.Fatal("EAGLE argumentation tool: ", err)
	}
}
