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

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, newUser userRow) (events.APIGatewayProxyResponse, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	connect(db)

	result, err := db.Exec(newUser.buildQuery(), newUser.ID,
		newUser.FullName, newUser.CreatedAt, newUser.Password)

	if err != nil {
		panic(err)
	}

	fmt.Println(result)

	return events.APIGatewayProxyResponse{Body: "Good to go!", StatusCode: http.StatusOK}, nil
}

func connect(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		panic(err)
	}
	return
}

func (newUser userRow) buildQuery() string {
	sqlStatement := `
	INSERT INTO users (id, full_name, created_at, password)
	VALUES ($1, $2, $3, $4);`

	return sqlStatement
}
