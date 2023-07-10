/*

* / wbnb
buy * for bnb on dex1
sell * for bnb on dex2

*/

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
	DexId       string `json:"dexId"`
	PriceNative string `json:"priceNative"`
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

	var prices []float64
	// prices := []float64{}

	var parsed map[string]any
	json.Unmarshal([]byte(resRawData), &parsed)
	pairsParsed := parsed["pairs"].([]any)
	for _, value := range pairsParsed {
		marshalled, err := json.Marshal(value.(map[string]any))
		if err != nil {
			log.Fatal(err)
		}
		pairData := &PairData{}
		if err := json.Unmarshal(marshalled, &pairData); err != nil {
			log.Fatal(err)
		}
		w := value.(map[string]any)

		pairData.Liq = w["liquidity"].(map[string]any)["usd"].(float64)
		fmt.Println(pairData)
		if price, err := strconv.ParseFloat(pairData.PriceNative, 64); err == nil {
			prices = append(prices, price)
		} else {
			log.Fatal(err)
		}
	}

	sort.Float64s(prices)
	fmt.Println("spread is ", (prices[len(prices)-1]/prices[0]-1)*100, "%")

	dt := time.Now()
	fmt.Println(dt.Format("15:04:05"))
}

func heartBeat() {
	for range time.Tick(time.Second * 30) {
		updPairs()
	}
}

func main() {
	go heartBeat()
	for true {

	}
}
