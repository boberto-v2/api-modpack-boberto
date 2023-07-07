-- ///TODO: Explain to Daniel why i started the database schema with users and what hell users can have a api key.

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
