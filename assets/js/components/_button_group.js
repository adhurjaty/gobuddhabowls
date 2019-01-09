
// + and - buttons that pop open modal or remove item respectively
export class ButtonGroup {
    constructor() {
        this.createButtonGroup();
    }

    createButtonGroup() {
        this.addButton = $(`<button class="btn btn-default" 
                type="button" data-toggle="modal"
                data-target="#add-item-modal" disabled>

                <span class="fa fa-plus">
            </button>`);
        this.removeButton = $(`<button class="btn btn-default"
                type="button" disabled>
                <span class="fa fa-minus">
            </button>`);
        this.$content = $(`<div class="input-group d-flex 
            justify-content-end"></div>`);
        this.addButton.appendTo(this.$content);
        this.removeButton.appendTo(this.$content);
    }

    setAddListenr(fn) {
        this.addButton.click(() => {
            fn();
        });
    }

    setRemoveListener(fn) {
        this.removeButton.click(() => {
            fn();
        });
    }

    enableAddButton() {
        this.enableButton(this.addButton);
    }

    disableAddButton() {
        this.disableButton(this.addButton);
    }

    enableRemoveButton() {
        this.enableButton(this.removeButton);
    }

    disableRemoveButton() {
        this.disableButton(this.removeButton);
    }

    enableButton(btn) {
        btn.removeAttr('disabled');
    }

    disableButton(btn) {
        btn.attr('disabled', 'disabled');
    }

    getGroup() {
        return this.group;
    }
}