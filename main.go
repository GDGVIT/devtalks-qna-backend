package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/rithikjain/LiveQnA/api/handler"
	"github.com/rithikjain/LiveQnA/api/middleware"
	"github.com/rithikjain/LiveQnA/api/websocket"
	"github.com/rithikjain/LiveQnA/pkg/question"
	"github.com/rithikjain/LiveQnA/pkg/user"
	"log"
	"net/http"
	"os"
)

func dbConnect(host, port, user, dbname, password, sslmode string) (*gorm.DB, error) {
	// In the case of heroku
	if os.Getenv("DATABASE_URL") != "" {
		return gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	}
	db, err := gorm.Open(
		"postgres",
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", host, port, user, dbname, password, sslmode),
	)

	return db, err
}

func GetPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		fmt.Println("INFO: No PORT environment variable detected, defaulting to 4000")
		return "localhost:4000"
	}
	return ":" + port
}

func main() {
	if os.Getenv("onServer") != "True" {
		// Loading the .env file
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// Setting up DB
	db, err := dbConnect(
		os.Getenv("dbHost"),
		os.Getenv("dbPort"),
		os.Getenv("dbUser"),
		os.Getenv("dbName"),
		os.Getenv("dbPass"),
		os.Getenv("sslmode"),
	)
	if err != nil {
		log.Fatalf("Error connecting to the database: %s", err.Error())
	}

	// Creating the tables
	db.AutoMigrate(&user.User{})
	db.AutoMigrate(&question.Question{})
	db.AutoMigrate(&question.UpVoteDetail{})

	defer db.Close()
	fmt.Println("Connected to DB...")
	//db.LogMode(true)

	// Setting up the router
	r := http.NewServeMux()

	// Setting up the hub
	hub := websocket.NewHub()
	go hub.Run()

	// Users
	userRepo := user.NewRepo(db)
	userSvc := user.NewService(userRepo)
	handler.MakeUserHandler(r, userSvc)

	// Questions
	questionRepo := question.NewRepo(db)
	questionSvc := question.NewService(questionRepo)
	handler.MakeQuestionHandler(r, questionSvc, hub)

	// To check if server up or not
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello There"))
		return
	})

	// Handling the websocket connection
	r.HandleFunc("/api/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWS(hub, w, r)
	})

	// Adding Cors middleware
	mwCors := middleware.CorsEverywhere(r)

	fmt.Println("Serving...")
	log.Fatal(http.ListenAndServe(GetPort(), mwCors))
}
