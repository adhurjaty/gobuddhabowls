import { formatMoney, parseModelJSON } from "../helpers/_helpers";
import { CategorizedItemsDisplay } from "../components/_categorized_items_display";

var _options = {
    breakdown: false
};

var _columnInfo = [
    {
        name: 'id',
        hidden: true,
        get_column: (recipe) => {
            return recipe.id;
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
        header: 'Category',
        get_column: (recipe) => {
            return recipe.Category.name;
        },
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
    }
]

$(() => {
    createDatagrid()
});

function createDatagrid() {
    var container = $('#datagrid-holder');
    // var items = parseModelJSON(container.attr('data'));
    new CategorizedItemsDisplay(container, _columnInfo, null, _options);
}

function calculateRecipeCost(items, ruc) {
    return 0;
}