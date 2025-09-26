swagger:
	rm ./internal/models/*.go ./internal/restapi/operations/*.go
	swagger generate server --exclude-main --exclude-spec -t internal/ -f api/swagger/file.yaml --name rest-server