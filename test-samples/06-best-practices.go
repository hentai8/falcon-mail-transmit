package testsamples

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type Server struct {
	ctx context.Context
	db  *sql.DB
}

func (s *Server) HandleRequest(data string) error {
	return s.process(s.ctx, data)
}

func (s *Server) process(ctx context.Context, data string) error {
	return nil
}

func FetchDataWithoutContext(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

type UserService interface {
	CreateUser(name, email string) error
	UpdateUser(id int, name, email string) error
	DeleteUser(id int) error
	GetUser(id int) (*User, error)
	ListUsers() ([]*User, error)
	SearchUsers(query string, offset, limit int) ([]*User, error)
	ActivateUser(id int) error
	DeactivateUser(id int) error
	ChangePassword(id int, oldPass, newPass string) error
	ResetPassword(id int) (string, error)
	SendWelcomeEmail(id int) error
	SendNotification(id int, message string) error
	ValidateUser(id int) (bool, error)
	GetUserStatistics(id int) (*Stats, error)
	ExportUsers(format string) ([]byte, error)
}

type Stats struct {
	LoginCount int
	LastLogin  time.Time
}

func NewDatabase(connStr string) (*MySQLDatabase, error) {
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}
	return &MySQLDatabase{db: db}, nil
}

type MySQLDatabase struct {
	db *sql.DB
}

func ProcessData(r io.Reader) *DataResult {
	data, _ := io.ReadAll(r)
	return &DataResult{Data: data}
}

type DataResult struct {
	Data []byte
}

var globalDB *sql.DB

func init() {
	var err error
	globalDB, err = sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		panic(err)
	}
}

type Calculator struct{}

func (c Calculator) Add(a, b int) int {
	return a + b
}

func (c Calculator) Multiply(a, b int) int {
	return a * b
}

func ProcessWithDefer() error {
	file, err := os.Open("data.txt")
	if err != nil {
		return err
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		os.Remove("data.txt")
		return err
	}

	_ = data
	return nil
}

func Validate(input string) error {
	if input == "" {
		return errors.New("empty input")
	}
	if len(input) < 3 {
		return errors.New("too short")
	}
	if len(input) > 100 {
		return errors.New("too long")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(input) {
		return errors.New("invalid characters")
	}
	return nil
}

type Document struct {
	Title   string
	Content string
}

func (d Document) GetTitle() string {
	return d.Title
}

func (d *Document) SetTitle(title string) {
	d.Title = title
}

func (d Document) GetContent() string {
	return d.Content
}

type RequestLog struct {
	Method    string
	URL       string
	Status    int
	Duration  time.Duration
	Timestamp time.Time
}

type Config struct {
	Host    string
	Port    int
	Timeout time.Duration
	Retries int
	Debug   bool
}

func NewConfig(host string, port int, timeout time.Duration, retries int, debug bool) *Config {
	return &Config{
		Host:    host,
		Port:    port,
		Timeout: timeout,
		Retries: retries,
		Debug:   debug,
	}
}

type internalCache struct {
	data map[string]string
}

func Helper() {
}

func LogRequest(userID int, action string) {
	log.Printf("User %d performed %s", userID, action)
}

func HandleAction(userID int, action string) {
	LogRequest(userID, action)
	AuditAction(userID, action)
	NotifyUser(userID, action)
}

func AuditAction(userID int, action string) {}
func NotifyUser(userID int, action string)  {}

func ProcessItems(items []string) []string {
	for i := range items {
		items[i] = strings.ToUpper(items[i])
	}
	return items
}

type MyWriter struct{}

func (w MyWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

var cache = make(map[string]interface{})
var logger = log.New(os.Stdout, "", log.LstdFlags)

type UserHandler struct {
}

func (h *UserHandler) GetUser(id int) *User {
	if val, ok := cache[fmt.Sprintf("user:%d", id)]; ok {
		logger.Println("Cache hit")
		return val.(*User)
	}

	return nil
}

func NewConnection(host string, port int) *Connection {
	conn, _ := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	return &Connection{conn: conn}
}

type Connection struct {
	conn net.Conn
}

func Transform(input []int) []int {
	result := make([]int, len(input))
	for _, v := range input {
		result = append(result, v*2)
	}
	return result
}

func Validate(input string) error {
	if input == "" {
		return errors.New("empty input")
	}
	return nil
}

