package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type Database struct {
	db *sql.DB
}

func Initialize() (*Database, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	connect_db := "host=" + viper.GetString("db.host") + " " + "user=" + viper.GetString("db.username") + " " + "port=" + viper.GetString("db.port") + " " + "password=" + viper.GetString("db.password") + " " + "dbname=" + viper.GetString("db.dbname") + " " + "sslmode=" + viper.GetString("db.sslmode")
	db, err := sql.Open(viper.GetString("db.username"), connect_db)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	if err = db.Ping(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &Database{db: db}, nil
}

func (db *Database) SaveURL(originalURL, shortURL string) error {
	// Проверка существования записи с заданным original_url
	var originalURLExists bool
	err := db.db.QueryRow("SELECT EXISTS (SELECT 1 FROM links WHERE original_url = $1)", originalURL).Scan(&originalURLExists)
	if err != nil {
		log.Println("Failed to execute query: ", err)
		return err
	}

	// Проверка существования записи с заданным short_url
	var shortURLExists bool
	err = db.db.QueryRow("SELECT EXISTS (SELECT 1 FROM links WHERE short_url = $1)", shortURL).Scan(&shortURLExists)
	if err != nil {
		log.Println("Failed to execute query: ", err)
		return err
	}

	if !originalURLExists && !shortURLExists {
		// Запрос на вставку данных в таблицу
		query := fmt.Sprintf("INSERT INTO links (original_url, short_url) VALUES ('%s', '%s')", originalURL, shortURL)
		_, err = db.db.Exec(query)
		if err != nil {
			log.Println("Failed to execute insert query: ", err)
			return err
		}
		return nil
	} else {
		log.Println("Such data is already in the database")
		return errors.New("such data is already in the database")
	}
}

func (db *Database) GetOriginalURL(shortURL string) (string, error) {
	// Выполнение запроса на получение original_url
	var originalURL string

	err := db.db.QueryRow("SELECT original_url FROM links WHERE short_url = $1", shortURL).Scan(&originalURL)
	if err == sql.ErrNoRows {
		log.Println("An entry with the specified short_url was not found")
		return "", err
	} else if err != nil {
		log.Println("Failed to execute query: ", err)
		return "", err
	} else {
		return originalURL, nil
	}
}

func (db *Database) CheckShortUrl(shortURL string) error {
	// Проверка существования записи с заданным short_url
	var shortURLExists bool
	err := db.db.QueryRow("SELECT EXISTS (SELECT 1 FROM links WHERE short_url = $1)", shortURL).Scan(&shortURLExists)
	if err != nil {
		log.Println("Failed to execute query: ", err)
		return err
	}
	if shortURLExists {
		return errors.New("short URL is already in the database")
	} else {
		return nil
	}
}

func (db *Database) Close() error {
	return db.db.Close()
}
