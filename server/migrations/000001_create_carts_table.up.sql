CREATE TABLE IF NOT EXISTS carts (
    id serial NOT NULL,
    user_id integer NOT NULL,
    CONSTRAINT carts_pkey PRIMARY KEY (id)
);