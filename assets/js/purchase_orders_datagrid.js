import { DataGrid, EditItem } from './datagrid';
import { getPurchaseOrderCost, formatMoney, formatSlashDate, replaceUrlId } from './helpers';
import { horizontalPercentageChart } from './horizontal_percentage_chart';

$(() => {
    var $container = $('#datagrid-holder');
    var purchaseOrders = JSON.parse($container.attr('data'));
    var editOrderPath = $container.attr('data-url');
    var openOrders = purchaseOrders.filter(x => !x.received_date);
    var recOrders = purchaseOrders.filter(x => x.received_date);
    var tableArea = "";

    if(openOrders.length > 0) {
        tableArea = getDataGrid("Open Orders", openOrders, editOrderPath);
    }
    if(recOrders.length > 0) {
        tableArea += getDataGrid("Received Orders", recOrders, editOrderPath);
    }

    $('#datagrid-holder').html(tableArea);
    
    $.each($('.datagrid'), function(i, grid) {
        new DataGrid(grid);
    });
});

function getDataGrid(title, purchaseOrders, editOrderPath) {
    var head = `
    <div class="row justify-content-center">
        <div class="col-6 text-center">
            <h4>${title}</h4>
        </div>
    </div>
    <div class="row">
    <table class="datagrid" update-url="<%= purchaseOrdersPath() %>">
        <thead>
            <th></th>
            <th>Vendor</th>
            <th>Order Date</th>
            <th>Cost</th>
            <th>&nbsp;</th>
        </thead>
    <tbody>`;

    var foot = `
    </tbody>
      </table>
    </div>`;
    
    var rows = purchaseOrders.map(po => {
        var total = getPurchaseOrderCost(po)
        return `
        <tr item-id="${po.id}">
            <td class="expander"><span class="fa fa-caret-right"></span></td>
            <td>${po.Vendor.name}</td>
            <td editable="true" data-type="date" field="OrderDate">${formatSlashDate(po.order_date)}</td>
            <td>${formatMoney(total)}</td>
            <td>
                <div class="dropdown show">
                    <button type="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                        ...
                    </button>
                    <div class="dropdown-menu">
                        <span class="dropdown-item" onclick="receiveItem(${po.id})">Received</span>
                        <a href="${replaceUrlId(editOrderPath, po.id)}" class="dropdown-item">Edit</a>
                        <span class="dropdown-item text-danger" onclick="deleteItem(${po.id})">Delete</span>
                    </div>
                </div>
            </td>
        </tr>
        <tr class="items-list" style="display: none;">
            <td colspan="100">
                ${horizontalPercentageChart(title, po.Items, total)}            
            </td> 
        </tr>`;
    }).join('');
    
    return head + rows + foot;
}