import { formatMoney, replaceUrlId } from '../_helpers';

$(() => {
    var $container = $('#datagrid-holder');
    var editOrderPath = $container.attr('data-url');

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
            name: 'shipping_cost',
            header: 'Shipping Cost',
            editable: true,
            column_func: (vendor) => {
                return formatMoney(vendor.shippingCost);
            }
        }
    ];

    var dataStr = $container.attr('data');
    var data = dataStr != '' ? JSON.parse(dataStr) : [];

});