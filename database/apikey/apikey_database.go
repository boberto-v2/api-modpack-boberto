package apikey_database

import (
	"time"

	"github.com/brutalzinn/boberto-modpack-api/database"
	entities_apikey "github.com/brutalzinn/boberto-modpack-api/database/apikey/entities"
)

func Insert(apiKey entities_apikey.ApiKey) (string, error) {
	conn, err, ctx := database.OpenConnection()
	if err != nil {
		return "", err
	}
	defer conn.Close(ctx)
	sql := `INSERT INTO users_api_key 
	(key, scopes, app_name, user_id, expire_at, duration, enabled) 
	VALUES ($1, $2, $3, $4, $5, $6, $7) 
	RETURNING id`
	err = conn.QueryRow(
		ctx,
		sql,
		apiKey.Key,
		apiKey.Scopes,
		apiKey.AppName,
		apiKey.UserId,
		apiKey.ExpireAt,
		apiKey.Duration,
		apiKey.Enabled,
	).Scan(&apiKey.ID)
	return apiKey.ID, err
}

func Update(apiKey entities_apikey.ApiKey) (int64, error) {
	conn, err, ctx := database.OpenConnection()
	if err != nil {
		return 0, err
	}
	defer conn.Close(ctx)
	sql := `
		UPDATE users_api_key set
		 key=$1,
		scopes=$2,
		expire_at=$3,
		update_at=$4
		where id=$5`
	res, err := conn.Exec(
		ctx,
		sql,
		apiKey.Key,
		apiKey.Scopes,
		apiKey.ExpireAt,
		time.Now(),
		apiKey.ID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

func Delete(apiKey entities_apikey.ApiKey) (int64, error) {
	conn, err, ctx := database.OpenConnection()
	if err != nil {
		return 0, nil
	}
	defer conn.Close(ctx)
	res, err := conn.Exec(ctx, "DELETE FROM users_api_key WHERE id=$1", apiKey.ID)
	if err != nil {
		return 0, nil
	}
	return res.RowsAffected(), nil
}

func Get(keyId string) (apiKey entities_apikey.ApiKey, err error) {
	conn, err, ctx := database.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close(ctx)
	row := conn.QueryRow(ctx, "SELECT * FROM users_api_key WHERE id=$1", apiKey.ID)
	err = row.Scan(
		apiKey.ID,
		apiKey.Key,
		apiKey.Scopes,
		apiKey.UserId,
		apiKey.ExpireAt,
		apiKey.CreateAt,
		apiKey.UpdateAt)
	return
}

func GetByAppName(appName string) (apiKey entities_apikey.ApiKey, err error) {
	conn, err, ctx := database.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close(ctx)
	row := conn.QueryRow(ctx,
		`SELECT id, key, app_name, scopes, expire_at, enabled, duration
		 FROM users_api_key
		 WHERE app_name=$1`,
		appName,
	)
	err = row.Scan(
		&apiKey.ID,
		&apiKey.Key,
		&apiKey.AppName,
		&apiKey.Scopes,
		&apiKey.ExpireAt,
		&apiKey.Enabled,
		&apiKey.Duration)
	return
}
