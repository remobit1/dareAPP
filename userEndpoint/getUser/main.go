package main

import (
	"context"
	"database/sql"
	"encoding/json"
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

func handle(ctx context.Context, dareUser userRow) (events.APIGatewayProxyResponse, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	connect(db)

	row := db.QueryRow(dareUser.buildQuery())

	err = row.Scan(&dareUser.ID, &dareUser.FullName, dareUser.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No such row exists.")
		} else {
			fmt.Printf("Unable to query row: %s \n", err.Error())
		}
	}

	resp, err := json.Marshal(dareUser)

	if err != nil {
		fmt.Printf("Unable to marshal struct into json: %s \n", err.Error())
	}

	return events.APIGatewayProxyResponse{Body: string(resp), StatusCode: http.StatusOK}, nil
}

func connect(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		panic(err)
	}
	return
}

func (user userRow) buildQuery() (string, int) {
	sqlStatement := `SELECT * FROM users WHERE id = $1;`

	return sqlStatement, user.ID
}
