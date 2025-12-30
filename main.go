package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/order_item", func(c *gin.Context) {
		c.HTML(http.StatusOK, "order_item.html", nil)
	})

	r.GET("/add_item", func(c *gin.Context) {
		c.HTML(http.StatusOK, "add_item.html", nil)
	})

	r.POST("/api/orders", createOrder)
	r.GET("/api/orders/:id", getOrder)
	r.POST("/api/orders/:id/items", addItem)

	log.Println("Server running on :8080")
	r.Run(":8080")
}

// Structs

type CreateOrderReq struct {
	Restaurant string `json:"restaurant"`
	Deadline   string `json:"deadline"`
	Password   string `json:"password"`
}

type Order struct {
	ID         string     `json:"id"`
	Restaurant string     `json:"restaurant"`
	Deadline   *time.Time `json:"deadline,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	Items      []Item     `json:"items"`
}

type Item struct {
	ID       int    `json:"id"`
	UserName string `json:"user_name"`
	ItemName string `json:"item_name"`
	Notes    string `json:"notes"`
}

type AddItemReq struct {
	UserName string `json:"user_name"`
	ItemName string `json:"item_name"`
	Notes    string `json:"notes"`
}

func createOrder(c *gin.Context) {
	var req CreateOrderReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if req.Restaurant == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "restaurant and password required"})
		return
	}

	var deadline *time.Time
	if req.Deadline != "" {
		t, err := time.Parse("2006-01-02T15:04", req.Deadline)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid deadline"})
			return
		}
		deadline = &t
	}

	id := uuid.New()
	_, err := db.Exec(
		"INSERT INTO orders (id, restaurant, deadline, password) VALUES ($1,$2,$3,$4)",
		id, req.Restaurant, deadline, req.Password,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id.String()})
}

func getOrder(c *gin.Context) {
	orderID := c.Param("id")
	password := c.Query("password")
	if password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password required"})
		return
	}

	var order Order
	row := db.QueryRow("SELECT id, restaurant, deadline, created_at FROM orders WHERE id=$1 AND password=$2", orderID, password)
	err := row.Scan(&order.ID, &order.Restaurant, &order.Deadline, &order.CreatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found or wrong password"})
		return
	}

	rows, err := db.Query("SELECT id, user_name, item_name, notes FROM items WHERE order_id=$1", orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	defer rows.Close()

	order.Items = []Item{}
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.UserName, &item.ItemName, &item.Notes)
		if err != nil {
			continue
		}
		order.Items = append(order.Items, item)
	}

	c.JSON(http.StatusOK, order)
}

func addItem(c *gin.Context) {
	orderID := c.Param("id")
	var req AddItemReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if req.UserName == "" || req.ItemName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_name and item_name required"})
		return
	}

	_, err := db.Exec(
		"INSERT INTO items (order_id, user_name, item_name, notes) VALUES ($1,$2,$3,$4)",
		orderID, req.UserName, req.ItemName, req.Notes,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
