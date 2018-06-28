export function addToRemaining(item) {
    var $container = $('#add-order-modal');
    var remainingItems = JSON.parse($container.attr('data'));

    remainingItems.push(item);
    remainingItems.sort((a, b) => {
        return a.index - b.index;
    });
    $container.attr('data', JSON.stringify(remainingItems));
    populateRemaining();
}

export function removeFromRemaining(item) {
    var $container = $('#add-order-modal');
    var remainingItems = JSON.parse($container.attr('data'));
    var idx = remainingItems.indexOf(item);

    remainingItems.splice(idx, 1);
    $container.attr('data', JSON.stringify(remainingItems));
    populateRemaining();
}

function populateRemaining() {
    var $container = $('#add-order-modal');
    var remainingItems = JSON.parse($container.attr('data'));
    var $select = $container.find('select');

    remainingItems.forEach((item) => {
        $('<option/>').val(item.id).html(item.name).appendTo($select);
    });
}