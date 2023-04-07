package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ViaCEP struct {
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	IBGE        string `json:"ibge"`
	GIA         string `json:"gia"`
	DDD         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type ApiCep struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

func main() {
	chanApiCep := make(chan ApiCep)
	chanViaCep := make(chan ViaCEP)
	ctx, cancel := context.WithCancel(context.Background())

	go findCepApiCep(ctx, "01001-000", chanApiCep)
	go findCepViaCep(ctx, "01001-000", chanViaCep)

	select {
	case c := <-chanApiCep:
		cancel()
		fmt.Println("ApiCep:", c)
	case c := <-chanViaCep:
		cancel()
		fmt.Println("ViaCep:", c)
	case <-time.After(1 * time.Second):
		cancel()
		fmt.Println("Timeout")
	}
}

func findCepApiCep(ctx context.Context, cep string, ch chan<- ApiCep) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://cdn.apicep.com/file/apicep/"+cep+".json", nil)
	if err != nil {
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	var c ApiCep
	err = json.Unmarshal(body, &c)
	if err != nil {
		return
	}
	ch <- c
}

func findCepViaCep(ctx context.Context, cep string, ch chan<- ViaCEP) {
	req, err := http.NewRequestWithContext(ctx, "GET", "http://viacep.com.br/ws/"+cep+"/json/", nil)
	if err != nil {
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	var c ViaCEP
	err = json.Unmarshal(body, &c)
	if err != nil {
		return
	}
	ch <- c
}
