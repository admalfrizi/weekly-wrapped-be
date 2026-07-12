package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func ptr(s string) *string {
	return &s
}

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		_ = godotenv.Load(".env")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set in your .env file")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	categories := []struct {
		Name     string
		Icon     *string
		ColorHex *string
	}{
		{Name: "Coding"},   
		{Name: "Reading"},
		{Name: "Exercise"},
		{Name: "Gaming"},  
		{Name: "Working"},   
	}

	fmt.Println("Starting database seeder...")

	for _, c := range categories {
		catID := uuid.New()
		now := time.Now()

		query := `
			INSERT INTO categories (id, name, created_at)
			VALUES ($1, $2, $3)
			ON CONFLICT (name) DO NOTHING;
		`

		tag, err := pool.Exec(ctx, query, catID, c.Name, now)
		if err != nil {
			log.Fatalf("Failed to insert category %s: %v\n", c.Name, err)
		}

		if tag.RowsAffected() > 0 {
			fmt.Printf("✅ Inserted: %s (ID: %s)\n", c.Name, catID)
		} else {
			fmt.Printf("⏭️  Skipped: %s (Already exists)\n", c.Name)
		}
	}

	fmt.Println("Seeding complete!")
}