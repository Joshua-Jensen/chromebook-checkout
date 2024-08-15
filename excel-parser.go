package main

import (
	"errors"
	"fmt"
	"log"
	"math"
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
	rowInt         int
}

func main() {

	var env envVariables

	var debug bool = true
	if debug {
		println("debug mode:", debug)
	}

	envChan := make(chan envVariables)
	go func() {
		envChan <- setupEnv()
	}()
	env = <-envChan

fmt.Println("")
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

		for rowInt, row := range rows {
			//TODO - send each row to be made into cb item and handle the error created
			//REVIEW - done
			item, err := newCbItem(row, sheet, rowInt)
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
			var breakStart int = 0
			var taskLen int = len(items)
			var breakSizeFloat float32 = float32(taskLen) / 20
			var breakSize int = int(math.Ceil(float64(breakSizeFloat)))
			var index int = 0
			var wg sync.WaitGroup

			for index <= 20 {
				var task []cbItem
				if breakStart+breakSize > taskLen {
					task = items[breakStart:taskLen]
				} else {
					task = items[breakStart : breakStart+breakSize]
				}
				wg.Add(1)
				go func() {
					defer wg.Done()
					searchWorkerAssetTag(index, task, search)
				}()
				index++
				breakStart += breakSize
			}
			wg.Wait()
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

func newCbItem(item []string, sheet string, rowInt int) (cbItem, error) {
	var newItem cbItem
	if len(item) >= 15 {

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
		if sheet == "general" {
			newItem.location = item[9]
		} else {
			newItem.location = item[10]
		}
		newItem.condition = item[11]
		newItem.inventoryTaken = item[12]
		newItem.disposalDate = item[13]
		newItem.disposalPrice = item[14]
		newItem.campus = item[15]
		newItem.sheetName = sheet
		newItem.rowInt = rowInt
		return newItem, nil
	}
	return newItem, errors.New("row did not contain enough values to be cbItem")
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

// NOTE - search function needs to work as a concurrent worker
func searchWorkerSN(id int, task []cbItem, keyword string) ([]cbItem, int) {
	var foundItem []cbItem
	for _, item := range task {
		if item.sn == keyword {
			foundItem = append(foundItem, item)
		}
	}
	return foundItem, id

}

func searchWorkerAssetTag(id int, task []cbItem, keyword string) ([]cbItem, int) {
	var foundItem []cbItem
	for _, item := range task {
		if item.assetTag == keyword {
			foundItem = append(foundItem, item)
		}
	}

	return foundItem, id
}

func changeRoom(selected cbItem, roomNumber string, file *excelize.File) error {
	var coordinates string
	if selected.sheetName == "general" {
		coordinates = fmt.Sprintf("%s_%s", "I", string(selected.rowInt))
	} else {
		coordinates = fmt.Sprintf("%s_%s", "K", string(selected.rowInt))

	}
	err := file.SetCellValue(selected.sheetName, coordinates, roomNumber)
	if err != nil {
		return err
	}
	return nil
}
