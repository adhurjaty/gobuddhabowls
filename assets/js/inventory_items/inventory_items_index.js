import { parseModelJSON, groupByCategory } from "../helpers/_helpers";

$(() => {
    insertCategoryLabels();
    enableDragging();
});

function insertCategoryLabels() {
    var container = $('#categories-movable');
    var data = parseModelJSON(container.attr('data'));
    var categorizedData = groupByCategory(data);

    var rows = container.find('li').toArray(); 

    var i = 0;
    categorizedData.forEach((catItem) => {
        var $row = $(rows[i]);
        var $labelRow = $(`<li class="list-group-item d-flex justify-content-between align-items-center" style="background-color: ${catItem.background};">
            <span>${catItem.name}</span>
            <span class="drag-handle" style="font-size: 20px;">☰</span>
        </li>`);

        $labelRow.insertBefore($row);
        i += catItem.value.length;
    });
}

function enableDragging() {
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

