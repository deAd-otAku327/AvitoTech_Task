CREATE TABLE IF NOT EXISTS "coin_transfers" 
(
    "id" SERIAL PRIMARY KEY,
    "from_user_id" INTEGER NOT NULL REFERENCES users(id),
    "to_user_id" INTEGER NOT NULL REFERENCES users(id),
    "amount" INTEGER NOT NULL,
    "timing" TIMESTAMP NOT NULL
);