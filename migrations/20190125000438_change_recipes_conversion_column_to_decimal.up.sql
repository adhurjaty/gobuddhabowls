ALTER TABLE recipes
ADD COLUMN tempCol decimal;
UPDATE recipes
SET tempCol = CAST(recipe_unit_conversion AS decimal);
ALTER TABLE recipes
DROP COLUMN recipe_unit_conversion;
ALTER TABLE recipes
RENAME COLUMN tempCol TO recipe_unit_conversion;