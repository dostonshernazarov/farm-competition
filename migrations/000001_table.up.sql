CREATE TYPE "product_union" AS ENUM (
    'kilogram',
    'liter',
    'piece'
);

CREATE TYPE "animal_gender" AS ENUM (
    'male',
    'female'
);

CREATE TYPE "store_category" AS ENUM (
    'drug',
    'food',
    'water'
);

CREATE TABLE IF NOT EXISTS category (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS animals (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    category_id UUID NOT NULL,
    gender animal_gender NOT NULL,
    birth_day DATE NOT NULL DEFAULT CURRENT_DATE,
    genus VARCHAR(100),
    weight FLOAT NOT NULL,
    description TEXT,
    is_health BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    FOREIGN KEY (category_id) REFERENCES category(id)
);

CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    union product_union NOT NULL,
    description TEXT,
    total_capacity BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ DEFAULT NULL

);

CREATE TABLE IF NOT EXISTS animal_products (
    id UUID PRIMARY KEY,
    animal_id UUID NOT NULL,
    product_id UUID NOT NULL,
    capacity BIGINT NOT NULL,
    get_time TIMESTAMP NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    FOREIGN KEY (animal_id) REFERENCES animals(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE TABLE IF NOT EXISTS foods (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    capacity BIGINT NOT NULL,
    union product_union NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS drugs (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    capacity BIGINT NOT NULL,
    union product_union NOT NULL,
    status VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS animal_eatable_info (
    id UUID PRIMARY KEY,
    animal_id UUID NOT NULL,
    eatables_id UUID NOT NULL,
    time TIME[],
    category store_category NOT NULL,
    capacity BIGINT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    FOREIGN KEY (animal_id) REFERENCES animals(id)
);

CREATE TABLE IF NOT EXISTS animal_given_eatables (
    id UUID PRIMARY KEY,
    animal_id UUID NOT NULL,
    eatables_id UUID NOT NULL,
    category store_category NOT NULL,
    capacity BIGINT NOT NULL,
    given_time TIMESTAMP NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    FOREIGN KEY (animal_id) REFERENCES animals(id)
);

CREATE TABLE IF NOT EXISTS into_store (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    category store_category NOT NULL,
    capacity BIGINT NOT NULL,
    union product_union NOT NULL,
    time TIMESTAMP NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);
