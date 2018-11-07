import { formatMoney, replaceUrlId } from '../_helpers';
import { DataGrid } from '../datagrid/_datagrid';

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
            name: 'phoneNumber',
            header: 'Phone Number',
            editable: true,
            column_func: (vendor) => {
                return vendor.phoneNumber;
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
            name: 'shippingCost',
            header: 'Shipping Cost',
            type: 'money',
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

    var datagrid = new DataGrid(data, columnObjects, sendUpdate);
    $container.html(datagrid.getTable());
});

function sendUpdate(updateObj) {
    var id = updateObj.id;

    submitUpdateForm(id, updateObj);
}


function submitUpdateForm(id, data) {
    var $form = $('#update-vendor-form');
    // TODO: repeated code from purchase_orders_datagrid.js - refactor
    $form.attr('action', replaceUrlId($form.attr('action'), id));

    for(var key in data) {
        var $field = $('<input />');
        $field.attr('type', 'hidden');
        $field.attr('name', key);
        $field.val(data[key]);
        $field.appendTo($form);
    }

    $form.submit();
}