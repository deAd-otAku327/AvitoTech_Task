CREATE TABLE IF NOT EXISTS "items"
(
    "id" SERIAL PRIMARY KEY,
    "type" VARCHAR(15) NOT NULL UNIQUE,
    "price" INTEGER NOT NULL
);

INSERT INTO "items" (type, price) VALUES
    ('t-shirt', 80),
    ('cup', 20),
    ('book', 50),
    ('pen', 10),
    ('powerbank', 200),
    ('hoody', 300),
    ('umbrella', 200),
    ('socks', 10),
    ('wallet', 50),
    ('pink-hoody', 500);