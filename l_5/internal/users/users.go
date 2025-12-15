package users

import (
	"errors"
	"os/exec"

	"github.com/newmedia27/R_D_G/l_5/internal/documentstore"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

//для reflect
//type User struct {
//	ID   string `json:"id"`
//	Name string `json:"name"`
//	Age  string `json:"age"`
//}
//
//type UserRequest struct {
//	Name string `json:"name"`
//	Age  string `json:"age"`
//}

// Без рефлекту
type User struct {
	ID string `json:"id"`
	UserRequest
}

type UserRequest struct {
	Name string `json:"name"`
	Age  string `json:"age"`
}

type UserService struct {
	coll *documentstore.Collection
}

func NewUserService(coll *documentstore.Collection) *UserService {
	return &UserService{
		coll: coll,
	}
}

func (s *UserService) CreateUser(u UserRequest) (*User, error) {
	var user = new(User)
	var uuid []byte
	var err error
	var doc *documentstore.Document

	uuid, err = exec.Command("uuidgen").Output()
	if err != nil {
		return nil, err
	}
	user.ID = string(uuid)
	user.Name = u.Name
	user.Age = u.Age

	doc, err = documentstore.MarshalDocument(user)
	if err != nil {
		return nil, err
	}
	err = s.coll.Put(*doc)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) ListUsers() ([]User, error) {
	docs := s.coll.List()
	var users []User
	users = make([]User, 0, len(docs))
	for _, doc := range docs {
		var user User
		err := documentstore.UnmarshalDocument(&doc, &user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *UserService) GetUser(userID string) (*User, error) {
	doc, err := s.coll.Get(userID)
	if err != nil {
		return nil, err
	}
	var user User
	err = documentstore.UnmarshalDocument(doc, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) DeleteUser(userID string) error {
	if ok := s.coll.Delete(userID); !ok {
		return ErrUserNotFound
	}
	return nil
}
