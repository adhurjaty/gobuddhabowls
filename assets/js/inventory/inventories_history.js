import { parseModelJSON, formatSlashDate, formatMoney, replaceUrlId } from '../helpers/_helpers';
import { CategorizedItemsDisplay } from '../components/_categorized_items_display';
import { datepicker } from '../_datepicker';
import { createInventoryDatagrid } from './_inventory_datagrid';
import { horizontalPercentageChart } from '../_horizontal_percentage_chart';

var _selectedInventory = null;

$(() => {
    var invList = $('#date-list');
    invList.find('li').first().addClass('active');

    _selectedInventory = parseModelJSON(invList.find('li.active')
                                .attr('data'));
    setClickInventory();
    setOnSubmit();

    var container = $('#categorized-items-display');
    createInventoryDatagrid(container, onDataGridEdit);
    createBreakdown()
});

function setClickInventory() {
    var listItems = $('#date-list').find('li');
    listItems.each((i, el) => {
        var $el = $(el);
        $el.click(() => {
            listItems.removeClass('active');
            $el.addClass('active');
            setSelectedInventory();
            clearItemsInput();
        });
    });
}

function setSelectedInventory() {
    var $form = $('#inventory-form');
    var invList = $('#date-list');
    var container = $('#categorized-items-display');
    _selectedInventory = parseModelJSON(invList.find('li.active')
                                .attr('data'));
    container.attr('data', JSON.stringify(_selectedInventory.Items));
    var dateInput = $form.find('input[name="Date"]');
    dateInput.val(formatSlashDate(_selectedInventory.time));

    var deleteLink = $('#delete-inventory');
    var url = replaceUrlId(deleteLink.attr('data-link'), _selectedInventory.id);
    deleteLink.attr('href', url);

    var table = createInventoryDatagrid(container, onDataGridEdit);
    createBreakdown();
}

function clearItemsInput() {
    var input = $('input[name="Items"]');
    input.val("");
}

function createBreakdown() {
    var container = $('#category-breakdown');
    var title = 'Inventory Breakdown';
    var items = parseModelJSON($('#categorized-items-display').attr('data'));
    items.forEach(item => {
        if(item.conversion > 0) {
            item.price = item.price / item.conversion;
        } else {
            item.price = 0;
        }
    });
    var total = items.reduce((total, item) => {
        return total + item.price * item.count;
    }, 0);
    if(total != 0) {
        container.html(horizontalPercentageChart(title, items, total));
    } else {
        bdContainer.html('');
    }
}

function setOnSubmit() {
    var form = $('#inventory-form');
    form.submit((event) => {
        var url = replaceUrlId(form.attr('action'), _selectedInventory.id);
        form.attr('action', url);
    });
}

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