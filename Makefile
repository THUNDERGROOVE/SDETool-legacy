all:
	go fmt github.com/THUNDERGROOVE/SDETool
	go build -v github.com/THUNDERGROOVE/SDETool
windows:
	go build -v -ldflags "-X main.BuildDate \"%date%  %time%\"" github.com/THUNDERGROOVE/SDETool
linux:
	go build -v -ldflags "-X main.BuildDate \"`date`\""
clean:
	rm -f *.db.zip
	rm -f *.db
	rm -f *.panic
	rm -f log.txt
	rm -f SDETool
	rm -f SDETool.exe
dep:
	go get github.com/joshlf13/term
	go install github.com/joshlf13/term
	go get github.com/mattn/go-sqlite3
	go install github.com/mattn/go-sqlite3
install: linux
	@echo Installing SDETool to /usr/local/bin
	@echo this will only work on Linux
	sudo cp -f SDETool /usr/local/bin/