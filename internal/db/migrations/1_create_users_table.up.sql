CREATE TABLE IF NOT EXISTS "users"
(
    "id" SERIAL PRIMARY KEY,
    "username" VARCHAR(10) NOT NULL UNIQUE,
    "password" TEXT NOT NULL,
    "balance" INTEGER NOT NULL DEFAULT 1000 CHECK("balance" >= 0),
    "inventory" JSONB NOT NULL DEFAULT '{}'
);