import { toGoName, replaceUrlId } from "./_helpers";

export function sendUpdate($form, updateObj, submitFn) {
    var id = updateObj.id;
    submitFn = submitFn || ((f) => f.submit());
    submitUpdateForm($form, id, convertUpdateObj(updateObj), submitFn);
}

function convertUpdateObj(updateObj) {
    var outObj = {}
    for(var key in updateObj) {
        outObj[toGoName(key)] = updateObj[key];
    }

    return outObj;
}

function submitUpdateForm($form, id, data, submitFn) {
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