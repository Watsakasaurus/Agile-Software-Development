var url = "results.html?";

/**
 * Function that gets the users ive location from the browser
 * 
 */
function get_location() {
    return new Promise((resolve, reject) => {
        if (navigator.geolocation) {
            navigator.geolocation.getCurrentPosition(position => {
                resolve(position);
            });
        } else {
            console.log("Error occured");
        }
    }).then((position) => {
        lat_long_to_zip(position.coords.latitude, position.coords.longitude).then((zip) => {
            console.log(zip)
            $("#zip-code").val(zip);
            $("#zip-code").stop().css("background-color", "#FFFF9C")
                .animate({ backgroundColor: "#FFFFFF" }, 1500);
        })
    });
}

function zip_code_info(zip) {
    return new Promise((resolve, reject) => {
        $.ajax({
            url: 'https://geocoder.ls.hereapi.com/6.2/geocode.json',
            type: 'GET',
            cached: false,
            data: {
                PostalCode: zip,
                country: 'usa',
                gen: '9',
                apiKey: 'wS3zBaE7wLZ5Im9u7TfCFwbDmPABPKlSCCg_7s4JM-U'
            },
            success: function (data) {
                let test = JSON.stringify(data);
                console.log(test)

                // Show validation error if no response was returned
                if (data.Response.View.length == 0) {
                    let validator = $("#search-form").validate()
                    validator.showErrors({ "zip-code": "Zip Code not found" })
                    return
                }
                var lat = (data.Response.View[0].Result[0].Location.DisplayPosition.Latitude);
                var lon = (data.Response.View[0].Result[0].Location.DisplayPosition.Longitude);

                // Store the results in the local storage
                localStorage.setItem("lat", lat);
                localStorage.setItem("lon", lon);

                resolve({ lat, lon })
            }
        });
    });
}

function lat_long_to_zip(lat, long) {
    return new Promise((resolve, reject) => {
        $.ajax({
            url: 'https://reverse.geocoder.ls.hereapi.com/6.2/reversegeocode.json',
            type: 'GET',
            cached: false,
            data: {
                bbox: lat + ',' + long + ';' + lat + ',' + long,
                mode: 'retrieveAddresses',
                country: 'usa',
                gen: '9',
                apiKey: 'wS3zBaE7wLZ5Im9u7TfCFwbDmPABPKlSCCg_7s4JM-U'
            },
            success: function (data) {
                if (data.Response.View.length == 0) {
                    alert("Failed to get Zip Code")
                    return
                }
                var zip = parseInt(data.Response.View[0].Result[0].Location.Address.PostalCode);

                resolve(zip)
            }
        });
    });
}

function get_data(query, lat, long, minPrice, maxPrice, proximity) {
    return new Promise((resolve, reject) => {
        $.ajax({
            url: 'http://99.81.88.54:5009/medicare/api/data',
            type: 'GET',
            cached: false,
            data: {
                query: query,
                lat: lat,
                long: long,
                min_price: minPrice,
                max_price: maxPrice,
                proximity: proximity,
                per_page: 2000,
            },
            success: function (data) {
                console.log(data)

                resolve(data)
            }
        });
    });
}

$(document).ready(function () {
    $("#search-form").on('submit', async function (e) {
        e.preventDefault();


        console.log($("#amount-top").val())
        let procedure = $("#injury-type").val() != "" ? $("#injury-type").val() : null;
        let minPrice = $("#amount-bottom").val() != "" ? parseInt($("#amount-bottom").val()) : null;
        let maxPrice = $("#amount-top").val() != "" ? parseInt($("#amount-top").val()) : null;
        let proximity = $("#distance-amount-top").val() != "" ? parseInt($("#distance-amount-top").val()) : null;
        let zipCode = $("#zip-code").val()

        var latLong = {}

        // Zip code should never be empty, because we get either live or manual
        if (zipCode != "") {
            latLong = await zip_code_info(zipCode)
        }
        
        let data = await get_data(procedure, latLong.lat, latLong.lon, minPrice, maxPrice, proximity)
        
        // Store the results in the local storage
        localStorage.setItem("data", JSON.stringify(data))
        localStorage.setItem("procedure", procedure)
        localStorage.setItem("minPrice", minPrice)
        localStorage.setItem("maxPrice", maxPrice)
        localStorage.setItem("proximity", proximity)
        localStorage.setItem("zipCode", zipCode)

        $(location).attr('href', 'results.html#map')
    });
})