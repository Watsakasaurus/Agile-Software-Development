// Disable searching on table
$.extend(true, $.fn.dataTable.defaults, {
    "searching": false
});

// Init table on document load
$(document).ready(function () {

    // Get data from browser.
    let data = localStorage.getItem("data")

    // Set data to mock if localstorage is empty
    if (data == null) {
        data = mockData
    }

    // Turn json string into object
    data = JSON.parse(data)

    let columns = []
    if (data.total > 0) {
        columns.push(
            {
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
                title: "Price",
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
            removeAllMarkersFromMap();
            data.map(row => {
                //let marker = new H.map.Marker({ lat: row.latitude, lng: row.longitude });
                addDomMarker(row.latitude, row.longitude, row.average_total_payments, row.provider_name);
                //map.addObject(marker);
            });
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