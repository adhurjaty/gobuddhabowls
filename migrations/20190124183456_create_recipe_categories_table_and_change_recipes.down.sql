ALTER TABLE public.recipes
    ADD COLUMN category CHARACTER varying(255) NULL;

UPDATE public.recipes
    SET category = c.name 
    FROM public.recipe_categories AS c
    WHERE recipe_category_id = c.id;

ALTER TABLE public.recipes
    DROP COLUMN recipe_category_id;

DROP TABLE public.recipe_categories;
    