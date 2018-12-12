import { parseModelJSON, formatSlashDate, formatMoney } from '../helpers/_helpers';
import { CategorizedItemsDisplay } from '../components/_categorized_items_display';
import { datepicker } from '../_datepicker';
import { createInventoryDatagrid } from './_inventory_datagrid';

$(() => {
    var dateInput = $('#inventory-form').find('input[name="Date"]');
    datepicker(dateInput[0]);

    var container = $('#categorized-items-display');
    
    createInventoryDatagrid(container, onDataGridEdit);
});

function onDataGridEdit(item) {
    var form = $('#inventory-form');
    var itemsInput = form.find('input[name="Items"]');
    var editedItems = [item];
    if(itemsInput.val()) {
        editedItems = JSON.parse(itemsInput.val());
        var idx = editedItems.findIndex(x => x.id == item.id);
        if(idx > -1) {
            editedItems[idx] = item;
        } else {
            editedItems.push(item);
        }
    }

    itemsInput.val(JSON.stringify(editedItems));
}