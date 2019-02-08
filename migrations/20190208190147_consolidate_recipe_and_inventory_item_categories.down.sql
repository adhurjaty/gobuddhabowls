CREATE TABLE public.recipe_categories (
    id uuid NOT NULL,
    name character varying(255) NOT NULL,
    background character varying(255) NOT NULL,
    index integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

ALTER TABLE public.recipe_categories OWNER TO postgres;

ALTER TABLE item_categories
    RENAME TO inventory_item_categories;