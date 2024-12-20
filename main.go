package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
)

// ตั้งค่า router และ middleware สำหรับเซิร์ฟเวอร์
func initServer() http.Handler {
	http.HandleFunc("/todo/register", register) // register user
	http.HandleFunc("/todo/login", login)       // login user
	http.HandleFunc("/todo/", handleTodo)       // handle todos

	// ตั้งค่า CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://127.0.0.1:5500"},         // อนุญาตให้เข้าถึงจาก localhost:5500
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},  // อนุญาตเฉพาะ HTTP Methods ที่ต้องการ
		AllowedHeaders: []string{"Content-Type", "Authorization"}, // อนุญาต headers ที่ต้องการ
	})

	// ใช้ CORS กับ handler
	return c.Handler(http.DefaultServeMux)
}

func main() {
	handler := initServer() // เรียกใช้ handler จาก initServer()
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
