INSERT INTO inventory_item_categories (id, name, background, index, created_at, updated_at)
    SELECT id, name, background, index, created_at, updated_at FROM
        recipe_categories;

DROP TABLE recipe_categories;

ALTER TABLE inventory_item_categories
    RENAME TO item_categories;