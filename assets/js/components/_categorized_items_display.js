import { formatMoney, unFormatMoney, sortItems } from '../helpers/_helpers';
import { CategorizedDatagrid } from '../datagrid/_categorized_datagrid';
import { horizontalPercentageChart } from '../_horizontal_percentage_chart';
import { Modal } from './_modal';

var columnInfo = [
    {
        name: 'id',
        hidden: true,
        column_func: (item) => {
            return item.id;
        }
    },
    {
        name: 'inventory_item_id',
        hidden: true,
        column_func: (item) => {
            return item.inventory_item_id;
        }
    },
    {
        name: 'index',
        hidden: true,
        column_func: (item) => {
            return item.index;
        }
    },
    {
        header: 'Name',
        column_func: (item) => {
            return item.name;
        }
    },
    {
        name: 'price',
        header: 'Price',
        editable: true,
        data_type: 'money',
        column_func: (item) => {
            return formatMoney(parseFloat(item.price));
        }
    },
    {
        name: 'count',
        header: 'Count',
        editable: true,
        data_type: 'number',
        column_func: (item) => {
            return item.count;
        }
    },
    {
        name: 'total_cost',
        header: 'Total Cost',
        column_func: (item) => {
            return formatMoney(parseFloat(item.price) * parseFloat(item.count));
        }
    }
];

export class CategorizedItemsDisplay {
    constructor(allItems, $container) {
        this.$container = $container;
        this.allItems = allItems;
        this.items = null;
        this.$selectedTr = null;
        this.datagrid = null;
        this.breakdown = null;
        this.columnInfo = columnInfo;
        this.buttonGroup = null;
        this.modal = null;

        this.datagridUpdated = this.datagridUpdated.bind(this);
        this.addItem = this.addItem.bind(this);

        this.initDisplay();
    }

    initDisplay() {
        this.initItems();
        this.writeItemsToDataAttr();

        if(this.allItems) {
            this.initAddRemoveButtons();
            this.insertModal();
        }

        this.updateTables();
    }

    initItems() {
        this.items = JSON.parse(this.$container.attr('data'));
        this.items.forEach(item => {
            if(!item.count) {
                item.count = 0;
            }
        });
    }

    initAddRemoveButtons() {
        var self = this;
        this.buttonGroup = new ButtonGroup();
        var remainingItems = this.getRemainingItems();
        if(remainingItems != null && remainingItems.length > 0) {
            this.buttonGroup.enableAddButton();
        }

        this.buttonGroup.setRemoveListener(() => {
            // sometimes triggers multiple times from UI
            // this check ensures this function happens once
            if(!self.$selectedTr) {
                return;
            }
            var selectedItem = self.datagrid.getItem(self.$selectedTr);
            self.$selectedTr = null;
            self.removeItem(selectedItem);
        });

        this.$container.find('div[name="button-group"]').html(this.buttonGroup.$content);
    }

    updateTables() {
        this.initDatagrid();
        this.initBreakdown();
    }

    initDatagrid() {
        this.datagrid = new CategorizedDatagrid(this.items, 
            this.columnInfo, this.datagridUpdated);
        this.$container.find('div[name="datagrid"]')
                       .html(this.datagrid.getTable());
    
        this.initSelection();
    }
    
    initSelection() {
        var self = this;
        this.datagrid.rows.forEach((row) => {
            var $tr = row.getRow();
            $tr.click(() => {
                if(self.buttonGroup) {
                    self.buttonGroup.enableRemoveButton();
                }
                self.$selectedTr = $tr;
            });
        });
    }

    initBreakdown() {
        var bdContainer = this.$container.find('div[name="breakdown"]');
        var title = 'Order Breakdown';
        var total = this.items.reduce((total, item) => total + item.count * item.price, 0);
        if(total != 0) {
            bdContainer.html(horizontalPercentageChart(title, this.items, total));
        } else {
            bdContainer.html('');
        }
    }

    writeItemsToDataAttr() {
        this.$container.attr('data', JSON.stringify(this.items));
    }

    addItem(item) {
        item.count = 0;
        this.items.push(item);
        this.items = sortItems(this.items);
        this.initDatagrid();
        this.initBreakdown();

        var remaining = this.getRemainingItems();
        if(remaining.length == 0) {
            this.buttonGroup.disableAddButton();
        }
    }

