import { sendUpdate, sendAjax } from './helpers/index_helpers';
import 'spectrum-colorpicker';

$(() => {
    setupSortable();
    setupEditName();
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

function setupEditName() {
    var nameSpans = $('#categories-movable').find('li span:first-child');
    nameSpans.on('click', (e) => {
        var span = $(e.target);
        var text = span.html().trim();
        var input = $(`<input type="text" value="${text}"></input>`);
        var li = span.parent();
        li.prepend(input);
        span.hide();
        
        input.on('blur', ((span) => {
            return (e) => {
                var el = $(e.target)
                var text = el.val();
                el.remove();
                span.html(text);
                span.show();
            }
        })(span));
    });
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
        data['name'] = $(el).find('span:first-child').first().html();
        data['index'] = i;

        sendUpdate($form, data, sendAjax);
    });

    location.replace('/settings');
}