import { replaceUrlId } from "./helpers";

var collapsedCaret = 'fa-caret-right';
var expandedCaret = 'fa-caret-down';

// EditItem represents an editable cell in a Datagrid. 
export class EditItem {
    constructor(datagrid, $td) {
        this.datagrid = datagrid;
        this.$td = $td;
        this.type = $td.attr('data-type');
        this.id = $td.parent().attr('item-id');
        this.field = $td.attr('field');
        this.contents = $td.text();
        this.errorMessage = "";
        this.isEditable = $td.attr('editable');
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
                    $date.datepicker({
                        autoclose: 'true',
                        format: 'mm/dd/yyyy',
                        defaultViewDate: startDate
                    }).on('changeDate', function(event) {
                        self.contents = event.format();
                        self.datagrid.sendUpdate(self);
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
                break;
            case 'selector':
                break;
            case 'number':
                self.$td.on('focus', function(event) {
                    self.clearError($(this));
                    $(this).selectText();
                });
                self.$td.on('blur', function(event) {
                    if(!isNaN($(this).text())) {
                        self.contents = $(this).text();
                        self.datagrid.sendUpdate(self);
                    } else {
                        $(this).text("0");
                    }
                });
                break;
            default:    // type 'text'
                self.$td.on('focus', function(event) {
                    self.clearError($(this));
                    $(this).selectText();                    
                });
                self.$td.on('blur', function(event) {
                    // var id = $(this).parent().attr('item-id');
                    // var field = $(this).attr('field');
                    self.content = $(this).text();
                    self.datagrid.sendUpdate(self);
                });
                break;
        }
    }

    showError(msg) {
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

    onUpdateSuccess() {
        switch(this.type) {
            case 'date':
                this.$td.text(this.contents);
                break;
            default:
                break;
        }
    }
}

// DataGrid is a class for creating a table that has editable cells that
// may update models on edit
export class DataGrid {
    // add data grid to grid (usually div tag). Give it a function to execute when a value is changed
    constructor(grid, updateFnc) {
        this.grid = grid;
        this.on_change_url = $(grid).attr('onchange-href');
        this.sendUpdate = updateFnc || this.defaultSendUpdate;
        this.initRows();

        if($(grid).find('td.expander') != undefined) {
            this.initCollapse();
        }
    }

    // initRows sets click highlighting for rows
    // TODO: remove highlighting when clicking off the table
    initRows() {
        self = this;
        $.each($(self.grid).find('tbody>tr'), function(i, tr) {
            $(tr).click(function(event) {
                if(!$(this).hasClass('active')) {
                    self.clearSelectedRow();
                    $(this).addClass('active');
                    self.setEditable($(this));
                }
            })
        });

        $(this.grid).on('focusout', function(event) {
            $(this).find('tr').each(function(i, tr) {
                $(tr).removeClass('active');
            });
        });
    }

    // initCollapse sets up tables with hidden rows to enable showing
    // row detail content. Add 'expander' class to the caret td tag
    initCollapse() {
        self = this;
        $.each($(this.grid).find('td.expander'), function(i, el) {
            $(el).click(function(event) {
                var $span = $(this).find('span');
                if($span.hasClass(collapsedCaret)) {
                    $span.removeClass(collapsedCaret);
                    $span.addClass(expandedCaret);
                    self.expandInfo($(this).parent().next());
                } else if($span.hasClass(expandedCaret)) {
                    $span.removeClass(expandedCaret);
                    $span.addClass(collapsedCaret);
                    self.collapseInfo($(this).parent().next());                    
                }
            });
        });
    }

    expandInfo($tr) {
        $tr.show();
    }

    collapseInfo($tr) {
        $tr.hide();
    }

    // clearSelectedRow unhighlighs a row
    clearSelectedRow() {
        var row = $(this.grid).find('tr.active');
        row.removeClass('active');
        this.removeEditable(row);
    }

    setEditable(row) {
        var editableRows = row.find('td[editable="true"]');
        editableRows.attr('contenteditable', true);

        editableRows.keydown(function(e) {
            if(e.keyCode == 13) {
                $(this).blur();
                if(e.getModifierState("Shift")) {
                    focusPrevRow($(this));
                } else {
                    focusNextRow($(this));
                }
                return false;
            }
            if(e.keyCode == 9) {
                $(this).blur();
                focusNextColumn($(this));
                return false;
            }
        })
    }

    removeEditable(row) {
        row.find('td[editable="true"]').removeAttr('contenteditable');
    }

    defaultSendUpdate(editItem) {
        console.log('default send update');
    }
}

function focusNextRow($el) {
    var colIdx = $el.index();
    var nextRow = $el.parent().next();
    if(nextRow.length == 0) {
        // $($el.parents('tr')[2]).blur();
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
        
        // nextRow.parents('tr').first().focus();
        // nextRow.focus();
        // debugger;
    }
    // debugger;
    nextRow.children().eq(colIdx).focus();
}

function focusNextColumn($el) {

}

// code from https://stackoverflow.com/questions/12243898/how-to-select-all-text-in-contenteditable-div/12244703
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