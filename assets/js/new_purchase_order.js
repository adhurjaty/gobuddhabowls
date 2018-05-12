// import { datepicker } from 'bootstrap-datepicker';
import { EditItem, DataGrid } from './datagrid.js';
import { formatMoney } from './helpers';

$(() => {
    $('#new-order-date').datepicker({
        autoclose: 'true',
        format: 'mm/dd/yyyy'
    });

    $('#new-order-vendor').change(function(d) {
        // remove none option
        $('#new-order-vendor option[value=""]').remove();
        var id = $(this).val();
        $.ajax({
            url: '/purchase_orders/order_vendor_changed/' + id,
            method: 'GET',
            error: function(xhr, status, err) {
                var errMessage = xhr.responseText;
                debugger;
            }
            // success: function(data, status, xhr) {
            //     editItem.onUpdateSuccess();
            // }
        });
    });

    $('#vendor-items-table').on('DOMNodeInserted', function(event) {
        if(event.target.parentNode.id == 'vendor-items-table') {
            $.each($('.datagrid'), function(i, grid) {
                var dg = new DataGrid(grid, orderCountChanged);
        
                $.each($(this).find('td[editable="true"]'), function(j, el) {
                    var ei = new EditItem(dg, $(el));
                });
            });
        }
    });
});

export function sendOrderItems() {
    var $input = $('form>input[name="Items"]');
    var data = $('#vendor-items-table').find('tr[item-id]').map(function(i, el) {
        return {
            'inventory_item_id': $(el).attr('item-id'),
            'price': $(el).find('td[name="price"]').attr('value'),
            'count': $(el).find('td[name="count"]').text()
        };
    }).get();
    data = JSON.stringify(data);
    $input.val(data);
    debugger;
}

export function orderCountChanged(editItem) {
    // change price extension for row
    var $tr = editItem.$td.parent();
    var price = parseFloat($tr.find('td[name="price"]').attr('value')),
      count = parseFloat($tr.find('td[name="count"]').text());
    var extension = price * count;
    $tr.find('td[name="extension"]').text(formatMoney(extension));

    // generate/update category breakdown
}