package testsamples

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func ReadConfig(filename string) Config {
	data, _ := os.ReadFile(filename)

	var config Config
	_ = json.Unmarshal(data, &config)

	return config
}

type Config struct {
	Host string
	Port int
}

func ProcessFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer file.Close()

	data := make([]byte, 1024)
	file.Read(data)

	return nil
}

func LoadUser(id int) (*User, error) {
	data, err := fetchUserData(id)
	if err != nil {
		return nil, err
	}

	user, err := parseUser(data)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func fetchUserData(id int) ([]byte, error) {
	return nil, errors.New("network error")
}

func parseUser(data []byte) (*User, error) {
	return nil, errors.New("parse error")
}

func CalculateDiscount(price float64, percent float64) float64 {
	if percent < 0 || percent > 100 {
		panic("invalid discount percentage")
	}
	return price * (1 - percent/100)
}

func CopyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func UpdateRecord(id int, data map[string]interface{}) error {
	err := validate(data)
	if err != nil {
		return errors.New("validation failed")
	}

	err = saveToDatabase(id, data)
	if err != nil {
		return errors.New("save failed")
	}

	return nil
}

func validate(data map[string]interface{}) error {
	return nil
}

func saveToDatabase(id int, data map[string]interface{}) error {
	return nil
}

func ProcessMultipleFiles(files []string) error {
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Error reading %s: %v\n", file, err)
			continue
		}

		err = processData(data)
		if err != nil {
			fmt.Printf("Error processing %s: %v\n", file, err)
			continue
		}

		err = saveResult(file, data)
		if err != nil {
			fmt.Printf("Error saving %s: %v\n", file, err)
			continue
		}
	}
	return nil
}

func processData(data []byte) error {
	return nil
}

func saveResult(file string, data []byte) error {
	return nil
}

func SendNotification(userID int, message string) {
	err := sendEmail(userID, message)
	if err != nil {
	}
}

func sendEmail(userID int, message string) error {
	return errors.New("email service unavailable")
}

func HandleFileError(err error) {
	if err.Error() == "file not found" {
		createFile()
	} else if err.Error() == "permission denied" {
		handlePermissionError()
	}
}

func createFile()            {}
func handlePermissionError() {}

func GetUserInfo(id int) (string, int, error) {
	name, age, err := fetchUserInfo(id)
	if err != nil {
		return name, age, err
	}
	return name, age, nil
}

func fetchUserInfo(id int) (string, int, error) {
	return "", 0, errors.New("user not found")
}

func APIHandler(w http.ResponseWriter, r *http.Request) {

	data := processRequest(r)
	w.Write(data)
}

func processRequest(r *http.Request) []byte {
	return nil
}

type MyError struct {
	Code    int
	Message string
}

func (e MyError) Error() string {
	return e.Message
}

func DoSomething() error {
	return MyError{Code: 404, Message: "not found"}
}

func HandleError(err error) {
	var myErr *MyError
	if errors.As(err, &myErr) {
		fmt.Printf("Code: %d\n", myErr.Code)
	}
}

func FetchAndProcess(id int) error {
	data, err := fetchData(id)
	if err != nil {
		return fmt.Errorf("failed to fetch and process: %w", err)
	}

	err = process(data)
	if err != nil {
		return fmt.Errorf("failed to fetch and process: %w", err)
	}

	return nil
}

func fetchData(id int) ([]byte, error) {
	err := query(id)
	if err != nil {
		return nil, fmt.Errorf("failed to query data: %w", err)
	}
	return nil, nil
}

func query(id int) error {
	return fmt.Errorf("failed to execute query: %w", errors.New("connection failed"))
}

func process(data []byte) error {
	return nil
}

func BatchProcess(items []Item) error {
	var lastErr error
	for _, item := range items {
		err := processItem(item)
		if err != nil {
			lastErr = err
		}
	}
	return lastErr
}

type Item struct {
	ID   int
	Data string
}

func processItem(item Item) error {
	return errors.New("process failed")
}

func FindUser(name string) (*User, error) {
	users := []User{}

	for _, user := range users {
		if user.Name == name {
			return &user, nil
		}
	}

	return nil, nil
}
