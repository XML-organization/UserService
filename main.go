package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"user_service/handler"
	"user_service/model"
	"user_service/repository"
	"user_service/service"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	connectionStr := "host=localhost user=postgres password=password dbname=UserServiceDatabase port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(connectionStr), &gorm.Config{})
	if err != nil {
		print(err)
		return nil
	}

	database.AutoMigrate(&model.User{}, &model.Address{})

	return database
}

func GetClient(host, user, password, dbname, port string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func initPostgresClient() *gorm.DB {
	client, err := GetClient(
		os.Getenv("USER_DB_HOST"), os.Getenv("USER_DB_USER"),
		os.Getenv("USER_DB_PASS"), os.Getenv("USER_DB_NAME"),
		os.Getenv("USER_DB_PORT"))
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func startServer(userHandler *handler.UserHandler) {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/update", userHandler.UpdateUser).Methods("PUT")
	router.HandleFunc("/changePassword", userHandler.ChangePassword).Methods("PUT")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://192.168.0.17:5173", "http://localhost:5173", "http://192.168.137.1:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	handler1 := corsHandler.Handler(router)

	router.Methods("OPTIONS").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.WriteHeader(http.StatusNoContent)
		})

	println("Server starting")
	log.Fatal(http.ListenAndServe(":8082", handler1))
}

func main() {
	//database := initDB()
	database := initPostgresClient()

	if database == nil {
		log.Fatal("FAILED TO CONNECT TO DB")
	}
	repoUser := &repository.UserRepository{DatabaseConnection: database}
	serviceUser := &service.UserService{UserRepo: repoUser}
	handlerUser := &handler.UserHandler{UserService: serviceUser}

	startServer(handlerUser)
}
