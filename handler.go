package function

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var db *gorm.DB
var err error

type User struct {
	Email    string `gorm:"size:255;not null;unique" json:"email"`
	Password string `gorm:"password" json:"password"`
}

type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Email   string `json:"email"`
}

func Handle(w http.ResponseWriter, r *http.Request) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", "postgres", 5432, "postgres", "postgres", "openfaas")
	db, err = gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)

	}
	db.Exec(`set search_path='roles'`)
	defer db.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Unable to read request body")
		w.WriteHeader(500) // Return 500 Internal Server Error.

	}
	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Println("unable to unmarshal data")
		w.WriteHeader(400) // Return 400 Bad Request.
		return
	}

	var resp Response
	res := db.Debug().Where("email = ?", user.Email).Find(&user)
	if res.RecordNotFound() {
		resp = Response{Message: "Email not found, Try again", Status: "404", Email: user.Email}
	} else {
		resp = Response{Message: "successfully logged in", Status: "200", Email: user.Email}
	}

	resBody, err := json.Marshal(resp)
	if err != nil {
		log.Println("Error in Marshalling Data")
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resBody)
}
