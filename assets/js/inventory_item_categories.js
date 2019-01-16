import { sendUpdate, sendAjax } from './helpers/index_helpers';
import 'spectrum-colorpicker';
import { replaceUrlId } from './helpers/_helpers';

$(() => {
    setupSortable();
    setupEditName();
    setupColorPicker();
    setupAddButton();
    // delete would destroy the data as implemented
    // setupDeleteButton();
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

function setupAddButton() {
    $('#add-category-button').click(() => {
        makeLi();
    });
}

function makeLi() {
    // var input = $('<input type="text" />');
    var li = $(`
        <li itemid="00000000-0000-0000-0000-000000000000"
            class="list-group-item d-flex justify-content-between 
            align-items-center">
            <span>NewCategory</span>
            <input name="color" type="text" value="" />
                    
            <span class="drag-handle" style="font-size: 20px;">☰</span>
        </li>`
    );
    var ul = $('#categories-movable');
    ul.prepend(li);

    setupEditName();
    setupColorPicker();
    setupSortable();
}

function setupDeleteButton() {
    var button = $('#delete-category-button');
    button.click(() => {
        var listEls = $('#categories-movable').find('li');
        if(button.hasClass('active')) {
            button.removeClass('active');
            listEls.each((i, el) => {
                var li = $(el);
                li.find('a').remove();
                li.find('span:last-child').show();
            });
        } else {
            button.addClass('active');
            listEls.each((i, el) => {
                var li = $(el);
                var span = li.find('span:last-child');
                span.hide();
                var delLink = makeDeleteLink(li.attr('itemid'));
                li.append(delLink);
            });
        }
    });
}

function makeDeleteLink(id) {
    var url = $('#delete-category-button').attr('data-link');
    return $(`<a href="${replaceUrlId(url, id)}" data-method="DELETE"
        data-confirm="Are you sure?">
            <span class="fa fa-minus-circle"
                style="font-size: 20px; color: rgb(200, 0, 0);"></span>
        </a>`);
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
        data['name'] = $(el).find('span:first-child').first().html().trim();
        data['index'] = i;

        sendUpdate($form, data, sendAjax);
    });

    location.replace('/settings');
}