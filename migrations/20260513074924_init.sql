-- +goose Up
-- +goose StatementBegin

-- ==========================================================
-- 0. ENUM TYPES (PostgreSQL)
-- ==========================================================
DO $do$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'payment_status_enum') THEN
        CREATE TYPE payment_status_enum AS ENUM ('pending', 'paid', 'expired', 'cancelled');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'order_status_enum') THEN
        CREATE TYPE order_status_enum AS ENUM ('waiting_payment', 'processed', 'shipped', 'delivered', 'cancelled', 'completed', 'pending');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'discount_type_enum') THEN
        CREATE TYPE discount_type_enum AS ENUM ('fixed', 'percentage');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role_enum') THEN
        CREATE TYPE user_role_enum AS ENUM ('customer', 'admin');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'provider_enum') THEN
        CREATE TYPE provider_enum AS ENUM ('local', 'google', 'github');
    END IF;
END $do$;

CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


-- ==========================================================
-- 1. MODUL: AUTH & USERS (Supports OAuth2)
-- ==========================================================
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    password TEXT, -- NULL jika login via OAuth
    full_name VARCHAR(100),
    role user_role_enum DEFAULT 'customer',
    provider provider_enum DEFAULT 'local',
    provider_id VARCHAR(255),
    avatar_url TEXT,
    session TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Index umum untuk soft delete
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- ==========================================================
-- 3. MODUL: PRODUCTS & INVENTORY
-- ==========================================================
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    slug VARCHAR(50) UNIQUE NOT NULL,
    parent_id INT REFERENCES categories(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_categories_deleted_at ON categories(deleted_at);

CREATE TABLE IF NOT EXISTS products (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255)   NOT NULL,
    description     TEXT           NOT NULL DEFAULT '',
    category_id     INT            NOT NULL REFERENCES categories(id),
    price           DECIMAL(15, 2) NOT NULL,
    original_price  DECIMAL(15, 2),
    image           TEXT           NOT NULL,
    rating          DECIMAL(2, 1)  NOT NULL DEFAULT 0,
    reviews         INT            NOT NULL DEFAULT 0,
    stock           INT            NOT NULL DEFAULT 0 CHECK (stock >= 0),
    in_stock        BOOLEAN        NOT NULL DEFAULT true,
    featured        BOOLEAN        NOT NULL DEFAULT false,
    tags            TEXT[]         NOT NULL DEFAULT '{}',
    created_at      TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_new_product_tbl_category_id ON products (category_id);
CREATE INDEX IF NOT EXISTS idx_new_product_tbl_featured    ON products (featured);
CREATE INDEX IF NOT EXISTS idx_new_product_tbl_deleted_at  ON products (deleted_at);
CREATE INDEX IF NOT EXISTS idx_new_product_tbl_tags        ON products USING GIN (tags);

CREATE TABLE IF NOT EXISTS product_variants (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id        UUID           NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    name              VARCHAR(100)   NOT NULL,
    type              VARCHAR(50)    NOT NULL,
    value             VARCHAR(100)   NOT NULL,
    price             DECIMAL(15, 2),
    price_adjustment  DECIMAL(15, 2),
    stock             INT            NOT NULL DEFAULT 0 CHECK (stock >= 0),
    created_at        TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at        TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_new_variants_tbl_product_id  ON product_variants (product_id);
CREATE INDEX IF NOT EXISTS idx_new_variants_tbl_type        ON product_variants (type);
CREATE INDEX IF NOT EXISTS idx_new_variants_tbl_deleted_at  ON product_variants (deleted_at);

CREATE TABLE IF NOT EXISTS articles (
    id VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    content TEXT,
    image VARCHAR(255),
    author VARCHAR(255),
    category VARCHAR(255),
    published_at TIMESTAMP WITH TIME ZONE,
    featured BOOLEAN DEFAULT FALSE,
    read_time INTEGER DEFAULT 0,
    excerpt TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (id)
);

-- Membuat index untuk optimasi query
CREATE INDEX IF NOT EXISTS idx_posts_category ON articles(category);
CREATE INDEX IF NOT EXISTS idx_posts_featured ON articles(featured);
CREATE INDEX IF NOT EXISTS idx_posts_published_at ON articles(published_at);

-- Fungsi untuk otomatis memperbarui updated_at (Postgres spesifik)
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger untuk menjalankan fungsi di atas sebelum data di-update
CREATE TRIGGER update_posts_modtime
    BEFORE UPDATE ON articles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE IF NOT EXISTS carts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id int NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Relasi ke tabel users (pastikan tabel users sudah ada)
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

---
-- 2. Tabel Cart Items
---
CREATE TABLE IF NOT EXISTS cart_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cart_id UUID NOT NULL,
    product_id UUID NOT NULL, -- Asumsi product_id juga menggunakan UUID
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    selected_variants JSONB DEFAULT '[]'::jsonb, -- Array of objects disimpan dalam JSONB
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Relasi ke tabel carts
    CONSTRAINT fk_cart FOREIGN KEY (cart_id) REFERENCES carts(id) ON DELETE CASCADE,
    CONSTRAINT fk_product_cart FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
    -- Anda juga bisa menambahkan FK ke tabel products di sini jika perlu
);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 3. Create Orders Table
CREATE TABLE orders (
    id VARCHAR(255) PRIMARY KEY,
    customer_name VARCHAR(255) NOT NULL,
    customer_email VARCHAR(255) NOT NULL,
    customer_phone VARCHAR(255) NOT NULL,
    total FLOAT8 NOT NULL, -- FLOAT8 is equivalent to double precision in Postgres
    status payment_status_enum DEFAULT 'pending',
    payment_url VARCHAR(255),
    user_id int,
    customer_address TEXT NOT NULL,
    city VARCHAR(255) NOT NULL,
    state VARCHAR(255) NOT NULL,
    zip_code VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_orders_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- 4. Create Order Items Table
CREATE TABLE order_items (
    id VARCHAR(255) PRIMARY KEY,
    order_id VARCHAR(255) NOT NULL,
    product_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    price FLOAT8 NOT NULL,
    quantity INT NOT NULL,
    image VARCHAR(255),
    selected_variants JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_order_items_order_id FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);


CREATE TRIGGER update_orders_updated_at BEFORE UPDATE ON orders FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_order_items_updated_at BEFORE UPDATE ON order_items FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
-- Drop tables (reverse dependency order)
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS new_carts;
DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS articles;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS product_variants CASCADE;
DROP TABLE IF EXISTS products  CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Drop types
DO $do$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'order_status_enum') THEN
        DROP TYPE order_status_enum;
    END IF;

    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'payment_status_enum') THEN
        DROP TYPE payment_status_enum;
    END IF;

    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'discount_type_enum') THEN
        DROP TYPE discount_type_enum;
    END IF;

    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'provider_enum') THEN
        DROP TYPE provider_enum;
    END IF;

    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role_enum') THEN
        DROP TYPE user_role_enum;
    END IF;
END $do$;

-- +goose StatementEnd