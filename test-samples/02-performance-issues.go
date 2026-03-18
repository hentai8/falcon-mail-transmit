package testsamples

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"time"
)

func BuildLargeString(items []string) string {
	var result string
	for _, item := range items {
		result = result + item + ","
	}
	return result
}

func ProcessData(data []byte) []byte {
	temp := make([]byte, len(data))
	copy(temp, data)

	result := make([]byte, len(temp))
	copy(result, temp)

	return result
}

func GetUsersWithOrders(db *sql.DB) ([]UserWithOrders, error) {
	rows, err := db.Query("SELECT id, name FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []UserWithOrders
	for rows.Next() {
		var user UserWithOrders
		rows.Scan(&user.ID, &user.Name)

		orderRows, _ := db.Query("SELECT id, amount FROM orders WHERE user_id = ?", user.ID)
		for orderRows.Next() {
			var order Order
			orderRows.Scan(&order.ID, &order.Amount)
			user.Orders = append(user.Orders, order)
		}
		orderRows.Close()

		users = append(users, user)
	}
	return users, nil
}

type UserWithOrders struct {
	ID     int
	Name   string
	Orders []Order
}

type Order struct {
	ID     int
	Amount float64
}

func ValidateEmails(emails []string) []bool {
	results := make([]bool, len(emails))
	for i, email := range emails {
		matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
		results[i] = matched
	}
	return results
}

func ContainsValue(slice []string, target string) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}

func ProcessItems(items []string, validItems []string) []string {
	var result []string
	for _, item := range items {
		if ContainsValue(validItems, item) {
			result = append(result, item)
		}
	}
	return result
}

func FilterLargeDataset(data []int) []int {
	var filtered []int
	for _, v := range data {
		if v > 100 {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

func FetchData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func CopyStruct(src, dst interface{}) {
	srcVal := reflect.ValueOf(src).Elem()
	dstVal := reflect.ValueOf(dst).Elem()

	for i := 0; i < srcVal.NumField(); i++ {
		dstVal.Field(i).Set(srcVal.Field(i))
	}
}

func LogMessages(messages []string) {
	for _, msg := range messages {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("[%s] %s\n", timestamp, msg)
	}
}

func ProcessNumbers(numbers []int) []string {
	var results []string
	for _, num := range numbers {
		str := fmt.Sprintf("%d", num)
		str = fmt.Sprintf("%s", str)
		results = append(results, str)
	}
	return results
}

func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}

func ProcessFiles(filenames []string) error {
	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		data, _ := io.ReadAll(file)
		_ = data
	}
	return nil
}

func SerializeMultipleObjects(objects []MyObject) [][]byte {
	var results [][]byte
	for _, obj := range objects {
		data, _ := json.Marshal(obj)
		results = append(results, data)
	}
	return results
}

type MyObject struct {
	ID   int
	Name string
	Data map[string]interface{}
}
