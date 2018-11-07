import { toGoName, replaceUrlId } from "./_helpers";

export function sendUpdate($form, updateObj) {
    var id = updateObj.id;
    submitUpdateForm($form, id, convertUpdateObj(updateObj));
}

function convertUpdateObj(updateObj) {
    var outObj = {}
    for(var key in updateObj) {
        outObj[toGoName(key)] = updateObj[key];
    }

    return outObj;
}

function submitUpdateForm($form, id, data) {
    $form.attr('action', replaceUrlId($form.attr('action'), id));

    for(var key in data) {
        var $field = $('<input />');
        $field.attr('type', 'hidden');
        $field.attr('name', key);
        $field.val(data[key]);
        $field.appendTo($form);
    }

    $form.submit();
}