package database

import (
    "context"
    "fmt"
    "github.com/jackc/pgx/v4"
)

// Функция для создания таблицы product, если её нет
func CreateProductTableIfNotExist(conn *pgx.Conn) error {
    createProductTable := `
    CREATE TABLE IF NOT EXISTS product (
        id SERIAL PRIMARY KEY,
        manufacturer_id SMALLINT NOT NULL,
        name VARCHAR(255) NOT NULL,
        alias VARCHAR(255) NOT NULL,
        short_description TEXT NOT NULL,
        description TEXT NOT NULL,
        price DECIMAL(20,2) NOT NULL,
        image VARCHAR(255) NOT NULL,
        available SMALLINT DEFAULT 1,
        meta_keywords VARCHAR(255),
        meta_description VARCHAR(255),
        meta_title VARCHAR(255)
    );`

    _, err := conn.Exec(context.Background(), createProductTable)
    if err != nil {
        return fmt.Errorf("failed to create product table: %v", err)
    }

    return nil
}

// Функция для создания таблицы product_properties, если её нет
func CreateProductPropertiesTableIfNotExist(conn *pgx.Conn) error {
    createProductPropertiesTable := `
    CREATE TABLE IF NOT EXISTS product_properties (
        id SERIAL PRIMARY KEY,
        product_id INT REFERENCES product(id),
        property_name VARCHAR(255) NOT NULL,
        property_value VARCHAR(255) NOT NULL,
        property_price DECIMAL(20,2)
    );`

    _, err := conn.Exec(context.Background(), createProductPropertiesTable)
    if err != nil {
        return fmt.Errorf("failed to create product_properties table: %v", err)
    }

    return nil
}

// Функция для создания таблицы product_images, если её нет
func CreateProductImagesTableIfNotExist(conn *pgx.Conn) error {
    createProductImagesTable := `
    CREATE TABLE IF NOT EXISTS product_images (
        id SERIAL PRIMARY KEY,
        product_id INT REFERENCES product(id),
        image VARCHAR(255) NOT NULL,
        title VARCHAR(255) NOT NULL
    );`

    _, err := conn.Exec(context.Background(), createProductImagesTable)
    if err != nil {
        return fmt.Errorf("failed to create product_images table: %v", err)
    }

    return nil
}
