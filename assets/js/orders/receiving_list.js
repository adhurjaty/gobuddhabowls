import { formatMoney } from "../helpers/_helpers";
import html2canvas from 'html2canvas';
import { DataGrid } from "../datagrid/_datagrid";


$(() => {
    var columnInfo = [
        {
            header: 'Pack Size',
            get_column: (item) => {
                return item.conversion;
            }
        },
        {
            header: 'Purchased Unit',
            get_column: (item) => {
                return item.purchased_unit;
            }
        },
        {
            header: 'Item Name',
            get_column: (item) => {
                return item.name;
            }
        },
        {
            header: 'Order Amt',
            get_column: (item) => {
                return item.count;
            }
        },
        {
            header: 'Current Price',
            get_column: (item) => {
                return formatMoney(item.price);
            }
        },
        {
            header: 'Extension',
            get_column: (item) => {
                return formatMoney(item.count * item.price);
            }
        },
        {
            header: 'Received?',
            get_column: (item) => {
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