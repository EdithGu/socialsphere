package service

import (
	"fmt"
	"reflect"

	"socialai/backend"
	"socialai/constants"
	"socialai/model"

	"github.com/olivere/elastic/v7"
	"golang.org/x/crypto/bcrypt"
)

// CheckUser looks up the user by username, then verifies the supplied
// password against the bcrypt hash stored in Elasticsearch.
// true / nil:  credentials are valid
// false / nil: user not found, or password does not match
func CheckUser(username, password string) (bool, error) {
	query := elastic.NewTermQuery("username", username)

	searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
	if err != nil {
		return false, err
	}

	if searchResult.TotalHits() == 0 {
		return false, nil
	}

	var utype model.User
	for _, item := range searchResult.Each(reflect.TypeOf(utype)) {
		u := item.(model.User)
		if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil {
			fmt.Printf("Login as %s\n", username)
			return true, nil
		}
	}

	return false, nil
}

// AddUser hashes the incoming password with bcrypt before persisting the
// user document, so plaintext passwords are never stored.
// false: user existed
// true: successfully added
func AddUser(user *model.User) (bool, error) {
	// check user existed or not
	query := elastic.NewTermQuery("username", user.Username)
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
	if err != nil {
		return false, err
	}

	if searchResult.TotalHits() > 0 {
		return false, nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}
	user.Password = string(hashedPassword)

	// call backend to save
	err = backend.ESBackend.SaveToES(user, constants.USER_INDEX, user.Username)
	if err != nil {
		return false, err
	}

	// construct response
	fmt.Printf("User is added: %s\n", user.Username)
	return true, nil
}
