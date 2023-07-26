// @author: @perpetualgod

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
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

func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func writeBigSpreadToFile(f *os.File, spread float64, path string, curTime string) {
	_, err := f.WriteString(fmt.Sprintf("spread: %.2f%%, path: %s\n", spread, path))
	HandleError(err)
	_, err = f.WriteString(fmt.Sprintf("time: %s\n", curTime))
	HandleError(err)
}

func updPairs(f *os.File) {
	chain := "bsc"
	// pairOders := []string{"LVL/WBNB"}
	pairs := [][]string{
		{"0x70f16782010fa7dDf032A6aaCdeed05ac6B0BC85", "0x7f9d307973CDAbe42769D9712DF8ee1cC1A28D10", "0x011077b8199cAb999e895AED7f8A78755A678106"},
		{"123"}}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.dexscreener.com/latest/dex/pairs/"+chain+"/"+strings.Join(pairs[0], ","), nil)

	HandleError(err)
	resp, err := client.Do(req)
	HandleError(err)
	defer resp.Body.Close()
	resRawData, err := ioutil.ReadAll(resp.Body)
	HandleError(err)

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

	curTime := time.Now().Format("15:04:05")
	fmt.Printf("time: %s\n", curTime)

	fmt.Println(parsedPairs)
	spread := (parsedPairs[len(parsedPairs)-1].PriceNative/parsedPairs[0].PriceNative - 1) * 100
	spreadPath := fmt.Sprintf("%s->%s", parsedPairs[0].DexId, parsedPairs[len(parsedPairs)-1].DexId)
	// spread += 5

	fmt.Printf("spread: %.2f%%, path: %s\n", spread, spreadPath)

	if spread > 1.5 && f != nil {
		writeBigSpreadToFile(f, spread, spreadPath, curTime)
	}
}

func heartBeat() {
	f, err := os.Create("big_spreads.txt")
	HandleError(err)
	defer f.Close()

	updRate := 30
	for range time.Tick(time.Second * time.Duration(updRate)) {
		updPairs(f)
	}
}

func decodeSqrtPrice(sqrtPriceDecimal string, factor int) float64 {
	num, _ := new(big.Float).SetPrec(256).SetString(sqrtPriceDecimal)
	exp := big.NewFloat(1)
	exp.SetMantExp(exp, -factor)
	num.Mul(num, exp)
	num.Mul(num, num)
	num.Quo(big.NewFloat(1e12), num)
	toFloat64, _ := num.Float64()
	return toFloat64
}

func main() {
	arg := flag.String("func", "nothing", "func choice to execute")
	flag.Parse()

	switch parsedFlag := *arg; parsedFlag {
	case "updPairs":
		updPairs(nil)
	case "runNodeProvider":
		runNodeProvider()
	case "heartBeat":
		heartBeat()
	default:
		return
	}
	// runNodeProvider()

	// res := decodeSqrtPrice("31365730092477674391378819923527281", 96)
	// fmt.Printf("%f", res)
}
