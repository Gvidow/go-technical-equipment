ALTER TABLE orders ADD COLUMN count int DEFAULT 1 CONSTRAINT orders_count_positive CHECK (count > 0);
