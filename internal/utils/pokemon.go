package utils

import (
	"truelayer/internal/models"
)

const _POKEMON_HABITAT_CAVE = "cave"

func IsYodaTranslation(pokemon *models.Pokemon) bool {
	return pokemon.IsLegendary || pokemon.Habitat == _POKEMON_HABITAT_CAVE
}