    removeItem(item) {
        var idx = this.items.indexOf(item);
        this.items.splice(idx, 1);
        this.writeItemsToDataAttr();
    
        this.updateTables();
    
        this.addToModal(item);
    
        this.buttonGroup.enableAddButton();
        this.buttonGroup.disableRemoveButton();
    }

    addToModal(item) {
        this.modal.addItem(item);
    }

    insertModal() {
        var remainingItems = this.getRemainingItems();
        this.modal = new Modal(remainingItems, this.addItem);
        this.$container.find('div[name="modal"]').html(this.modal.$content);
    }
    
    getRemainingItems() {
        return this.allItems.filter(x => {
            return this.items.findIndex(y => y.inventory_item_id ==
                                            x.inventory_item_id) == -1;
        }, this);
    }

    datagridUpdated(updateObj) {
        var price = parseFloat(unFormatMoney(updateObj.price));
        var count = parseFloat(updateObj.count);
        var $tr = $('.datagrid').find(`tr td:contains(${updateObj.id})`).parent();
        $tr.find('td[name="total_cost"]').text(formatMoney(price * count));

        // update items and breakdown
        var idx = this.items.findIndex((x) => x.id == updateObj.id);
        this.items[idx].price = price;
        this.items[idx].count = count;
        this.writeItemsToDataAttr();

        this.initBreakdown();
    }
}

class ButtonGroup {
    constructor() {
        this.createButtonGroup();
    }

    createButtonGroup() {
        this.addButton = $(`<button class="btn btn-default" 
                type="button" data-toggle="modal"
                data-target="#add-item-modal" disabled>

                <span class="fa fa-plus">
            </button>`);
        this.removeButton = $(`<button class="btn btn-default"
                type="button" disabled>
                <span class="fa fa-minus">
            </button>`);
        this.$content = $(`<div class="input-group d-flex 
            justify-content-end"></div>`);
        this.addButton.appendTo(this.$content);
        this.removeButton.appendTo(this.$content);
    }

    setAddListenr(fn) {
        this.addButton.click(() => {
            fn();
        });
    }

    setRemoveListener(fn) {
        this.removeButton.click(() => {
            fn();
        });
    }

    enableAddButton() {
        this.enableButton(this.addButton);
    }

    disableAddButton() {
        this.disableButton(this.addButton);
    }

    enableRemoveButton() {
        this.enableButton(this.removeButton);
    }

    disableRemoveButton() {
        this.disableButton(this.removeButton);
    }

    enableButton(btn) {
        btn.removeAttr('disabled');
    }

    disableButton(btn) {
        btn.attr('disabled', 'disabled');
    }

    getGroup() {
        return this.group;
    }
}

// export function initOrderItemsArea(items) {
//     _$container = $('#vendor-items-table');
//     _items = items;
//     _items.forEach((item) => {
//         if(!item.count) {
//             item.count = 0;
//         }
//     });
//     writeItemsToDataAttr();

//     initDatagrid();
//     initBreakdown();
//     initAddRemoveButtons();
// }

// export function addToDatagrid(item) {
//     // need to re-initialize globals because this is called from outside
//     _$container = $('#vendor-items-table');
//     _items = JSON.parse(_$container.attr('data'));

//     if(item.count == undefined) {
//         item.count = 0;
//     }

//     // add item to datagrid data
//     _items.push(item);
//     _items = sortItems(_items);
//     writeItemsToDataAttr();

//     initDatagrid();
//     initBreakdown();

//     // remove item from modal remaining
//     removeFromRemaining(item);
// }


// function datagridUpdated(updateObj) {
//     var price = parseFloat(unFormatMoney(updateObj.price));
//     var count = parseFloat(updateObj.count);
//     var $tr = $('.datagrid').find(`tr td:contains(${updateObj.id})`).parent();
//     $tr.find('td[name="total_cost"]').text(formatMoney(price * count));

//     // update items and breakdown
//     var idx = _items.findIndex((x) => x.id == updateObj.id);
//     _items[idx].price = price;
//     _items[idx].count = count;
//     writeItemsToDataAttr();

//     initBreakdown();
// }
