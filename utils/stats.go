package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"sync"
)

var PercentileMap = map[string]float64{}

type ResultStats struct {
	StatMap     map[string]int
	ResultStore map[uint32][]int
}

func calc(schema, total int) float64 {
	return math.Round((float64(schema) / float64(total) * 100))
}

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
func AppendResult(stat ResultStats, offset int64, schemaId uint32, lock *sync.RWMutex) {
	lock.Lock()
	defer lock.Unlock()
	stat.ResultStore[schemaId] = append(stat.ResultStore[schemaId], int(offset))
}

// CalcStat keep on track the stats
func CalcStat(stat ResultStats, schemaId uint32, lock *sync.RWMutex) {
	lock.Lock()
	defer lock.Unlock()
	// 4294967295 represents error
	if schemaId == 4294967295 {
		stat.StatMap["ERROR"] += 1
		return
	}
	stat.StatMap[fmt.Sprint(schemaId)] += 1
	stat.StatMap["TOTAL"] += 1
}

func DumpStats(stat ResultStats, path string) {
	j, err := json.Marshal(stat.ResultStore)
	if err != nil {
		log.Printf("Error: %s", err.Error())
	} else {
		err := os.WriteFile(path, j, 0600)
		if err != nil {
			log.Fatalf("Error: %s", err.Error())
		}
		log.Printf("Results saved to %v\n", path)
	}
}
