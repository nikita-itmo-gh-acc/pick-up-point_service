package main

import (
	"fmt"
	"net/http"
	"os"
	"pvz_service/database"
	"pvz_service/handlers"
	"pvz_service/repos"
	"pvz_service/services"

	"github.com/gorilla/mux"
	_ "github.com/gorilla/schema"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
    	
    }
}

func main() {
	var (
		dbUrl = os.Getenv("DB_URL")
		baseUrl = os.Getenv("POSTGRES_URL")
		dbName = os.Getenv("DB_NAME")
		port = os.Getenv("SERVE_PORT")
		migrationsDir = os.Getenv("MIGRATION_DIR")
	)

	conn := &database.DBConnection{
		DbName: dbName,
		URL: dbUrl,
		BaseURL: baseUrl,
	}

	migrator := &database.Migrator{}

	if err := conn.InitPostgresConn(); err != nil {
		fmt.Println("Can't create database connection: ", err)
		return 
	}

	if err := migrator.Init(migrationsDir, dbUrl); err != nil {
		fmt.Println("Can't create migrator: ", err)
		return
	}

	if err := migrator.Apply(); err != nil {
		return
	}

	defer conn.DB.Close()

	router := mux.NewRouter()

	userRepo := repos.NewUserRepo(conn.DB)
	pvzRepo := repos.NewPvzRepo(conn.DB)
	productRepo := repos.NewProductRepo(conn.DB)
	receptionRepo := repos.NewReceptionRepo(conn.DB)

	userService := services.NewUserService(userRepo)
	pvzService := services.NewPvzService(pvzRepo)
	productService := services.NewProductService(productRepo, receptionRepo)
	receptionService := services.NewReceptionService(receptionRepo)

	userHandler := handlers.NewUserHandler(userService)
	userHandler.SetUpRoutes(router)

	pvzRoutes := router.PathPrefix("/pvz").Subrouter()
	pvzHandler := handlers.NewPvzHandler(pvzService, receptionService, productService)
	pvzHandler.SetUpRoutes(pvzRoutes)

	Employee := handlers.RequireRole("employee")
	productsHandler := handlers.NewProductHandler(productService)
	handleProducts := Employee(handlers.AuthMiddleware(http.HandlerFunc(productsHandler.AddHandler)))
	router.Handle("/products", handleProducts).Methods("POST")

	receptionHandler := handlers.NewReceptionHandler(receptionService)
	handleReceptions := Employee(handlers.AuthMiddleware(http.HandlerFunc(receptionHandler.CreateHandler)))
	router.Handle("/receptions", handleReceptions).Methods("POST")

	http.ListenAndServe(":" + port, router)
}
