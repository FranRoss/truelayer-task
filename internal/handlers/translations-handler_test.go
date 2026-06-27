//go:build unit

package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"truelayer/internal/clients"
	"truelayer/internal/models"
)

type MockPokemonService struct {
	GetPokemonSpecieFunc func(name string) (*models.Pokemon, error)
}

func (m *MockPokemonService) GetPokemonSpecie(name string) (*models.Pokemon, error) {
	return m.GetPokemonSpecieFunc(name)
}

type MockTranslationClient struct {
	TranslateFunc func(text string, translationType clients.TranslationType) (string, error)
}

func (m *MockTranslationClient) Translate(text string, translationType clients.TranslationType) (string, error) {
	return m.TranslateFunc(text, translationType)
}

func TestTranslationsHandler_Handle(t *testing.T) {
	tests := []struct {
		name               string
		mockService        func(name string) (*models.Pokemon, error)
		mockClient         func(text string, translationType clients.TranslationType) (string, error)
		expectedStatusCode int
	}{
		{
			name: "GetPokemonSpecie returns an error",
			mockService: func(name string) (*models.Pokemon, error) {
				return &models.Pokemon{}, errors.New("database or API down")
			},
			mockClient:         nil, // Shouldn't be called
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Successful Yoda Translation",
			mockService: func(name string) (*models.Pokemon, error) {
				return &models.Pokemon{Name: "mewtwo", Description: "A powerful pokemon.", Habitat: "rare", IsLegendary: true}, nil
			},
			mockClient: func(text string, translationType clients.TranslationType) (string, error) {
				if translationType != clients.Yoda {
					t.Errorf("expected Yoda translation, got %v", translationType)
				}
				return "A powerful pokemon, it is.", nil
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Successful Shakespeare Translation",
			mockService: func(name string) (*models.Pokemon, error) {
				return &models.Pokemon{Name: "pidgey", Description: "A small bird.", Habitat: "forest", IsLegendary: false}, nil
			},
			mockClient: func(text string, translationType clients.TranslationType) (string, error) {
				if translationType != clients.Shakespeare {
					t.Errorf("expected Shakespeare translation, got %v", translationType)
				}
				return "A small bird, eke.", nil
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Translation fails",
			mockService: func(name string) (*models.Pokemon, error) {
				return &models.Pokemon{Name: "pidgey", Description: "A small bird.", Habitat: "forest", IsLegendary: false}, nil
			},
			mockClient: func(text string, translationType clients.TranslationType) (string, error) {
				return "", errors.New("translation API rate limited")
			},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewTranslationsHandler(nil)

			handler.service = &MockPokemonService{GetPokemonSpecieFunc: tt.mockService}
			handler.client = &MockTranslationClient{TranslateFunc: tt.mockClient}

			req := httptest.NewRequest(http.MethodGet, "/pokemon/charizard", nil)
			req.SetPathValue("name", "charizard")

			rec := httptest.NewRecorder()

			handler.Handle(rec, req)

			if rec.Code != tt.expectedStatusCode {
				t.Errorf("expected status code %d, got %d", tt.expectedStatusCode, rec.Code)
			}
		})
	}
}
