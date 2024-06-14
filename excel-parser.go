package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)
	type envVariables struct {
		path          string
		worksheetNames []string
	}

func main() {

	var env envVariables


	var debug bool = true
	println("debug mode:", debug)


	envChan := make(chan envVariables)
	go func() {
		envChan <- setupEnv()
	}()
	env = <-envChan


	file, err := excelize.OpenFile(env.path)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	defer file.Close()

	rows, err := file.GetRows("csp")
	if err != nil {
		fmt.Println(err)
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
			var loopCount int = 0 
			for _, row := range rows {
				loopCount++
				if len(row) >3 {
					if row[2] == search {
					fmt.Println(row[10])
					}
				}
			}
			search = ""
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

func setupEnv() envVariables {
	var env envVariables
	var worksheetStr string
	fmt.Println("enter file path: ")
	fmt.Scanln(&env.path)
	re := regexp.MustCompile(`\\`)
	env.path = re.ReplaceAllLiteralString(env.path, "/")
fmt.Println("Enter worksheet names comma separated:")
fmt.Scanln(&worksheetStr)
env.worksheetNames = strings.Split(worksheetStr, ",")
	return env
}
