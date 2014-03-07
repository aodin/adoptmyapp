package robot

import (
	"encoding/json"
	"net/http"
)

type API struct {
	db DB
}

// The sole HTTP function
// Can be tested with curl:
// curl localhost:8081 -d "url=git@github.com:kkochis/adoptmyapp.git"
func (api *API) Root(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Check for the parameter `url`
		rawurl := r.FormValue("url")

		// Normalize
		url, err := NormalizeGitHubURL(rawurl)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		// Remove the scheme: `https://`
		normalized := url.String()[8:]

		// Check if the url already exists
		meta, err := api.db.Get(normalized)
		if err == nil {
			// TODO JSON response?
			w.Header().Set("Content-Type", "application/json")
			// TODO JSON errors ignored
			b, _ := json.Marshal(meta)
			w.Write(b)
			return
		}
		// TODO What if it was a database error?

		// Create it if it does not
		meta, err = api.db.Save(normalized)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		// Return the meta information
		w.Header().Set("Content-Type", "application/json")
		// TODO JSON errors ignored
		b, _ := json.Marshal(meta)
		w.Write(b)
		return
	}
	w.Write([]byte("Not Implemented\n"))
}

// Start the API
func New(db DB) (*API, error) {
	api := &API{db: db}
	http.HandleFunc("/", api.Root)
	return api, nil
}
