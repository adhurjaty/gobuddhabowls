ALTER TABLE public.inventory_items
    ADD COLUMN category CHARACTER varying(255) NULL;

UPDATE public.inventory_items
    SET category = c.name 
    FROM public.inventory_item_categories AS c
    WHERE inventory_item_category_id = c.id;

ALTER TABLE public.inventory_items
    DROP COLUMN inventory_item_category_id;

DROP TABLE public.inventory_item_categories;
    