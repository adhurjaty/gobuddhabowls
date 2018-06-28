import { DataGrid } from "./datagrid";
import { groupByCategory } from "./helpers";

export class CategorizedDatagrid extends DataGrid {
    constructor(data, columnInfo, updateFnc) {
        super(data, columnInfo, updateFnc);
        this.categorizedData = groupByCategory(data);
        this.insertCategoryLabels();
    }

    insertCategoryLabels() {
        var i = 0;
        this.categorizedData.forEach((catItem) => {
            var $row = this.rows[i].getRow();
            var $labelRow = $(`<tr class="category-header" style="background-color: ${catItem.background};">
                <td colspan="100">${catItem.name}</td>
            </tr>`);

            $labelRow.insertBefore($row);
            i += catItem.value.length;
        });
    }

    // add from outside for now. May want to revisit this
    // addRow(item) {
    //     this.data.push(item);
    //     this.data.sort((a, b) => {
    //         return a.index - b.index;
    //     });
    //     super(this.data, this.columnInfo, this.sendUpdate)
    //     this.categorizedData = groupByCategory(data);
    //     this.insertCategoryLabels();
    // }

    // removeRow($tr) {
    //     var idx = this.rows.findIndex((row) => row.getRow() == $tr);
    //     this.rows.splice(idx, 1);
    //     super(this.data, this.columnInfo, this.sendUpdate)
    //     this.categorizedData = groupByCategory(data);
    //     this.insertCategoryLabels();
    // }
}