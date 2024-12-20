package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

type Todo struct {
	ID       int    `json:"id"`
	Text     string `json:"text"`
	Complete bool   `json:"complete"`
}

var db *sql.DB
var jwtSecretKey string

func init() {
	var err error
	// Database connection string
	connStr := "host=localhost user=postgres password=1234 dbname=postgres sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping the database:", err)
	}

	log.Println("Connected to the database!")

	// Fetching JWT secret key from environment variable or setting a default
	jwtSecretKey = os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		log.Println("Warning: JWT_SECRET_KEY is not set, using default secret key")
		jwtSecretKey = "defaultSecretKey"
	}
}

// ฟังก์ชันการลงทะเบียนผู้ใช้
func register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user User
	// อ่านข้อมูล JSON จาก request body
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// แฮชรหัสผ่านก่อนบันทึกลงฐานข้อมูล
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// สร้าง UUID สำหรับ user
	user.ID = uuid.New() // สร้าง UUID ใหม่

	// แทรกข้อมูลผู้ใช้ลงในฐานข้อมูล
	err = db.QueryRow("INSERT INTO todo_users1 (id, email, password) VALUES ($1, $2, $3) RETURNING id", user.ID, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// กำหนดสถานะการตอบกลับให้เป็น "Created"
	w.WriteHeader(http.StatusCreated)

	// เข้ารหัสข้อมูลเป็น JSON และส่งกลับไปยัง client
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

// ฟังก์ชันสำหรับการล็อกอิน (login)
func login(w http.ResponseWriter, r *http.Request) {
	// Set content type as JSON
	w.Header().Set("Content-Type", "application/json")

	var user User
	// Decode incoming JSON request body
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Check user in the database
	var dbUser User
	err := db.QueryRow("SELECT id, email, password FROM todo_users1 WHERE email = $1", user.Email).Scan(&dbUser.ID, &dbUser.Email, &dbUser.Password)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Compare hashed password from DB with input password
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Create JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = dbUser.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix() // 72 hours expiration time

	// Sign the JWT token
	t, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Send success response with token
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
		"token":   t, // Send JWT token in the response body
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecretKey), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ฟังก์ชันสำหรับการจัดการ Todo
func handleTodo(w http.ResponseWriter, r *http.Request) {
	// ตรวจสอบ JWT token
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Authorization token missing", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// คำสั่งการดำเนินการตาม HTTP method (GET, POST, PUT, DELETE)
	switch r.Method {
	case "GET":
		// ดึง Todo ทั้งหมด
		idStr := strings.TrimPrefix(r.URL.Path, "/todo/")
		if idStr == "" {
			rows, err := db.Query("SELECT id, text, complete FROM todos ORDER BY id ")
			if err != nil {
				http.Error(w, "Failed to fetch todos", http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			var todos []Todo
			for rows.Next() {
				var todo Todo
				if err := rows.Scan(&todo.ID, &todo.Text, &todo.Complete); err != nil {
					http.Error(w, "Failed to scan todo", http.StatusInternalServerError)
					return
				}
				todos = append(todos, todo)
			}

			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(todos); err != nil {
				http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
			}
			return
		}

		// ดึง Todo รายบุคคล
		id, err := strconv.Atoi(idStr)
		if err != nil || id < 1 {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		var todo Todo
		err = db.QueryRow("SELECT id, text, complete FROM todos WHERE id = $1", id).Scan(&todo.ID, &todo.Text, &todo.Complete)
		if err == sql.ErrNoRows {
			http.Error(w, "Todo Not Found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Failed to fetch todo", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(todo); err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		}

	case "POST":
		var newTodo Todo
		if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		err := db.QueryRow("INSERT INTO todos (text, complete) VALUES ($1, $2) RETURNING id", newTodo.Text, newTodo.Complete).Scan(&newTodo.ID)
		if err != nil {
			http.Error(w, "Failed to create todo", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(newTodo); err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		}

	case "DELETE":
		idStr := strings.TrimPrefix(r.URL.Path, "/todo/")
		id, err := strconv.Atoi(idStr)
		if err != nil || id < 1 {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("DELETE FROM todos WHERE id = $1", id)
		if err != nil {
			http.Error(w, "Failed to delete todo", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	case "PUT":
		idStr := strings.TrimPrefix(r.URL.Path, "/todo/")
		id, err := strconv.Atoi(idStr)
		if err != nil || id < 1 {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		var updateTodo Todo
		if err := json.NewDecoder(r.Body).Decode(&updateTodo); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("UPDATE todos SET text = $1, complete = $2 WHERE id = $3", updateTodo.Text, updateTodo.Complete, id)
		if err != nil {
			http.Error(w, "Failed to update todo", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
