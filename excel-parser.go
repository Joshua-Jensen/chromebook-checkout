package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

func main() {
	type envVariables struct {
		path          string
		worksheetName string
	}
	var debug bool = true
	println("debug mode:", debug)
	pathChan := make(chan string)
	go func() {
		pathChan <- setupEnv()
	}()
	path := <-pathChan

	file, err := excelize.OpenFile(path)
	if err != nil {
		println(err)
		log.Fatal(err)
	}
	defer file.Close()

	rows, err := file.GetRows("csp")
	if err != nil {
		println(err)
		log.Fatal(err)
	}

	fmt.Println(rows[10][10])
	var room string
	var checkedItems [][]string

	fmt.Println("Enter room number: ")
	fmt.Scanln(&room)
	fmt.Println(room)

	for _, row := range rows {
		for _, colCell := range row {
			if colCell == room {
				fmt.Println(row)
				checkedItems = append(checkedItems, row)
			}
		}
	}

	var itemsCount int = len(checkedItems)
	var itemsString string = strconv.FormatInt(int64(itemsCount), 10)
	msg := fmt.Sprintf("\n Checked items count: %s", itemsString)
	fmt.Println(msg)

	var loop bool = true
	var search string
	for loop {
		fmt.Println("input sn/id: ")
		fmt.Scanln(&search)
		if search != "" {
			for _, row := range rows {
				if row[2] == search {
					fmt.Println(row[10])
				}
			}
		} else {
			println("exiting")
			time.Sleep(time.Second)
			c := exec.Command("clear")
			c.Stdout = os.Stdout
			c.Run()
			loop = false
		}
	}

}

func setupEnv() string {
	var path string
	fmt.Println("enter file path: ")
	fmt.Scanln(&path)
	re := regexp.MustCompile(`\\`)
	path = re.ReplaceAllLiteralString(path, "/")

	return path
}
