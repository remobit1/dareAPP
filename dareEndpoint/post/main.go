package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	_ "github.com/lib/pq"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	host     = "postgres.ccfpdckmolcy.us-east-1.rds.amazonaws.com"
	port     = 5432
	user     = "postgres"
	password = "Livelife4love"
	dbname   = "postgres"
)

type userRow struct {
	ID        int       `json:"id"`
	FullName  string    `json:"fullName"`
	CreatedAt time.Time `json:"createdAt"`
	Password  string    `json:"password"`
}

type dareRow struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	OwnerID   int       `json:"owner_id"`
}

type request struct {
	DareUser userRow `json:"user"`
	Dare     dareRow `json:"dare"`
}

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, req request) (events.APIGatewayProxyResponse, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	connect(db)

	result, err := db.Exec(buildQuery(), req.Dare.ID, req.Dare.Title,
		req.Dare.CreatedAt, req.DareUser.ID)

	if err != nil {
		panic(err)
	}

	fmt.Println(result)

	return events.APIGatewayProxyResponse{Body: "Dare created.", StatusCode: http.StatusOK}, nil
}

func connect(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		panic(err)
	}
	return
}

func buildQuery() string {
	sqlStatement := `
	INSERT INTO dares(id, title, created_at, owner_id) VALUES ($1, $2, $3, $4);`

	return sqlStatement
}
