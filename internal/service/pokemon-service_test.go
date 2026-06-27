//go:build unit

package service

import (
	"errors"
	"testing"

	"truelayer/internal/clients"
)

type mockCache struct {
	getFunc func(key string) (*clients.PokémonData, bool)
	addFunc func(key string, value *clients.PokémonData)
}

func (m *mockCache) Get(key string) (*clients.PokémonData, bool) { return m.getFunc(key) }
func (m *mockCache) Add(key string, value *clients.PokémonData)  { m.addFunc(key, value) }

type mockApiClient struct {
	getSpecieFunc func(idOrName string) (*clients.PokémonData, error)
}

func (m *mockApiClient) GetSpecie(idOrName string) (*clients.PokémonData, error) {
	return m.getSpecieFunc(idOrName)
}

func TestPokemonService_GetPokemonSpecie(t *testing.T) {
	apiData := &clients.PokémonData{
		Name:    "mewtwo",
		Habitat: clients.Habitat{Name: "rare"},
		FlavorTextEntries: []clients.FlavorTextEntry{
			{FlavorText: "It was created by\na scientist.", Language: clients.NamedAPIResource{Name: "en"}},
		},
		IsLegendary: true,
	}

	t.Run("cache hit", func(t *testing.T) {
		cacheCalled := false
		clientCalled := false

		mockC := &mockCache{
			getFunc: func(key string) (*clients.PokémonData, bool) {
				cacheCalled = true
				return apiData, true
			},
			addFunc: func(key string, value *clients.PokémonData) {},
		}

		mockCli := &mockApiClient{
			getSpecieFunc: func(idOrName string) (*clients.PokémonData, error) {
				clientCalled = true
				return nil, nil
			},
		}

		svc := &PokemonService{cache: mockC, client: mockCli}
		result, err := svc.GetPokemonSpecie("mewtwo")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !cacheCalled {
			t.Error("expected cache to be checked")
		}
		if clientCalled {
			t.Error("expected API client NOT to be called on cache hit")
		}
		if result.Name != "mewtwo" || !result.IsLegendary {
			t.Errorf("unexpected conversion output: %+v", result)
		}
	})

	t.Run("cache miss", func(t *testing.T) {
		clientCalled := false
		cacheAdded := false

		mockC := &mockCache{
			getFunc: func(key string) (*clients.PokémonData, bool) {
				return nil, false
			},
			addFunc: func(key string, value *clients.PokémonData) {
				cacheAdded = true
			},
		}

		mockCli := &mockApiClient{
			getSpecieFunc: func(idOrName string) (*clients.PokémonData, error) {
				clientCalled = true
				return apiData, nil
			},
		}

		svc := &PokemonService{cache: mockC, client: mockCli}
		result, err := svc.GetPokemonSpecie("mewtwo")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !clientCalled {
			t.Error("expected API client to be called on cache miss")
		}
		if !cacheAdded {
			t.Error("expected fetched data to be saved to the cache")
		}
		if result.Description != "It was created by a scientist." {
			t.Errorf("expected clean description, got %q", result.Description)
		}
	})

	t.Run("cache miss and error", func(t *testing.T) {
		mockC := &mockCache{
			getFunc: func(key string) (*clients.PokémonData, bool) {
				return nil, false
			},
			addFunc: func(key string, value *clients.PokémonData) {
				t.Error("should not add to cache if client fails")
			},
		}

		expectedErr := errors.New("api failure")
		mockCli := &mockApiClient{
			getSpecieFunc: func(idOrName string) (*clients.PokémonData, error) {
				return nil, expectedErr
			},
		}

		svc := &PokemonService{cache: mockC, client: mockCli}
		result, err := svc.GetPokemonSpecie("mewtwo")

		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
		if result != nil {
			t.Errorf("expected result to be nil on error, got %+v", result)
		}
	})
}
