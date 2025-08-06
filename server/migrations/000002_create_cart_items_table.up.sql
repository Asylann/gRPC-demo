CREATE TABLE IF NOT EXISTS cart_items (
    id serial NOT NULL,
    cart_id integer NOT NULL,
    product_id integer NOT NULL,
    CONSTRAINT cart_items_pkey PRIMARY KEY (id)
);