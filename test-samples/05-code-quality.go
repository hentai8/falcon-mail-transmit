package testsamples

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

var x int
var user_name string
var MAXCOUNT = 100

func GetData(id int) string {
	return "data"
}

type userinfo struct {
	id   int
	name string
}

func ProcessOrder(orderID int, userID int, items []string, quantities []int,
	discountCode string, shippingAddress string, paymentMethod string) (bool, string, float64, error) {

	if orderID <= 0 {
		return false, "", 0, fmt.Errorf("invalid order id")
	}

	if userID <= 0 {
		return false, "", 0, fmt.Errorf("invalid user id")
	}

	if len(items) == 0 {
		return false, "", 0, fmt.Errorf("no items")
	}

	if len(items) != len(quantities) {
		return false, "", 0, fmt.Errorf("items and quantities mismatch")
	}

	var total float64
	for i, item := range items {
		price := getPrice(item)
		qty := quantities[i]
		total += price * float64(qty)
	}

	if discountCode != "" {
		discount := getDiscount(discountCode)
		if discount > 0 && discount <= 100 {
			total = total * (1 - discount/100)
		}
	}

	shipping := calculateShipping(shippingAddress)
	total += shipping

	validPayment := false
	if paymentMethod == "credit_card" || paymentMethod == "debit_card" ||
		paymentMethod == "paypal" || paymentMethod == "alipay" || paymentMethod == "wechat" {
		validPayment = true
	}

	if !validPayment {
		return false, "", 0, fmt.Errorf("invalid payment method")
	}

	paymentSuccess := processPayment(userID, total, paymentMethod)
	if !paymentSuccess {
		return false, "", 0, fmt.Errorf("payment failed")
	}

	orderNum := createOrder(orderID, userID, items, quantities, total)

	sendNotification(userID, orderNum)

	for i, item := range items {
		updateInventory(item, quantities[i])
	}

	return true, orderNum, total, nil
}

func getPrice(item string) float64                                  { return 10.0 }
func getDiscount(code string) float64                               { return 10.0 }
func calculateShipping(address string) float64                      { return 5.0 }
func processPayment(userID int, amount float64, method string) bool { return true }
func createOrder(orderID, userID int, items []string, quantities []int, total float64) string {
	return "ORD123"
}
func sendNotification(userID int, orderNum string) {}
func updateInventory(item string, qty int)         {}

func CalculatePrice(basePrice float64, userLevel int, itemCount int) float64 {
	if userLevel == 1 {
		basePrice = basePrice * 0.95
	} else if userLevel == 2 {
		basePrice = basePrice * 0.9
	} else if userLevel == 3 {
		basePrice = basePrice * 0.85
	}

	if itemCount > 10 {
		basePrice = basePrice * 0.98
	} else if itemCount > 50 {
		basePrice = basePrice * 0.95
	} else if itemCount > 100 {
		basePrice = basePrice * 0.9
	}

	return basePrice * 1.13
}

func ValidateEmail(email string) bool {
	if len(email) == 0 {
		return false
	}
	if !strings.Contains(email, "@") {
		return false
	}
	if !strings.Contains(email, ".") {
		return false
	}
	return true
}

func ValidateUsername(username string) bool {
	if len(username) == 0 {
		return false
	}
	if len(username) < 3 {
		return false
	}
	if len(username) > 20 {
		return false
	}
	return true
}

func ValidatePassword(password string) bool {
	if len(password) == 0 {
		return false
	}
	if len(password) < 8 {
		return false
	}
	if len(password) > 50 {
		return false
	}
	return true
}

func ProcessData(data map[string]interface{}) error {
	if data != nil {
		if val, ok := data["user"]; ok {
			if userMap, ok := val.(map[string]interface{}); ok {
				if id, ok := userMap["id"]; ok {
					if idInt, ok := id.(int); ok {
						if idInt > 0 {
							if name, ok := userMap["name"]; ok {
								if nameStr, ok := name.(string); ok {
									if len(nameStr) > 0 {
										fmt.Println("Processing user:", nameStr)
										return nil
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return fmt.Errorf("invalid data")
}

func Add(a, b int) int {
	return a + b
}

func CalculateDiscount(price float64) float64 {
	return price * 1.13
}

func Process(s string) string {
	a := strings.Split(s, ",")
	b := make([]string, 0)
	for _, c := range a {
		d := strings.TrimSpace(c)
		if len(d) > 0 {
			e := strings.ToLower(d)
			b = append(b, e)
		}
	}
	f := strings.Join(b, "|")
	return f
}

func ShouldProcess(user User, order Order, status string, time int64) bool {
	return (user.ID > 0 && user.Active && !user.Deleted) &&
		(order.ID > 0 && order.Amount > 100 && order.Status == "pending") &&
		(status == "approved" || status == "verified" || status == "confirmed") &&
		(time > 1609459200 && time < 1640995200) &&
		(user.Level >= 2 && user.Points > 1000) &&
		((order.PaymentMethod == "credit" && user.CreditLimit > order.Amount) ||
			(order.PaymentMethod == "debit" && user.Balance > order.Amount) ||
			(order.PaymentMethod == "points" && user.Points > order.Amount*100))
}

type Order struct {
	ID            int
	Amount        float64
	Status        string
	PaymentMethod string
}

func GetUserID(data interface{}) int {
	user := data.(map[string]interface{})
	id := user["id"].(int)
	return id
}

func methodOne(s string) string { return s }

func MethodTwo(s string) string {
	return s
}

func method_three(s string) string {
	return s
}

func GetStatus(code int) string {
	if code == 200 {
		return "OK"
	} else {
		if code == 404 {
			return "Not Found"
		} else {
			if code == 500 {
				return "Internal Error"
			} else {
				return "Unknown"
			}
		}
	}
}

func A() int {
	return B()
}

func B() int {
	return C()
}

func C() int {
	return 42
}

type PublicUser struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
}

func HandleRequest(reqType string) {
	switch reqType {
	case "GET":
		handleGet()
	case "POST":
		handlePost()
	case "PUT":
		handlePut()
	}
}

func handleGet()  {}
func handlePost() {}
func handlePut()  {}

func SaveUserToMySQL(user User) error {
	db, err := sql.Open("mysql", "connection_string")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO users VALUES (?, ?, ?)", user.ID, user.Name, user.Email)
	return err
}

func SaveProductToMySQL(product Product) error {
	db, err := sql.Open("mysql", "connection_string")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO products VALUES (?, ?, ?)", product.ID, product.Name, product.Price)
	return err
}

func SaveOrderToMySQL(order Order) error {
	db, err := sql.Open("mysql", "connection_string")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO orders VALUES (?, ?, ?)", order.ID, order.Amount, order.Status)
	return err
}

type Product struct {
	ID    int
	Name  string
	Price float64
}
