package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PokémonData struct {
	FlavorTextEntries []FlavorTextEntry `json:"flavor_text_entries"`
	Name              string            `json:"name"`
	Habitat           Habitat           `json:"habitat"`
	IsLegendary       bool              `json:"is_legendary"`
}

type FlavorTextEntry struct {
	FlavorText string           `json:"flavor_text"`
	Language   NamedAPIResource `json:"language"`
}

type NamedAPIResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Habitat struct {
	Name string `json:"name"`
}

type PokeApiClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewPokeApiClient() *PokeApiClient {
	return &PokeApiClient{
		baseURL: "https://pokeapi.co/api/v2",
		httpClient: &http.Client{
			Timeout: CLIENTS_TIMEOUT * time.Second,
		},
	}
}

func (c *PokeApiClient) GetSpecie(idOrName string) (*PokémonData, error) {
	url := fmt.Sprintf("%s/pokemon-species/%s", c.baseURL, idOrName)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed making request to PokeAPI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("pokemon specie '%s' not found", idOrName)
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from PokeAPI: %d", resp.StatusCode)
	}

	species := &PokémonData{}
	if err := json.NewDecoder(resp.Body).Decode(&species); err != nil {
		return nil, fmt.Errorf("failed decoding JSON response: %w", err)
	}

	return species, nil
}
