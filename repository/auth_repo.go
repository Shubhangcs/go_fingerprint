package repository

import (
	"database/sql"
	"encoding/json"
	"io"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"vsensetech.in/go_fingerprint_server/models"
)

type Auth struct{
	db *sql.DB
	mut *sync.Mutex
}

func NewAdminAuth(db *sql.DB , mut *sync.Mutex) *Auth {
	return &Auth{
		db,
		mut,
	}
}

func(a *Auth) Register(reader *io.ReadCloser , urlPath string) (string , error) {
	//Creating a new variable of type AdminAuthDetails
	var newUser models.AuthDetails
	
	//Decoding the json from reader to the newly created variale
	if err := json.NewDecoder(*reader).Decode(&newUser); err != nil {
		return "",err
	}
	
	//Hashing the password
	hashpass , err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return "",err
	}
	//Execuiting the query and Creating new UUID and returning error if present
	var newUID = uuid.New().String()
	if _ , err := a.db.Exec("INSERT INTO "+urlPath+"(user_id , username , password) VALUES($1 , $2 , $3)", &newUID , &newUser.Name , hashpass); err != nil {
		return "",err
	}
	
	
	//Creating JWT token and Setting Cookie
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":newUID,
		"username":newUser.Name,
		"expiry": time.Now().Add(365 * 24 * time.Hour).Unix(),
	})
	tokenString , err := token.SignedString([]byte("vsense"))
	if err != nil {
		return "",err
	}

	//Return JWT token if No error
	return tokenString,nil
}

func(a *Auth) Login(reader *io.ReadCloser , urlPath string)  (string , error) {
	//Creating a new variable of type AdminAuthDetails
	var userIns models.AuthDetails
	var dbUser models.AuthDetails
	var UID string
	
	//Decoding the json from reader to the newly created variale
	if err := json.NewDecoder(*reader).Decode(&userIns); err !=  nil {
		return "",err
	}
	
	//Querying User from Database
		err := a.db.QueryRow("SELECT user_id , username , password FROM "+urlPath+" WHERE username=$1", &userIns.Name).Scan(&UID, &dbUser.Name , &dbUser.Password)
		if err != nil {
			return "",err
		}

	
	//Comparing HashedPassword with Normal Password
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(userIns.Password))
	if err != nil {
		return "",err
	}
	
	//Creating JWT token and Setting Cookie
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":UID,
		"username":dbUser.Name,
		"expiry": time.Now().Add(365 * 24 * time.Hour).Unix(),
	})
	tokenString , err := token.SignedString([]byte("vsense"))
	if err != nil {
		return "",err
	}
	
	//JWT token if No Error
	return tokenString,nil
}