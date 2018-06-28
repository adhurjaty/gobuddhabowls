import { formatMoney, unFormatMoney } from "./helpers";
import { datepicker } from "./datepicker";


class Cell {
    constructor($td) {
        this.$td = $td;
        this.contents = $td.text();
    }

    getCell() {
        return this.$td;
    }

    getText() {
        return this.$td.text();
    }

    setText(text) {
        if(text) {
            this.$td.text(text);
        }
        this.contents = this.$td.text();
    }

    isEditable() {
        return this.constructor.name == "EditCell";
    }
}

// EditCell represents an editable cell in a Datagrid. 
class EditCell extends Cell {
    constructor($td) {
        super($td);
        this.type = $td.attr('data-type');
        this.errorMessage = "";
    }

    onError(msg) {
        this.errorMessage = msg;
        var $span = $('<span></span>');
        $span.text(msg);
        $span.css('color', 'red');
        this.$td.append($('<br>'));
        this.$td.append($span);
    }

    clearError() {
        this.$td.find('span').remove();
    }

    onSuccess() {
        switch(this.type) {
            case 'date':
                this.$td.text(this.contents);
                break;
            default:
                break;
        }
    }
}


class Row {
    constructor(cells, updateFnc) {
        this.cells = cells;
        this.updateFnc = updateFnc;
        this.$tr = null;
        this.initRow();
    }

    initRow() {
        var self = this;
        this.$tr = $(`<tr></tr>`);
        this.cells.forEach((cell) => {
            cell.getCell().appendTo(self.$tr);
            if(cell.isEditable()) {
                self.setListener(cell);
            }
        });
    }

    getRow() {
        return this.$tr;
    }

    // setListener initializes the behaviors of the different types of cell data types
    setListener(cell) {
        var self = this;
        switch(cell.type) {
        case 'date':
            cell.$td.on('focus', function(event) {
                cell.clearError();
                var $date = $('<input data-provide="datepicker" value="' + cell.contents + '">');
                $date.css('width', '85px');
                $(this).empty();
                $(this).append($date);
                var startDate = cell.contents ? cell.contents : new Date().toLocaleDateString("en-US");
                // debugger;
                // $date.datepicker({
                datepicker($date, {
                    autoclose: 'true',
                    format: 'mm/dd/yyyy',
                    defaultViewDate: startDate
                }).on('changeDate', function(event) {
                    cell.contents = event.format();
                    cell.setText(cell.contents);
                    self.sendUpdate();
                }).on('hide', function(event) {
                    cell.$td.text(cell.contents);
                    if(cell.errorMessage) {
                        cell.showError(cell.errorMessage);
                    }
                });
                $date.focus();
            });
            break;
        // TODO: fill these options in
        case 'money':
            cell.$td.on('focusin', function(event) {
                cell.clearError($(this));
                cell.setText(unFormatMoney(cell.getText()));
                $(this).selectText();
            });
            cell.$td.on('focusout', function(event) {
                // HACK: event firing multiple times causes
                // text to go to $0.00 without this
                var text = cell.getText().replace('$', '')
                if(text == undefined) {
                    return;
                }
                // debugger;

                if(!isNaN(text)) {
                    var amt = parseFloat(text);
                    $(this).attr('value', amt);
                    cell.setText(formatMoney(amt));
                    self.sendUpdate();
                } else {
                    debugger;
                    cell.setText("$0.00");
                }
            });
            break;
        case 'selector':
            break;
        case 'number':
            cell.$td.on('focusin', function(event) {
                cell.clearError();
                $(this).selectText();
            });
            cell.$td.on('focusout', function(event) {
                if(!isNaN(cell.getText())) {
                    cell.setText()
                    self.sendUpdate();
                } else {
                    $(this).text("0");
                }
            });
            break;
        default:    // type 'text'
            cell.$td.on('focusin', function(event) {
                self.clearError($(this));
                $(this).selectText();                    
            });
            cell.$td.on('focusout', function(event) {
                cell.setText();
                self.sendUpdate();
            });
            break;
        }
    }

    sendUpdate() {
        var updateObj = {}
        this.$tr.find('td[name]').each((_, td) => {
            updateObj[$(td).attr('name')] = $(td).text()
        });
        this.updateFnc(updateObj);
    }

    getEditableCells() {
        return this.cells.filter((cell) => cell.isEditable());
    }
}

// DataGrid is a class for creating a table that has editable cells that
// may update models on edit
export class DataGrid {
    // data: data for the datagrid
    // columnInfo: array of objects that describe the columns of the form
    // {
    //     header: 'header',
    //     name: 'name',
    //     editable: optional, default false,
    //     data_type: 'text',
    //     hidden: optional, default false
    //     column_func: function(data_item)
    // }
    // updateFnc: function to execute after updating an editable cell
    constructor(data, columnInfo, updateFnc) {
        this.data = data;
        this.columnInfo = columnInfo;
        this.sendUpdate = updateFnc || this.defaultSendUpdate;
        this.$table = null;
        this.rows = null;
        this.initTable();
    }

    getTable() {
        return this.$table;
    }

