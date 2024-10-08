package initilize

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-redis/redis/v8"
	"vsensetech.in/go_fingerprint_server/payload"
)

type Init struct{
	db *sql.DB
	rdb *redis.Client
	ctx context.Context
}

func NewInitInstance(db *sql.DB , rdb *redis.Client , ctx context.Context) *Init{
	return &Init{
		db,
		rdb,
		ctx,
	}
}

func(i *Init) InitilizeTables(w http.ResponseWriter , r *http.Request){
	if _ , err := i.db.Exec("CREATE TABLE admin(user_id VARCHAR(100) PRIMARY KEY, user_name VARCHAR(50) NOT NULL, password VARCHAR(100) NOT NULL)"); err != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
		return
	}
	if _ , err := i.db.Exec("CREATE TABLE users(user_id VARCHAR(100) PRIMARY KEY, user_name VARCHAR(50) NOT NULL, password VARCHAR(100) NOT NULL)"); err != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
		return
	}
	if _ , err := i.db.Exec("CREATE TABLE biometric(user_id VARCHAR(100), unit_id VARCHAR(50) PRIMARY KEY , online BOOLEAN NOT NULL, FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE)"); err != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
		return
	}
	if _ , err := i.db.Exec("CREATE TABLE fingerprintdata(student_id VARCHAR(100) PRIMARY KEY, student_unit_id VARCHAR(100) , unit_id VARCHAR(50) , fingerprint VARCHAR(1000), FOREIGN KEY (unit_id) REFERENCES biometric(unit_id) ON DELETE CASCADE)"); err != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
		return
	}
	if _ , err := i.db.Exec("CREATE TABLE attendance(student_id VARCHAR(100), student_unit_id VARCHAR(100) , unit_id VARCHAR(50), date VARCHAR(20), login VARCHAR(20), logout VARCHAR(20), FOREIGN KEY (unit_id) REFERENCES biometric(unit_id) ON DELETE CASCADE , FOREIGN KEY (student_id) REFERENCES fingerprintdata(student_id) ON DELETE CASCADE)"); err != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
		return
	}
	if _ , err := i.db.Exec("CREATE TABLE times(user_id VARCHAR(200) , morning_start VARCHAR(20) , morning_end VARCHAR(20) , afternoon_start VARCHAR(20) , afternoon_end VARCHAR(20) , evening_start VARCHAR(20) , evening_end VARCHAR(20) , FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE)"); err != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
		return 
	}
	if _ , err := i.rdb.Do(i.ctx , "JSON.SET" , "deletes" , "$" , "{}").Result(); err != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
		return
	}
	if _ , err := i.rdb.Do(i.ctx , "JSON.SET" , "inserts" , "$" , "{}").Result(); err != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload.SimpleSuccessPayload{Message: "Success"})
}