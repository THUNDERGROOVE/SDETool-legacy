all:
	go fmt github.com/THUNDERGROOVE/SDETool
	go build -v github.com/THUNDERGROOVE/SDETool
	go install github.com/THUNDERGROOVE/SDETool/category
	go install github.com/THUNDERGROOVE/SDETool/util
	go install github.com/THUNDERGROOVE/SDETool/args
clean:
	rm dustSDE.db.zip
	rm dustSDE.db
test:
	go test -v
	go test -bench .
dep:
	go get github.com/joshlf13/term
	go install github.com/joshlf13/term
	go get github.com/mattn/go-sqlite3
	go install github.com/mattn/go-sqlite3