WITH category_to_keep AS (
    SELECT id, ic.name FROM item_categories AS ic
        JOIN (
            SELECT MAX(name) AS name, MIN(created_at) AS created_at FROM 
            (
                SELECT ic.* FROM item_categories AS ic
                    JOIN (
                        SELECT name FROM item_categories
                            GROUP BY name HAVING COUNT(*) > 1
                    ) q ON ic.name = q.name
            ) q
        ) q ON q.created_at = ic.created_at AND q.name = ic.name
),
categories_to_remove AS (
    SELECT ic.id, ck.name, ck.id AS id_to_keep FROM item_categories AS ic
        JOIN category_to_keep AS ck ON ck.name = ic.name
    WHERE ic.id != ck.id
)
UPDATE recipes AS r SET category_id = ctr.id_to_keep
    FROM categories_to_remove AS ctr
    WHERE r.category_id = ctr.id;

WITH category_to_keep AS (
    SELECT id, ic.name FROM item_categories AS ic
        JOIN (
            SELECT MAX(name) AS name, MIN(created_at) AS created_at FROM 
            (
                SELECT ic.* FROM item_categories AS ic
                    JOIN (
                        SELECT name FROM item_categories
                            GROUP BY name HAVING COUNT(*) > 1
                    ) q ON ic.name = q.name
            ) q
        ) q ON q.created_at = ic.created_at AND q.name = ic.name
),
categories_to_remove AS (
    SELECT ic.id, ck.name, ck.id AS id_to_keep FROM item_categories AS ic
        JOIN category_to_keep AS ck ON ck.name = ic.name
    WHERE ic.id != ck.id
)
UPDATE inventory_items AS i SET category_id = ctr.id_to_keep
    FROM categories_to_remove AS ctr
    WHERE i.category_id = ctr.id;

WITH category_to_keep AS (
    SELECT id, ic.name FROM item_categories AS ic
        JOIN (
            SELECT MAX(name) AS name, MIN(created_at) AS created_at FROM 
            (
                SELECT ic.* FROM item_categories AS ic
                    JOIN (
                        SELECT name FROM item_categories
                            GROUP BY name HAVING COUNT(*) > 1
                    ) q ON ic.name = q.name
            ) q
        ) q ON q.created_at = ic.created_at AND q.name = ic.name
),
categories_to_remove AS (
    SELECT ic.id, ck.name, ck.id AS id_to_keep FROM item_categories AS ic
        JOIN category_to_keep AS ck ON ck.name = ic.name
    WHERE ic.id != ck.id
)
DELETE FROM item_categories WHERE id IN (
    SELECT id FROM categories_to_remove
);