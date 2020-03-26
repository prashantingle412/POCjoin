package function

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	// POCenv "github.com/prashantingle412/POC_ENV/env-function"
)

var err error
var db *gorm.DB

type User struct {
	UserID   int    `gorm:"not null;primary_key" json:"user_id"`
	Name     string `gorm:"name" json:"name"`
	Email    string `gorm:"size:255;not null;unique" json:"email"`
	Password string `gorm:"password" json:"password,omitempty"`
	Book     Book   `json:"user" gorm:"foreignkey:UserID"`
}
type Book struct {
	UserID int    `json:"_"`
	BookID int    `gorm:"not null;primary_key" json:"book_id"`
	BName  string `gorm:"bname" json:"seqno"`
}
type Result struct {
	Email string `gorm:"email" json:"user email"`
	BName string `gorm:"b_name" json:"book name"`
}

type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	User    Result `json:"result"`
	Book    Book   `json:"book"`
}

// type Response struct {
// 	Message string `json:"message"`
// 	Status  string `json:"status"`
// 	User    User   `json:"user"`
// }
/*

func Handle(w http.ResponseWriter, r *http.Request) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", "postgres", 5432, "postgres", "postgres", "openfaas")
	db, err = gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	db.Exec(`set search_path='s1'`)
	var users User
	var books Book
	//db.DropTableIfExists(&users)
	db.AutoMigrate(&users)
	db.Exec(`set search_path='s2'`)
	db.AutoMigrate(&books)
	defer db.Close()
	db.Exec(`set search_path='s1'`)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Unable to read request body")
	}
	var user User
	// var book Book
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Println("unable to unmarshal data")
		return
	}
	// user.CreatedAt = time.Now().Local()

	var resp Response

	// var m = map[string]interface{}{
	// 	"Name": user.Name, "Email": user.Email,
	// }
	// res := db.Find(&user, &User{Email: user.Email})
	// if res.RecordNotFound() {
	db.Exec(`set search_path='s1'`)
	db.Debug().Create(&user)
	db.Exec(`set search_path='s2'`)
	// b :=
	db.Debug().Create(&Book{BName: "golang", UserID: 1})
	resp = Response{Message: "User Created Successfully", Status: "200", User: user}
	// } else {
	// resp = Response{Message: "Email already exist", Status: "200", User: user}
	// }
	// db.Debug().Create(&user)
	resBody, err := json.Marshal(resp)
	if err != nil {
		log.Println("Error in marshalling response")
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resBody)
}
*/

func Handle(w http.ResponseWriter, r *http.Request) {
	// POCenv.Handle()
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", "postgres", 5432, "postgres", "postgres", "openfaas")
	db, err = gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	db.Exec(`set search_path='s1','s2'`)
	defer db.Close()
	var result Result
	var resp Response
	res := db.Table("users").Select("users.email,books.b_name").Joins("left join books on books.user_id = users.user_id").Scan(&result)
	if res.RecordNotFound() {
		resp = Response{Message: "not found, Try again", Status: "404", User: result}
	} else {
		resp = Response{Message: "data retrive successfully", Status: "200", User: result}
	}
	resBody, err := json.Marshal(resp)
	if err != nil {
		log.Println("Error in Marshalling Data")
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resBody)
}
