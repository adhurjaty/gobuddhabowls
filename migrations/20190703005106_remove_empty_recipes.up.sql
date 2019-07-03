DELETE FROM recipes WHERE id NOT IN
(SELECT DISTINCT(recipe_id) FROM recipe_items)
AND is_batch = TRUE;

DELETE FROM prep_items WHERE batch_recipe_id NOT IN
(SELECT id FROM recipes)