package repository

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/neo4j"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetClient(host, user, password, dbname, port string) (*gorm.DB, error) {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func GetNeo4jClient(uri, username, password string) (neo4j.Driver, error) {
	config := func(conf *neo4j.Config) {
		conf.Encrypted = false // Postavite na "true" ako koristite enkripciju
	}

	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""), config)
	if err != nil {
		return nil, err
	}

	return driver, nil
}
