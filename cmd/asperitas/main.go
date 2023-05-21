package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Totus-Floreo/asperitas-on-go/pkg/application"
	route "github.com/Totus-Floreo/asperitas-on-go/pkg/delivery/http"
	"github.com/Totus-Floreo/asperitas-on-go/pkg/middleware"
	repository "github.com/Totus-Floreo/asperitas-on-go/pkg/repository/inmemory"
	repositoryRedis "github.com/Totus-Floreo/asperitas-on-go/pkg/repository/redis"
	"github.com/redis/go-redis/v9"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:32768",
		Password: "",
		DB:       0,
	})

	userRepository := repository.NewUserStorage()
	tokenRepository := repositoryRedis.NewTokenRepository(rdb)
	JWTService := application.NewJWTService(os.Getenv("signature"))
	authService := application.NewAuthService(userRepository, tokenRepository, JWTService)

	userHandler := &route.UserHandler{
		Logger:      logger,
		AuthService: authService,
	}

	postStorage := repository.NewPostStorage()
	postService := application.NewPostService(postStorage)

	postHandler := &route.PostHandler{
		Logger:      logger,
		PostService: postService,
	}

	router := mux.NewRouter()
	router.Use(middleware.Panic)
	router.Use(middleware.AccessLog(logger))
	router.PathPrefix("/static/").Handler(route.StaticHandler())
	router.HandleFunc("/", route.WebHandler)

	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/register", userHandler.SignUp).Methods("POST")
	api.HandleFunc("/login", userHandler.Login).Methods("POST")
	api.HandleFunc("/posts/", postHandler.GetAllPosts).Methods("GET")
	api.HandleFunc("/posts/{category}", postHandler.GetPostsByCategory).Methods("GET")
	api.HandleFunc("/post/{postID}", postHandler.GetPostByID).Methods("GET")
	api.HandleFunc("/user/{user}", postHandler.GetPostsByUser).Methods("GET")

	apiAuth := router.PathPrefix("/api").Subrouter()

	apiAuth.Use(middleware.Auth(JWTService, tokenRepository))
	apiAuth.HandleFunc("/posts", postHandler.AddPost).Methods("POST")
	apiAuth.HandleFunc("/post/{id}", postHandler.DeletePost).Methods("DELETE")
	apiAuth.HandleFunc("/post/{id}", postHandler.AddComment).Methods("POST")
	apiAuth.HandleFunc("/post/{postID}/{commentID}", postHandler.DeleteComment).Methods("DELETE")
	apiAuth.HandleFunc("/post/{postID}/upvote", postHandler.Vote).Methods("GET")
	apiAuth.HandleFunc("/post/{postID}/unvote", postHandler.Vote).Methods("GET")
	apiAuth.HandleFunc("/post/{postID}/downvote", postHandler.Vote).Methods("GET")

	router.NotFoundHandler = http.HandlerFunc(route.WebHandler)

	log.Println("Starting server on port" + os.Getenv("port"))
	http.ListenAndServe(os.Getenv("port"), router)
}
