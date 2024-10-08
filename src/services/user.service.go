package services

import (
	"errors"
	"fmt"
	"line/src/common/token"
	"line/src/models"
	"log"
	"sync"

	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

type PaginatedUsers struct {
	Users      []models.User `json:"users"`
	TotalCount int64         `json:"total_count"`
}

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (us *UserService) CreateUser(name, email, password string) (map[string]interface{}, error) {
	user := User{
		Name:     name,
		Email:    email,
		Password: password,
	}

	err := us.DB.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"user":    user,
		"message": "User signed up successfully",
	}, nil
}

func (us *UserService) SigninUser(email, password string) (map[string]interface{}, error) {
	var user models.User

	err := us.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if user.Password != password {
		return nil, errors.New("invalid credentials")
	}

	token, err := token.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return nil, errors.New("token generation error")
	}
	log.Print("User login successfully... ", user)
	return map[string]interface{}{
		"id":      user.ID,
		"name":    user.Name,
		"email":   user.Email,
		"token":   token,
		"message": "login successfully",
	}, nil
}

func (us *UserService) FindAllUsers(searchQuery string, limit, offset int) (PaginatedUsers, error) {
	var users []models.User
	var totalCount int64
	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	countQuery := "SELECT COUNT(*) FROM users"
	query := "SELECT * FROM users"

	if searchQuery != "" {
		searchCondition := fmt.Sprintf(" WHERE name LIKE '%%%s%%' OR email LIKE '%%%s%%'", searchQuery, searchQuery)
		countQuery += searchCondition
		query += searchCondition
	}

	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := us.DB.Raw(countQuery).Scan(&totalCount).Error
		if err != nil {
			errChan <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := us.DB.Raw(query).Scan(&users).Error
		if err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return PaginatedUsers{}, err
		}
	}

	return PaginatedUsers{
		Users:      users,
		TotalCount: totalCount,
	}, nil

}
