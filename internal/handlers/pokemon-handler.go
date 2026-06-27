package handlers

import (
	"net/http"
	c "truelayer/internal/cache"
	"truelayer/internal/clients"
	s "truelayer/internal/service"
	"truelayer/internal/utils"
)

type PokemonHandler struct {
	service s.PokemonService
}

func NewPokemonHandler(cache *c.Cache[clients.PokémonData]) *PokemonHandler {
	return &PokemonHandler{
		service: *s.NewPokemonService(cache, clients.NewPokeApiClient()),
	}
}

func (p *PokemonHandler) Handle(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	pokemon, err := p.service.GetPokemonSpecie(name)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "error retrieving the pokemon data")
		return
	}

	utils.WriteJsonResponse(w, pokemon)
}
