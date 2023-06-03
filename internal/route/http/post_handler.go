package route

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/Totus-Floreo/asperitas-on-go/internal/application"
	"github.com/Totus-Floreo/asperitas-on-go/internal/model"
	"github.com/Totus-Floreo/asperitas-on-go/internal/route/helpers"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type PostHandler struct {
	Logger      *zap.SugaredLogger
	PostService *application.PostService
}

func (h *PostHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	posts, err := h.PostService.GetAllPosts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helpers.SendResponse(w, http.StatusOK, posts)
}

func (h *PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	postID, found := vars["postID"]
	if !found {
		http.Error(w, model.ErrPostInvalidHTTP.Error(), http.StatusUnprocessableEntity)
		return
	}

	post, err := h.PostService.GetPostByID(r.Context(), postID)
	if err != nil {
		http.Error(w, helpers.HTTPError(err), http.StatusNotFound)
		return
	}

	helpers.SendResponse(w, http.StatusOK, post)
}

func (h *PostHandler) GetPostsByCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	postCategory, found := vars["category"]
	if !found {
		http.Error(w, model.ErrPostCategoryInvalidHTTP.Error(), http.StatusUnprocessableEntity)
		return
	}

	posts, err := h.PostService.GetPostsByCategory(r.Context(), postCategory)
	if err != nil {
		http.Error(w, helpers.HTTPError(err), http.StatusInternalServerError)
		return
	}

	helpers.SendResponse(w, http.StatusOK, posts)
}

func (h *PostHandler) GetPostsByUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	userName, found := vars["user"]
	if !found {
		http.Error(w, model.ErrUserInvalidHTTP.Error(), http.StatusUnprocessableEntity)
		return
	}

	posts, err := h.PostService.GetPostsByUser(r.Context(), userName)
	if err != nil {
		http.Error(w, helpers.HTTPError(err), http.StatusInternalServerError)
		return
	}

	helpers.SendResponse(w, http.StatusOK, posts)
}

func (h *PostHandler) AddPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	post := model.NewPost()
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, helpers.HTTPError(err), http.StatusBadRequest)
		return
	}
	if post.Text == "" || post.Title == "" {
		http.Error(w, helpers.HTTPError(model.ErrPostInvalidHTTP), http.StatusBadRequest)
		return
	}

	response, err := h.PostService.AddPost(r.Context(), post)
	if err == model.ErrInvalidUrl {
		msg, err := model.NewErrorStack("body", "url", post.Url, "is invalid")
		if err != nil {
			http.Error(w, helpers.HTTPError(err), http.StatusUnprocessableEntity)
			return
		}
		http.Error(w, msg, http.StatusUnprocessableEntity)
		return
	}
	if err != nil {
		http.Error(w, helpers.HTTPError(err), http.StatusBadRequest)
		return
	}

	helpers.SendResponse(w, http.StatusCreated, response)
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	postID, found := vars["id"]
	if !found {
		http.Error(w, model.ErrPostInvalidHTTP.Error(), http.StatusUnprocessableEntity)
		return
	}

	err := h.PostService.DeletePost(r.Context(), postID)
	if err != nil {
		http.Error(w, helpers.HTTPError(err), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"success"}`))
}

func (h *PostHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	postID, found := vars["id"]
	if !found {
		http.Error(w, model.ErrPostInvalidHTTP.Error(), http.StatusUnprocessableEntity)
		return
	}

	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	comment, ok := data["comment"]
	if !ok {
		msg, err := model.NewErrorStack("body", "comment", "", "is required")
		if err != nil {
			http.Error(w, helpers.HTTPError(err), http.StatusInternalServerError)
			return
		}
		http.Error(w, msg, http.StatusUnprocessableEntity)
		return
	}

	post, err := h.PostService.AddComment(r.Context(), postID, comment)
	if err == model.ErrCommentTooLong {
		msg, err := model.NewErrorStack("body", "comment", comment, "must be at most 2000 characters long")
		if err != nil {
			http.Error(w, helpers.HTTPError(err), http.StatusUnprocessableEntity)
			return
		}
		http.Error(w, msg, http.StatusUnprocessableEntity)
		return
	}
	if err != nil {
		http.Error(w, helpers.HTTPError(err), http.StatusNotFound)
		return
	}

	helpers.SendResponse(w, http.StatusCreated, post)
}

func (h *PostHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	postID, found := vars["postID"]
	if !found {
		http.Error(w, model.ErrPostInvalidHTTP.Error(), http.StatusUnprocessableEntity)
		return
	}

	commentID, found := vars["commentID"]
	if !found {
		http.Error(w, model.ErrCommentInvalidHTTP.Error(), http.StatusUnprocessableEntity)
		return
	}

	post, err := h.PostService.DeleteComment(r.Context(), postID, commentID)
	if err != nil {
		http.Error(w, helpers.HTTPError(err), http.StatusNotFound)
		return
	}

	helpers.SendResponse(w, http.StatusOK, post)
}

func (h *PostHandler) Vote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	postID, found := vars["postID"]
	if !found {
		http.Error(w, model.ErrPostInvalidHTTP.Error(), http.StatusUnprocessableEntity)
		return
	}

	post, err := h.PostService.Vote(r.Context(), postID, filepath.Base(filepath.Clean(r.URL.Path)))
	if err != nil {
		http.Error(w, helpers.HTTPError(err), http.StatusNotFound)
		return
	}

	helpers.SendResponse(w, http.StatusOK, post)
}
