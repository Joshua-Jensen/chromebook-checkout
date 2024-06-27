package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
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

type task func(data map[string]string) (map[string]string)

type worker struct{
	id int
	taskQueue <-chan task
	resultChan <-chan Result
}

//Result is the output from each worker it should have a new map of all the found rows and what sheet it was found from
type Result struct{
	workerID int
	data map[string]string
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
	sheetsLen := len(env.worksheetNames)
	data := make(map[string][][]string, sheetsLen)

	for _, sheet := range env.worksheetNames {
		rows, err := file.GetRows(sheet)
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
		data[sheet] = rows
	}

	// depreciated testing learning code
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
			
			var searchWG sync.WaitGroup
			var bufferSize int = 10	


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

func newCbItem(item []string) cbItem {
	var newItem cbItem
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
	return newItem
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




func searchWorker(data [][]string, keyword string) ([][]string){
	var foundItem [][]string
	for _, item := range data {
		if item[10] == keyword {
foundItem = append(foundItem, item)
		}
	}
	return foundItem
} 