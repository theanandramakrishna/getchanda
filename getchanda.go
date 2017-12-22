package main

import "fmt"
import "net/http"
import "io/ioutil"
import "log"
import "os"
import "sync"

type YearSet struct {
	beginYear int
	endYear int
}
func main() {
//	yearsets := []YearSet{{1954, 1956}, {1957, 1959}, {1970, 1973}, {1974, 1976},{1977, 1979}, {1980, 1983}, {1984, 1986}}
	yearsets := []YearSet{{1954, 1956}}
	var wg sync.WaitGroup

	for _, yearset := range yearsets {
		for year := yearset.beginYear; year <= yearset.endYear; year++ {
			for month := 1; month <= 12; month++ {
				var fileName string = 
					fmt.Sprintf("Chandamama-%d-%d.pdf", year, month)
				var urlString string = 
					fmt.Sprintf("http://chandamama.in/resources/english/%d-%d/%s", 
						yearset.beginYear, yearset.endYear, fileName)
				
				if _, err := os.Stat(fileName); err == nil {
					// filename already exists, just skip it.
					log.Printf("%s already exists, continuing.", fileName)
					continue
				}

				wg.Add(1)
				go func (urlString string) {
					defer wg.Done()
					log.Printf("Fetching %s\n", urlString)
					resp, err := http.Get(urlString)
					if err != nil {
						log.Printf("ERROR: %s\n", err.Error())
						return
					}
					defer resp.Body.Close()
	
					if resp.StatusCode != http.StatusOK {
						log.Printf("%s not found, continuing.\n", urlString)
						return
					}
					buf, err := ioutil.ReadAll(resp.Body)
					if (err != nil) {
						log.Printf("ERROR: %s\n", err.Error())
						return
					}
	
					err = ioutil.WriteFile(fileName, buf, 0777)
					if err != nil {
						log.Printf("ERROR: %s\n", err.Error())
						return
					} 
				}(urlString)
			}
		}
	}

	wg.Wait()	// Wait until all goroutines are done
}