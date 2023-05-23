package route

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/Totus-Floreo/asperitas-on-go/pkg/application"
	"github.com/Totus-Floreo/asperitas-on-go/pkg/delivery/helpers"
	"github.com/Totus-Floreo/asperitas-on-go/pkg/middleware"
	"github.com/Totus-Floreo/asperitas-on-go/pkg/model"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type PostHandler struct {
	Logger      *zap.SugaredLogger
	PostService *application.PostService
}

func (h *PostHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	posts, err := h.PostService.GetAllPosts()
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

	post, err := h.PostService.GetPostByID(postID)
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

	posts, err := h.PostService.GetPostsByCategory(postCategory)
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

	posts, err := h.PostService.GetPostsByUser(userName)
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
	post.Author = r.Context().Value(middleware.AuthorContextKey).(*model.Author)
	post, err := h.PostService.AddPost(post)
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

	helpers.SendResponse(w, http.StatusCreated, post)
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	postID, found := vars["id"]
	if !found {
		http.Error(w, model.ErrPostInvalidHTTP.Error(), http.StatusUnprocessableEntity)
		return
	}
	author := r.Context().Value(middleware.AuthorContextKey).(*model.Author)
	err := h.PostService.DeletePost(postID, author)
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
		}
		http.Error(w, msg, http.StatusUnprocessableEntity)
		return
	}

	author := r.Context().Value(middleware.AuthorContextKey).(*model.Author)

	post, err := h.PostService.AddComment(postID, comment, author)
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

	author := r.Context().Value(middleware.AuthorContextKey).(*model.Author)

	post, err := h.PostService.DeleteComment(postID, commentID, author)
	if err != nil {
		http.Error(w, helpers.HTTPError(err), http.StatusBadRequest)
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

	author := r.Context().Value(middleware.AuthorContextKey).(*model.Author)

	post, err := h.PostService.Vote(postID, author, filepath.Base(filepath.Clean(r.URL.Path)))
	if err != nil {
		http.Error(w, helpers.HTTPError(err), http.StatusBadRequest)
		return
	}

	helpers.SendResponse(w, http.StatusOK, post)
}
