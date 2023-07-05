package apikey_database

import (
	"time"

	"github.com/brutalzinn/boberto-modpack-api/database"
	"github.com/brutalzinn/boberto-modpack-api/database/apikey/entities"
)

func Insert(apiKey entities.ApiKey) {
	conn, err, ctx := database.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close(ctx)
	sql := `INSERT INTO users_api_key 
	(key, scopes, user_id, expire_at) 
	VALUES ($1, $2, $3, $4, $5) 
	RETURNING id`
	err = conn.QueryRow(ctx, sql,
		apiKey.Key,
		apiKey.Scopes,
		apiKey.UserId,
		apiKey.ExpireAt).Scan(&apiKey.Id)
	return
}

func Update(apiKey entities.ApiKey) (int64, error) {
	conn, err, ctx := database.OpenConnection()
	if err != nil {
		return 0, err
	}
	defer conn.Close(ctx)
	sql := `UPDATE users_api_key set 
	key=$1,
	scopes=$2,
	expire_at=$3,
	update_at=$4
	where id=$5`
	res, err := conn.Exec(ctx, sql,
		apiKey.Key,
		apiKey.Scopes,
		apiKey.ExpireAt,
		time.Now(),
		apiKey.Id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

func Delete(apiKey entities.ApiKey) (int64, error) {
	conn, err, ctx := database.OpenConnection()
	if err != nil {
		return 0, nil
	}
	defer conn.Close(ctx)
	res, err := conn.Exec(ctx, "DELETE FROM users_api_key WHERE id=$1", apiKey.Id)
	if err != nil {
		return 0, nil
	}
	return res.RowsAffected(), nil
}

func Get(keyId string) (apiKey entities.ApiKey, err error) {
	conn, err, ctx := database.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close(ctx)
	row := conn.QueryRow(ctx, "SELECT * FROM users_api_key WHERE id=$1", apiKey.Id)
	err = row.Scan(apiKey.Id,
		apiKey.Key,
		apiKey.Scopes,
		apiKey.UserId,
		apiKey.ExpireAt,
		apiKey.CreateAt,
		apiKey.UpdateAt)
	return
}
