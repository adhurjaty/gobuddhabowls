import { sendUpdate, sendAjax } from './helpers/index_helpers';
import 'spectrum-colorpicker';

$(() => {
    setupSortable();
    setupColorPicker();
    setupSubmitButton();
});

function setupSortable() {
    var el = document.getElementById('categories-movable');
    if(el != undefined) {
        var sortable = Sortable.create(el, {
            group: {
                name: "components",
                pull: function(to, from, dragEl, evt) {
                  if(evt.type === 'dragstart') {
                    return false;
                  }
                  return true;
                }
            },
            animation: 150,
            handle: '.drag-handle'
        });
    }
}

function setupColorPicker() {
    var components = $('input[name="color"]');
    components.each((i, el) => {
        $(el).spectrum({
            hideAfterPaletteSelect: true,
            color: $(el).val(),
            preferredFormat: "hex"
        });
    });
    components.on('hide.spectrum', (e, color) => {
        $(e.currentTarget).attr('value', color.toHexString());
    });
}

function setupSubmitButton() {
    $('#save-inv-item-categories').click(saveInvItemsCategories);
}

function saveInvItemsCategories() {
    var $form = $('#update-category-form');

    $('#categories-movable').find('li').each(function(i, el) {
        var data = {};
        data['id'] = $(el).attr('itemid');
        data['background'] = $(el).find('input[name="color"]').val();
        data['index'] = i;

        sendUpdate($form, data, sendAjax);
    });

    location.replace('/settings');
}