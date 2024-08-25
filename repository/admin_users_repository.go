package repository

import (
	"database/sql"
	"sync"

	"vsensetech.in/go_fingerprint_server/models"
)

type UsersRepo struct{
	db *sql.DB
	mut *sync.Mutex
}

func NewUsersRepo(db *sql.DB , mut *sync.Mutex) *UsersRepo{
	return &UsersRepo{
		db,
		mut,
	}
}

func(ur *UsersRepo) FetchAllUsers() ([]models.Users , error) {
	
	res , err := ur.db.Query("SELECT username , user_id FROM users")
	if err != nil {
		return nil , err
	}
	defer res.Close()
	
	var userList []models.Users
	var user models.Users
	
	for res.Next() {
		err := res.Scan(&user.UserName , &user.UserID)
		if err != nil {
			return nil , err
		}
		userList = append(userList, user)
	}
	if res.Err() != nil {
		return nil , res.Err()
	}
	return userList , nil
}