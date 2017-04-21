
# RESTful API for managing domain models

## Creating Domain Models

	$ curl -X PUT "http://<domain-name>:<port>/create-domain?apikey=<uuid>&user=<userid>" -H "Content-Type: application/json" -d @<file>.json
	
where the domain name (or IP number) and port number are for the EAGLE argumentation tool (not the Couchdb server).
	
### Example

    $ curl -X PUT "http://127.0.0.1:8080/create-domain?apikey=123&user=arun" -H "Content-Type: application/json" -d @test1.json 

## Reading Domain Models

	$ curl http://<domain-name>:<port>/read-domain/<domainId>

### Example

	$ curl http://127.0.0.1:8080/read-domain/6ddbc7bf2e9a54415d6c84a2ba004a65

## Updating Domain Models

	$ curl -X PUT "http://<domain-name>:<port>/update-domain?apikey=<uuid>&user=<userid>" -H "Content-Type: application/json" -d @<file>.json
	
### Example

	$ curl -X PUT "http://127.0.0.1:8080/update-domain?apikey=123&user=arun" -H "Content-Type: application/json" -d @animals.json

## Deleting Domain Models

	$ curl "http://<domain-name>:<port>/delete-domain/<domainId>?apikey=<uuid>&user=<userid>"
	
### Example

	$ curl "http://127.0.0.1:8080/delete-domain/6ddbc7bf2e9a54415d6c84a2ba0041b7?apikey=123&user=arun"

