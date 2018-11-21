import { replaceUrlId, parseModelJSON } from '../helpers/_helpers';
import { sendUpdate } from '../helpers/index_helpers';
import { CollapsibleDatagrid } from '../datagrid/_collapsible_datagrid';

$(() => {
    var $container = $('#datagrid-holder');
    var editInventoryPath = $container.attr('data-url');

    var columnObjects = [
        {
            name: 'id',
            hidden: true,
            column_func: (inventory) => {
                return inventory.id;
            }
        },
        {
            name: 'date',
            header: 'Date',
            editable: true,
            column_func: (inventory) => {
                return inventory.data;
            }
        },
        {
            name: 'dropdown',
            column_func: ((editInventoryPath) => {
                return (inventory) => {
                    return `<div class="dropdown show">
                            <button type="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                                ...
                            </button>
                            <div class="dropdown-menu">
                                <a href="${replaceUrlId(editInventoryPath, inventory.id)}" class="dropdown-item">Edit</a>
                            </div>
                        </div>`
                }
            })(editInventoryPath)
        }
    ];

    var dataStr = $container.attr('data');
    var data = parseModelJSON(dataStr) || [];

    var datagrid = new CollapsibleDatagrid(data, columnObjects, getHiddenRow, sendDatagridUpdate);
    $container.html(datagrid.getTable());
});

// TODO: repeated code from purchase_orders_datagrid.js - refactor
function sendDatagridUpdate(updateObj) {
    sendUpdate($('#update-vendor-form'), updateObj);
}

function getHiddenRow(inventory) {
    debugger;
    return `<tr class="items-list">
            <td colspan="100">
                <table>
                    ${inventory.Items.map((item) => {
                        return `<tr>
                            <td>${item.name}</td>
                            <td>${item.count}</td>
                        </tr>`
                    }).join('\n')}
                </table>
            </td>
        </tr>`
}