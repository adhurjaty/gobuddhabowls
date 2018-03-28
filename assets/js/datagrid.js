$.each($('.datagrid').find('td[contenteditable="true"]'), function(i, el) {
    $(el).on('focus', function(event) {
        console.log('before');
    });
    $(el).on('blur', function(event) {
        console.log('after');
    });
    // switch(type) {
    //     case 'date':
    //         el = '<input data-provide="datepicker" format="mm/dd/yyyy" startDate="' + $td.text + '">'
    //         break;
    //     case 'money':
    //         break;
    //     case 'selector':
    //         break;
    //     default:    // type 'text'
    //         el = '<input type="text" class="form-control" value="' + prevValue + '" onfocus="this.select()">'
    //         $td.append(el);
    //         $(el).focus().select();
    // }
});