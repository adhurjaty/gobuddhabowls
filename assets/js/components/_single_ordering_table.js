import { OrderingTable } from "./_ordering_table";

export class SingleOrderingTable extends OrderingTable {
    constructor(items, movingItem, onChange) {
        super(items, onChange);
        this.itemElement = this.insertItem(movingItem);
    }

    initTable() {
        super.initTable();
        this.ul.find('.drag-handle').remove();
    }

    insertItem(item) {
        var lis = this.ul.find('li');
        var idx = lis.toArray().findIndex(el => {
            return item.index <= parseInt($(el).attr('index'));
        });

        var listItem = this.getListElement(item);
        if(idx == -1) {
            this.ul.append(listItem);
        } else {
            lis.eq(idx).before(listItem);
        }

        return listItem
    }
    
    updateItemName(name) {
        this.itemElement.find('span').first().html(name);
    }
}