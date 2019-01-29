export class OrderingTable {
    constructor(items) {
        this.items = items;
        this.ul = null;
        this.initTable();
    }

    initTable() {
        this.ul = $(`<ul class="list-group" name="movable-list"></ul>`);
        this.items.forEach(item => {
            var li = this.getListElement(item);
            this.ul.append(li);
        }, this);
    }

    getListElement(item) {
        return $(
            `<li itemid="${item.id}" index="${item.index}"
                class="list-group-item d-flex justify-content-between 
                align-items-center">
                    <span>${item.name}</span>
                    <span class="drag-handle" style="font-size: 20px;">
                        ☰
                    </span>
            </li>`
        );
    }

    enableDragging() {
        var sortable = Sortable.create(this.ul.get(0), {
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

    attach(element) {
        element.empty();
        element.append(this.ul);

        this.enableDragging();
    }
}