package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type BrasilAPI struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {
	brasilAPIChannel := make(chan BrasilAPI)
	viaCEPChannel := make(chan ViaCEP)

	go func() {
		res, err := getCEPInfo("https://brasilapi.com.br/api/cep/v1/01153000")
		checkErr(err)
		var data BrasilAPI
		err = json.Unmarshal(res, &data)
		checkErr(err)
		brasilAPIChannel <- data
	}()

	go func() {
		res, err := getCEPInfo("http://viacep.com.br/ws/01153000/json/")
		checkErr(err)
		var data ViaCEP
		err = json.Unmarshal(res, &data)
		checkErr(err)
		viaCEPChannel <- data
	}()

	select {
	case cep := <-viaCEPChannel:
		fmt.Printf("ViaCEPAPI is faster. The CEP info is: \n-CEP: %s\n-City: %s\n-Street: %s", cep.Cep, cep.Localidade, cep.Logradouro)
	case cep := <-brasilAPIChannel:
		fmt.Printf("BrasilAPI is faster. The CEP info is: \n-CEP: %s\n-City: %s\n-Street: %s", cep.Cep, cep.City, cep.Street)
	case <-time.After(1 * time.Second):
		panic("Timeout")
	}
}

func getCEPInfo(url string) ([]byte, error) {
	req, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	res, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
