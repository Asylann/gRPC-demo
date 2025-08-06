ALTER TABLE cart_items
    ADD CONSTRAINT fk_carts
    FOREIGN KEY (cart_id) REFERENCES carts(id) ON DELETE CASCADE;