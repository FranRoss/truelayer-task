//go:build unit

package clients

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"truelayer/internal/utils"
)

func setupPokeAPIMockServer(handler http.HandlerFunc) (*httptest.Server, *PokeApiClient) {
	server := httptest.NewServer(handler)
	client := &PokeApiClient{
		baseURL:    server.URL,
		httpClient: server.Client(),
	}
	return server, client
}

func TestGetSpecie_Success(t *testing.T) {
	expectedPokemon := &PokémonData{
		Name:        "mewtwo",
		IsLegendary: true,
		Habitat:     Habitat{Name: "rare"},
		FlavorTextEntries: []FlavorTextEntry{
			{
				FlavorText: "It was created by a scientist.",
				Language:   NamedAPIResource{Name: "en"},
			},
		},
	}

	server, client := setupPokeAPIMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.WriteJsonResponse(w, expectedPokemon)
	})
	defer server.Close()

	result, err := client.GetSpecie("mewtwo")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result.Name != expectedPokemon.Name {
		t.Errorf("expected name %s, got %s", expectedPokemon.Name, result.Name)
	}
	if result.IsLegendary != expectedPokemon.IsLegendary {
		t.Errorf("expected legendary to be %t", expectedPokemon.IsLegendary)
	}
}

func TestGetSpecie_NetworkFailureError(t *testing.T) {
	server, client := setupPokeAPIMockServer(func(w http.ResponseWriter, r *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			t.Fatal("webserver doesn't support hijacking")
		}
		conn, _, _ := hj.Hijack()
		conn.Close()
	})
	defer server.Close()

	_, err := client.GetSpecie("mew")
	if err == nil {
		t.Fatal("expected a network failure error, but got nil")
	}

	if !strings.Contains(err.Error(), "failed making request to PokeAPI") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestGetSpecie_NotFound(t *testing.T) {
	server, client := setupPokeAPIMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	defer server.Close()

	_, err := client.GetSpecie("missingno")
	if err == nil {
		t.Fatal("expected a 404 error, but got nil")
	}

	expectedErr := "pokemon specie 'missingno' not found"
	if err.Error() != expectedErr {
		t.Errorf("expected error %q, got %q", expectedErr, err.Error())
	}
}

func TestGetSpecie_UnexpectedStatusCode(t *testing.T) {
	server, client := setupPokeAPIMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer server.Close()

	_, err := client.GetSpecie("ditto")
	if err == nil {
		t.Fatal("expected an unexpected status code error, but got nil")
	}

	expectedErr := fmt.Sprintf("unexpected status code from PokeAPI: %d", http.StatusInternalServerError)
	if err.Error() != expectedErr {
		t.Errorf("expected error %q, got %q", expectedErr, err.Error())
	}
}

func TestGetSpecie_InvalidJSON(t *testing.T) {
	server, client := setupPokeAPIMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{invalid-json`))
	})
	defer server.Close()

	_, err := client.GetSpecie("pikachu")
	if err == nil {
		t.Fatal("expected a JSON decoding error, but got nil")
	}

	if !strings.Contains(err.Error(), "failed decoding JSON response") {
		t.Errorf("unexpected error message: %v", err)
	}
}
