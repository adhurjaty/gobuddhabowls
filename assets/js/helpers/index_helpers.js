import { toGoName, replaceUrlId } from "./_helpers";

export function sendUpdate($form, updateObj, submitFn) {
    var id = updateObj.id;
    submitUpdateForm($form, id, convertUpdateObj(updateObj), submitFn);
}

function convertUpdateObj(updateObj) {
    var outObj = {}
    for(var key in updateObj) {
        outObj[toGoName(key)] = updateObj[key];
    }

    return outObj;
}

export function submitUpdateForm($form, id, data, submitFn) {
    submitFn = submitFn || ((f) => f.submit());    
    $form.attr('action', replaceUrlId($form.attr('action'), id));
    clearForm($form);

    for(var key in data) {
        var $field = $('<input />');
        $field.attr('type', 'hidden');
        $field.attr('name', key);
        $field.val(data[key]);
        $field.appendTo($form);
    }

    submitFn($form);

    // HACK: need to reset the form action URL in the case of 
    // repeated calls
    resetFormAction($form, id);
}

function clearForm($form) {
    $form.find('*')
         .not('input[name="authenticity_token"]')
         .not('input[name="_method"]')
         .remove();
}

function resetFormAction($form, id) {
    var url = $form.attr('action');
    var resetUrl = url.replace(id, '{id}');

    $form.attr('action', resetUrl);
}

export function sendAjax($form, isSync) {
    isSync = isSync || false;
    var data = {};
    $form.find('input').each((i, el) => {
        data[$(el).attr('name')] = $(el).val();
    });

    $.ajax({
        url: $form.attr('action'),
        method: 'POST',
        async: !isSync,
        data: data
    });
}

export function showError(msg) {
    var errContainer = $('#error-container');
    if(errContainer.length == 0) {
        errContainer = $('<div id="error-container" class="alert alert-danger"></div>');
        $('main').append(errContainer);
    }
    errContainer.html(msg);
}