package database

import (
	"context"
	"database/sql"
	"fmt"
	"pvz_service/logger"
	"time"

	_ "github.com/lib/pq"
)

type DBConnection struct {
	DB *sql.DB
	DbName string
	URL string
	BaseURL string
}

func (c *DBConnection) createDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
	defer cancel()

	db, err := sql.Open("postgres", c.BaseURL)
	if err != nil {
		return err
	}

	defer db.Close()

	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = '%s')", c.DbName)
	if err = db.QueryRowContext(ctx, query).Scan(&exists); err != nil {
		return err
	}

	if !exists {
		if _, err := db.Exec("CREATE DATABASE " + c.DbName); err != nil {
			return err
		}
	}

	return nil
}

func (c *DBConnection) InitPostgresConn() error {
	if err := c.createDB(); err != nil {
		logger.Err.Println("error occured during db creation - ", err)
		return err
	}

	db, err := sql.Open("postgres", c.URL)

	if err != nil {
		logger.Err.Println("can't establish connection with database - ", err)
		return nil
	}

	c.DB = db

	logger.Debug.Println("succesfully established database connection")
	return db.Ping()
}
