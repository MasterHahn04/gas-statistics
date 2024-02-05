package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type DatabaseMariaDB struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	db       *sql.DB
}

func (s *DatabaseMariaDB) GetIDs() ([]string, error) {
	rows, err := s.db.Query("SELECT id from stations")
	if err != nil {
		return nil, err
	}
	var ids []string
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if len(ids) > 10 {
		return nil, fmt.Errorf("too many ids in database /// Max 10 ids are allowed")
	}
	return ids, nil
}

func (s *DatabaseMariaDB) Connect() error {
	var err error
	s.db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", s.User, s.Password, s.Host, s.Port, s.Database))
	if err != nil {
		return err
	}
	return nil
}
