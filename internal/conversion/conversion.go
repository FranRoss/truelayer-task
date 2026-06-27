package conversion

import (
	"strings"

	"truelayer/internal/clients"
	"truelayer/internal/models"
)

func ConvertPokemonData(src *clients.PokémonData) *models.Pokemon {
	if src == nil {
		return nil
	}

	var description string
	for _, entry := range src.FlavorTextEntries {
		if entry.Language.Name == "en" {
			description = cleanFlavorText(entry.FlavorText)
			break
		}
	}

	return &models.Pokemon{
		Name:        src.Name,
		Description: description,
		Habitat:     src.Habitat.Name,
		IsLegendary: src.IsLegendary,
	}
}

// assuming the additional text only contains this chars as seen in a few responses from the pokemon apis,
// I'm unaware if more need to be sanitized as I couldn't find more examples.
func cleanFlavorText(text string) string {
	r := strings.NewReplacer("\n", " ", "\f", " ", "\t", " ")
	return strings.TrimSpace(r.Replace(text))
}
