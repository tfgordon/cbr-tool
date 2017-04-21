// Copyright © 2015 Thomas F. Gordon
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
	"sort"
)

type SimilarCase struct {
	Id         string
	Title      string
	Year       string // year decided
	Similarity float64
}

/*
The Cases of a CaseTable map the id of an option to cases
which were decided in favor of this option, ordered by similarity
with the current case.
*/
type CaseTable struct {
	DomainId string
	Options  map[string][]SimilarCase
}

// to satisfy the sorting interface
type BySimilarity []SimilarCase

func (a BySimilarity) Len() int      { return len(a) }
func (a BySimilarity) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a BySimilarity) Less(i, j int) bool {
	return a[i].Similarity >= a[j].Similarity
}

func makeSearchHandler(cdb CouchDBConfig, tc TemplatesConfig) func(http.ResponseWriter, *http.Request) {
	couchdbURL := cdb.URL
	domainsDBName := cdb.Domains
	casesDBName := cdb.Cases
	errorTemplate := tc.errorTemplate
	templatesDir := tc.templatesDir
	return func(w http.ResponseWriter, req *http.Request) {
		// this is the only reference to the animals model,
		// making it easier to generalize this code to all domain models later

		domainId := req.FormValue("domainId")
		domain, err := model.GetDomain(couchdbURL, domainsDBName, domainId)
		if err != nil {
			errorTemplate.Execute(w, err.Error())
			return
		}

		var currentCase = make(model.CaseFactors)

		// read the factors of the current case, as selected in the search form
		for dimension, _ := range domain.Dimensions {
			currentCase[dimension] = req.FormValue(dimension)
		}

		resp, err := http.Get(couchdbURL + casesDBName + "/_design/cases/_view/statement")
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
				Id    string // case id
				Value struct {
					Title      string
					Year       string              // year of the decision
					Statements map[string]struct { // factor id
						Meta    interface{}
						Text    string
						Assumed bool
						Label   string
					}
				}
			}
		}{}

		err = json.Unmarshal(body, &data)

		if err != nil {
			errorTemplate.Execute(w, err.Error())
			return
		}

		// Compute case similarities

		// Table of simlilar cases, with one column for each option of the issue.
		// The keys are the statement ids of the options of the issue
		caseTable := CaseTable{
			DomainId: domainId,
			Options:  make(map[string][]SimilarCase),
		}

		// Initialize the table by adding a column for each position of the model
		for _, p := range domain.Options {
			caseTable.Options[p] = []SimilarCase{}
		}

		// Compare each case with the current case and append it to
		// the column of the case table for the position it supports.
		for _, c := range data.Rows {
			var precedent model.CaseFactors = make(model.CaseFactors)
			var decision string // id of the chosen position
			for sid, s := range c.Value.Statements {
				if s.Label == "in" {
					// assumption: only one position of the issue is in
					if domain.IsOption(sid) {
						decision = sid
					} else {
						dimension, ok := domain.FactorDimension(sid)
						// if the in statement is a factor of a dimension
						// assign it as the value of the dimension in
						// the precedent case description
						// assumption: no more than one factor of the
						// dimension is in
						if ok {
							precedent[dimension] = sid
						}
					}
				}
			}

			// Compare the similarity of the precedent to the current
			// case and record the result in the case table
			if decision != "" {
				v := domain.ArithmetricMeanSimilarity(currentCase, precedent)
				sc := SimilarCase{
					Id:         c.Id,
					Title:      c.Value.Title,
					Year:       c.Value.Year,
					Similarity: v,
				}
				caseTable.Options[decision] = append(caseTable.Options[decision], sc)
			}
		}

		// sort the columns of the case table, with the most similar cases
		// first (at the top)

		for _, cases := range caseTable.Options {
			sort.Sort(BySimilarity(cases))
		}

		t := template.Must(template.ParseFiles(filepath.Join(templatesDir, "search-results.html")))
		t.Execute(w, caseTable)
	}
}
