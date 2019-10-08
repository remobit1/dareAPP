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

type dareRow struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	OwnerID   int       `json:"owner_id"`
}

type dareRows struct {
	Rows []dareRow `json:"dares"`
}

type request struct {
	DareUser userRow `json:"user"`
	Offset   int     `json:"offset"`
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

	var dares dareRows

	rows, err := db.Query(buildQuery(), req.DareUser.FullName, req.Offset)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No such row exists.")
		} else {
			fmt.Printf("Unable to query row: %s \n", err.Error())
		}
	}
	defer rows.Close()

	for rows.Next() {
		var row dareRow
		err = rows.Scan(&row.ID, &row.Title, &row.CreatedAt, &row.OwnerID)

		if err != nil {
			fmt.Printf("Unable to continue scanning rows: %s \n", err.Error())
		}
		dares.Rows = append(dares.Rows, row)
	}

	err = rows.Err()
	if err != nil {
		fmt.Printf("Request failed: %s \n", err.Error())
		panic(err)
	}

	resp, err := json.Marshal(dares)

	if err != nil {
		fmt.Printf("Unable to marshal struct in json: %s \n", err.Error())
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

func buildQuery() string {
	sqlStatement := `
	SELECT dares.id, dares.title 
	FROM users INNER JOIN dares ON users.id = dares.owner_id 
	WHERE full_name = $1 LIMIT 20 OFFSET $2;`

	return sqlStatement
}
