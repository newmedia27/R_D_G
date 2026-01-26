package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"routines/internal/documentstore"
	"routines/internal/users"
	"routines/pkg/printobject"
)

var names = []string{
	"Oleh",
	"Anna",
	"Ivan",
	"Maria",
	"Sviatoslav",
	"Dmytro",
	"Olena",
}

func generateRandomUser(rnd *rand.Rand) users.UserRequest {
	name := names[rnd.Intn(len(names))]
	age := rnd.Intn(60) + 18
	return users.UserRequest{Name: fmt.Sprintf("%s-%d", name, age), Age: strconv.Itoa(age)}
}

var numberOffNeededRoutines = 1000

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
	if err != nil {
		slog.Default().Warn("Error create index", slog.Any("err", err))
	}

	userService := users.NewUserService(userCollection)

	wg := sync.WaitGroup{}
	ch := make(chan string)
	for i := 0; i < numberOffNeededRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
			usr := generateRandomUser(rnd)
			newUser, err := userService.CreateUser(usr)
			if err != nil {
				slog.Default().Warn("Error create new user from file", slog.Any("err", err))
			}
			ch <- newUser.ID
		}()
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	ids := make([]string, 0, numberOffNeededRoutines)
	idsMu := sync.Mutex{}
	wg2 := sync.WaitGroup{}
	for id := range ch {
		idsMu.Lock()
		ids = append(ids, id)
		idsMu.Unlock()
		wg2.Add(1)
		go func(userId string) {
			defer wg2.Done()
			_, err = userService.GetUser(userId)
			if err != nil {
				slog.Default().Warn("Error get user from file", slog.Any("err", err))
			}
			err = userService.DeleteUser(userId)
			if err != nil {
				slog.Default().Warn("Error delete user from file", slog.Any("err", err))
			}
			idsMu.Lock()
			for i, v := range ids {
				if v == userId {
					ids = append(ids[:i], ids[i+1:]...)
					break
				}
			}
			idsMu.Unlock()
		}(id)
	}
	wg2.Wait()

	printobject.PrintObject("Users", ids)
	fmt.Println("Length", len(ids))
	printobject.PrintObject("collection", userService.GetDocumentsSize())

}
