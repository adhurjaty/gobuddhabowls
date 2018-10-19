CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE public.inventory_item_categories(
    id uuid NOT NULL,
	name character varying(255) NOT NULL,
	background character varying(255) NOT NULL,
    index integer NOT NULL, 
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);
ALTER TABLE public.inventory_item_categories OWNER TO postgres;
ALTER TABLE ONLY public.inventory_item_categories
    ADD CONSTRAINT inventory_item_categories_pkey PRIMARY KEY (id);

ALTER TABLE public.inventory_items
    ADD COLUMN inventory_item_category_id uuid NULL;

INSERT INTO public.inventory_item_categories (name, id, background, index, created_at, updated_at)
    SELECT c.category,
        uuid_generate_v4() as id,
        '#FFFFFF' as background,
        0 as index,
        current_timestamp,
        current_timestamp 
    FROM (SELECT DISTINCT(category) FROM public.inventory_items) AS c;

UPDATE public.inventory_items
    SET inventory_item_category_id = c.id
    FROM public.inventory_item_categories AS c
    WHERE category = c.name;

ALTER TABLE public.inventory_items
    DROP COLUMN category;