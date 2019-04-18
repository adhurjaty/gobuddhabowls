import { daterange } from '../_datepicker';
import { CategorizedItemsDisplay } from '../components/_categorized_items_display';
import { formatMoney, replaceUrlId } from '../helpers/_helpers';


$(() => {
    var vendorItemsMap = JSON.parse($('#vendor-items-map').attr('data'));
    _tableContainer = $('#categorized-items-display');

    setDateRange();
    setDropdown(vendorItemsMap);
    setFormOnSubmit();
    setDateCheckbox();
});

var _tableContainer;
var _table = null;

var _columnInfo = [
    {
        name: 'id',
        hidden: true,
        get_column: (item) => {
            return item.id;
        }
    },
    {
        name: 'inventory_item_id',
        hidden: true,
        get_column: (item) => {
            return item.inventory_item_id;
        }
    },
    {
        name: 'index',
        hidden: true,
        get_column: (item) => {
            return item.index;
        }
    },
    {
        header: 'Name',
        get_column: (item) => {
            return item.name;
        }
    },
    {
        name: 'price',
        header: 'Price',
        editable: true,
        data_type: 'money',
        get_column: (item) => {
            return formatMoney(parseFloat(item.price));
        },
        set_column: (item, value) => {
            item.price = parseFloat(value);
        }
    },
    {
        name: 'count',
        header: 'Count',
        editable: true,
        data_type: 'number',
        get_column: (item) => {
            return item.count;
        },
        set_column: (item, value) => {
            item.count = parseFloat(value);
        },
        default: 0
    },
    {
        name: 'total_cost',
        header: 'Total Cost',
        get_column: (item) => {
            return formatMoney(parseFloat(item.price) * parseFloat(item.count));
        }
    }
];

var _startDate = null;
var _endDate = null;

function setDateRange() {
    var $inputs = $('#new-order-date, #new-received-date');
    var dates = daterange($inputs, onDateChanged);
    _startDate = dates[0];
    _endDate = dates[1];
}

function onDateChanged(date) {
    if(_startDate == null || _endDate == null) {
        return;
    }

    _endDate.setMinDate(_startDate.getDate());
}

function setDropdown(vendorItemsMap) {
    var dropdownId = '#new-order-vendor';
    var selectedId = $(`${dropdownId} option:selected`).val();
    var options = {
        breakdown: true,
        breakdownTitle: 'Order Breakdown'
    };

    if(selectedId) {
        var allItems = vendorItemsMap[selectedId];
        _table = new CategorizedItemsDisplay(_tableContainer, _columnInfo,
            allItems, options);
    }
    $(dropdownId).change(function(d) {
        // remove none option
        $(`${dropdownId} option[value=""]`).remove();

        var id = $(this).val();
        var items = vendorItemsMap[id];
        _tableContainer.attr('data', JSON.stringify(items));

        _table = new CategorizedItemsDisplay(_tableContainer, _columnInfo,
            null, options);

        var editVendorButton = $('#edit-vendor-button');
        editVendorButton.show();
        var url = replaceUrlId(editVendorButton.attr('base_href'), id);
        editVendorButton.attr('href', url);
    });
}

function setFormOnSubmit() {
    $('#purchase-order-form-submit').closest('form').submit(function(event) {
        debugger;
        if(!isValidOptionSelected()) {
            event.preventDefault();
            showError('must select a vendor');
            return;
        }

        var data = getDataFromTable();
        if(data.length == 0) {
            event.preventDefault();
            showError('must order at least one item');
            return;
        }
        
        removeUncheckedRecDate();
        
        sendOrderItems(data)
    });
}

function isValidOptionSelected() {
    return $('#new-order-vendor option:selected').val().length > 0;
}

function getDataFromTable() {
    return _table.items.filter(x => x.count > 0);
}

function removeUncheckedRecDate() {
    if(!$('#received-order-checkbox').is(':checked')) {
        $('#received-date-input').remove();
    }
}

function sendOrderItems(data) {
    var $input = $('form input[name="Items"]');
    $input.val(JSON.stringify(data));
}

function showError(msg) {
    $('#form-errors').text(msg);
    $('#form-errors').show();
}

function setDateCheckbox() {
    $('#received-order-checkbox').change(function() {
        if(this.checked) {
            $('#received-date-input').show();
        } else {
            $('#received-date-input').hide();
        }
    });
}
