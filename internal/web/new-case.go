// Copyright © 2015 Thomas F. Gordon
// This Source Code Form is subject to the terms of the
// Mozilla Public License, v. 2.0. If a copy of the MPL
// was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.

package web

import (
	"bytes"
	"encoding/json"
	"git.list.lu/eagle/argumentation-tool/internal/model"
	"github.com/carneades/carneades-4/src/engine/caes"
	// "github.com/carneades/carneades-4/src/engine/caes/encoding/dot"
	cj "github.com/carneades/carneades-4/src/engine/caes/encoding/json"
	"github.com/pborman/uuid"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
)

func makeNewCaseHandler(cdb CouchDBConfig, tc TemplatesConfig) func(http.ResponseWriter, *http.Request) {
	couchdbURL := cdb.URL
	domainsDBName := cdb.Domains
	casesDBName := cdb.Cases
	errorTemplate := tc.errorTemplate
	templatesDir := tc.templatesDir
	newCaseCreatedTemplate := template.Must(template.ParseFiles(filepath.Join(templatesDir, "new-case-created.html")))

	return func(w http.ResponseWriter, req *http.Request) {
		ag := caes.NewArgGraph()
		// ag.Metadata["domain"] = req.FormValue("domainId")
		ag.Metadata["title"] = req.FormValue("title")
		ag.Metadata["citation"] = req.FormValue("citation")
		ag.Metadata["year"] = req.FormValue("year")
		ag.Metadata["court"] = req.FormValue("court")
		ag.Metadata["majority"] = req.FormValue("majority")
		ag.Metadata["minority"] = req.FormValue("minority")
		ag.Metadata["keywords"] = req.FormValue("keywords")
		ag.Metadata["language"] = req.FormValue("language")
		ag.Metadata["description"] = req.FormValue("description")

		domainId := req.FormValue("domainId")

		domain, err := model.GetDomain(couchdbURL, domainsDBName, domainId)
		if err != nil {
			errorTemplate.Execute(w, err.Error())
			return
		}

		decision := req.FormValue("decision")

		// add the premises to the statement table
		// and assume them to be true
		for dimension, _ := range domain.Dimensions {
			factor := req.FormValue(dimension)
			ag.Statements[factor] =
				&caes.Statement{Id: factor, Text: domain.FactorStatement(factor),
					Label: caes.In}
			ag.Assumptions[factor] = true
		}

		// take note of the premises, which are all the statements
		// added thus far to the statements table
		premises := []caes.Premise{}
		for _, s := range ag.Statements {
			premises = append(premises, caes.Premise{Stmt: s})
		}

		// construct the issue

		i1 := caes.NewIssue()
		i1.Id = "i1"
		ag.Issues[i1.Id] = i1

		// add the conclusions of the arguments, constructed next below,
		// to the statements table

		for i, option := range domain.Options {
			label := caes.Out
			if option == decision {
				label = caes.In
			}
			conclusion := caes.Statement{
				Id:    option,
				Text:  domain.FactorStatement(option),
				Label: label,
				Args:  []*caes.Argument{},
			}
			ag.Statements[option] = &conclusion
			arg := caes.NewArgument()
			arg.Id = "a" + strconv.Itoa(i)
			arg.Scheme = caes.BasicSchemes["cumulative"]
			arg.Premises = premises
			arg.Conclusion = &conclusion
			ag.Arguments[arg.Id] = arg
			conclusion.Args = append(conclusion.Args, arg)
			i1.Positions = append(i1.Positions, &conclusion)
		}

		// validate the argument graph
		err = validate(ag)
		if err != nil {
			errorTemplate.Execute(w, err.Error())
			return
		}

		// export the argument graph to JSON

		cag, err := cj.Caes2Json(ag)
		if err != nil {
			errorTemplate.Execute(w, err.Error())
			return
		}

		// upload the JSON document to Couchdb

		id := uuid.NewUUID().String()
		req2, err := http.NewRequest("PUT", couchdbURL+casesDBName+"/"+id, bytes.NewBufferString(cag.String()))
		if err != nil {
			errorTemplate.Execute(w, err.Error())
			return
		}

		client := http.Client{}
		resp, err := client.Do(req2)
		if err != nil {
			errorTemplate.Execute(w, err.Error())
			return
		}
		defer resp.Body.Close()

		// check that uploading was successful

		m := make(map[string]interface{})
		v, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errorTemplate.Execute(w, err.Error())
			return
		}
		err = json.Unmarshal(v, &m)
		if err != nil {
			errorTemplate.Execute(w, err.Error())
			return
		}

		// m should have the form:
		// {"ok":true,"id":"6e1295ed6c29495e54cc05947f18c8af","rev":"1-2902191555"}

		if m["ok"] != true {
			errorTemplate.Execute(w, "Saving the case to Couchdb failed")
			return
		}

		data := map[string]string{
			"DomainId": domainId,
			"CaseId":   id,
		}

		newCaseCreatedTemplate.Execute(w, data)

	}
}
