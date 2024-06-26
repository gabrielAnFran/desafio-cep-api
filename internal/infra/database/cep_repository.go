package database

import (
	"cep-gin-clean-arch/models"
	"encoding/json"
	"errors"
	"os"

	"github.com/supabase-community/supabase-go"
)

type CEPRepository struct{}

func NewCEPRepository() *CEPRepository {
	return &CEPRepository{}
}

type Address struct {
	CEP    string `json:"CEP"`
	Estado string `json:"Estado"`
	Cidade string `json:"Cidade"`
	Bairro string `json:"Bairro"`
	Rua    string `json:"Rua"`
}

func (r *CEPRepository) Buscar(cep string) (models.CEPResponse, error) {

	// Primeiro busca os dados "em memória"
	dados := []byte(jsonDados)

	var addresses []Address

	// Unmarshal JSON data
	err := json.Unmarshal(dados, &addresses)
	if err != nil {
		return models.CEPResponse{}, errors.New("Erro ao acessar dados de CEP")
	}

	// Create a map for quick lookup
	addressMap := make(map[string]Address)
	for _, address := range addresses {
		addressMap[address.CEP] = address
	}
	// Se o CEP estiver na memória, retorna os dados
	desiredAddress, found := addressMap[cep]
	if found {
		return models.CEPResponse{
			Estado: desiredAddress.Estado,
			Cidade: desiredAddress.Cidade,
			Bairro: desiredAddress.Bairro,
			Rua:    desiredAddress.Rua,
		}, nil
	}

	// Caso nao encontre, busca na tabela Supabase
	// Se der erro ao buscar na tabela Supabase, retorna um objeto vazio e o erro, que sera logado no sentry
	// Exemplo: CEP 99150000 nao se encontra no json mas está na tabela Supabase
	address, err := buscarDadosNaTabelaSupabase(cep)
	if err != nil {
		return models.CEPResponse{}, errors.New(err.Error())
	}

	// Se encontrou na tabela Supabase, retorna os dados
	return models.CEPResponse{
		Estado: address.Estado,
		Cidade: address.Cidade,
		Bairro: address.Bairro,
		Rua:    address.Rua,
	}, nil

}

func buscarDadosNaTabelaSupabase(cep string) (Address, error) {
	// Cliente do Supabase
	client, err := supabase.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"), nil)
	if err != nil {
		return Address{}, errors.New("Erro ao acessar dados de CEP: " + err.Error())
	}

	data, _, err := client.From("cep").
		Select("rua, bairro, cidade, estado", "exact", false).
		Eq("cep", cep).
		Execute()
	if err != nil {
		return Address{}, errors.New("Erro ao acessar dados de CEP: " + err.Error())
	}

	var addresses []Address
	err = json.Unmarshal([]byte(data), &addresses)
	if err != nil {
		return Address{}, errors.New("Erro ao acessar dados de CEP: " + err.Error())
	}

	// Caso retorne um array vazio, retorna um erro
	if len(addresses) == 0 {
		return Address{}, errors.New("CEP não encontrado")
	}

	// se o len de addresses for maior que 0, retorna o primeiro endereco
	return addresses[0], nil
}
