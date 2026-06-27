//go:build unit

package conversion

import (
	"reflect"
	"testing"

	"truelayer/internal/clients"
	"truelayer/internal/models"
)

func TestConvertPokemonData(t *testing.T) {
	tests := []struct {
		name     string
		input    *clients.PokémonData
		expected *models.Pokemon
	}{
		{
			name:     "Nil case",
			input:    nil,
			expected: nil,
		},
		{
			name: "No flavor text entries present",
			input: &clients.PokémonData{
				Name:              "mewtwo",
				FlavorTextEntries: []clients.FlavorTextEntry{},
				Habitat:           clients.Habitat{Name: "rare"},
				IsLegendary:       true,
			},
			expected: &models.Pokemon{
				Name:        "mewtwo",
				Description: "",
				Habitat:     "rare",
				IsLegendary: true,
			},
		},
		{
			name: "English description not present (only other languages)",
			input: &clients.PokémonData{
				Name: "pikachu",
				FlavorTextEntries: []clients.FlavorTextEntry{
					{
						FlavorText: "Pika Pika!",
						Language:   clients.NamedAPIResource{Name: "fr"},
					},
					{
						FlavorText: "Pika Chu!",
						Language:   clients.NamedAPIResource{Name: "ja"},
					},
				},
				Habitat:     clients.Habitat{Name: "forest"},
				IsLegendary: false,
			},
			expected: &models.Pokemon{
				Name:        "pikachu",
				Description: "",
				Habitat:     "forest",
				IsLegendary: false,
			},
		},
		{
			name: "Character replacement and trimming in description",
			input: &clients.PokémonData{
				Name: "charizard",
				FlavorTextEntries: []clients.FlavorTextEntry{
					{
						FlavorText: "\n Spits fire\fthat is hot\t",
						Language:   clients.NamedAPIResource{Name: "en"},
					},
				},
				Habitat:     clients.Habitat{Name: "mountain"},
				IsLegendary: false,
			},
			expected: &models.Pokemon{
				Name:        "charizard",
				Description: "Spits fire that is hot",
				Habitat:     "mountain",
				IsLegendary: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ConvertPokemonData(tt.input)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("ConvertPokemonData() = %v, want %v", actual, tt.expected)
			}
		})
	}
}
