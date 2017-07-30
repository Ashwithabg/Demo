package handlers

import (
	"fmt"
	"net/http"
	"encoding/json"

	"ashwitha/workspace/demo/api"
)

func GetReposHandler(api api.GithubAPI) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo, err := api.GetRepo(r.Context())
		if err != nil {
			fmt.Printf("GetRepo: unable to get github repositories: %s", err)
			http.Error(w, "Failed to get github repositories", 500)
			return
		}

		err = json.NewEncoder(w).Encode(repo)
		if err != nil {
			fmt.Printf("GetRepo: unable to write to response writer: %s",err)
			http.Error(w, "unable to write to response writer", 500)
			return
		}

	}
}