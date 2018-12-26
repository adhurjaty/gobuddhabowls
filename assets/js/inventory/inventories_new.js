import { createInventoryDatagrid } from './_inventory_datagrid';

$(() => {
    setOnSubmit();

    var container = $('#categorized-items-display');
    createInventoryDatagrid(container);
});

function setOnSubmit() {
    var form = $('#inventory-form');
    form.submit((event) => {
        var itemsInput = form.find('input[name="Items"]');
        var datagrid = $('#categorized-items-display');
        itemsInput.val(datagrid.attr('data'));
    });
}
