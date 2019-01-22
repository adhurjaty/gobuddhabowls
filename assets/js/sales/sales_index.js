import { transactionURL } from "../helpers/_square_helper";
import { showError } from "../helpers/index_helpers";

var _sales = []

$(() => {
    getSquareData();
});


function getSquareData() {
    var infoContainer = $('#square-info');
    var token = infoContainer.attr('square-token');
    var locationID = infoContainer.attr('square-location');

    var dateStrs = getDateRange();
    var startTime = new Date(dateStrs[0]);
    var endTime = new Date(dateStrs[1]);
    var url = transactionURL(locationID, startTime, endTime);

    getTransactionAjax(url, token);
}

function getDateRange() {

}

function getTransactionAjax(url, token) {
    $.ajax(url, {
        method: 'GET',
        headers: {Authorization: `Bearer: ${token}`},
        error: (xhr, status, err) => {
            showError(err);
        },
        success: (data, status, xhr) => {
            populateDatagrid(data);
        },
    });
}

function populateDatagrid(data) {
    
}