package main

import (
	"log/slog"
	"os"

	"github.com/newmedia27/R_D_G/l_6/internal/documentstore"
	"github.com/newmedia27/R_D_G/l_6/internal/users"
	"github.com/newmedia27/R_D_G/l_6/pkg/printobject"
)

// Привіт. Логую тільки в мейні. По суті всі помилки виходять сюди.
//Про тести ще читаю, розбираюся, дякую за лінк.

func main() {
	l := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	slog.SetDefault(slog.New(l))
	store := documentstore.NewStore()
	userCollection, err := store.CreateCollection("users", &documentstore.CollectionConfig{PrimaryKey: "id"})

	if err != nil {
		slog.Default().Warn("Error create new collection")
	}

	userService := users.NewUserService(userCollection)

	newUser := users.UserRequest{Name: "Sviatoslav", Age: "40"}
	newUser2 := users.UserRequest{Name: "John", Age: "42"}

	var usr *users.User
	var usr2 *users.User

	usr, err = userService.CreateUser(newUser)
	if err != nil {
		slog.Default().Warn("Error create new user from file", slog.Any("err", err))
	}
	usr2, err = userService.CreateUser(newUser2)
	if err != nil {
		slog.Default().Warn("Error create new user from file", slog.Any("err", err))
	}

	//dump, err := store.Dump()
	//if err != nil {
	//	fmt.Println(err)
	//}

	_ = userService.DeleteUser(usr.ID)

	//printobject.PrintObject("Store", store)

	//store, err = documentstore.NewStoreFromDump(dump)

	err = store.DumpToFIle("temp/dump/file.json")
	if err != nil {
		slog.Default().Warn("Error create new store from dump", slog.Any("err", err))
	}
	_ = userService.DeleteUser(usr2.ID)
	//printobject.PrintObject("storeFromDump", store)

	store, err = documentstore.NewStoreFromFile("temp/dump/file.json")
	if err != nil {
		slog.Default().Warn("Error create new store from file", slog.Any("err", err))
	}
	printobject.PrintObject("NewStoreFromFile", store)

}
