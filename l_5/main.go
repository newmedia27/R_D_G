package main

import (
	"encoding/json"
	"fmt"

	"github.com/newmedia27/R_D_G/l_5/internal/documentstore"
	"github.com/newmedia27/R_D_G/l_5/internal/users"
)

func PrintObject(message string, object any) {
	if message == "" {
		message = "Object: "
	}

	p, err := json.MarshalIndent(object, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(message, string(p))
}

func main() {
	store := documentstore.NewStore()
	userCollection, err := store.CreateCollection("users", &documentstore.CollectionConfig{PrimaryKey: "id"})

	if err != nil {
		fmt.Println(err)
	}

	userService := users.NewUserService(userCollection)

	newUser := users.UserRequest{Name: "Sviatoslav", Age: "40"}
	newUser2 := users.UserRequest{Name: "John", Age: "42"}

	var usr *users.User

	usr, err = userService.CreateUser(newUser)
	if err != nil {
		fmt.Println(err)
	}
	_, err = userService.CreateUser(newUser2)
	if err != nil {
		fmt.Println(err)
	}
	list, err := userService.ListUsers()
	if err != nil {
		fmt.Println(err)
	}

	usr, _ = userService.GetUser(usr.ID)

	PrintObject("List", list)
	PrintObject("get_user", usr)
	_ = userService.DeleteUser(usr.ID)
	PrintObject("store", store)
}
