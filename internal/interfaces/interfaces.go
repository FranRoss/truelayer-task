package interfaces

import (
	"truelayer/internal/clients"
	"truelayer/internal/models"
)

type PokemonCache interface {
	Get(key string) (*clients.PokémonData, bool)
	Add(key string, value *clients.PokémonData)
}

type PokemonApiClient interface {
	GetSpecie(idOrName string) (*clients.PokémonData, error)
}

type PokemonService interface {
	GetPokemonSpecie(idOrName string) (*models.Pokemon, error)
}

type TranslationClient interface {
	Translate(text string, translationType clients.TranslationType) (string, error)
}
