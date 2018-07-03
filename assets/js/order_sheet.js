import { groupByCategory, formatMoney } from "./helpers";
import { PrintTable } from "./print_table";

function makeTable(items) {
    var categorizedItems = groupByCategory(items);

    var $table = $(`<table></table>`);

    categorizedItems.forEach((catItem) => {
        addSubTable(catItem, $table);
    });
    var $thead = $('<thead></thead>');
    $thead.appendTo($table);

     
}

function addSubTable(categoryItems, $table) {
    var category = categoryItems[0].Category;
    var header = ['Pack Size', 'Purchased Unit', category.name, 'Order Amt', 'Current Price', 'Extension'];

    var $tr = $('<tr></tr>');
    $tr.style(`background-color: ${category.background};`);
    header.forEach((heading) => {
        var $th = $(`<td>${heading}</td>`);
        $th.appendTo($tr);
    });
    $tr.appendTo($table);

    categoryItems.forEach((item) => {

    })
}


$(() => {
    var columnInfo = [
        {
            name: 'Pack Size',
            column_func: (item) => {
                return item.conversion;
            }
        },
        {
            name: 'Purchased Unit',
            column_func: (item) => {
                return item.purchased_unit;
            }
        },
        {
            name: 'category',
            column_func: (item) => {
                return item.name;
            }
        },
        {
            name: 'Order Amt',
            column_func: (item) => {
                return item.count;
            }
        },
        {
            name: 'Current Price',
            column_func: (item) => {
                return formatMoney(item.price);
            }
        },
        {
            name: 'Extension',
            column_func: (item) => {
                return formatMoney(item.count * item.price);
            }
        }
    ];

    var $container = $('#order-table');
    var data = JSON.parse($container.attr('data'));
    var table = new PrintTable(data, columnInfo);
    $container.append(table.getTable());
});