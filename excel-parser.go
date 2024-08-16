package main

import (
	"fmt"
	"log"
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
	// itemDesc       string
	sn             string
	assetTag       string
	// funding        string
	// award          string
	// fain           string
	// titleHolder    string
	// aqcDate        string
	// cost           string
	// fedPartPercent string
	location       string
	// condition      string
	// inventoryTaken string
	// disposalDate   string
	// disposalPrice  string
	// campus         string
	sheetName      string
	rowInt         int
}

type searchWorkerResult struct {
	goodResult []cbItem
	id         int
	keyword    string
	task       []cbItem
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
			if len(row)> 4 {
			item := newCbItem(row, sheet, rowInt)
				items = append(items, item)
			}else{
				fmt.Println("unsupported item", row)
			}
		}
	}
for _, printRow := range items{
	fmt.Println(printRow)
}


	var newRoomNum string
	fmt.Println("Enter room number to add to")
	fmt.Scanln(&newRoomNum)
	fmt.Println("entered room number: ", newRoomNum)
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
		fmt.Println("searching ....")
		if search != "e" {
			var wg sync.WaitGroup
			for _, item := range items {
				if item.assetTag == search {
					fmt.Println(item)
					wg.Add(1)
					go func() {
						err := changeRoom(item, newRoomNum, *&file)
						if err != nil {
							fmt.Println(err)
						}
					}()
					fmt.Println(item)
				}
			}

			//TODO - add search function with workers
			//NOTE - turn this into working code but commented out to do first run
			// var breakStart int = 0
			// var taskLen int = len(items)
			// var breakSizeFloat float32 = float32(taskLen) / 20
			// var breakSize int = int(math.Ceil(float64(breakSizeFloat)))
			// var index int = 0
			// workerChan := make(chan searchWorkerResult)
			// go func() {

			// 	for index <= 20 {
			// 		var task []cbItem
			// 		if breakStart+breakSize > taskLen {
			// 			task = items[breakStart:taskLen]
			// 		} else {
			// 			task = items[breakStart : breakStart+breakSize]
			// 		}
			// 		var message searchWorkerResult
			// 		message.task = task
			// 		message.id = index
			// 		message.keyword = search
			// 		message.goodResult = nil
			// 		workerChan <- message
			// 		// wg.Add(1)
			// 		// go func() {
			// 		// 	defer wg.Done()
			// 		// 	searchWorkerAssetTag(index, task, search)
			// 		// }()
			// 		index++
			// 		breakStart += breakSize
			// 	}
			// }()
			// close(workerChan)
			search = ""
		} else {
			println("exiting")
			time.Sleep(time.Second)
			// c := exec.Command("clear")
			// c.Stdout = os.Stdout
			// c.Run()
			loop = false
		}
	}

}

func newCbItem(item []string, sheet string, rowInt int) (cbItem) {
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

// NOTE - search function needs to work as a concurrent worker
func searchWorkerSN(ch chan searchWorkerResult, incoming searchWorkerResult) {
	for _, item := range incoming.task {
		if item.sn == incoming.keyword {
			incoming.goodResult = append(incoming.goodResult, item)
		}
	}
	ch <- incoming
}

func searchWorkerAssetTag(ch chan searchWorkerResult, incoming searchWorkerResult) {
	for _, item := range incoming.task {
		if item.assetTag == incoming.keyword {
			incoming.goodResult = append(incoming.goodResult, item)
		}
	}
	ch <- incoming
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
