CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER TABLE count_prep_items;
ALTER TABLE DROP COLUMN count;
ALTER TABLE DROP COLUMN inventory_id;
ALTER TABLE ADD COLUMN line_count DECIMAL NOT NULL DEFAULT 0;
ALTER TABLE ADD COLUMN walk_in_count DECIMAL NOT NULL DEFAULT 0;

ALTER TABLE prep_items ADD COLUMN inventory_item_id NOT NULL DEFAULT uuid_generate_nil();