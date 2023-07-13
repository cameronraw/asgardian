package asgardian

import (
	"net/http"
)

func (app *Application) healthCheck(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"status":      "available",
		"version":     "1.0.0",
		"environment": "dev",
	}

	err := WriteJSON(w, r, http.StatusOK, Envelope{"response": data})
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
}

