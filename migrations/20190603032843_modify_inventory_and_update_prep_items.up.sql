TRUNCATE TABLE count_prep_items;

ALTER TABLE count_prep_items ADD COLUMN count DECIMAL NOT NULL DEFAULT 0;
ALTER TABLE count_prep_items ADD COLUMN inventory_id uuid NOT NULL;
ALTER TABLE count_prep_items DROP COLUMN line_count;
ALTER TABLE count_prep_items DROP COLUMN walk_in_count;
ALTER TABLE count_prep_items ADD CONSTRAINT inventory_id_fk 
    FOREIGN KEY (inventory_id) REFERENCES inventories (id);

ALTER TABLE prep_items DROP COLUMN inventory_item_id;
UPDATE prep_items SET count_unit = 'EACH' WHERE count_unit IS NULL;
