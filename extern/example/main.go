package main

import (
	"fmt"
	"github.com/THUNDERGROOVE/SDETool/extern"
)

func main() {
	c, err := extern.DTDKGetFittingByID(669)
	if err != nil {
		fmt.Println("Error parsing JSON", err.Error())
	}
	fmt.Println(c)
}
