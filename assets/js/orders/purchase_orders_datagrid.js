import { CollapsibleDatagrid } from '../datagrid/_collapsible_datagrid';
import { getPurchaseOrderCost, formatMoney, formatSlashDate, replaceUrlId } from '../_helpers';
import { horizontalPercentageChart } from '../_horizontal_percentage_chart';

$(() => {
    var $container = $('#datagrid-holder');
    // var purchaseOrders = JSON.parse($container.attr('data'));
    var editOrderPath = $container.attr('data-url');
    var orderSheetPath = $container.attr('order-sheet-url');
    var receivingListPath = $container.attr('receiving-list-url');

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

    var $openDatagridContainer = $('#open-order-datagrid')
    var openOrdersData = $openDatagridContainer.attr('data')
    if(openOrdersData) {
        var openOrders = JSON.parse(openOrdersData);
        var openColumnObjects = baseColumnObjects.concat([
            {
                name: 'dropdown',
                column_func: ((editOrderPath, orderSheetPath, receivingListPath) => {
                    return (purchaseOrder) => {
                        return `<div class="dropdown show">
                            <button type="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                                ...
                            </button>
                            <div class="dropdown-menu">
                                <span name="receive" class="dropdown-item" onclick="receiveItem(${purchaseOrder.id})">Received</span>
                                <a href="${replaceUrlId(editOrderPath, purchaseOrder.id)}" class="dropdown-item">Edit</a>
                                <a name="delete" class="dropdown-item text-danger" data-method="DELETE">Delete</a>
                                <a href="${replaceUrlId(orderSheetPath, purchaseOrder.id)}" class="dropdown-item">Order Sheet<a>
                                <a href="${replaceUrlId(receivingListPath, purchaseOrder.id)}" class="dropdown-item">Receiving List<a>
                            </div>
                        </div>`
                    }
                })(editOrderPath, orderSheetPath, receivingListPath)
            }
        ]);
        var openDatagrid = new CollapsibleDatagrid(openOrders, openColumnObjects, getHiddenRow, sendUpdate)
        $openDatagridContainer.html(openDatagrid.getTable());

        initDropdownActions(openDatagrid);
    }

    var $recDatagridContainer = $('#rec-order-datagrid');
    var recOrdersData = $recDatagridContainer.attr('data');
    if(recOrdersData) {
        var recOrders = JSON.parse(recOrdersData);
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
                            <span name="reopen" class="dropdown-item" >Re-open</span>
                            <a href="${replaceUrlId(editOrderPath, purchaseOrder.id)}" class="dropdown-item">Edit</a>
                            <a name="delete" class="dropdown-item text-danger" data-method="DELETE">Delete</a>
                            <a href="${replaceUrlId(orderSheetPath, purchaseOrder.id)}" class="dropdown-item">Order Sheet<a>
                            <a href="${replaceUrlId(receivingListPath, purchaseOrder.id)}" class="dropdown-item">Receiving List<a>
                        </div>
                    </div>`
                }
            })(editOrderPath)
        });

        var recDatagrid = new CollapsibleDatagrid(recOrders, recColumnObjects, getHiddenRow, sendUpdate);
        $recDatagridContainer.html(recDatagrid.getTable());

        initDropdownActions(recDatagrid);
    }
});

function sendUpdate(updateObj) {
    var id = updateObj.id;
    
    submitUpdateForm(id, convertUpdateObj(updateObj));
}

function convertUpdateObj(updateObj) {
    var outObj = {}
    if(updateObj['order_date']) {
        outObj['OrderDate'] = updateObj['order_date']
    }
    if(updateObj['received_date']) {
        outObj['ReceivedDate'] = updateObj['received_date']
    }

    return outObj;
}

function getHiddenRow(purchaseOrder) {
    return `<tr class="items-list">
                <td colspan="100">
                    ${horizontalPercentageChart(purchaseOrder.Vendor.name + ' Order', purchaseOrder.Items, getPurchaseOrderCost(purchaseOrder))}            
                </td> 
            </tr>`;
}

function initDropdownActions(datagrid) {
    var updatePath = $('#datagrid-holder').attr('update-url');
    datagrid.rows.forEach((row) => {
        var details = row.getInfo();
        var id = details.id;
        var $row = row.getRow();

        $row.find('td[name="dropdown"] span[name="receive"]').click(() => {
            var date = new Date();
            date.setHours(0, 0, 0, 0);
            submitUpdateForm(id, { ReceivedDate: date.toISOString() });
        });
        $row.find('td[name="dropdown"] span[name="reopen"]').click(() => {
            submitUpdateForm(id, { ReceivedDate: '' });
        });
        $row.find('td[name="dropdown"] a[name="delete"]').attr('href', replaceUrlId(updatePath, id));
    });
}

function submitUpdateForm(id, data) {
    var $form = $('#update-po-form');
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