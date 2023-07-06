package user_database

import (
	"github.com/brutalzinn/boberto-modpack-api/database"
	"github.com/brutalzinn/boberto-modpack-api/database/user/entities"
)

func Delete(id string) (int64, error) {
	conn, err, ctx := database.OpenConnection()
	if err != nil {
		return 0, err
	}
	defer conn.Close(ctx)
	res, err := conn.Exec(ctx, "DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

func Get(id string) (user *entities.User, err error) {
	conn, err, ctx := database.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close(ctx)
	row := conn.QueryRow(ctx, "SELECT * FROM users WHERE id=$1", id)
	err = row.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.CreateAt,
		&user.UpdateAt)
	return
}

func Insert(user entities.User) (id string, err error) {
	conn, err, ctx := database.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close(ctx)
	sql := `INSERT INTO users 
	(email, password, username) 
	VALUES ($1, $2, $3) 
	RETURNING id`
	err = conn.QueryRow(ctx, sql,
		user.Email,
		user.Password,
		user.Username).Scan(&id)
	return
}
func Update(id int64, user entities.User) (int64, error) {
	conn, err, ctx := database.OpenConnection()
	if err != nil {
		return 0, err
	}
	defer conn.Close(ctx)
	sql := `UPDATE users SET 
	email=$1,
	password=$2,
	username=$3,
	update_at=$4 
	WHERE id=$5`
	res, err := conn.Exec(ctx, sql,
		user.Email,
		user.Password,
		user.Username,
		user.UpdateAt,
		id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

func FindByEmail(email string) (user entities.User, err error) {
	conn, err, ctx := database.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close(ctx)
	row := conn.QueryRow(ctx, "SELECT id, password FROM users WHERE email=$1", email)
	err = row.Scan(user.ID, user.Password)
	return
}
