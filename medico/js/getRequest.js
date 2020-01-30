

function getLatLongs(latitude, longitude, prox, que, priceMax, priceMin) {
    console.log("Executing: API Request:")
    console.log("API Address: http://99.81.88.54:5009/medicare/api/data");
    console.log("type: GET")
    console.log("dataType: json")

    return new Promise((resolve, reject) => {

        $.ajax({
            url: 'http://99.81.88.54:5009/medicare/api/data',
            type: 'GET',
            dataType: 'json',
            data: {
                lat: latitude,
                long: longitude,
                proximity: prox,
                query: que,
                page: 1,
                per_page: 2000,
                max_price: priceMax,
                min_price: priceMin
            },

            success: function (data) {
                JSON.stringify(data);
                console.log("Success: Data returned.")
                console.log(data)

                var lats = new Array;
                var longs = new Array;

                for (var i = 0; i < data.objects.length; i++) {
                    lats[i] = data.objects[i].latitude;
                    longs[i] = data.objects[i].longitude;
                }

                console.log("Map_Data_Lat_Longs:");
                console.log(lats, longs);

                resolve({ lats, longs });
            }
        });
    });

}