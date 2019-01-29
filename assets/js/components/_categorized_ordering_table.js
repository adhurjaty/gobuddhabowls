import { groupByCategory } from "../helpers/_helpers";
import { OrderingTable } from "./_ordering_table";

export class CategorizedOrderingTable {
    constructor(items) {
        this.categorizedItems = groupByCategory(items);
        this.initOrderTables();
        this.initTable();
    }

    initOrderTables() {
        this.orderTables = this.categorizedItems.map((catItem) => {
            return {
                name: catItem.name,
                background: catItem.background,
                table: new OrderingTable(catItem.value)
            }
        });
    }

    initTable() {
        this.table = $('<ul></ul>');
        this.orderTables.forEach((item) => {
            var header = $(`<li class="text-center list-group-item category-header"
                style="background: ${item.background};">
                ${item.name}
            </li>`);
            var catTable = $('<li name="reorder-li" class="list-group-item"></li>');
            this.table.append(header);
            this.table.append(catTable);
            item.table.attach(catTable);
        });
    }

    attach(container) {
        container.empty();
        container.append(this.table);
    }
}