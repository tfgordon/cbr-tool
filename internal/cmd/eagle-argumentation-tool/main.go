// Copyright © 2015 Thomas F. Gordon
// This Source Code Form is subject to the terms of the
// Mozilla Public License, v. 2.0. If a copy of the MPL
// was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"git.list.lu/eagle/argumentation-tool/internal/web"
	"log"
	"os"
	"path/filepath"
)

const help = `
usage: eagle-argumentation-tool [-p port] [-s couchdb-url]

Starts the Eagle argumentation tool as a web service.

The -p flag ("port") specifies the port to use for the web service (default: 8080)

The -c flag ("couchdb") specifies the URL of the Couchdb server (default: http://127.0.0.1:5984/)
`

var goPath = os.Getenv("GOPATH")

const defaultHttpPort = "8080"
const defaultCouchdbURL = "http://127.0.0.1:5984/"
const casesDBName = "cases"
const domainsDBName = "domains" // domain models

var templatesDir = filepath.Join(goPath, "/src/git.list.lu/eagle/argumentation-tool/internal/web/templates/")

func main() {
	flags := flag.NewFlagSet("flags", flag.ContinueOnError)
	httpPort := flags.String("p", defaultHttpPort, "the port number of the web service")
	couchdbURL := flags.String("c", defaultCouchdbURL, "the URL of the Couchdb server")
	apiKey := flags.String("k", "", "The API key of the EAGLE platform with access to this service")
	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	cdb := web.CouchDBConfig{
		URL:     *couchdbURL,
		Domains: domainsDBName,
		Cases:   casesDBName,
		APIKey:  *apiKey,
	}

	web.ArgumentationToolServer(*httpPort, cdb, templatesDir)
}
