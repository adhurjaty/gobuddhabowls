DELETE FROM recipe_items WHERE inventory_item_id IS NULL AND batch_recipe_id IS NULL;
UPDATE recipe_items SET measure = 'NA'
WHERE measure IS NULL;
