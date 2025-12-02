-- +goose Up
CREATE TABLE IF NOT EXISTS cart_items (
    user_id    BIGINT NOT NULL REFERENCES users(id)   ON DELETE CASCADE,
    product_id INT    NOT NULL REFERENCES product(id) ON DELETE CASCADE,
    quantity   INT    NOT NULL CHECK (quantity > 0),
    PRIMARY KEY (user_id, product_id)
);

-- +goose Down
DROP TABLE IF EXISTS cart_items;
