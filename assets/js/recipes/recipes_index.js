import { formatMoney, parseModelJSON, replaceUrlId, dblclickEdit } from "../helpers/_helpers";
import { CollapseCategorizedDatagrid } from "../datagrid/_collapse_categorized_datagrid";
import { CategorizedDatagrid } from "../datagrid/_categorized_datagrid";
import { sendUpdate, sendAjax } from "../helpers/index_helpers";

var _options = {
    breakdown: false
};

var _editPath = $('#data-holder').attr('data-url');
var _recipePath = $('#update-recipe-form').attr('action');

var _columnInfo = [
    {
        name: 'id',
        hidden: true,
        get_column: (recipe) => {
            return recipe.id;
        }
    },
    {
        name: 'index',
        hidden: true,
        get_column: (recipe) => {
            return recipe.index;
        }
    },
    {
        name: 'name',
        header: 'Name',
        editable: true,
        get_column: (recipe) => {
            return recipe.name;
        },
        set_column: (recipe, name) => {
            recipe.name = name;
        }
    },
    {
        name: 'recipe_unit',
        header: 'RU',
        editable: true,
        get_column: (recipe) => {
            return recipe.recipe_unit;
        },
        set_column: (recipe, ru) => {
            recipe.recipe_unit = ru;
        }
    },
    {
        name: 'recipe_unit_conversion',
        header: 'Yield',
        editable: true,
        data_type: 'number',
        get_column: (recipe) => {
            return recipe.recipe_unit_conversion;
        },
        set_column: (recipe, ruc) => {
            recipe.recipe_unit_conversion = ruc;
        }
    },
    {
        header: 'Cost',
        get_column: (recipe) => {
            return formatMoney(calculateRecipeCost(recipe.Items),
                recipe.recipe_unit_conversion);
        }
    },
    {
        name: 'dropdown',
        get_column: ((editPath, recipePath) => {
            return (recipe) => {
                return `<div class="dropdown show">
                    <button type="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                        ...
                    </button>
                    <div class="dropdown-menu">
                        <a href="${replaceUrlId(editPath, recipe.id)}"
                            class="dropdown-item">Edit</a>
                        <a name="delete" class="dropdown-item text-danger"
                            data-method="DELETE"
                            href="${replaceUrlId(recipePath, recipe.id)}">
                            Delete
                            </a>
                    </div>
                </div>`
            };
        })(_editPath, _recipePath)
    }
];

var _subColumnInfo = [
    {
        name: 'id',
        hidden: true,
        get_column: (item) => {
            return item.id;
        }
    },
    {
        header: 'Name',
        get_column: (item) => {
            return item.name;
        }
    },
    {
        header: 'RU',
        get_column: (item) => {
            return item.recipe_unit;
        }
    },
    {
        header: 'Count',
        get_column: (item) => {
            return item.count;
        }
    },
    {
        header: 'RU Cost',
        get_column: (item) => {
            return formatMoney(item.price);
        }
    },
    {
        header: 'Ext',
        get_column: (item) => {
            return formatMoney(item.price * item.count);
        }
    }
];

var _collapseInfo = (recipe) => {
    var dg = new CategorizedDatagrid(recipe.Items, _subColumnInfo);
    var row = `<tr><td colspan="100"><div>
            ${dg.$table.prop('outerHTML')}
        </div></td></tr>`;
    return row;
};

$(() => {
    createDatagrid()
});

function createDatagrid() {
    var dataHolder = $('#data-holder');
    var items = parseModelJSON(dataHolder.attr('data'));
    var batchContainer = $('#batch-datagrid');
    var menuContainer = $('#menu-datagrid');

    [batchContainer, menuContainer].forEach((container, i) => {
        var isBatch = i == 0;
        var recItems = items.filter(x => x.is_batch == isBatch);
        var dg = new CollapseCategorizedDatagrid(recItems, _columnInfo,
            _collapseInfo, onRecipeUpdated);
        $(container).html(dg.$table);
    });
}

function calculateRecipeCost(items, ruc) {
    return items.reduce((total, item) => {
        return total + (item.price * item.count);
    }, 0)
}

function onRecipeUpdated(updateObj) {
    var copyObj = JSON.parse(JSON.stringify(updateObj));
    if(copyObj.Items) {
        delete copyObj["Items"];
    }
    if(copyObj.Category) {
        copyObj['category_id'] = copyObj.Category.id;
        delete copyObj["Category"];
    }

    sendUpdate($('#update-recipe-form'), copyObj, sendAjax);
}