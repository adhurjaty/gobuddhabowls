import { sortItems } from "../helpers/_helpers";

export class Modal {
    constructor(items, addFn) {
        this.items = items;
        this.$content = null;
        this.$select = null;
        this.addFn = addFn;

        this.initModal();
    }

    initModal() {

        this.$content = $(`
            <div id="add-item-modal" class="modal fade" tabindex="-1" role="dialog" aria-labelledby="AddItem" aria-hidden="true">
                <div class="modal-dialog modal-dialog-centered" role="document">
                    <div class="modal-content">
                        <div class="modal-header">
                            <h5 class="modal-title" id="exampleModalLongTitle">Add item</h5>
                            <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                                <span aria-hidden="true">&times;</span>
                            </button>
                        </div>
                        <form>
                            <div class="modal-body">
                            </div>
                            <div class="modal-footer">
                                <button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button>
                                
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        `);


        this.setSelect();
        this.setSubmit();
    }

    setSelect() {
        this.$select = $(
            `<select>
                ${this.items.map(item => {
                    return `<option value="${item.id}">${item.name}</option>`;
                })}
            </select>`
        );
        
        this.$content.find('div.modal-body').html(this.$select);
    }

    setSubmit() {
        var submit = $(`<button type="button" class="btn btn-success" data-dismiss="modal" id="add-po-item-submit">OK</button>`);
        var self = this;
        submit.click(() => {
            var id = self.$select.find('option:selected').val();
            var item = self.removeItem(id);
            self.addFn(item);
        });
        submit.appendTo(this.$content.find('div.modal-footer'));
    }

    addItem(item) {
        this.items.push(item);
        this.sortItems();
        this.setSelect();
    }

    removeItem(id) {
        var index = this.items.findIndex(x => x.id == id);
        var item = this.items[index];
        this.items.splice(index, 1);
        this.setSelect();
        return item;
    }

    sortItems() {
        this.items = sortItems(this.items);
    }
    
    getContent() {
        return this.$content;
    }
}