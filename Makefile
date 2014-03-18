all:
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
