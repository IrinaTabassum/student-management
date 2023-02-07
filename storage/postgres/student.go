package postgres

import (
	"fmt"
	"log"
	"student-management/storage"

	"golang.org/x/crypto/bcrypt"
)

const listQuery = `SELECT * from students WHERE deleted_at IS NULL ORDER BY id DESC`

func (ps PostgresStorage) ListofStudent() ([]storage.Student, error) {
	var listStudent []storage.Student

	if err := ps.DB.Select(&listStudent, listQuery); err != nil {
		log.Println(err)
		return nil, err
	}
	return listStudent, nil
}

const insertSQuery = `
	INSERT INTO students(
		first_name,
		last_name,
		username,
		email,
		password
	) VALUES (
		:first_name,
		:last_name,
		:username,
		:email,
		:password
	) RETURNING * ;	
	`

func (ps PostgresStorage) CreateStudent(s storage.Student) (*storage.Student, error) {

	stmt, err := ps.DB.PrepareNamed(insertSQuery)
	if err != nil {
		log.Fatal(err)
	}
	hashPass, err := bcrypt.GenerateFromPassword([]byte(s.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	s.Password = string(hashPass)
	if err = stmt.Get(&s, s); err != nil {
		log.Fatal(err)
		return nil, err
	}
	if s.ID == 0 {
		return nil, fmt.Errorf("Unable to insert")
	}

	return &s, nil
}

const updatStudenrQuery = `
	UPDATE students SET
		first_name = :first_name,
		last_name = :last_name,
		status = :status
	WHERE id = :id AND deleted_at IS NULL
	RETURNING *;
	`

func (ps PostgresStorage) UpdateStudent(s storage.Student) (*storage.Student, error) {

	stmt, err := ps.DB.PrepareNamed(updatStudenrQuery)
	if err != nil {
		log.Fatal(err)
	}
	if err = stmt.Get(&s, s); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &s, nil
}

const getStudentByIDQuery = `
SELECT * from students WHERE id = $1 AND deleted_at IS NULL`

func (ps PostgresStorage) GetStudentByID(id int) (*storage.Student, error) {
	var s storage.Student
	err := ps.DB.Get(&s, getStudentByIDQuery, id)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &s, nil
}

const deleteStudentByIDQuery = `
UPDATE students SET
	deleted_at = CURRENT_TIMESTAMP
	WHERE id = $1 AND deleted_at IS NULL
	RETURNING *;
	`

func (ps PostgresStorage) DeleteStudentByID(id int) error {
	res, err := ps.DB.Exec(deleteStudentByIDQuery, id)
	if err != nil {
		log.Fatal(err)
	}
	rowCounr, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	if rowCounr <= 0 {
		return fmt.Errorf("unable to delate student")
	}
	return nil
}
