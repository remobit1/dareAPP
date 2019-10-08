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

type dareRow struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	OwnerID   int       `json:"owner_id"`
}

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, dare dareRow) (events.APIGatewayProxyResponse, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	connect(db)

	result, err := db.Exec(buildQuery(), dare.ID)

	if err != nil {
		fmt.Printf("Unable to delete row: %s \n", err.Error())
	}

	fmt.Println(result)

	return events.APIGatewayProxyResponse{Body: "Successfully deleted", StatusCode: http.StatusOK}, nil
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
	DELETE FROM table dares WHERE id = $1;`

	return sqlStatement
}
