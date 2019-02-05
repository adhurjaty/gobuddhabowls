CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE public.recipe_categories(
    id uuid NOT NULL,
	name character varying(255) NOT NULL,
	background character varying(255) NOT NULL,
    index integer NOT NULL, 
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);
ALTER TABLE public.recipe_categories OWNER TO postgres;
ALTER TABLE ONLY public.recipe_categories
    ADD CONSTRAINT recipe_categories_pkey PRIMARY KEY (id);

ALTER TABLE public.recipes
    ADD COLUMN recipe_category_id uuid NULL;

INSERT INTO public.recipe_categories (name, id, background, index, created_at, updated_at)
    SELECT c.category,
        uuid_generate_v4() as id,
        '#FFFFFF' as background,
        0 as index,
        current_timestamp,
        current_timestamp 
    FROM (SELECT DISTINCT(category) FROM public.recipes) AS c;

UPDATE public.recipes
    SET recipe_category_id = c.id
    FROM public.recipe_categories AS c
    WHERE category = c.name;

ALTER TABLE public.recipes
    DROP COLUMN category;