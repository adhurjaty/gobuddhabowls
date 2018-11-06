import 'bootstrap-colorpicker';

$(() => {
    var el = document.getElementById('categories-movable');
    if(el != undefined) {
        debugger;
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
    $('#save-inv-item-categories').click(SaveInvItemsCategories);
});

function SaveInvItemsCategories() {
    var url = $('#categories-movable').attr('on-save-url');
    $('#categories-movable').find('li').each(function(i, el) {
        var id = $(el).attr('itemid');
        var data = {};
        data['Background'] = $(el).find('input[name="Background"]').val();
        data['Index'] = i;

        $.ajax({
            url: url + '/' + id,
            data: data,
            method: 'PUT',
            error: function(xhr, status, err) {
                var errMessage = xhr.responseText;
                debugger;
            },
            success: function(data, status, xhr) {
                
            }
        });
    });

    location.replace('/settings');
}