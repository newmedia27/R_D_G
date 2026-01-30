package handler

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"strings"

	"tcp/internal/users"
)

//example
// put {"name":"Alex","age":"40"}

const (
	putCommand    string = "put"
	getCommand    string = "get"
	deleteCommand string = "delete"
	listCommand   string = "list"
)

type Handler interface {
	Handle()
	Put(user *users.UserRequest) (string, error)
	Get(id string) (string, error)
	Delete(id string) (string, error)
	List() (string, error)
}

type HandleConnection struct {
	conn net.Conn
	usr  *users.UserService
}

func NewHandleConnection(conn net.Conn, usr *users.UserService) *HandleConnection {
	return &HandleConnection{
		conn: conn,
		usr:  usr,
	}
}

func (h *HandleConnection) Handle() {
	defer func() {
		err := h.conn.Close()
		if err != nil {
			slog.Error("Failed to close connection", slog.Any("error", err))
		}
	}()

	scanner := bufio.NewScanner(h.conn)
	writer := bufio.NewWriter(h.conn)

	for scanner.Scan() {
		message := scanner.Text()
		var err error
		request, err := base64.StdEncoding.DecodeString(message)
		if err != nil {
			slog.Error("Failed to decode origin message", slog.Any("error", err))
			_, _ = writer.WriteString("Failed to decode message\n")
			_ = writer.Flush()
			continue
		}

		input := strings.Split(string(request), " ")
		command := input[0]
		msg := ""
		if command != listCommand {
			if len(input) != 2 {
				slog.Error("Invalid command")
				_, _ = writer.WriteString("Invalid format command\n")
				_ = writer.Flush()
				continue
			}
			enc, err := base64.StdEncoding.DecodeString(input[1])
			if err != nil {
				slog.Error("Failed to decode message", slog.Any("error", err))
				_, _ = writer.WriteString("Failed to decode message\n")
				_ = writer.Flush()
				continue
			}
			msg = string(enc)
			fmt.Println("Message received: ", msg)
		}

		var response string

		switch command {
		case putCommand:
			response, err = h.Put(msg)
		case getCommand:
			response, err = h.Get(msg)
		case deleteCommand:
			response, err = h.Delete(msg)
		case listCommand:
			response, err = h.List()
		default:
			err = errors.New("invalid command")
		}

		if err != nil {
			_, _ = writer.WriteString(err.Error() + "\n")
		}

		_, _ = writer.WriteString(response + "\n")
		_ = writer.Flush()

		fmt.Println("Message received: ", message)
	}
	fmt.Println("Connection closed: ", h.conn.RemoteAddr())
}

func (h *HandleConnection) Put(usr string) (string, error) {
	var user users.UserRequest
	err := json.Unmarshal([]byte(usr), &user)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling payload: %w", err)
	}
	var responseUser *users.User

	responseUser, err = h.usr.CreateUser(user)
	if err != nil {
		return "", err
	}
	resp := make(map[string]string)
	resp["user_id"] = responseUser.ID
	resp["status"] = strconv.Itoa(http.StatusCreated)
	r, err := json.Marshal(resp)
	if err != nil {
		return "", fmt.Errorf("error marshalling response: %w", err)
	}
	return string(r), nil
}

func (h *HandleConnection) Get(id string) (string, error) {
	userRequest := struct {
		ID string `json:"id"`
	}{}
	err := json.Unmarshal([]byte(id), &userRequest)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling payload: %w", err)
	}
	user, err := h.usr.GetUser(userRequest.ID)
	if err != nil {
		return "", err
	}
	resp := make(map[string]string)
	resp["id"] = user.ID
	resp["name"] = user.Name
	resp["age"] = user.Age
	resp["status"] = strconv.Itoa(http.StatusOK)
	r, err := json.Marshal(resp)
	if err != nil {
		return "", fmt.Errorf("error marshalling response: %w", err)
	}
	return string(r), nil
}

func (h *HandleConnection) Delete(id string) (string, error) {
	userRequest := struct {
		ID string `json:"id"`
	}{}
	err := json.Unmarshal([]byte(id), &userRequest)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling payload: %w", err)
	}

	err = h.usr.DeleteUser(userRequest.ID)
	if err != nil {
		return "", err
	}
	resp := make(map[string]string)
	resp["status"] = strconv.Itoa(http.StatusOK)
	r, err := json.Marshal(resp)
	if err != nil {
		return "", fmt.Errorf("error marshalling response: %w", err)
	}
	return string(r), nil
}

func (h *HandleConnection) List() (string, error) {
	usrs, err := h.usr.ListUsers()
	if err != nil {
		return "", err
	}
	resp := make(map[string]interface{})
	resp["data"] = usrs
	resp["status"] = strconv.Itoa(http.StatusOK)
	r, err := json.Marshal(resp)
	if err != nil {
		return "", fmt.Errorf("error marshalling response: %w", err)
	}
	return string(r), nil
}
