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

func makeEditCaseFormHandler(cdb CouchDBConfig, tc TemplatesConfig) func(http.ResponseWriter, *http.Request) {
	couchdbURL := cdb.URL
	domainsDBName := cdb.Domains
	casesDBName := cdb.Cases
	errorTemplate := tc.errorTemplate
	templatesDir := tc.templatesDir
	return func(w http.ResponseWriter, req *http.Request) {

		domainId := path.Base(path.Dir(req.URL.Path))
		caseId := path.Base(req.URL.Path)
		domain, err := model.GetDomain(couchdbURL, domainsDBName, domainId)
		if err != nil {
			errorTemplate.Execute(w, err.Error())
			return
		}

		case1, err := model.GetCase(couchdbURL, casesDBName, caseId)

		if err != nil {
			errorTemplate.Execute(w, err.Error())
			return
		}

		currentCase := map[string]interface{}{
			"DomainId":    domainId,
			"Id":          case1.Id,
			"Revision":    case1.Revision,
			"Title":       case1.Meta.Title,
			"Citation":    case1.Meta.Citation,
			"Court":       case1.Meta.Court,
			"Year":        case1.Meta.Year,
			"Description": case1.Meta.Description,
			"Keywords":    case1.Meta.Keywords,
			"Language":    case1.Meta.Language,
			"Majority":    case1.Meta.Majority,
			"Minority":    case1.Meta.Minority,
			"Options":     make(map[string]Option),
			"Dimensions":  make(map[string]Dimension),
		}

		// get the decision and dimensions of the current case
		for k, s := range case1.Statements {
			d, ok := domain.FactorDimension(k)

			// assumption: at most one factor of each dimension is in
			if ok && s.Label == "in" {
				dimensions := currentCase["Dimensions"].(map[string]Dimension)
				factors := []Factor{}
				for _, f := range domain.Dimensions[d].Factors {
					factors = append(factors, Factor{Id: f, Text: domain.FactorStatement(f), Selected: f == k})
				}
				dimensions[d] = Dimension{Id: d, Description: domain.Dimensions[d].Description, Factors: factors}
			}
			for _, p := range domain.Options {
				if k == p {
					// the statement is an option of the issue of the model
					in := false
					if s.Label == "in" {
						in = true
					}
					options := currentCase["Options"].(map[string]Option)
					options[p] = Option{Statement: s.Text, In: in}
				}
			}
		}

		t := template.Must(template.ParseFiles(filepath.Join(templatesDir, "case-editor.html")))
		t.Execute(w, currentCase)
	}
}
