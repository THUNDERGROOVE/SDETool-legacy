all:
	go build github.com/THUNDERGROOVE/SDETool
clean:
	rm dustSDE.db.zip
	rm dustSDE.db
test:
	go test -v
	go test -bench .