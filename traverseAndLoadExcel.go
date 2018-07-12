package loaddata

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

//var wg sync.WaitGroup
var count int

func ReadExcelFileAndLoad(name, fullPath string, container chan interface{}, count int, rowContainer chan []string) {
	var flag int
	defer func() {
		log.Println("Finish", fullPath)
		return
	}()
	if _, err := os.Stat(fullPath); err == nil {
		container <- name
		xlsx, err := excelize.OpenFile(fullPath)
		if err != nil {
			fmt.Println("Failed because of ", err)
			return
		}
		rows := xlsx.GetRows(xlsx.GetSheetName(xlsx.GetActiveSheetIndex()))

		for _, row := range rows {
			rowContainer <- row
		}
	} else {
		fmt.Println("Not workign", flag)
	}
	//time.Sleep(1)
	return
}

/*
	track the channel and load the data into the table
*/

func DumpInsideDB(conn *sql.DB, container chan interface{}, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	for {
		select {
		case details := <-container:
			log.Println("File Detail :", details)
		case <-time.After(time.Second * 59):
			//wg.Done()
			log.Println("Waiting for file name", len(container))
			return
		}
	}
}

func insertDetail(conn *sql.DB, query string, details []interface{}) {
	_, err := conn.Exec(query, details...)
	if err != nil {
		log.Println("Got error at the inserttime :", err)
	}
}

func RowDumpInsideDB(conn *sql.DB, rowContainer chan []string, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	for {
		select {
		case details := <-rowContainer:
			//count += 1
			query := "insert into studentDetail(id, name, details) values(?,?,?);"
			var data []interface{}
			for _, val := range details {
				data = append(data, val)
			}
			//log.Println(query, data)
			insertDetail(conn, query, data)
		case <-time.After(time.Second * 59):
			log.Println("Waiting for Row name", len(rowContainer))
			return
		}
	}
}

func TraverseAllDirectoryForExcel(conn *sql.DB, path string) {
	//change directory to the given directory
	runtime.GOMAXPROCS(3)
	var wg sync.WaitGroup
	//log.Println(path)
	err := os.Chdir(path)
	if nil != err {
		log.Panic(err)
	}

	// for reading the current direcotry we need
	listDir, err := ioutil.ReadDir("./")
	if nil != err {
		log.Fatal(err)
	}
	//var chnlCount int = 1
	chnl := make(chan interface{})
	stringChnl := make(chan []string)
	wg.Add(2)
	go DumpInsideDB(conn, chnl, &wg)
	go RowDumpInsideDB(conn, stringChnl, &wg)
	for _, fileInfo := range listDir {
		switch mode := fileInfo.Mode(); {
		case mode.IsDir():
			log.Println("Not working inside directory", fileInfo.Name())
		case mode.IsRegular():
			fullPath, _ := filepath.Abs(fileInfo.Name())
			count += 1
			go ReadExcelFileAndLoad(fileInfo.Name(), fullPath, chnl, count, stringChnl)
		}
	}
	wg.Wait()
	close(chnl)
	close(stringChnl)
	fmt.Println("list is has len", len(chnl))
}
