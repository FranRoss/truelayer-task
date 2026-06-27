package service

import (
	"truelayer/internal/clients"
	"truelayer/internal/conversion"
	i "truelayer/internal/interfaces"
	"truelayer/internal/models"
)

type PokemonService struct {
	cache  i.PokemonCache
	client i.PokemonApiClient
}

func NewPokemonService(cache i.PokemonCache, client i.PokemonApiClient) *PokemonService {
	return &PokemonService{
		cache:  cache,
		client: client,
	}
}

func (s *PokemonService) GetPokemonSpecie(idOrName string) (*models.Pokemon, error) {
	var pokemonData *clients.PokémonData

	// checking the cache, otherwise using the client
	if cachedData, found := s.cache.Get(idOrName); found {
		pokemonData = cachedData
	} else {
		data, err := s.client.GetSpecie(idOrName)
		if err != nil {
			return nil, err
		}
		s.cache.Add(idOrName, data)

		pokemonData = data
	}

	pokemon := conversion.ConvertPokemonData(pokemonData)

	return pokemon, nil
}
