package main

import (
	"fmt"
)

func main() {
	r, err := LoadCsv("test.csv")
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	fmt.Println(r)
	fmt.Printf("Field1 Col: %v", r.GetColumn("Field"))
}
