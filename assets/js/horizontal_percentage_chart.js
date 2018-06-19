import { formatMoney } from "./helpers";

export function horizontalPercentageChart(title, items, total) {
    var head = `
    <div class="container horizontal rounded" style="width: 100%">

        <h4>${title}</h4>
        <div class="bar-holder">
            <div class="percentage-bar horizontal">
                <div class="progress-track">
                    <div class="progress-fill" style="background: #666;width: 100%; text-align: center;">
                        <span>Total</span>
                    </div>
                </div>
            </div>
            <div class="bar-label">
                <span>${formatMoney(total)}</span>
            </div>
        </div>
        <div id="horizontal-percentage-bars">`;

    var foot = `
        </div>
    </div>`;

    return getStyle() + head + getBars(categorize(items), total) + foot;

}

function getStyle() {
    return `
    <style>
    .horizontal>h4 {
        text-align: center;
    }

    .horizontal .percentage-bar {
        float: left;
        height: 35px;
        width: 90%;
    }

    .bar-holder {
        height: 35px;
        width: 100%;
        padding: 12px 0;
    }

    .horizontal .progress-track {
        position: relative;
        width: 100%;
        height: 30px;
        background: #ebebeb;
    }

    .horizontal .progress-fill {
        position: relative;
        height: 30px;
        width: 50%;
        color: #fff;
        text-align: right;
        padding-right: 20px;
        font-family: "Lato","Verdana",sans-serif;
        font-size: 16px;
        line-height: 30px;
        float: left;
    }

    .progress-empty {
        height: 30px;
        line-height: 30px;
        overflow: hidden;
        font-size: 16px;
        text-align: left;
        padding-left: 20px;
        font-family: "Lato","Verdana",sans-serif;        
    }

    .rounded .progress-track,
    .rounded .progress-fill {
        border-radius: 3px;
        box-shadow: inset 0 0 5px rgba(0,0,0,.2);
    }

    .bar-label {
        float: right;
        width: 8%;
        text-align: left;
        height: 100%;
    }

    .bar-label>span {
        line-height: 35px;
    }
    </style>`;
}

function getBars(categorizedItems, total) {
    return categorizedItems.map((category) => {
        var proportion = category.value / total * 100;
        debugger;
        return `
        <div class="bar-holder">
            <div class="percentage-bar horizontal">
                <div class="progress-track">
                    <div class="progress-fill" style="background: ${category.background};width: ${proportion}%;">
                        ${(proportion >= 50) ? `<span>${category.name}</span>` : ''}
                    </div>
                    <div class="progress-empty">
                        ${(proportion < 50) ? `<span>${category.name}</span>` : ''}
                    </div>
                </div>
            </div>
            <div class="bar-label">
                <span>${formatMoney(category.value)}</span>
            </div>
        </div>`
    }).join('');
}

function categorize(items) {
    return items.reduce((categorizedItems, item) => {
        var value = item.price * item.count;
        var category = categorizedItems.find((x) => x.name == item.InventoryItem.Category.name);

        if(category) {
            category.value += value;
        } else {
            categorizedItems.push({
                index: item.InventoryItem.Category.index,
                name: item.InventoryItem.Category.name,
                value: value,
                background: item.InventoryItem.Category.background
            });
        }

        return categorizedItems;
    }, []).sort((a, b) => {
        return a.index - b.index;
    });
}
