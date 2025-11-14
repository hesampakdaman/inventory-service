CREATE TABLE IF NOT EXISTS products (
    product_id uuid PRIMARY KEY,
    available int,
    reserved int,
    title text,
    description text
);

CREATE TABLE IF NOT EXISTS reservations (
    product_id uuid,
    reservation_id uuid,
    qty int,
    state text,
    PRIMARY KEY (product_id, reservation_id)
);

CREATE TABLE IF NOT EXISTS requests (
    product_id uuid,
    request_id uuid,
    created_at timestamp,
    PRIMARY KEY (product_id, request_id)
);
