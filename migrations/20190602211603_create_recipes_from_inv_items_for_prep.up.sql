INSERT INTO recipes (id, name, recipe_unit, is_batch, category_id, index, 
    recipe_unit_conversion, created_at, updated_at)
SELECT uuid_generate_v1(), ii.name, 'EACH' AS recipe_unit, TRUE AS is_batch,
    ii.category_id, 0 AS index, ii.recipe_unit_conversion, NOW(), NOW()
    FROM prep_items AS pi
JOIN inventory_items AS ii ON ii.id = pi.inventory_item_id
WHERE pi.inventory_item_id IS NOT NULL;
