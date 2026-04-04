package users

import (
	"errors"

	"github.com/google/uuid"
	"index/internal/documentstore"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	ID string `json:"id"`
	UserRequest
}

type UserRequest struct {
	Name string `json:"name"`
	Age  string `json:"age"`
}

type UserService struct {
	coll documentstore.Collector
}

func NewUserService(coll documentstore.Collector) *UserService {
	return &UserService{
		coll: coll,
	}
}

func (s *UserService) CreateUser(u UserRequest) (*User, error) {
	var user = new(User)
	var err error
	var doc *documentstore.Document

	id := uuid.New()

	user.ID = id.String()
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
func (s *UserService) Query(fieldName string, params documentstore.QueryParams) ([]documentstore.Document, error) {
	return s.coll.Query(fieldName, params)
}
