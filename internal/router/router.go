package router

import (
	"net/http"

	c "truelayer/internal/cache"
	"truelayer/internal/clients"
	"truelayer/internal/handlers"
)

func New(cache *c.Cache[clients.PokémonData]) *http.ServeMux {
	mux := http.NewServeMux()

	pokemonHandler := handlers.NewPokemonHandler(cache)
	translationsHandler := handlers.NewTranslationsHandler(cache)

	mux.HandleFunc("GET /health", handlers.HealthHandler)
	mux.HandleFunc("GET /pokemon/{name}", pokemonHandler.Handle)
	mux.HandleFunc("GET /pokemon/translated/{name}", translationsHandler.Handle)

	return mux
}
