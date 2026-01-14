package main

import (
	"log/slog"
	"os"

	"index/internal/documentstore"
	"index/internal/users"
	"index/pkg/printobject"
)

func main() {
	l := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	slog.SetDefault(slog.New(l))
	store := documentstore.NewStore()
	userCollection, err := store.CreateCollection("users", &documentstore.CollectionConfig{PrimaryKey: "id"})

	if err != nil {
		slog.Default().Warn(err.Error())
	}

	if err != nil {
		slog.Default().Warn("Error create new collection")
	}
	err = userCollection.CreateIndex("age")

	userService := users.NewUserService(userCollection)

	newUser := users.UserRequest{Name: "Sviatoslav", Age: "40"}
	newUser2 := users.UserRequest{Name: "John", Age: "42"}
	newUser3 := users.UserRequest{Name: "Jane", Age: "18"}

	//var usr *users.User
	//var usr2 *users.User

	_, err = userService.CreateUser(newUser)
	if err != nil {
		slog.Default().Warn("Error create new user from file", slog.Any("err", err))
	}
	_, err = userService.CreateUser(newUser2)
	if err != nil {
		slog.Default().Warn("Error create new user from file", slog.Any("err", err))
	}
	_, err = userService.CreateUser(newUser3)
	if err != nil {
		slog.Default().Warn("Error create new user from file", slog.Any("err", err))
	}

	//dump, err := store.Dump()
	//if err != nil {
	//	fmt.Println(err)
	//}

	//_ = userService.DeleteUser(usr.ID)

	//printobject.PrintObject("Store", store)

	//store, err = documentstore.NewStoreFromDump(dump)

	//err = store.DumpToFIle("temp/dump/file.json")
	//if err != nil {
	//	slog.Default().Warn("Error create new store from dump", slog.Any("err", err))
	//}
	//_ = userService.DeleteUser(usr2.ID)

	//printobject.PrintObject("storeFromDump", store)

	//printobject.PrintObject("index", q)
	store, err = documentstore.NewStoreFromFile("temp/dump/file.json")
	if err != nil {
		slog.Default().Warn("Error create new store from file", slog.Any("err", err))
	}

	minValue := "018"
	maxValue := "42"

	q, err := userService.Query("age", documentstore.QueryParams{
		Desc:     true,
		MinValue: &minValue,
		MaxValue: &maxValue,
	})
	printobject.PrintObject("NewStoreFromFile", q)

}
