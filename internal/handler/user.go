package handler

import (
	"encoding/json"
	"net/http"
	"sodam/internal/service"

	"github.com/matryer/way"
)

type createUserInput struct {
	Email, Username string
}

func (h *handler) createUser(w http.ResponseWriter, r *http.Request) {
	var in createUserInput
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.CreateUser(r.Context(), in.Email, in.Username)
	if err == service.ErrInvalidEmail || err == service.ErrInvalidUsername {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err == service.ErrEmailTaken || err == service.ErrUsernameTaken {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	if err != nil {
		respondError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// 팔로우 핸들러
func (h *handler) toggleFollow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := way.Param(ctx, "username")

	out, err := h.ToggleFollow(ctx, username)
	// 인증되지 않은 유저일 때
	if err == service.ErrUnauthenticated {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// 유저 이름이 맞지 않을 때
	if err == service.ErrInvalidUsername {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	// 유저를 찾을 수 없을 때
	if err == service.ErrUserNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// 스스소를 팔로우 할 때
	if err != nil {
		respondError(w, err)
		return
	}

	respond(w, out, http.StatusOK)
}
