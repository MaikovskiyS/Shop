CREATE TABLE IF NOT EXISTS products(
    id BIGSERIAL PRIMARY KEY,
    category_id BIGSERIAL NOT NULL,
    sku VARCHAR(255) NOT NULL,
    name varchar NOT NULL,
    price decimal(18, 2) NOT NULL,
    image varchar,
    created_at timestamptz NOT NULL DEFAULT (now())
    );
CREATE TABLE IF NOT EXISTS categories(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255)
);

INSERT INTO categories(name) VALUES('toys');
INSERT INTO categories(name) VALUES('sneakers');
INSERT INTO categories(name) VALUES('shorts');
INSERT INTO categories(name) VALUES('shoes');



-- CREATE INDEX "products_category_id" ON "products" ("category_id");

-- CREATE INDEX "products_name" ON "products" ("name");