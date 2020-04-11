package handler

import (
	"encoding/json"
	"github.com/rithikjain/LiveQnA/api/middleware"
	"github.com/rithikjain/LiveQnA/api/view"
	"github.com/rithikjain/LiveQnA/api/websocket"
	"github.com/rithikjain/LiveQnA/pkg/question"
	"net/http"
	"strconv"
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

		// Getting the user from claims
		claims, err := middleware.ValidateAndGetClaims(r.Context(), "user")
		if err != nil {
			view.Wrap(err, w)
			return
		}
		user, err := s.GetUser(claims["id"].(float64))

		// Setting email in question as user email
		queRequest.CreatedByEmail = user.Email

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

// Protected Request
func viewAllQuestions(s question.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			view.Wrap(view.ErrMethodNotAllowed, w)
			return
		}

		questions, err := s.ViewAllQuestions()
		if err != nil {
			view.Wrap(err, w)
			return
		}
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message":   "All questions fetched",
			"questions": questions,
			"status":    http.StatusOK,
		})
	})
}

// Protected Request
func increaseUpVote(s question.Service, hub *websocket.Hub) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			view.Wrap(view.ErrMethodNotAllowed, w)
			return
		}

		// Getting the user from claims
		claims, err := middleware.ValidateAndGetClaims(r.Context(), "user")
		if err != nil {
			view.Wrap(err, w)
			return
		}
		user, err := s.GetUser(claims["id"].(float64))

		// Get questionID from url
		questionIDStr := r.URL.Query().Get("questionID")
		if questionIDStr != "" {
			questionID, _ := strconv.Atoi(questionIDStr)
			que, err := s.IncreaseUpVote(float64(questionID), user)
			if err != nil {
				view.Wrap(err, w)
				return
			}
			// Send message on websocket channel
			hub.Broadcast <- que
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Successfully Up voted",
				"status":  http.StatusOK,
			})
		} else {
			view.Wrap(view.ErrNoParameter, w)
			return
		}
	})
}

// Handler
func MakeQuestionHandler(r *http.ServeMux, s question.Service, hub *websocket.Hub) {
	r.Handle("/api/question/create", middleware.Validate(sendQuestion(s, hub)))
	r.Handle("/api/question/view", middleware.Validate(viewAllQuestions(s)))
	r.Handle("/api/question/upvote", middleware.Validate(increaseUpVote(s, hub)))
}
