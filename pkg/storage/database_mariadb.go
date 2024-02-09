package storage

import (
	"database/sql"
	"errors"
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

func (s *DatabaseMariaDB) StoreResponse(body string) error {
	_, err := s.db.Exec("INSERT INTO responses(respond) VALUE (?)", body)
	if err != nil {
		return err
	}
	return nil
}

func (s *DatabaseMariaDB) StorePrices(id string, status string, e5 float32, e10 float32, diesel float32) error {
	var (
		statusDB string
		e5DB     float32
		e10DB    float32
		dieselDB float32
	)

	//Check if the last inserted row has the same values as the new one
	err := s.db.QueryRow(fmt.Sprintf("SELECT status, e5, e10, diesel FROM `%s` ORDER BY `index` DESC LIMIT 1", id)).Scan(&statusDB, &e5DB, &e10DB, &dieselDB)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}
	if statusDB == status && e5DB == e5 && e10DB == e10 && dieselDB == diesel {
		return nil
	}

	sqlCmd := fmt.Sprintf("INSERT INTO `%s` (status, e5, e10, diesel) VALUES (?, ?, ?, ?)", id)
	_, err = s.db.Exec(sqlCmd, status, e5, e10, diesel)
	if err != nil {
		sqlCreateCmd := fmt.Sprintf("CREATE TABLE `%s` (`index` bigint(20) NOT NULL AUTO_INCREMENT, `time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP(), `status` varchar(255) NOT NULL, `e5` float(4,3), `e10` float(4,3), `diesel` float(4,3), PRIMARY KEY (`index`))", id)
		fmt.Printf("The Command to create the Table: %s\n", sqlCreateCmd)
		_, err = s.db.Exec(sqlCreateCmd)
		if err != nil {
			return err
		}
		_, err = s.db.Exec(sqlCmd, status, e5, e10, diesel)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *DatabaseMariaDB) Close() error {
	err := s.db.Close()
	if err != nil {
		return err
	}
	return nil
}
