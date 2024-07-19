package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type envVariables struct {
	path           string
	worksheetNames []string
}
type room struct {
	roomNum      string
	roomContents []cbItem
}
type cbItem struct {
	itemDesc       string
	sn             string
	assetTag       string
	funding        string
	award          string
	fain           string
	titleHolder    string
	aqcDate        string
	cost           string
	fedPartPercent string
	location       string
	condition      string
	inventoryTaken string
	disposalDate   string
	disposalPrice  string
	campus         string
	sheetName      string
}



func main() {

	var env envVariables

	var debug bool = true
	if debug{
	println("debug mode:", debug)
	}

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
	var items []cbItem
	for _, sheet := range env.worksheetNames {
		rows, err := file.GetRows(sheet)
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}

		for _, row:= range rows{
//TODO - send each row to be made into cb item and handle the error created  
//REVIEW - done
item, err :=newCbItem(row, sheet)
if err == nil {
	items = append(items, item)
}
		}
	}

//NOTE depreciated testing learning code- 
	// fmt.Println(rows[10][10])
	// fmt.Println("Enter room number: ")
	// fmt.Scanln(&room)
	// fmt.Println(room)

	var loop bool = true
	var search string

	for loop {
		fmt.Println("input sn/id: ")
		fmt.Scanln(&search)

		if search != "e" {
//TODO - add search function with workers
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

func newCbItem(item []string, sheet string) (cbItem, error)  {
	var newItem cbItem
if len(item) >= 15{


		newItem.itemDesc = item[0]
		newItem.sn = item[1]
		newItem.assetTag = item[2]
		newItem.funding = item[3]
		newItem.award = item[4]
		newItem.fain = item[5]
		newItem.titleHolder = item[6]
		newItem.aqcDate = item[7]
		newItem.cost = item[8]
		newItem.fedPartPercent = item[9]
		newItem.location = item[10]
		newItem.condition = item[11]
		newItem.inventoryTaken = item[12]
		newItem.disposalDate = item[13]
		newItem.disposalPrice = item[14]
		newItem.campus = item[15]
	return newItem,nil
	}
	return newItem , errors.New("row did not contain enough values to be cbItem")
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


//NOTE - search function needs to work as a concurrent worker
func searchWorker(id int, task [][]string, keyword string) {
	var foundItem [][]string
	for _, item := range task {
		if item[10] == keyword {
			foundItem = append(foundItem, item)
		}	
	}
}
