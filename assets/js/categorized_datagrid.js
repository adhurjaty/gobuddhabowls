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
}