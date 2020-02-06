// Disable searching on table
$.extend(true, $.fn.dataTable.defaults, {
    "searching": false
});

// Init table on document load
$(document).ready(function () {

    // Get data from browser.
    let data = localStorage.getItem("data")

    // Turn json string into object
    data = JSON.parse(data)

    // If there is no data, show 404
    if (data == null || data.total == 0) {
        $("#main-content").replaceWith(`
            <div class="jumbotron jumbotron-fluid">
                <div class="container justify-content-ceter">
                    <h1 class="display-4 d-flex justify-content-center">Not Found</h1>
                    <p class="lead d-flex justify-content-center">Please try again with different search parameters.</p>
                </div>
            </div>
        `)
        $("#custom-spinner").attr('style', 'display: none !important')
        return
    }

    let columns = []
    if (data.total > 0) {
        columns.push(
            {
                // Build and populate table.
                class_name: "full-span-col",
                data: null,
                title: "Name",
                render: function (data, type, row, meta) {
                    return `
                        <div class="card">
                            <div class="card-body" style="padding: 0.75rem;">
                                <h5 class="card-title">${row.provider_name}</h5>
                                <h6 class="card-subtitle mb-2 text-muted">${row.hrr_description}</h6>
                                <div class="row">
                                    <div class="col">
                                        <i class="fas fa-road"></i> Street - ${row.provider_street}</i>
                                    </div>
                                    <div class="col">
                                        <i class="fas fa-city"></i> City - ${row.provider_city}</i>
                                    </div>
                                    <div class="w-100"></div>
                                    <div class="col">
                                        <i class="fas fa-location-arrow"></i> Zip Code - ${row.provider_zip_code}</i>
                                    </div>
                                    <div class="col">
                                        <i class="fas fa-dollar-sign"> Price - $${row.average_total_payments}</i>
                                    </div>
                                </div>
                            </div>
                        </div>
                    `
                }
            }
        )

        columns.push(
            {
                data: "average_total_payments",
                title: '<span style="background-color: white;">Price</span>',
            }
        )

        // Check if distance is present
        if ("distance" in data.objects[0]) {
            columns.push(
                {
                    data: "distance",
                    title: "Distance",
                }
            )
        }
    }
    $('#resultsTable').DataTable({
        "scrollY": "80vh",
        "scrollCollaps": true,
        info: false,
        select: true,
        paging: true,
        data: data.objects,
        processing: true,
        columns: columns,
        createdRow: function (row, data, dataIndex) {
            // Add COLSPAN attribute
            $('td:eq(0)', row).attr('colspan', 3);

            // Hide required number of columns
            // next to the cell with COLSPAN attribute
            $('td:eq(1)', row).css('display', 'none');
            $('td:eq(2)', row).css('display', 'none');
        },
        "initComplete": function (settings, json) {
            $("#custom-spinner").attr('style', 'display: none !important')
            $("#main-content").css('visibility', 'visible')
        },
        "fnDrawCallback": function (oSettings) {
            let table = $("#resultsTable").DataTable()
            let data = table.rows({ page: "current" }).data();

            // Remove all the markers from the list that are already there.
            removeAllMarkersFromMap();
            // Add marker to map to show where user is.
            addMarkerToMap( localStorage.getItem("lat"), localStorage.getItem("lon")  );
            // Add new items to the map.
            data.map(row => {addDomMarker(row.latitude, row.longitude, row.average_total_payments, row.provider_name);});

            // Center map on the location of the first item in the list.
            map.setCenter({ lat: data[0].latitude, lng: data[0].longitude });
            map.setZoom(14);
        }
    });

    $('#resultsTable tbody').on('click', 'tr', function () {
        // Create data from the table based on clicked element.
        let clickData = $('#resultsTable').DataTable().row(this).data();

        // Center the map when user clicks an entry in the list
        map.setCenter({ lat: clickData.latitude, lng: clickData.longitude });
        map.setZoom(14);
    })
});