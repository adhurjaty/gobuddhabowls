import { formatMoney } from "../_helpers";
import html2canvas from 'html2canvas';
import { DataGrid } from "../datagrid/_datagrid";


$(() => {
    var columnInfo = [
        {
            header: 'Pack Size',
            column_func: (item) => {
                return item.conversion;
            }
        },
        {
            header: 'Purchased Unit',
            column_func: (item) => {
                return item.purchased_unit;
            }
        },
        {
            header: 'Item Name',
            column_func: (item) => {
                return item.name;
            }
        },
        {
            header: 'Order Amt',
            column_func: (item) => {
                return item.count;
            }
        },
        {
            header: 'Current Price',
            column_func: (item) => {
                return formatMoney(item.price);
            }
        },
        {
            header: 'Extension',
            column_func: (item) => {
                return formatMoney(item.count * item.price);
            }
        },
        {
            header: 'Received?',
            column_func: (item) => {
                return '';
            }
        }
    ];

    var $container = $('#receiving-table');
    var data = JSON.parse($container.attr('data'));
    var table = new DataGrid(data, columnInfo);
    $container.append(table.getTable());

    html2canvas(document.body).then((canvas) => {
        document.body.appendChild(canvas);
        $('#receiving-sheet').remove();
    });
});