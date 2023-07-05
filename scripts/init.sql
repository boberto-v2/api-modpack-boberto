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
	"scopes" text NOT NULL,
	"user_id" uuid NOT NULL,
    "description" text NOT NULL,
    "expire_at" timestamp default NULL,
	"create_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
	"update_at" timestamp default NULL,
    CONSTRAINT users_api_key_pkey PRIMARY KEY (id)
);
