$(document).ready(function() {
    $.ajax({
        url: "http://99.81.88.54:5009/medicare/api/filtering",
        type: "GET",
        success: function(response) {
        console.log(response) 
            $("#injury-type").autocomplete({
                source: response.procedure_definitions
            });

            $( "#slider-range" ).slider({
                range: true,
                min: 0,
                max: response.price_max,
                values: [ 0, response.price_max ],
                slide: function( event, ui ) {
                    $( "#amount-top" ).val(ui.values[ 0 ]);
                    $( "#amount-bottom" ).val(ui.values[ 1 ]);
                }
            // }).slider("pips", {
            //     first: "label",
            //     last: "label",
            //     rest: "label",
            //     step: 50000,
            //     labels: true,
            //     prefix: "$",
            //     suffix: ""
            // }).slider("float", {
            //     handle: true,
            //     pips: true,
            //     labels: true,
            //     prefix: "",
            //     suffix: ""
            // });
            });

            $( "#distance-slider-range" ).slider({
                value: 200,
                max: 20000,
                range: "max",
                animate: true,
                slide: function( event, ui ) {
                    $( "#distance-amount-top" ).val(ui.value + " " + "Miles");
                }
            });
            
            $( "#amount-top" ).val($( "#slider-range" ).slider( "values", 0 ));
            $( "#amount-bottom" ).val($( "#slider-range" ).slider( "values", 1 ));
            $( "#distance-amount-top" ).val($( "#distance-slider-range" ).slider( "value" ));
        }
    });
});
