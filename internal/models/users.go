package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserModelInterface interface {
	Insert(name, email, password string, userType int) error
	InsertCompany(name, logoSvg, logoBg, website string) error
	Authenticate(email, password string) (int, error)
	UserExists(email string) (bool, error)
	GetLastUserCompanyCreated(email, name string) (int, int, error)
	InsertCompanyUser(usrId, compId int) error
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Type           string
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string, userType int) error {
	// return nil

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created, type) 
  VALUES(?, ?, ?, UTC_TIMESTAMP(), ?)`

	_, err = m.DB.Exec(stmt, name, email, hashedPassword, userType)

	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}

		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}

func (m *UserModel) UserExists(email string) (bool, error) {
	return false, nil
}

func (m *UserModel) InsertCompany(name, logoSvg, logoBg, website string) error {

	stmt := `INSERT INTO companies (name, logo_svg, logo_bg_color, website) 
	VALUES(?, ?, ?, ?)`

	_, err := m.DB.Exec(stmt, name, logoSvg, logoBg, website)

	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "companies_uc_name") {
				return ErrDuplicateCompanyName
			}
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "companies_uc_website") {
				return ErrDuplicateCompanyWebsite
			}
		}

		return err
	}

	return nil
}

func (m *UserModel) GetLastUserCompanyCreated(email, name string) (int, int, error) {

	stmtUsers := `SELECT id FROM users WHERE email = ?`
	stmtCompany := `SELECT company_id FROM companies WHERE name = ?`

	rowUsr := m.DB.QueryRow(stmtUsers, email)
	rowComp := m.DB.QueryRow(stmtCompany, name)

	var usrId int
	var compId int
	err := rowUsr.Scan(&usrId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 0, ErrNoRecord
		} else {
			return 0, 0, err
		}
	}
	err = rowComp.Scan(&compId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 0, ErrNoRecord
		} else {
			return 0, 0, err
		}
	}

	return usrId, compId, nil
}

func (m *UserModel) InsertCompanyUser(usrId, compId int) error{
	stmt := `INSERT INTO users_employers (user_id, company_id) VALUES (?, ?)`

	_, err := m.DB.Exec(stmt, usrId, compId)
	if err != nil {
		return err
	}

	return nil
}
