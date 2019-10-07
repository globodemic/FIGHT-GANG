package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	authorization "../auths"
	db "../database"
	model "../models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
)

//CreatePlayer to registering new players
func CreatePlayer(w http.ResponseWriter, r *http.Request) {

	var player model.Player
	_ = json.NewDecoder(r.Body).Decode(&player)

	if player.Name == "" || player.Alias == "" || player.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resultID := db.CreatePlayer(player)

	if resultID == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	validToken, err := authorization.GenerateJWT(player.Name, player.Password)

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("X-API-Token", validToken)
	json.NewEncoder(w).Encode(resultID)
}

//GetPlayer /api/player/{id}
func GetPlayer(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	if len(params) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	decoded := r.Context().Value("player")

	var cred model.Credentials
	mapstructure.Decode(decoded.(jwt.MapClaims), &cred)
	resultPlayer := db.GetPlayer(cred, params)

	if resultPlayer == (model.Player{}) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resultPlayer)
}

//GetPlayers /api/players
func GetPlayers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	players := db.GetPlayers()
	if len(players) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(players)
}

//DeletePlayer /api/player/{id}
func DeletePlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//Get Player ID
	params := mux.Vars(r)
	if len(params) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	decoded := r.Context().Value("player")

	var cred model.Credentials
	mapstructure.Decode(decoded.(jwt.MapClaims), &cred)
	deletePlayer := db.DeletePlayer(params)
	if deletePlayer == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

//PutPlayer update Player's alias and password /api/player/{id}
func PutPlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	if len(params) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	decoded := r.Context().Value("player")
	var player model.Player
	_ = json.NewDecoder(r.Body).Decode(&player)

	var cred model.Credentials
	mapstructure.Decode(decoded.(jwt.MapClaims), &cred)
	putPlayer := db.ChangePlayerAuth(cred, params, player)
	if putPlayer == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//refresh another token
	refreshToken, err := authorization.GenerateJWT(cred.Name, player.Password)

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("X-API-Token", refreshToken)

	w.WriteHeader(http.StatusOK)
}
