import { formatMoney, replaceUrlId } from '../helpers/_helpers';
import { DataGrid } from '../datagrid/_datagrid';
import { sendUpdate } from '../helpers/index_helpers';
import { CollapsibleDatagrid } from '../datagrid/_collapsible_datagrid';

$(() => {
    var $container = $('#datagrid-holder');
    var editVendorPath = $container.attr('data-url');

    var columnObjects = [
        {
            name: 'id',
            hidden: true,
            column_func: (vendor) => {
                return vendor.id;
            }
        },
        {
            name: 'name',
            header: 'Name',
            editable: true,
            column_func: (vendor) => {
                return vendor.name;
            }
        },
        {
            name: 'email',
            header: 'Email',
            editable: true,
            column_func: (vendor) => {
                return vendor.email;
            }
        },
        {
            name: 'phone_number',
            header: 'Phone Number',
            editable: true,
            column_func: (vendor) => {
                return vendor.phone_number;
            }
        },
        {
            name: 'contact',
            header: 'Contact',
            editable: true,
            column_func: (vendor) => {
                return vendor.contact;
            }
        },
        {
            name: 'shipping_cost',
            header: 'Shipping Cost',
            data_type: 'money',
            editable: true,
            column_func: (vendor) => {
                return formatMoney(vendor.shipping_cost);
            }
        },
        {
            name: 'dropdown',
            column_func: ((editVendorPath) => {
                return (vendor) => {
                    return `<div class="dropdown show">
                            <button type="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                                ...
                            </button>
                            <div class="dropdown-menu">
                                <a href="${replaceUrlId(editVendorPath, vendor.id)}" class="dropdown-item">Edit</a>
                            </div>
                        </div>`
                }
            })(editVendorPath)
        }
    ];

    var dataStr = $container.attr('data');
    var data = dataStr != '' ? JSON.parse(dataStr) : [];

    var datagrid = new CollapsibleDatagrid(data, columnObjects, getHiddenRow, sendDatagridUpdate);
    $container.html(datagrid.getTable());
});

// TODO: repeated code from purchase_orders_datagrid.js - refactor
function sendDatagridUpdate(updateObj) {
    sendUpdate($('#update-vendor-form'), updateObj);
}

function getHiddenRow(vendor) {
    return `<tr class="items-list">
            <td colspan="100">
                <table>
                    ${vendor.Items.map((item) => {
                        return `<tr>
                            <td>${item.name}</td>
                            <td>${formatMoney(item.price)}</td>
                        </tr>`
                    }).join('\n')}
                </table>
            </td>
        </tr>`
}