CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    "id" uuid DEFAULT uuid_generate_v4(),
	"email" varchar(100),
    "password" varchar(100),
    "username" TEXT,
    "create_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    "update_at" timestamp default NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS modpacks (
    "id" uuid NOT NULL,
	"name" varchar(100),
    "client_manifest" text,
    "server_manifest" text,
    "user_id" uuid NOT NULL,
    "create_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    "update_at" timestamp default NULL,
    CONSTRAINT modpacks_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS users_api_key (
    "id" uuid DEFAULT uuid_generate_v4(),
	"key" text NOT NULL,
    "enabled" boolean default true,
	"scopes" text NOT NULL,
    "app_name" text NOT NULL,
	"user_id" uuid NOT NULL,
    "duration" int default 0,
    "expire_at" timestamp default NULL,
	"create_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
	"update_at" timestamp default NULL,
    CONSTRAINT users_api_key_pkey PRIMARY KEY (id)
);


ALTER TABLE users_api_key
ADD FOREIGN KEY ("user_id") REFERENCES users (id)
ON DELETE CASCADE
DEFERRABLE INITIALLY DEFERRED;

INSERT INTO users ("id", "email", "password", "username") VALUES('ab7d7136-6c24-4cd0-ba30-97ff0110ecac'::uuid, 'test', '$2a$15$zqJBJTKH7LZbTSQHhnNzeOx9VjcGwv3HamUksu8VQ81E/WbRJCLPW', 'usertest');