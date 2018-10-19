import pandas as pd
import os
import re
import arrow
import datetime
import uuid
import numpy as np
import psycopg2 as sql

conn = sql.connect("dbname='buddhabowls_development' user='postgres' host='localhost' port='15432' password='mysecretpassword'")


def insert_df_contents(file_df, table, schema):
    cols = file_df.columns.values
    # sql_cols = [mapping[key] for key in cols]
    if 'new_id' not in cols:
        file_df['new_id'] = create_new_id(file_df)
    if file_df['new_id'].str.isdigit().any():
        raise Exception('Wrong ID format')
    file_df['created_at'] = datetime.datetime.now()
    file_df['updated_at'] = datetime.datetime.now()
    new_schema = schema + ['created_at', 'updated_at']
    values = file_df[['new_id'] + new_schema].fillna('')
    bulk_insert(['id'] + new_schema, values, table)


def bulk_insert(col_names, items, table):
    cursor = conn.cursor()
    for _, vals in items.iterrows():
        query = "INSERT INTO %s (%s) VALUES ('%s')" %\
                (table, ','.join(col_names), make_val_string(vals.values))
        cursor.execute(query)
        print('Ran %s' % query)


def make_val_string(values):
    s = '\',\''.join(map(lambda a: clean_string(str(a)), values))
    return s.replace('\'\'', 'NULL')


def clean_string(s):
    return s.replace('\'', '\'\'')


def get_db_df(df, mapping):
    df = df[list(mapping.keys())].rename(index=str, columns=mapping)
    # df['id'] = df.apply(lambda a: uuid.uuid4(), axis=1)
    return df


def conv_date(date_str):
    for form in ['M/D/YYYY HH:mm', 'M/D/YYYY h:mm:ss A', 'MM-DD-YYYY']:
        try:
            return arrow.get(date_str, form).datetime
        except:
            pass


def get_closest_date_vi(joined_order_df, a, c_date):
    filtered = joined_order_df.loc[(pd.notnull(joined_order_df['formatted_date'])) &
                                   (joined_order_df['inventory_item_id'] == a['inventory_item_id'])]
    diff = (filtered['formatted_date'] - c_date)
    date_df = diff[(diff < pd.to_timedelta(0))]
    if not date_df.empty:
        return filtered.ix[[date_df.idxmax()]]['vendor_id'].item()


def create_new_id(df):
    return df.apply(lambda _: uuid.uuid4(), axis=1)


