package main

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"log"
)


var (
	db *sql.DB
    err error
)

func main() {
	db, err = sql.Open("mysql", "[user]:[pass]@/[dbname]")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("About to close db connection")
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	app := fiber.New()
	app.Static("/", "index.html")
	app.Get("/signup", func (c *fiber.Ctx) error {
		username := c.FormValue("username")
		password := c.FormValue("password")
		var user string
		err := db.QueryRow("SELECT username FROM users WHERE username=?", username).Scan(&user)
		switch {
		case err == sql.ErrNoRows:
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				return c.Status(500).SendString("Server error, unable to create your account.")
			}

			_, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, hashedPassword)
			if err != nil {
				return c.Status(500).SendString("Server error, unable to create your account.")
			}

			return c.SendString("User created.")
		case err != nil:
			return c.Status(500).SendString("Server error, unable to create your account.")
		}
		return c.Status(301).Redirect("/")
	})
	app.Get("/login", func (c *fiber.Ctx) error {
		username := c.FormValue("username")
		password := c.FormValue("password")
		var (
			databaseUsername string
			databasePassword string
		)
		err := db.QueryRow("SELECT username, password FROM users WHERE username=?", username).Scan(&databaseUsername, &databasePassword)

		if err != nil {
			return c.Redirect("/login", 301)
		}

		err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
		if err != nil {
			return c.Redirect("/login", 301)
		}
		return c.SendString("Hello" + databaseUsername)
	})
	log.Fatal(app.Listen(":80"))

	if err != nil { panic(err.Error()) }
}
