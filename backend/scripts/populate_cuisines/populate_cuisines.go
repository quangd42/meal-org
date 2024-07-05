package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// DBConfig holds the database configuration
type DBConfig struct {
	DBName   string
	User     string
	Password string
	Host     string
	Port     string
}

// getParentID retrieves the ID of the parent row by its name
func getParentID(db *sql.DB, parentName, tableName string) (sql.NullString, error) {
	query := fmt.Sprintf("SELECT id FROM %s WHERE name = $1", tableName)
	var parentID sql.NullString
	err := db.QueryRow(query, parentName).Scan(&parentID)
	if err == sql.ErrNoRows {
		return parentID, nil
	}
	return parentID, err
}

// importCSVToPostgres imports data from a CSV file into a PostgreSQL table
func importCSVToPostgres(csvFilePath string, dbConfig DBConfig, tableName string) error {
	connStr := fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s sslmode=disable",
		dbConfig.DBName, dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	file, err := os.Open(csvFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for i, row := range records {
		if i == 0 {
			continue // Skip header row
		}
		rowID := uuid.New().String()
		name := row[0]
		parentName := row[1]
		createdAt := time.Now()
		updatedAt := createdAt

		var parentID sql.NullString
		if parentName != "" {
			parentID, err = getParentID(db, parentName, tableName)
			if err != nil {
				return err
			}
		}

		insertQuery := fmt.Sprintf(
			"INSERT INTO %s (id, name, parent_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
			tableName,
		)
		_, err = db.Exec(insertQuery, rowID, name, parentID, createdAt, updatedAt)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatal("error loading env file: database")
	}

	dbConfig := DBConfig{
		DBName:   os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
	}

	csvFilePath := "data/cuisines.csv"
	tableName := "cuisines"

	err := importCSVToPostgres(csvFilePath, dbConfig, tableName)
	if err != nil {
		log.Fatalf("error importing CSV to PostgreSQL: %v", err)
	}

	fmt.Println("cuisines data successfully imported!")
}
