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
defer file.Close()

for _, row := range excelize.Rows {
	for _, colCell
}
}
