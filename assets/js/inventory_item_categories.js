import 'bootstrap-colorpicker';
import { sendUpdate } from './helpers/index_helpers';

$(() => {
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

    $('.colorpicker-component').colorpicker({
        useAlpha: false
    });
    $('#save-inv-item-categories').click(saveInvItemsCategories);
});

function saveInvItemsCategories() {
    var $form = $('#update-category-form');
    $('#categories-movable').find('li').each(function(i, el) {
        var data = {};
        data['id'] = $(el).attr('itemid');
        data['background'] = $(el).find('input[name="Background"]').val();
        data['index'] = i;

        sendUpdate($form, data);
    });

    location.replace('/settings');
}