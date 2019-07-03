import { createInventoryDatagrid } from './_inventory_datagrid';
import { createPrepItemDatagrid } from './_prep_item_datagrid';

$(() => {
    setOnSubmit();

    var container = $('#categorized-items-display');
    createInventoryDatagrid(container);
    var prepItemsContainer = $('#categorized-prep-items-display');
    createPrepItemDatagrid(prepItemsContainer);
});

function setOnSubmit() {
    var form = $('#inventory-form');
    form.submit((event) => {
        var itemsInput = form.find('input[name="Items"]');
        var datagrid = $('#categorized-items-display');
        itemsInput.val(datagrid.attr('data'));
    });
}
