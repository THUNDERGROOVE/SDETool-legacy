all:
	go fmt github.com/THUNDERGROOVE/SDETool
	go build -v github.com/THUNDERGROOVE/SDETool
date:
	go build -v -ldflags "-X main.BuildDate \"%date%  %time%\"" github.com/THUNDERGROOVE/SDETool
clean:
	rm dustSDE.db.zip
	rm dustSDE.db
dep:
	go get github.com/joshlf13/term
	go install github.com/joshlf13/term
	go get github.com/mattn/go-sqlite3
	go install github.com/mattn/go-sqlite3