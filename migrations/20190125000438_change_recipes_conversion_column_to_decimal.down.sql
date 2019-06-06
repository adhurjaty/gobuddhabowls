ALTER TABLE recipes
ADD COLUMN tempCol varchar(255);
UPDATE recipes
SET tempCol = CAST(recipe_unit_conversion AS varchar(255));
ALTER TABLE recipes
DROP COLUMN recipe_unit_conversion;
ALTER TABLE recipes
RENAME COLUMN tempCol TO recipe_unit_conversion;