package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func main() {
	file, err := excelize.OpenFile("C:/Users/User/source/go-repos/chromebook-checkout/cli-tool-test.xlsx")
	if err != nil {
		println(err)
		log.Fatal(err)
	}
	defer file.Close()

rows,err := file.GetRows("csp") 
if err != nil {
	println(err)
	log.Fatal(err)	
}
fmt.Println(rows[10][9])
reader := 	bufio.NewReader(os.Stdin)
fmt.Println("Enter colum num: ")
text, err := reader.ReadString("\n")
num, err := strconv.Atoi(text)
if err != nil {
	fmt.Println(err)
	log.Fatal(err)
}

fmt.Println(rows[10][num])

}