    initTable() {
        var self = this;
        this.$table = $('<table class="datagrid"></table>');
        var headersArr = this.columnInfo.map((info) => {
            var $header = $('<th></th>');
            if(info.hidden) {
                $header.hide();
            }
            if(info.header) {
                $header.text(info.header);
            }
            return $header;
        });
        var $header = $('<thead></thead>');
        headersArr.forEach(($h) => {
            $h.appendTo($header);
        });
        $header.appendTo(this.$table);

        this.rows = this.data.map((item) => {
            var cells = self.columnInfo.map((info) => {
                var $td = $('<td></td>');
                $td.attr('name', info.name);
                $td.html(info.column_func(item));

                if(info.hidden) {
                    $td.hide();
                }
                if(info.editable) {
                    $td.attr('editable', true);
                    $td.attr('data-type', info.data_type);
                    return new EditCell($td);
                }

                return new Cell($td);
            });
            return new Row(cells, this.sendUpdate);
        });

        var $tbody = $('<tbody></tbody>');
        this.rows.forEach((row) => {
            self.setRowClick(row);
            var $row = row.getRow();
            // TODO: remove highlighting when clicking off the table
            $row.appendTo($tbody);
        });
        $tbody.appendTo(this.$table);
    }

    setRowClick(row) {
        var self = this;
        row.getRow().click(function() {
            if(!$(this).hasClass('active')) {
                self.clearSelectedRow();
                $(this).addClass('active');
                self.setEditable(row);
            }
        });
    }

    // clearSelectedRow unhighlighs a row
    clearSelectedRow() {
        var $row = this.$table.find('tr.active');
        $row.removeClass('active');
        this.removeEditable($row);
    }

    setEditable(row) {
        var self = this;
        var $row = row.getRow();
        var editableCells = $row.find('td[editable="true"]');
        editableCells.attr('contenteditable', true);
        var fired = false;

        editableCells.keydown(function(e) {
            var colIdx = $(this).index();
            var rowIdx = self.rows.indexOf(row);
            if(e.keyCode == 13 || e.keyCode == 9) {
                if(fired) {
                    return false;
                }
                fired = true;
                setTimeout(function() { fired = false; }, 200);
                $(this).blur();

                // ENTER key
                if(e.keyCode == 13) {
                    if(window.event.getModifierState("Shift")) {
                        self.focusPrevRow(rowIdx, colIdx);
                    } else {
                        self.focusNextRow(rowIdx, colIdx);
                    }
                    return false;
                }
                // TAB key
                if(e.keyCode == 9) {
                    if(window.event.getModifierState("Shift")) {
                        self.focusPrevColumn(rowIdx, colIdx);
                    } else {
                        self.focusNextColumn(rowIdx, colIdx);
                    }
                    return false;
                }
            }
            
        });
    }

    removeEditable($row) {
        $row.find('td[editable="true"]').removeAttr('contenteditable');
    }

    removeItem($row) {
        $row.remove();
    }

    addItem($row, idx) {
        $.each($row.find('td[editable="true"]'), function(i, el) {
            return new EditItem(self, $(el));
        });

    }

    defaultSendUpdate($row, editCell) {
        console.log('default send update');
    }

    focusPrevRow(rowIdx, colIdx) {
        if(rowIdx > 0) {
            rowIdx--;
        }
        var $el = this.rows[rowIdx].cells[colIdx].getCell();
        $el.click();
        $el.focus();
    }

    focusNextRow(rowIdx, colIdx) {
        if(rowIdx < this.rows.length - 1) {
            rowIdx++;
        }
        var $el = this.rows[rowIdx].cells[colIdx].getCell();
        $el.click();
        $el.focus();
    }

    focusPrevColumn(rowIdx, colIdx) {
        var row = this.rows[rowIdx];
        var colNum = row.cells.length;
        var i = colIdx-1;
        if(i < 0) {
            i = colNum - 1;
        }
        
        while(i != colIdx) {
            var cell = row.cells[i];
            if(cell.isEditable()) {
                cell.getCell().click();
                cell.getCell().focus();
                break;
            }
            i--;
            if(i < 0) {
                i = colNum - 1;
            }
        }
    }

    focusNextColumn(rowIdx, colIdx) {
        var row = this.rows[rowIdx];
        var colNum = row.cells.length;
        var i = colIdx+1 % colNum;
        while(i != colIdx) {
            var cell = row.cells[i];
            if(cell.isEditable()) {
                cell.getCell().click();
                cell.getCell().focus();
                break;
            }
            i++;
            i %= colNum;
        }
    }

}

// code from https://stackoverflow.com/questions/12243898/how-to-select-all-text-in-contenteditable-div/12244703
// for highlighting all text when clicking in a cell
$.fn.selectText = function(){
    var doc = document;
    var element = this[0];
    if (doc.body.createTextRange) {
        var range = document.body.createTextRange();
        range.moveToElementText(element);
        range.select();
    } else if (window.getSelection) {
        var selection = window.getSelection();        
        var range = document.createRange();
        range.selectNodeContents(element);
        selection.removeAllRanges();
        selection.addRange(range);
    }
 };