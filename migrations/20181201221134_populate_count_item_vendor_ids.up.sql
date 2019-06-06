WITH items AS (
    SELECT c.id AS count_inventory_id, po.order_date AS order_date,
        item.name AS item_name, po.vendor_id AS vendor_id
        FROM count_inventory_items AS c
    LEFT JOIN inventory_items AS item ON c.inventory_item_id = item.id
    INNER JOIN inventories AS i ON i.id = c.inventory_id
    LEFT JOIN order_items AS o ON c.inventory_item_id = o.inventory_item_id
    INNER JOIN purchase_orders AS po ON po.id = o.order_id
    WHERE i.date > po.order_date
),
grouped AS (
    SELECT count_inventory_id, MAX(order_date) AS order_date FROM items
    GROUP BY count_inventory_id
)
UPDATE count_inventory_items AS c SET selected_vendor_id = r.vendor_id
FROM (
    SELECT i.item_name, v.name, i.vendor_id AS vendor_id,
        i.count_inventory_id AS count_inventory_id
        FROM items AS i
    INNER JOIN grouped AS g ON i.count_inventory_id = g.count_inventory_id
        AND i.order_date = g.order_date
    INNER JOIN vendors AS v ON v.id = vendor_id
) AS r
WHERE r.count_inventory_id = c.id;