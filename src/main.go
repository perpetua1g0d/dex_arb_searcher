// @author: @perpetualgod

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type PairData struct {
	DexId       string
	PriceNative float64
	Liq         float64
}

func updPairs() {
	chain := "bsc"
	// pairOders := []string{"LVL/WBNB"}
	pairs := [][]string{
		{"0x70f16782010fa7dDf032A6aaCdeed05ac6B0BC85", "0x7f9d307973CDAbe42769D9712DF8ee1cC1A28D10", "0x011077b8199cAb999e895AED7f8A78755A678106"},
		{"123"}}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.dexscreener.com/latest/dex/pairs/"+chain+"/"+strings.Join(pairs[0], ","), nil)

	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	resRawData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var parsedPairs []PairData

	var parsed map[string]any
	json.Unmarshal([]byte(resRawData), &parsed)
	pairsMapArray := parsed["pairs"].([]any)
	for _, value := range pairsMapArray {
		curMap := value.(map[string]any)
		pairData := &PairData{}

		pairData.DexId = curMap["dexId"].(string)
		pairData.Liq = curMap["liquidity"].(map[string]any)["usd"].(float64)
		if pairData.PriceNative, err = strconv.ParseFloat(curMap["priceNative"].(string), 64); err != nil {
			log.Fatal(err)
		}

		parsedPairs = append(parsedPairs, *pairData)
	}

	sort.Slice(parsedPairs, func(i int, j int) bool {
		return parsedPairs[i].PriceNative < parsedPairs[j].PriceNative
	})

	dt := time.Now()
	fmt.Printf("time: %s\n", dt.Format("15:04:05"))

	fmt.Println(parsedPairs)
	spread := (parsedPairs[len(parsedPairs)-1].PriceNative/parsedPairs[0].PriceNative - 1) * 100
	fmt.Printf("spread: %.2f%%, path: %s->%s\n", spread, parsedPairs[0].DexId, parsedPairs[len(parsedPairs)-1].DexId)
}

func heartBeat() {
	updRate := 30
	for range time.Tick(time.Second * time.Duration(updRate)) {
		updPairs()
	}
}

func main() {
	go heartBeat()
	for true {

	}
}
