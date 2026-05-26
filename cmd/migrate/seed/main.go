package main

import (
	"log"

	"github.com/yusufnuru/gopher_social/internal/db"
	"github.com/yusufnuru/gopher_social/internal/env"
	"github.com/yusufnuru/gopher_social/internal/store"
)

func main() {
	addr := env.GetString(
		"DB_ADDR",
		"postgres://admin:adminpassword@localhost/gopher_social?sslmode=disable",
	)
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	store := store.NewStorage(conn)

	db.Seed(store, conn)
}
