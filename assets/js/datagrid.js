import { formatMoney, unFormatMoney } from "./helpers";
import { datepicker } from "./datepicker";

// EditCell represents an editable cell in a Datagrid. 
class EditCell {
    constructor($td, sendUpdate) {
        this.$td = $td;
        this.sendUpdate = sendUpdate;
        this.type = $td.attr('data-type');
        this.contents = $td.text();
        this.errorMessage = "";
        this.setListener();
    }

    // setListener initializes the behaviors of the different types of cell data types
    setListener() {
        var self = this;

        switch(self.type) {
        case 'date':
            self.$td.on('focus', function(event) {
                self.clearError($(this));
                var $date = $('<input data-provide="datepicker" value="' + self.contents + '">');
                $date.css('width', '85px');
                $(this).empty();
                $(this).append($date);
                var startDate = self.contents ? self.contents : new Date().toLocaleDateString("en-US");
                // debugger;
                // $date.datepicker({
                datepicker($date, {
                    autoclose: 'true',
                    format: 'mm/dd/yyyy',
                    defaultViewDate: startDate
                }).on('changeDate', function(event) {
                    self.contents = event.format();
                    self.$td.text(self.contents);
                    self.sendUpdate(self);
                }).on('hide', function(event) {
                    self.$td.text(self.contents);
                    if(self.errorMessage) {
                        self.showError(self.errorMessage);
                    }
                });
                $date.focus();
            });
            break;
        // TODO: fill these options in
        case 'money':
            self.$td.on('focusin', function(event) {
                self.clearError($(this));
                $(this).text(unFormatMoney($(this).text()));
                $(this).selectText();
            });
            self.$td.on('focusout', function(event) {
                // HACK: event firing multiple times causes
                // text to go to $0.00 without this
                var text = $(this).text().replace('$', '')
                if(text == undefined) {
                    return;
                }
                // debugger;

                if(!isNaN(text)) {
                    var amt = parseFloat(text);
                    $(this).attr('value', amt);
                    self.contents = formatMoney(amt);
                    $(this).text(self.contents);
                    self.sendUpdate(self);
                } else {
                    debugger;
                    $(this).text("$0.00");
                }
            });
            break;
        case 'selector':
            break;
        case 'number':
            self.$td.on('focusin', function(event) {
                self.clearError($(this));
                $(this).selectText();
            });
            self.$td.on('focusout', function(event) {
                if(!isNaN($(this).text())) {
                    self.contents = $(this).text();
                    self.sendUpdate(self);
                } else {
                    $(this).text("0");
                }
            });
            break;
        default:    // type 'text'
            self.$td.on('focusin', function(event) {
                self.clearError($(this));
                $(this).selectText();                    
            });
            self.$td.on('focusout', function(event) {
                // var id = $(this).parent().attr('item-id');
                // var field = $(this).attr('field');
                self.contents = $(this).text();
                self.sendUpdate(self);
            });
            break;
        }
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

// may want to store cells and rows as classes, not sure yet
// class Cell {
//     constructor($cell) {
//         this.$cell = $cell;
//     }
//     getCell() {
//         return this.$cell;
//     }
// }

// class Row {
//     constructor($cells) {
//         this.$cells = $cells;
//         this.$row = null;
//         this.initRow();
//     }

//     initRow() {
//         var self = this;
//         this.$row = $(`<tr></tr>`);
//         this.$cells.forEach(($cell) => {
//             $cell.appendTo(self.$row);
//         });
//     }

//     getRow() {
//         return this.$row;
//     }
// }

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
        this.$rows = null;
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

        this.$rows = this.data.map((item) => {
            var $row = $('<tr></tr>');
            self.columnInfo.forEach((info) => {
                var $td = $('<td></td>');
                $td.attr('name', info.name);
                $td.html(info.column_func(item));

                if(info.hidden) {
                    $td.hide();
                }
                if(info.editable) {
                    $td.attr('editable', true);
                    $td.attr('data-type', info.data_type);
                    new EditCell($td, ((item, info) => {
                        var updateObj = {
                            id: item.id,
                            name: info.name,
                        }
                        return (editCell) => {
                            updateObj.value = editCell.contents;
                            updateObj.onError = editCell.onError;
                            updateObj.onSuccess = editCell.onSuccess;
                            return self.sendUpdate(updateObj);
                        }
                    })(item, info));
                }
                $td.appendTo($row);
            });
            return $row;
        });

        var $tbody = $('<tbody></tbody>');
        this.$rows.forEach(($row) => {
            self.setRowClick($row);
            // TODO: remove highlighting when clicking off the table
            $row.appendTo($tbody);
        });
        $tbody.appendTo(this.$table);
    }

    setRowClick($row) {
        var self = this;
        $row.click(function() {
            if(!$(this).hasClass('active')) {
                self.clearSelectedRow();
                $(this).addClass('active');
                self.setEditable($(this));
            }
        });
    }

    // clearSelectedRow unhighlighs a row
    clearSelectedRow() {
        var $row = this.$table.find('tr.active');
        $row.removeClass('active');
        this.removeEditable($row);
    }

    setEditable($row) {
        var editableCells = $row.find('td[editable="true"]');
        editableCells.attr('contenteditable', true);
        var fired = false;

        editableCells.keydown(function(e) {
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
                        focusPrevRow($(this));
                    } else {
                        focusNextRow($(this));
                    }
                    return false;
                }
                // TAB key
                if(e.keyCode == 9) {
                    if(window.event.getModifierState("Shift")) {
                        focusPrevColumn($(this));
                    } else {
                        focusNextColumn($(this));
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
}

function focusPrevRow($el) {
    var colIdx = $el.index();
    var nextRow = $el.parent().prev();
    if(nextRow.length == 0) {
        var $td = $el.parents('.datagrid')
                     .first()
                     .parents('tr')
                     .prev().prev()
                     .find('td')
                     .first();

        $td.click();
        nextRow = $td.find('.datagrid')
                        .find('tr')
                        .last();

    }
    if(nextRow.length == 0) {
        $el.focus();
    } else {
        nextRow.children().eq(colIdx).click();
        nextRow.children().eq(colIdx).focus();
    }
}

function focusNextRow($el) {
    var colIdx = $el.index();
    var nextRow = $el.parent().next();
    if(nextRow.length == 0) {
        var $td = $el.parents('.datagrid')
                     .first()
                     .parents('tr')
                     .next().next()
                     .find('td')
                     .first();
        $td.click();
        nextRow = $td.find('.datagrid')
                     .find('tr')
                     .eq(1);
    }
    if(nextRow.length == 0) {
        $el.focus();
    } else {
        nextRow.children().eq(colIdx).click();
        nextRow.children().eq(colIdx).focus();
    }
}

function focusPrevColumn($el) {
    var cols = $el.parent().children();
    var colNum = cols.length;
    var colIdx = $el.index();
    var i = colIdx-1;
    if(i < 0) {
        i = colNum - 1;
    }
    
    while(i != colIdx) {
        if(cols.eq(i).attr('editable') == 'true') {
            cols.eq(i).focus();
            cols.eq(i).focusin();
            break;
        }
        i--;
        if(i < 0) {
            i = colNum - 1;
        }
    }
}

function focusNextColumn($el) {
    var cols = $el.parent().children();
    var colNum = cols.length;
    var colIdx = $el.index();
    var i = colIdx+1 % colNum;
    while(i != colIdx) {
        if(cols.eq(i).attr('editable') == 'true') {
            cols.eq(i).focus();
            cols.eq(i).focusin();
            break;
        }
        i++;
        i %= colNum;
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