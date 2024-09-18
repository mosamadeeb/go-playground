package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/mosamadeeb/chirpy/internal/chirpydb"
)

func (s serverState) handleWebhooks() {
	s.Mux.HandleFunc("POST /api/polka/webhooks", func(w http.ResponseWriter, r *http.Request) {
		apiKey, ok := strings.CutPrefix(r.Header.Get("Authorization"), "ApiKey ")
		if !ok || apiKey != s.ApiCfg.polkaApi {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var req struct {
			Event string `json:"event"`
			Data  struct {
				UserId int `json:"user_id"`
			} `json:"data"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Error decoding request body: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch req.Event {
		case "user.upgraded":
			if err := s.DB.SetUserChirpyRed(req.Data.UserId, true); err != nil {
				if errors.Is(err, chirpydb.ErrNotExist) {
					w.WriteHeader(http.StatusNotFound)
					return
				}

				log.Printf("Error updating user in database: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusNoContent)
		}
	})
}
