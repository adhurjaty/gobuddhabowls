import "bootstrap-datepicker";
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
            },
            success: function(data, status, xhr) {
                initDatagrid();
            }
        });
    });

    if($('#vendor-items-table').children().length > 0) {
        initDatagrid();
    }

    $('#purchase-order-form>button[role="submit"]').click(function(event) {
        if(!$('#received-order-checkbox').is(':checked')) {
            $('#received-date-input').remove();
        }
        sendOrderItems();
    });

    // $('#received-date-input').hide();
    $('#received-order-checkbox').change(function() {
        if(this.checked) {
            $('#received-date-input').show();
        } else {
            $('#received-date-input').hide();
        }
    });
    
    $('#modal-form button[role="submit"]').click(function(event) {
        $('#modal-form').submit();
    });

    // $('#remove-po-item').click(function(event) {
    //     var selected = $('.datagrid .datagrid .tr:active');
    //     selected.remove();
    // });
});

function initDatagrid() {
    var grid = $('.datagrid .datagrid').get();
    // $.each($('.datagrid'), function(i, grid) {
    var dg = new DataGrid(grid, orderCountChanged);

    $.each($('.datagrid td[editable="true"]'), function(j, el) {
        var ei = new EditItem(dg, $(el));
    });

    $.each($('.datagrid .datagrid tr'), function(i, el) {
        $(el).click(function(event) {
            $('#remove-po-item').removeAttr('disabled');
        })
    });

}

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

function orderCountChanged(editItem) {
    // change price extension for row
    var $tr = editItem.$td.parent();
    var price = parseFloat($tr.find('td[name="price"]').attr('value')),
      count = parseFloat($tr.find('td[name="count"]').text());
    var extension = price * count;
    $tr.find('td[name="extension"]').text(formatMoney(extension));

    // generate/update category breakdown
    // use backend to generate percentage chart
    var on_change_url = '/purchase_orders/count_changed'
    var itemsJSON = $('#vendor-items-table').find('tr[item-id]').map(function(i, el) {
        return {
            'inventory_item_id': $(el).attr('inv-item-id'),
            'price': $(el).find('td[name="price"]').attr('value'),
            'count': $(el).find('td[name="count"]').text()
        };
    }).get();
    var data = {};
    data['Items'] = JSON.stringify(itemsJSON);

    $.ajax({
        url: on_change_url,
        data: data,
        method: 'POST',
        error: function(xhr, status, err) {
            var errMessage = xhr.responseText;
            debugger;
        },
        success: function(data, status, xhr) {
            
        }
    });
}