import { datepicker } from './datepicker.js';

var _selected_tr;

$(() => {
    datepicker($('#new-order-date'), {
        autoclose: 'true',
        format: 'mm/dd/yyyy'
    });

    $('#new-order-vendor').change(function(d) {
        // remove none option
        // $('#new-order-vendor option[value=""]').remove();
        // var id = $(this).val();
        // $.ajax({
        //     url: '/purchase_orders/order_vendor_changed/' + id,
        //     method: 'GET',
        //     error: function(xhr, status, err) {
        //         var errMessage = xhr.responseText;
        //         debugger;
        //     },
        //     success: function(data, status, xhr) {
        //         initDatagrid();
        //     }
        // });
    });

    $('#purchase-order-form>button[role="submit"]').click(function(event) {
        if(!$('#received-order-checkbox').is(':checked')) {
            $('#received-date-input').remove();
        }
        sendOrderItems();
    });
;
    $('#received-order-checkbox').change(function() {
        if(this.checked) {
            $('#received-date-input').show();
        } else {
            $('#received-date-input').hide();
        }
    });

    $.each($('.datagrid .datagrid tr'), function(i, el) {
        $(el).click(function(event) {
            $('#remove-po-item').removeAttr('disabled');
            _selected_tr = $(this);
        })
    });

    $('#remove-po-item').click((event) => {
        removeItem();
    });
});

function sendOrderItems() {
    var $input = $('form>input[name="Items"]');
    var data = $('#vendor-items-table').find('tr[item-id]').map(function(i, el) {
        return {
            'id': $(el).attr('item-id'),
            'inventory_item_id': $(el).attr('inv-item-id'),
            'price': $(el).find('td[name="price"]').attr('value'),
            'count': $(el).find('td[name="count"]').text()
        };
    }).get();
    data = JSON.stringify(data);
    $input.val(data);
}

function removeItem() {
    if(_selected_tr.siblings().length == 0) {
        var $categoryTable = _selected_tr.parent().closest('tr');

        // remove category header
        $categoryTable.prev().remove();
        $categoryTable.remove();
    } else {
        _selected_tr.remove();
    }
    // remove item from display
        // remove entire category if necessary

    // add item to the remaining items holder

    // enable + button
}
