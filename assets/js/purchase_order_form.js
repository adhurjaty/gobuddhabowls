import { datepicker } from './datepicker';
import { addToDatagrid, initDatagrid } from './inventory_items_datagrid'

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

    $('#received-order-checkbox').change(function() {
        if(this.checked) {
            $('#received-date-input').show();
        } else {
            $('#received-date-input').hide();
        }
    });

    // $('#remove-po-item').click((event) => {
    //     removeItem();
    // });
    
});

function sendOrderItems() {
    var $input = $('form>input[name="Items"]');
    var data = $('#vendor-items-table').find('tr[item-id]').map(function(i, el) {
        return createItemFromRow($(el));
    }).get();
    data = JSON.stringify(data);
    $input.val(data);
}

// function removeItem() {
//     var item = createItemFromRow(_selected_tr);
//     // remove item from display
//     if(_selected_tr.siblings().length == 0) {
//         // remove entire category if necessary
//         var $categoryTable = _selected_tr.parent().closest('tr');

//         // remove category header
//         $categoryTable.prev().remove();
//         $categoryTable.remove();
//     } else {
//         _selected_tr.remove();
//     }

//     // add item to the remaining items holder
//     addToRemaining(item)
    
//     $('#add-po-item').removeAttr('disabled');
//     $('#remove-po-item').attr('disabled', true);
// }

// function addItem(id) {
//     var itemToAdd = removeFromRemaining(id);
//     addToDatagrid(itemToAdd, _datagrid);

//     if(JSON.parse($('#add-order-modal').attr('data')).length == 0) {
//         $('#add-po-item').attr('disabled', true)
//     }
// }

// function addToRemaining(item) {
//     var $container = $('#add-order-modal');
//     var remainingItems = JSON.parse($container.attr('data'));

//     remainingItems.push(item);
//     var newRemainingString = JSON.stringify(remainingItems.sort((a, b) => {
//         return a.index - b.index;
//     }));
//     $container.attr('data', newRemainingString);
//     populateRemaining();
// }

// function removeFromRemaining(id) {
//     var $container = $('#add-order-modal');
//     var remainingItems = JSON.parse($container.attr('data'));
//     var idx = remainingItems.findIndex((x) => x.id == id);

//     var itemToRemove = remainingItems[idx]

//     remainingItems.splice(idx, 1);
//     $container.attr('data', JSON.stringify(remainingItems));
//     populateRemaining();

//     return itemToRemove;
// }

// function populateRemaining() {
//     var $container = $('#add-order-modal');
//     var remainingItems = JSON.parse($container.attr('data'));
//     var $select = $container.find('select');

//     remainingItems.forEach((item) => {
//         $('<option/>').val(item.id).html(item.name).appendTo($select);
//     });
// }

// function createItemFromRow($tr) {
//     return {
//         'id': $tr.attr('item-id'),
//         'inventory_item_id': $tr.attr('inv-item-id'),
//         'name': $tr.find('td[name="name"]').html(),
//         'price': $tr.find('td[name="price"]').attr('value'),
//         'count': $tr.find('td[name="count"]').text(),
//         'index': $tr.attr('data-index')
//     };
// }
