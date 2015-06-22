package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
)

type LottoResults struct {
	Results []LottoItem
}

type LottoItem struct {
	Game        string
	Combination string
	Date        string
	Jackpot     string
	Winners     string
}

const CLR_0 = "\x1b[39;1m"
const CLR_R = "\x1b[31;1m"
const CLR_G = "\x1b[32;1m"
const CLR_Y = "\x1b[33;1m"
const CLR_B = "\x1b[34;1m"
const CLR_M = "\x1b[35;1m"
const CLR_C = "\x1b[36;1m"
const CLR_W = "\x1b[37;1m"
const CLR_X = "\x1b[38;1m"
const CLR_N = "\x1b[0m"

func main() {
	wg := &sync.WaitGroup{}
	if _, err := os.Stat("results"); os.IsNotExist(err) {
		os.Mkdir("."+string(filepath.Separator)+"results", 0777)
	}
	runtime.GOMAXPROCS(4)
	wg.Add(8);
	go StartScrape("http://pcso-lotto-results-and-statistics.webnatin.com/6-55results.asp", CLR_0, "results/6-55.json", wg )
	go StartScrape("http://pcso-lotto-results-and-statistics.webnatin.com/6-49results.asp", CLR_R, "results/6-49.json", wg )
	go StartScrape("http://pcso-lotto-results-and-statistics.webnatin.com/6-45results.asp", CLR_G, "results/6-45.json", wg )
	go StartScrape("http://pcso-lotto-results-and-statistics.webnatin.com/6-42results.asp", CLR_Y, "results/6-42.json", wg )
	go StartScrape("http://pcso-lotto-results-and-statistics.webnatin.com/6-dresults.asp", CLR_B, "results/6-d.json", wg )
	go StartScrape("http://pcso-lotto-results-and-statistics.webnatin.com/4-dresults.asp", CLR_M, "results/4-d.json", wg )
	go StartScrape("http://pcso-lotto-results-and-statistics.webnatin.com/3-dresults.asp", CLR_C, "results/3-d.json", wg )
	go StartScrape("http://pcso-lotto-results-and-statistics.webnatin.com/2-dresults.asp", CLR_W, "results/2-d.json", wg )
	wg.Wait()
}

func StartScrape(url string, color string, filename string,wg *sync.WaitGroup ) {
	fmt.Println(color + "Creating " + filename + CLR_N)
	var item = LottoItem{Game: "", Combination: "", Date: "", Jackpot: "", Winners: ""}
	var results = LottoResults{}
	doc, queryError := goquery.NewDocument(url)
	if queryError != nil {
		fmt.Println("\nError fetching data..\n" + queryError.Error())
		return
	}
	insertedCounter := 0
	doc.Find("table tr").Each(func(i int, s *goquery.Selection) {
		s.Find("td").Each(func(j int, t *goquery.Selection) {
			if j == 0 {
				item.Game = t.Text()
			}
			if j == 2 {
				item.Date = t.Text()
			}
			if j == 3 {
				item.Combination = t.Text()
			}
			if j == 4 {
				item.Jackpot = t.Text()
			}
			if j == 5 {
				item.Winners = t.Text()
			}
		})
		results.Results = append(results.Results, item)
		fmt.Print(color + "●" + CLR_N)
		insertedCounter++
	})
	jsonResults, _ := json.Marshal(results)
	ioutil.WriteFile(filename, []byte(jsonResults), 0644)
	fmt.Println(color + "\nCreated json file :" + filename + " -> " + strconv.Itoa(insertedCounter) + " records " + CLR_N)
	wg.Done()
}
