run:
	go run .
tidy:
	go mod tidy
migrate:
	goose -dir ./database/migrations sqlite3 ./database/url-shortly.db up