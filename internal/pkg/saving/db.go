package saving

import (
	"database/sql"
	"log"
	"time"
)

const (
	queryCreateTable = `CREATE TABLE IF NOT EXISTS core (
		version bigserial PRIMARY KEY,
		timestamp bigint NOT NULL,
		payload JSONB NOT NULL
	)`

	queryDeleteOld = `DELETE FROM core
		WHERE version NOT IN (
			SELECT version FROM core
			ORDER BY timestamp DESC
			LIMIT 5
		)`
	querySave = `INSERT INTO core (timestamp, payload) VALUES ($1, $2)`

	queryVacuum = `VACUUM core`
)

type StorageDB struct {
	Db *sql.DB
}

func NewStorageDB(url string) (*StorageDB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(queryCreateTable)
	if err != nil {
		return nil, err
	}

	return &StorageDB{Db: db}, nil
}

func (s *StorageDB) SaveVersion(data []byte) error {
	timestamp := time.Now().Unix()
	_, err := s.Db.Exec(querySave, timestamp, data)
	if err != nil {
		log.Println("Ошибка сохранения версии:", err)
		return err
	}

	_, err = s.Db.Exec(queryVacuum)
	if err != nil {
		log.Println("Ошибка очистки:", err)
		return err
	}

	_, err = s.Db.Exec(queryDeleteOld)
	return err
}
