package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type user struct {
	ID        int       `json:"id"`
	FullName  string    `json:"fullName"`
	CreatedAt time.Time `json:"createdAt"`
	Password  string    `json:"password"`
}

type dare struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	OwnerID   int       `json:"owner_id"`
}

type request struct {
	DareUser user `json:"user"`
	Dare     dare `json:"dare"`
}

func main() {
	testUser := user{
		ID:        00001,
		FullName:  "Joseph Morgan",
		CreatedAt: time.Now(),
		Password:  "Livelife4Love",
	}

	testDare := dare{
		ID:        9,
		Title:     "Best dare evah!",
		CreatedAt: time.Now(),
	}

	testReq := request{
		DareUser: testUser,
		Dare:     testDare,
	}

	msg, err := json.Marshal(testReq)

	if err != nil {
		fmt.Printf("Unable to marshal given struct: %s \n", err.Error())
	}

	body := bytes.NewReader(msg)

	response, err := http.Post("https://3ne30pte12.execute-api.us-east-1.amazonaws.com/test/user/dares", "application/json", body)

	if err != nil {
		fmt.Printf("Unable to POST to given url: %s \n", err.Error())
	}

	fmt.Println(response)
}
