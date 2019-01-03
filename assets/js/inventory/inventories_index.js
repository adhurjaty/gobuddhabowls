import { createInventoryDatagrid } from './_inventory_datagrid';
import { parseModelJSON, getObjectDiff, replaceUrlId } from '../helpers/_helpers';
import { sendAjax } from '../helpers/index_helpers';

$(() => {
    var container = $('#categorized-items-display');
    createInventoryDatagrid(container, onDataGridEdit);
});

function onDataGridEdit(item) {
    var form = $('#inventory-item-form');
    var gridContainer = $('#categorized-items-display');

    var allItems = parseModelJSON(gridContainer.attr('data'));
    var oldItemIdx = allItems.findIndex(x => x.id == item.id);
    var oldItem = allItems[oldItemIdx];
    debugger;
    var props = getObjectDiff(item, oldItem);

    form.html('');
    props.forEach(prop => {
        addInput(form, prop);
    });

    submitForm(form, item.inventory_item_id);

    allItems[oldItemIdx] = item;
    gridContainer.attr('data', JSON.stringify(allItems));
}

function addInput(form, name) {
    var input = $(`<input type="text" name="${name}" style="display: none;" />`);
    form.append(input);
}

function submitForm(form, id) {
    var templatePath = form.attr('action');
    var actionPath = replaceUrlId(templatePath, id);

    form.attr('action', actionPath);
    sendAjax(form);
}