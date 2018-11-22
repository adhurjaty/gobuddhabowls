import { replaceUrlId, parseModelJSON, formatSlashDate } from '../helpers/_helpers';
import { sendUpdate } from '../helpers/index_helpers';
import { CollapsibleDatagrid } from '../datagrid/_collapsible_datagrid';


$(() => {
    var invList = $('#date-list');
    invList.find('li').first().addClass('active');

    setClickInventory();
    setSelectedInventory();

    // var $container = $('#datagrid-holder');
    // var editInventoryPath = $container.attr('data-url');

    // var columnObjects = [
    //     {
    //         name: 'id',
    //         hidden: true,
    //         column_func: (inventory) => {
    //             return inventory.id;
    //         }
    //     },
    //     {
    //         name: 'date',
    //         header: 'Date',
    //         editable: true,
    //         data_type: 'date',
    //         column_func: (inventory) => {
    //             return formatSlashDate(inventory.time);
    //         }
    //     },
    //     {
    //         name: 'dropdown',
    //         column_func: ((editInventoryPath) => {
    //             return (inventory) => {
    //                 return `<div class="dropdown show">
    //                         <button type="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
    //                             ...
    //                         </button>
    //                         <div class="dropdown-menu">
    //                             <a href="${replaceUrlId(editInventoryPath, inventory.id)}" class="dropdown-item">Edit</a>
    //                         </div>
    //                     </div>`
    //             }
    //         })(editInventoryPath)
    //     }
    // ];

    // var dataStr = $container.attr('data');
    // var data = parseModelJSON(dataStr) || [];

    // var datagrid = new CollapsibleDatagrid(data, columnObjects, getHiddenRow, sendDatagridUpdate);
    // $container.html(datagrid.getTable());
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
    var selectedInventory = JSON.parse(invList.find('li.active').attr('data'));
    $form.find('input[name="Date"]').val(formatSlashDate(selectedInventory.time));

}

// TODO: repeated code from purchase_orders_datagrid.js - refactor
// function sendDatagridUpdate(updateObj) {
//     sendUpdate($('#update-inventory-form'), updateObj);
// }

