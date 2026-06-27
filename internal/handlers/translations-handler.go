package handlers

import (
	"log"
	"net/http"

	c "truelayer/internal/cache"
	"truelayer/internal/clients"
	i "truelayer/internal/interfaces"
	s "truelayer/internal/service"
	"truelayer/internal/utils"
)

type TranslationsHandler struct {
	service i.PokemonService
	client  i.TranslationClient
}

func NewTranslationsHandler(
	cache *c.Cache[clients.PokémonData],
) *TranslationsHandler {
	return &TranslationsHandler{
		service: s.NewPokemonService(cache, clients.NewPokeApiClient()),
		client:  clients.NewTranslationClient(),
	}
}

func (p *TranslationsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	pokemon, err := p.service.GetPokemonSpecie(name)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "error retrieving the pokemon data")
		return
	}

	var translationType = clients.Shakespeare
	if utils.IsYodaTranslation(pokemon) {
		translationType = clients.Yoda
	}
	newDescription, err := p.client.Translate(pokemon.Description, translationType)
	if err != nil {
		log.Printf("error while translating: %s", err)
	} else {
		pokemon.Description = newDescription
	}

	utils.WriteJsonResponse(w, pokemon)
}
