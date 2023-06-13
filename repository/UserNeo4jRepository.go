package repository

import (
	neo4j "github.com/neo4j/neo4j-go-driver/neo4j"
)

type Neo4jUserRepository struct {
	Session neo4j.Session
}

func NewNeo4jUserRepository(driver neo4j.Driver) *Neo4jUserRepository {
	session, err := driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return nil
	}

	return &Neo4jUserRepository{
		Session: session,
	}
}

func (repo *Neo4jUserRepository) Close() {
	repo.Session.Close()
}

func (repo *Neo4jUserRepository) SaveUser(userId string) error {
	_, err := repo.Session.Run("CREATE (:User {idInPostgre: $userId})",
		map[string]interface{}{
			"userId": userId,
		})

	if err != nil {
		println("EROR PRILIKOM UPISA U NEO4J: USER REPOSITORY")
		println(err.Error())
		return err
	}

	println("USPJESNO SACUVAN USER U NEO4J")

	return nil
}
