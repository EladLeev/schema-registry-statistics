package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/fatih/color"
)

type ResultStats struct {
	StatMap     map[string]int
	ResultStore map[uint32][]int
}

func CalcPercentile(k string, v, consumedMessages int) {
	idPerc := math.Round((float64(v) / float64(consumedMessages) * 100))
	c := color.New(color.FgGreen, color.BgWhite)
	c.Printf("Schema ID %v => %v%%\n", k, idPerc)
}

// appendResult will map the results to a storeable map
func AppendResult(stat ResultStats, offset int64, schemaId uint32) {
	stat.ResultStore[schemaId] = append(stat.ResultStore[schemaId], int(offset))
}

// calcStat keep on track the stats
func CalcStat(stat ResultStats, schemaId uint32) {
	stat.StatMap[fmt.Sprint(schemaId)] += 1
	stat.StatMap["TOTAL"] += 1
}

func DumpStats(stat ResultStats, path string) {
	j, err := json.Marshal(stat.ResultStore)
	if err != nil {
		log.Printf("Error: %s", err.Error())
	} else {
		err := os.WriteFile(path, j, 0644)
		if err != nil {
			log.Fatalf("Error: %s", err.Error())
		}
		log.Printf("Results saved to %v\n", path)
	}
}
