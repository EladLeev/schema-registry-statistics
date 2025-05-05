package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"sync"
)

// PercentileMap ...
var PercentileMap = map[string]float64{}

// ResultStats ...
type ResultStats struct {
	StatMap     map[string]int
	ResultStore map[uint32][]int
}

func calc(schema, total int) float64 {
	return math.Round((float64(schema) / float64(total) * 100))
}

// BuildPercentileMap will calculate the percentile into a map
func BuildPercentileMap(s map[string]int) {
	for k, v := range s {
		if k == "TOTAL" {
			continue
		} else if k == "ERROR" {
			defer log.Printf("Unable to decode schema in %v messages. They might be empty, or do not contains any schema.", v)
		} else {
			PercentileMap[k] = calc(v, s["TOTAL"])
		}
	}
}

// AppendResult will map the results to a storeable map
func AppendResult(stat ResultStats, offset int64, schemaID uint32, lock *sync.RWMutex) {
	lock.Lock()
	defer lock.Unlock()
	stat.ResultStore[schemaID] = append(stat.ResultStore[schemaID], int(offset))
}

// CalcStat keep on track the stats
func CalcStat(stat ResultStats, schemaID uint32, lock *sync.RWMutex) {
	lock.Lock()
	defer lock.Unlock()
	// 4294967295 represents error
	if schemaID == 4294967295 {
		stat.StatMap["ERROR"]++
		return
	}
	stat.StatMap[fmt.Sprint(schemaID)]++
	stat.StatMap["TOTAL"]++
}

// DumpStats will print the results into a file
func DumpStats(stat ResultStats, path string) {
	j, err := json.Marshal(stat.ResultStore)
	if err != nil {
		log.Printf("Error: %s", err.Error())
	} else {
		err := os.WriteFile(path, j, 0o600)
		if err != nil {
			log.Fatalf("Error: %s", err.Error())
		}
		log.Printf("Results saved to %v\n", path)
	}
}
