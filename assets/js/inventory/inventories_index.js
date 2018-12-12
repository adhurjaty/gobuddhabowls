import { parseModelJSON, formatSlashDate, formatMoney, replaceUrlId } from '../helpers/_helpers';
import { CategorizedItemsDisplay } from '../components/_categorized_items_display';
import { datepicker } from '../_datepicker';

var _selectedInventory = null;

$(() => {
    var invList = $('#date-list');
    invList.find('li').first().addClass('active');

    _selectedInventory = parseModelJSON(invList.find('li.active')
                                .attr('data'));
    setClickInventory();
    setOnSubmit();
});

function setClickInventory() {
    var listItems = $('#date-list').find('li');
    listItems.each((i, el) => {
        var $el = $(el);
        $el.click(() => {
            listItems.removeClass('active');
            $el.addClass('active');
            setSelectedInventory();
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

    var table = new CategorizedItemsDisplay(container, _columns,
        null, _categorizedOptions);
}

function setOnSubmit() {
    var form = $('#inventory-form');
    form.submit((event) => {
        var url = replaceUrlId(form.attr('action'), _selectedInventory.id);
        form.attr('action', url);
    });
}