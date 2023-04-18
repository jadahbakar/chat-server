package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"
)

type Message struct {
	ID         int    `json:"id"`
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Message    string `json:"message"`
	CreatedAt  string `json:"created_at"`
}

type Conn struct {
	db *sql.DB
}

func main() {
	db, err := sql.Open("sqlite3", "./chat.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	app := fiber.New()
	app.Use(logger)

	// Define the routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, world!")
	})

	app.Post("/messages", func(c *fiber.Ctx) error {
		var message Message
		if err := c.BodyParser(&message); err != nil {
			return err
		}

		// Insert the message into the database
		_, err := db.Exec("INSERT INTO messages(sender_id, receiver_id, message) VALUES(?, ?, ?)", message.SenderID, message.ReceiverID, message.Message)
		if err != nil {
			return err
		}

		// // Get the ID of the new message
		// id, err := result.LastInsertId()
		// if err != nil {
		// 	return err
		// }

		rows, err := db.Query("SELECT * FROM messages ORDER BY created_at")
		if err != nil {
			return err
		}

		var history []Message
		for rows.Next() {
			var t Message
			if err := rows.Scan(&t.ID, &t.SenderID, &t.ReceiverID, &t.Message, &t.CreatedAt); err != nil {
				return err
			}
			history = append(history, t)
		}

		return c.JSON(history)

	})

	log.Fatal(app.Listen(":3000"))

}

func logger(c *fiber.Ctx) error {
	fmt.Printf("%s %s\n", c.Method(), c.Path())
	return c.Next()
}
