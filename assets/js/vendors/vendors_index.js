import { formatMoney, replaceUrlId, parseModelJSON } from '../helpers/_helpers';
import { DataGrid } from '../datagrid/_datagrid';
import { sendUpdate, sendAjax } from '../helpers/index_helpers';
import { CollapsibleDatagrid } from '../datagrid/_collapsible_datagrid';

$(() => {
    var $container = $('#datagrid-holder');
    var editVendorPath = $container.attr('data-url');

    var columnObjects = [
        {
            name: 'id',
            hidden: true,
            get_column: (vendor) => {
                return vendor.id;
            }
        },
        {
            name: 'name',
            header: 'Name',
            editable: true,
            get_column: (vendor) => {
                return vendor.name;
            }
        },
        {
            name: 'email',
            header: 'Email',
            editable: true,
            get_column: (vendor) => {
                return vendor.email;
            },
            set_column: (item, value) => {
                item.email = value;
            }
        },
        {
            name: 'phone_number',
            header: 'Phone Number',
            editable: true,
            get_column: (vendor) => {
                return vendor.phone_number;
            },
            set_column: (item, value) => {
                item.phone_number = value;
            }
        },
        {
            name: 'contact',
            header: 'Contact',
            editable: true,
            get_column: (vendor) => {
                return vendor.contact;
            },
            set_column: (item, value) => {
                item.contact = value;
            }
        },
        {
            name: 'shipping_cost',
            header: 'Shipping Cost',
            data_type: 'money',
            editable: true,
            get_column: (vendor) => {
                return formatMoney(vendor.shipping_cost);
            },
            set_column: (item, value) => {
                item.shipping_cost = parseFloat(value);
            }
        },
        {
            name: 'dropdown',
            get_column: ((editVendorPath) => {
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
    var data = parseModelJSON(dataStr) || [];

    var datagrid = new CollapsibleDatagrid(data, columnObjects, getHiddenRow, sendDatagridUpdate);
    $container.html(datagrid.getTable());
});

// TODO: repeated code from purchase_orders_datagrid.js - refactor
function sendDatagridUpdate(updateObj) {
    var copyObj = JSON.parse(JSON.stringify(updateObj));
    delete copyObj["Items"];
    var form = $('#update-vendor-form');
    sendUpdate($('#update-vendor-form'), copyObj, ((form) => {
        return () => sendAjax(form)
    })(form));
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