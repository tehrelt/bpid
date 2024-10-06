package main

import (
	"fmt"
	"os"
)

func main() {
	var filename string

	fmt.Print("enter the filename: ")
	fmt.Scan(&filename)

	var file *os.File
	var err error

	file, err = os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fmt.Println("file opened successfully")
}
