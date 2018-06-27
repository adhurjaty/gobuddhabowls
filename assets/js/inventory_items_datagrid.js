import { groupByCategory, formatMoney } from './helpers';
import { DataGrid } from './datagrid';

export function initDatagrid() {
    var grid = $('.datagrid .datagrid').get();

    return new DataGrid(grid, orderCountChanged);
}

function orderCountChanged(editItem) {
    // change price extension for row
    var $tr = editItem.$td.parent();
    var price = parseFloat($tr.find('td[name="price"]').attr('value')),
      count = parseFloat($tr.find('td[name="count"]').text());
    var extension = price * count;
    $tr.find('td[name="extension"]').text(formatMoney(extension));
}

function createDatagrid(items) {
    var head = `
    <div class="row justify-content-center">
        <div class="col-6 text-center">
            <h4>Order Items</h4>
        </div>
    </div>
    <div class="row">
        <table class="datagrid">
            <thead>
                <th width="40%">Name</th>
                <th width="22%">Cost</th>
                <th width="15%">Count</th>
                <th width="22%">Total Cost</th>
            </thead>
            <tbody>`;
    var foot = `
            </tbody>
        </table>
    </div>`

    var categorized = groupByCategory(items);

    var categoryRows = categorized.map((categoryGroup) => {
        return `
        <tr class="category-header" style="background-color: ${categoryGroup.background}">
            <td colspan="100">
                ${categoryGroup.name}
            </td>
        </tr>
        <tr>
            <td colspan="100" style="padding: 0;">
                <table class="datagrid">
                    <thead style="display: none;">
                        <th>Name</th>
                        <th>Cost</th>
                        <th>Count</th>
                        <th>Total Cost</th>
                    </thead>
                    <tbody>
                        ${categoryGroup.value.map((item) => {
                            return createRow(item);
                        }).join('')}
                    </tbody>
                </table>
            </td>
        </tr>`
    }).join('');

    return head + categoryRows + foot;
}

export function addToDatagrid(item, datagrid) {
    // move this stuff to datagrid
    var tr = createRow(item);
    var $rows = $('.datagrid .datagrid tr')

    var idx = -1;
    $.each($rows, (i, x) => {
        if(parseInt(item.index) < parseInt($(x).attr('data-index'))) {
            idx = i;
            return false;
        }
    });

    // check whether this is the first item of its category
    // add new category
    // add row to datagrid

    // datagrid.initRow(tr);
    // $.each($(tr).find('td[editable="true"]'), function(i, el) {
    //     new EditItem(datagrid, $(el));
    // });

    if(idx == -1) {
        $rows.parent().append(tr)
    } else {
        $(tr).insertBefore($rows.get(idx));
    }
}

function createRow(item) {
    var price = parseFloat(item.price);
    var count = parseFloat(item.count);
    return `
    <tr item-id="${item.id}" inv-item-id="${item.inventory_item_id}" data-index="${item.index}">
        <td name="name" width="40%">${item.name}</td>
        <td name="price" width="22%" editable="true" data-type="money" value="${price}">${formatMoney(price)}</td>
        <td name="count" width="15%" editable="true" data-type="number">${count}</td>
        <td name="extension" width="22%">${formatMoney(price * count)}</td>
    </tr>`
}

$(() => {
    var $container = $('#vendor-items-table');
    var items = JSON.parse($container.attr('data'));

    $container.html(createDatagrid(items));
});

