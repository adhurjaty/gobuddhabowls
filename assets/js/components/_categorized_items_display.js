import { formatMoney, unFormatMoney, sortItems, parseModelJSON } from '../helpers/_helpers';
import { CategorizedDatagrid } from '../datagrid/_categorized_datagrid';
import { horizontalPercentageChart } from '../_horizontal_percentage_chart';
import { Modal } from './_modal';
import { ButtonGroup } from './_button_group';

export class CategorizedItemsDisplay {
    constructor(container, columnInfo, allItems, options) {
        this.$container = container;
        this.allItems = allItems;
        this.options = options || {};
        this.items = parseModelJSON(this.$container.attr('data')) || [];
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
        this.initContainer();
        this.updateTables();
    }

    initContainer() {
        if(this.options.breakdown) {
            this.$container.html(
                `<div class="col-md-6">
                    <div name="datagrid"></div>
                </div>
                <div class="col-md-6" name="breakdown"></div>`);
        } else {
            this.$container.html('<div class="col-md-12"><div name="datagrid"></div></div>');
        }
        if(this.allItems) {
            this.initAddRemoveButtons();
            this.insertModal();
        }
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

        var colDiv = this.$container.find('> div').first();
        var buttonContainer = $('<div></div>');
        buttonContainer.append(this.buttonGroup.$content);
        colDiv.prepend(buttonContainer);
    }

    updateTables() {
        this.initDatagrid();
        if(this.options.breakdown) {
            this.initBreakdown();
        }
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
        this.setDefaults(item);
        this.items.push(item);
        this.items = sortItems(this.items);
        this.updateTables();

        var remaining = this.getRemainingItems();
        if(remaining.length == 0) {
            this.buttonGroup.disableAddButton();
        }
    }

    setDefaults(item) {
        this.columnInfo.forEach((info) => {
            if(Object.keys(info).indexOf('default') > -1) {
                item[info.name] = info.default;
            }
        });
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
        var modalContainer = $(`<div name="modal"></div>`);
        modalContainer.html(this.modal.$content);
        this.$container.append(modalContainer);
    }
    
    getRemainingItems() {
        return this.allItems.filter(x => {
            return this.items.findIndex(y => y.inventory_item_id ==
                                            x.inventory_item_id) == -1;
        }, this);
    }

    datagridUpdated(updateObj) {
        if(this.options.datagridUpdated) {
            this.options.datagridUpdated(updateObj);
        }
        this.writeItemsToDataAttr();

        if(this.options.breakdown)
            this.initBreakdown();
    }
}
