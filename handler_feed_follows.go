package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/k124k3n/rssagg/internal/database"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create feed follow: %s", err))
		return
	}

	respondWithJSON(w, 201, databaseFeedFollowToFeedFollow(feed))
}

func (apiCfg *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {

	feed, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get feed follow: %s", err))
		return
	}

	respondWithJSON(w, 200, databaseFeedFollowsToFeedFollows(feed))
}

func (apiCfg *apiConfig) handlerDeleteFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedParams := chi.URLParam(r, "feedFollowID")

	feedFollow, err := uuid.Parse(feedParams)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't parse feedFollowID: %v", err))
		return
	}
	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollow,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't dlete feedFollowID: %v", err))
		return
	}

	respondWithJSON(w, 204, struct{}{})
}
