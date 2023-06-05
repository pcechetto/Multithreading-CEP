package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CDNApiCep struct {
	Cep        string `json:"cep"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

type ViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func makeAPIRequest(url string, ch chan<- []byte) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("method 'get' error:", err)
		ch <- nil
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read body error", err)
		ch <- nil
		return
	}

	ch <- body
}

func main() {
	cep := "69915-630"
	apicep := "https://cdn.apicep.com/file/apicep/" + cep + ".json"
	viacep := "http://viacep.com.br/ws/" + cep + "/json/"

	ch1 := make(chan []byte)
	ch2 := make(chan []byte)

	go makeAPIRequest(apicep, ch1)
	go makeAPIRequest(viacep, ch2)

	select {
	case res1 := <-ch1:
		if res1 != nil {
			var address1 CDNApiCep
			err := json.Unmarshal(res1, &address1)
			if err != nil {
				fmt.Println("Unmarshal error on cdn.apicep.com: ", err)
			} else {
				fmt.Println("cdn.apicep.com:")
				fmt.Printf("%+v\n", address1)
			}
		}
	case res2 := <-ch2:
		if res2 != nil {
			var address2 ViaCep
			err := json.Unmarshal(res2, &address2)
			if err != nil {
				fmt.Println("Unmarshall error viacep.com.br: ", err)
			} else {
				fmt.Println("viacep.com.br: ")
				fmt.Printf("%+v\n", address2)
			}
		}
	case <-time.After(1 * time.Second):
		fmt.Println("timeout")
	}
}
