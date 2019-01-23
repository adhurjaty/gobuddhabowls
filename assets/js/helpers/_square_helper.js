import { rfc3339 } from "./_helpers";

export function transactionURL(locationID, startTime, endTime) {
    var formattedStart = rfc3339(startTime);
    var formattedEnd = rfc3339(endTime);
    return `https://connect.squareup.com/v2/locations/${locationID}`
        + `/transactions?begin_time=${formattedStart}&end_time=${formattedEnd}`;
}
