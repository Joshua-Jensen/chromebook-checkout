package main

import (
	"log"

	"github.com/xuri/excelize/v2"
)

func main() {
	file, err := excelize.OpenFile("C:/Users/User/source/go-repos/chromebook-checkout/div-test.xlsx")
	if err != nil {
		println(err)
		log.Fatal(err)
	}
	rows, err := file.GetRows("High VC")
	if err != nil {
		println(err)
		log.Fatal(err)
	}
	println(rows)
	log.Println(rows[8][1])

}
