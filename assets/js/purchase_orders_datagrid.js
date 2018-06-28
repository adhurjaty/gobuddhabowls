import { CollapsibleDatagrid } from './collapsible_datagrid';
import { getPurchaseOrderCost, formatMoney, formatSlashDate, replaceUrlId } from './helpers';
import { horizontalPercentageChart } from './horizontal_percentage_chart';

$(() => {
    var $container = $('#datagrid-holder');
    // var purchaseOrders = JSON.parse($container.attr('data'));
    var editOrderPath = $container.attr('data-url');
    var $openDatagridContainer = $('#open-order-datagrid')
    var openOrders = JSON.parse($openDatagridContainer.attr('data'));
    var $recDatagridContainer = $('#rec-order-datagrid');
    var recOrders = JSON.parse($recDatagridContainer.attr('data'));

    var baseColumnObjects = [
        {
            name: 'id',
            hidden: true,
            column_func: (purchaseOrder) => {
                return purchaseOrder.id;
            }
        },
        {
            name: 'vendor',
            header: 'Vendor',
            column_func: (purchaseOrder) => {
                return purchaseOrder.Vendor.name;
            }
        },
        {
            name: 'order_date',
            header: 'Order Date',
            editable: true,
            data_type: 'date',
            column_func: (purchaseOrder) => {
                return formatSlashDate(purchaseOrder.order_date);
            }
        },
        {
            name: 'cost',
            header: 'Cost',
            column_func: (purchaseOrder) => {
                return formatMoney(getPurchaseOrderCost(purchaseOrder));
            }
        }
    ];

    var openColumnObjects = baseColumnObjects.concat([
        {
            name: 'dropdown',
            column_func: ((editOrderPath) => {
                return (purchaseOrder) => {
                    return `<div class="dropdown show">
                        <button type="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                            ...
                        </button>
                        <div class="dropdown-menu">
                            <span class="dropdown-item" onclick="receiveItem(${purchaseOrder.id})">Received</span>
                            <a href="${replaceUrlId(editOrderPath, purchaseOrder.id)}" class="dropdown-item">Edit</a>
                            <span class="dropdown-item text-danger" onclick="deleteItem(${purchaseOrder.id})">Delete</span>
                        </div>
                    </div>`
                }
            })(editOrderPath)
        }
    ]);

    var recColumnObjects = baseColumnObjects.slice();
    recColumnObjects.splice(3, 0, {
        name: 'received_date',
        header: 'Received Date',
        editable: true,
        data_type: 'date',
        column_func: (purchaseOrder) => {
            return formatSlashDate(purchaseOrder.received_date);
        }
    });
    recColumnObjects.push({
        name: 'dropdown',
        column_func: ((editOrderPath) => {
            return (purchaseOrder) => {
                return `<div class="dropdown show">
                    <button type="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                        ...
                    </button>
                    <div class="dropdown-menu">
                        <span class="dropdown-item" onclick="reopen(${purchaseOrder.id})">Re-open</span>
                        <a href="${replaceUrlId(editOrderPath, purchaseOrder.id)}" class="dropdown-item">Edit</a>
                        <span class="dropdown-item text-danger" onclick="deleteItem(${purchaseOrder.id})">Delete</span>
                    </div>
                </div>`
            }
        })(editOrderPath)
    });

    var openDatagrid = new CollapsibleDatagrid(openOrders, openColumnObjects, getHiddenRow, sendUpdate)
    var recDatagrid = new CollapsibleDatagrid(recOrders, recColumnObjects, getHiddenRow, sendUpdate);

    $openDatagridContainer.html(openDatagrid.getTable());
    $recDatagridContainer.html(recDatagrid.getTable());
});

function sendUpdate(updateObj) {
    var id = updateObj.id;
    var updateUrl = replaceUrlId($('#datagrid-holder').attr('update-url'), id);
    var data = {};
    data[updateObj.name] = updateObj.value;
    $.ajax({
        url: updateUrl,
        data: data,
        method: 'PUT',
        error: function(xhr, status, err) {
            var errMessage = xhr.responseText;
            debugger;
            updateObj.onError(errMessage);
        },
        success: function(data, status, xhr) {
            updateObj.onSuccess();
        }
    });
}

function getHiddenRow(purchaseOrder) {
    return `<tr class="items-list">
                <td colspan="100">
                    ${horizontalPercentageChart(purchaseOrder.Vendor.name + ' Order', purchaseOrder.Items, getPurchaseOrderCost(purchaseOrder))}            
                </td> 
            </tr>`;
}