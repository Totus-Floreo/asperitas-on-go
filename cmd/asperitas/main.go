package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Totus-Floreo/asperitas-on-go/internal/application"
	"github.com/Totus-Floreo/asperitas-on-go/internal/middleware"
	mongo_repository "github.com/Totus-Floreo/asperitas-on-go/internal/repository/mongo"
	pgx_repository "github.com/Totus-Floreo/asperitas-on-go/internal/repository/pgx"
	redis_repository "github.com/Totus-Floreo/asperitas-on-go/internal/repository/redis"
	route "github.com/Totus-Floreo/asperitas-on-go/internal/route/http"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()

	// TODO: commandline settings for choose needed method of storage info
	// method := flag.String("method", "db", "Specify the storage method")
	// flag.Parse()
	// if *method != "db" && *method != "inmemory" {
	// 	flag.Usage()
	// 	logger.Fatalln("Invalid method specified, use db or im")
	// }
	// userRepository := inmemory.NewUserStorage()
	// postStorage := inmemory.NewPostStorage()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost" + os.Getenv("redis"),
		Password: "",
		DB:       0,
	})
	if redisStatus := rdb.Ping(context.Background()); redisStatus.Err() != nil {
		logger.Panicln("Redis connection error: ", redisStatus.Err().Error())
	}

	postgreUrl := "postgres://" + os.Getenv("pg_uri")
	pgxdb, err := pgxpool.New(context.Background(), postgreUrl)
	if err != nil {
		logger.Panicln("Postgre connection error: ", err.Error())
	}

	poolScheduler, err := mongo_repository.NewDBReadersPool(os.Getenv("mongo_uri"), 10)
	if err != nil {
		logger.Panicln("Mongo connection error: ", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(os.Getenv("mongo_uri"))
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	userRepository := pgx_repository.NewUserStorage(pgxdb)
	tokenRepository := redis_repository.NewTokenRepository(rdb)
	JWTService := application.NewJWTService(os.Getenv("signature"))
	authService := application.NewAuthService(userRepository, tokenRepository, JWTService)

	userHandler := &route.UserHandler{
		Logger:      logger,
		AuthService: authService,
	}

	postStorage := mongo_repository.NewPostStorage(client, poolScheduler)
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
