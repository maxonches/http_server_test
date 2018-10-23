package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type M map[string]interface{}

type response struct {
    Rates []M `json:"Rates"`
}

var cb_url = "https://www.cbr-xml-daily.ru/daily_json.js"

func handler(w http.ResponseWriter, r *http.Request) {
	c := http.Client{}
	resp, err := c.Get(cb_url)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	bodyMap := map[string]interface{}{}

	json.Unmarshal(body, &bodyMap)

	valutes, ok := bodyMap["Valute"].(map[string]interface{})
	if !ok {
		fmt.Println("Incorrect value in Valute, exiting.")
		return
	}

	valutesUsd, ok := valutes["USD"].(map[string]interface{})
	if !ok {
		fmt.Println("Incorrect value in Valute->USD, exiting.")
		return
	}

	valutesEur, ok := valutes["EUR"].(map[string]interface{})
	if !ok {
		fmt.Println("Incorrect value in Valute->EUR, exiting.")
		return
	}

	var mMapSlice []M

    m1 := M{"From": "Рубль", "To": "Доллар США", "Value": valutesUsd["Value"]}
    m2 := M{"From": "Рубль", "To": "Евро", "Value": valutesEur["Value"]}

    mMapSlice = append(mMapSlice, m1, m2)

     res := &response{
         Rates: mMapSlice}
     resed, _ := json.MarshalIndent(res, "", "    ")

	w.WriteHeader(200)
	w.Write(resed)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":9000", nil)
}