def insert_inventory_items():
    def get_index(item):
        try:
            return inv_item_order_list.index(item['name'])
        except ValueError:
            high_index[0] += 1
            return high_index[0]

    def get_prep_index(item):
        if not pd.isna(item['inventory_item_id']):
            name = db_inv_item_df.loc[db_inv_item_df['id'] == int(item['inventory_item_id']), 'name']
        else:
            name = db_inv_item_df.loc[db_inv_item_df['id'] == int(item['batch_recipe_id']), 'name']
        try:
            return prep_order_list.index(name)
        except ValueError:
            high_index[0] += 1
            return high_index[0]

    # Inventory Items
    inv_items_df = pd.read_csv(os.path.join(os.path.curdir, 'Data', 'InventoryItem.csv'))
    inv_items_mapping = {
        'Id': 'id',
        'Name': 'name',
        'Category': 'category',
        'CountUnit': 'count_unit',
        'RecipeUnit': 'recipe_unit',
        'RecipeUnitConversion': 'recipe_unit_conversion',
        'Yield': 'yield'
    }

    inv_item_schema = [
        'name',
        'category',
        'count_unit',
        'recipe_unit',
        'recipe_unit_conversion',
        'yield',
        'index'
    ]

    db_inv_item_df = get_db_df(inv_items_df, inv_items_mapping)

    with open(os.path.join(os.path.curdir, 'Data', 'Settings', 'InventoryOrder.txt')) as f:
        inv_item_order_list = f.read().split('\n')
    high_index = [len(inv_item_order_list)]
    db_inv_item_df['index'] = db_inv_item_df.apply(get_index, axis=1)
    db_inv_item_df['new_id'] = create_new_id(db_inv_item_df)
    db_inv_item_df['yield'] = db_inv_item_df['yield'].fillna(1)

    insert_df_contents(db_inv_item_df, 'inventory_items', inv_item_schema)

    # Vendors
    vendor_df = pd.read_csv(os.path.join(os.path.curdir, 'Data', 'Vendor.csv'))
    vendor_mapping = {
        'Id': 'id',
        'Name': 'name',
        'Email': 'email',
        'PhoneNumber': 'phone_number',
        'Contact': 'contact',
        'ShippingCost': 'shipping_cost'
    }

    vendor_schema = [
        'name',
        'email',
        'phone_number',
        'contact',
        'shipping_cost'
    ]
    db_vendor_df = get_db_df(vendor_df, vendor_mapping)
    db_vendor_df['new_id'] = create_new_id(db_vendor_df)

    insert_df_contents(db_vendor_df, 'vendors', vendor_schema)

    # Vendor Inventory Items
    vendor_inventory_item_mapping = {
        'Id': 'inventory_item_id',
        'PurchasedUnit': 'purchased_unit',
        'Conversion': 'conversion',
        'LastPurchasedPrice': 'price'
    }

    vendor_item_schema = [
        'inventory_item_id',
        'vendor_id',
        'purchased_unit',
        'conversion',
        'price'
    ]

    dirpath, _, filenames = next(os.walk(os.path.join(os.path.curdir, 'Data', 'Vendors')))
    for filename in filenames:
        vendor_name = re.findall(r'(.*)_.*', filename)[0]
        vendor_item_df = pd.read_csv(os.path.join(dirpath, filename))
        db_vendor_item_df = get_db_df(vendor_item_df, vendor_inventory_item_mapping)
        db_vendor_item_df['vendor_id'] = db_vendor_df.loc[db_vendor_df['name'] == vendor_name, 'new_id'].item()
        db_vendor_item_df['inventory_item_id'] = db_vendor_item_df.apply(
            lambda a: db_inv_item_df.loc[db_inv_item_df['id'] == a['inventory_item_id'], 'new_id'].item(), axis=1)

        insert_df_contents(db_vendor_item_df, 'vendor_items', vendor_item_schema)

    # Inventories
    inv_df = pd.read_csv(os.path.join(os.path.curdir, 'Data', 'Inventory.csv'))
    inv_mapping = {
        'Id': 'id',
        'Date': 'date'
    }

    inv_schema = [
        'date'
    ]
    db_inv_df = get_db_df(inv_df, inv_mapping)
    db_inv_df['new_id'] = create_new_id(db_inv_df)
    db_inv_df['formatted_date'] = db_inv_df.apply(lambda a: conv_date(a['date']).date(), axis=1)

    insert_df_contents(db_inv_df, 'inventories', inv_schema)

    # Orders
    order_df = pd.read_csv(os.path.join(os.path.curdir, 'Data', 'PurchaseOrder.csv'))
    order_mapping = {
        'Id': 'id',
        'OrderDate': 'order_date',
        'ReceivedDate': 'received_date',
        'VendorName': 'vendor_id'
    }

    order_schema = [
        'vendor_id',
        'order_date',
        'received_date'
    ]
    db_order_df = get_db_df(order_df, order_mapping)
    db_order_df['vendor_id'] = db_order_df.apply(
        lambda a: db_vendor_df.loc[db_vendor_df['name'] == a['vendor_id'], 'new_id'][0], axis=1)
    db_order_df['new_id'] = create_new_id(db_order_df)
    db_order_df['formatted_date'] = db_order_df.apply(lambda a: conv_date(a['received_date']), axis=1)

    insert_df_contents(db_order_df, 'purchase_orders', order_schema)

    # Order Items
    order_items_mapping = {
        'Id': 'inventory_item_id',
        'LastPurchasedPrice': 'price',
        'LastOrderAmount': 'count'
    }

    order_item_schema = [
        'inventory_item_id',
        'order_id',
        'price',
        'count'
    ]

    joined_order_df = pd.DataFrame()
    dirpath, _, filenames = next(os.walk(os.path.join(os.path.curdir, 'Data', 'Orders')))
    for filename in filenames:
        order_id = int(re.findall(r'_(.*)\.', filename)[0])
        order_item_df = pd.read_csv(os.path.join(dirpath, filename))
        db_order_item_df = get_db_df(order_item_df, order_items_mapping)
        # throw away data from defunct inventory items
        db_order_item_df = db_order_item_df[(db_order_item_df['inventory_item_id'].isin(db_inv_item_df['id']))]
        db_order_item_df['inventory_item_id'] = db_order_item_df.apply(
            lambda a: db_inv_item_df.loc[db_inv_item_df['id'] == a['inventory_item_id'], 'new_id'].item(), axis=1)
        matching_order = db_order_df.loc[db_order_df['id'] == order_id]
        db_order_item_df['order_id'] = matching_order['new_id'].item()
        date, v_id = matching_order[['formatted_date', 'vendor_id']].values.tolist()[0]
        db_order_item_df['formatted_date'] = date
        db_order_item_df['vendor_id'] = v_id
        joined_order_df = joined_order_df.append(db_order_item_df, ignore_index=True)

        insert_df_contents(db_order_item_df, 'order_items', order_item_schema)

    # Recipes
    recipe_mapping = {
        'Id': 'id',
        'Name': 'name',
        'RecipeUnit': 'recipe_unit',
        'RecipeUnitConversion': 'recipe_unit_conversion',
        'Category': 'category',
        'IsBatch': 'is_batch'
    }

    recipe_schema = [
        'name',
        'recipe_unit',
        'recipe_unit_conversion',
        'category',
        'is_batch',
        'index'
    ]

    recipe_df = pd.read_csv(os.path.join(os.path.curdir, 'Data', 'Recipe.csv'))
    db_recipe_df = get_db_df(recipe_df, recipe_mapping)
    db_recipe_df['index'] = list(range(db_recipe_df.shape[0]))
    db_recipe_df['new_id'] = create_new_id(db_recipe_df)

    insert_df_contents(db_recipe_df, 'recipes', recipe_schema)

    # Recipe Items
    recipe_item_mapping = {
        'Id': 'inventory_item_id',
        'Name': 'name',
        'Measure': 'measure',
        'Quantity': 'count'
    }

    recipe_item_schema = [
        'recipe_id',
        'inventory_item_id',
        'batch_recipe_id',
        'measure',
        'count'
    ]

    def select_or_nan(df, cond, col='new_id'):
        try:
            return df.loc[cond, col].item()
        except:
            return np.nan

    dirpath, _, filenames = next(os.walk(os.path.join(os.path.curdir, 'Data', 'Recipes')))
    for filename in filenames:
        recipe_name = filename.split('.')[0]
        recipe_item_df = pd.read_csv(os.path.join(dirpath, filename))
        db_recipe_item_df = get_db_df(recipe_item_df, recipe_item_mapping)
        recipe_id = db_recipe_df.loc[db_recipe_df['name'] == recipe_name, 'new_id'].iloc[0]
        db_recipe_item_df['recipe_id'] = recipe_id
        db_recipe_item_df['inventory_item_id'] = db_recipe_item_df.apply(
            lambda a: select_or_nan(db_inv_item_df, db_inv_item_df['id'] == a['inventory_item_id']), axis=1)
        db_recipe_item_df['batch_recipe_id'] = db_recipe_item_df.apply(
            lambda a: select_or_nan(db_recipe_df, (db_recipe_df['is_batch']) & (db_recipe_df['name'] == a['name'])
                                    & pd.isna(a['inventory_item_id'])), axis=1)

        insert_df_contents(db_recipe_item_df, 'recipe_items', recipe_item_schema)

    # Prep Items
    prep_items_df = pd.read_csv(os.path.join(os.path.curdir, 'Data', 'PrepItem.csv'))
    prep_items_mapping = {
        'Id': 'id',
        'InventoryItemId': 'inventory_item_id',
        'RecipeItemId': 'batch_recipe_id',
        'Conversion': 'conversion',
        'CountUnit': 'count_unit',
    }

    prep_item_schema = [
        'inventory_item_id',
        'batch_recipe_id',
        'conversion',
        'count_unit',
        'index'
    ]
    db_prep_df = get_db_df(prep_items_df, prep_items_mapping)
    with open(os.path.join(os.path.curdir, 'Data', 'Settings', 'PrepItemOrder.txt')) as f:
        prep_order_list = f.read().split('\n')
    high_index[0] = len(prep_order_list)
    db_prep_df['index'] = db_prep_df.apply(get_prep_index, axis=1)
    db_prep_df['inventory_item_id'] = db_prep_df.apply(
        lambda a: select_or_nan(db_inv_item_df, db_inv_item_df['id'] == a['inventory_item_id']), axis=1)
    db_prep_df['batch_recipe_id'] = db_prep_df.apply(
        lambda a: select_or_nan(db_recipe_df, db_recipe_df['id'] == a['batch_recipe_id']), axis=1)
    db_prep_df['new_id'] = create_new_id(db_prep_df)

    insert_df_contents(db_prep_df, 'prep_items', prep_item_schema)

    # Count Prep Items
    count_prep_items_mapping = {
        'Id': 'prep_item_id',
        'LineCount': 'line_count',
        'WalkInCount': 'walk_in_count'
    }

    count_prep_item_schema = [
        'prep_item_id',
        'line_count',
        'walk_in_count'
    ]

    db_count_prep_df = get_db_df(prep_items_df, count_prep_items_mapping)
    db_count_prep_df['prep_item_id'] = db_count_prep_df.apply(
        lambda a: select_or_nan(db_prep_df, db_prep_df['id'] == a['prep_item_id']), axis=1)

    insert_df_contents(db_count_prep_df, 'count_prep_items', count_prep_item_schema)

    # Count Inventory Items
    count_items_mapping = {
        'Id': 'inventory_item_id',
        'Count': 'count',
    }

    count_inv_item_schema = [
        'inventory_item_id',
        'inventory_id',
        'selected_vendor_id',
        'count'
    ]

    dirpath, _, filenames = next(os.walk(os.path.join(os.path.curdir, 'Data', 'Inventory History')))
    for filename in filenames:
        c_date = conv_date(re.findall(r'_(.*)\.', filename)[0])
        count_items_df = pd.read_csv(os.path.join(dirpath, filename))
        db_count_items_df = get_db_df(count_items_df, count_items_mapping)
        inv_id = db_inv_df.loc[db_inv_df['formatted_date'] == c_date.date(), 'new_id'].item()
        # throw away data from defunct inventory items
        db_count_items_df = db_count_items_df[(db_count_items_df['inventory_item_id'].isin(db_inv_item_df['id']))]
        db_count_items_df['inventory_id'] = inv_id
        db_count_items_df['selected_vendor_id'] = db_count_items_df.apply(
            lambda a: get_closest_date_vi(joined_order_df, a, c_date), axis=1)
        db_count_items_df['inventory_item_id'] = db_count_items_df.apply(
            lambda a: select_or_nan(db_inv_item_df, db_inv_item_df['id'] == a['inventory_item_id']), axis=1)

        insert_df_contents(db_count_items_df, 'count_inventory_items', count_inv_item_schema)

    conn.commit()


insert_inventory_items()
