package handler

import (
	"encoding/json"
	"github.com/rithikjain/LiveQnA/api/middleware"
	"github.com/rithikjain/LiveQnA/api/view"
	"github.com/rithikjain/LiveQnA/api/websocket"
	"github.com/rithikjain/LiveQnA/pkg/question"
	"net/http"
)

// Protected Request
func sendQuestion(s question.Service, hub *websocket.Hub) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			view.Wrap(view.ErrMethodNotAllowed, w)
			return
		}

		queRequest := &question.Question{}
		if err := json.NewDecoder(r.Body).Decode(queRequest); err != nil {
			view.Wrap(err, w)
			return
		}

		que, err := s.CreateQuestion(queRequest)
		if err != nil {
			view.Wrap(err, w)
			return
		}

		// Send message on websocket channel
		hub.Broadcast <- que
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Successfully created question",
			"status":  http.StatusOK,
		})
	})
}

// Handler
func MakeQuestionHandler(r *http.ServeMux, s question.Service, hub *websocket.Hub) {
	r.Handle("/api/question/create", middleware.Validate(sendQuestion(s, hub)))
}
