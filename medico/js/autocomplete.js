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
                min: response.price_min,
                max: response.price_max,
                values: [ response.price_min, response.price_max ],
                slide: function( event, ui ) {
                    for (var i=0; i<ui.values.length; i++)
                    {
                        $("input.amount[data-index=" + i + "]").val(ui.values[i]);
                    }
                    //$( "#amount-bottom" ).val(ui.values[ 0 ]);
                    //$( "#amount-top" ).val(ui.values[ 1 ]);
                }
            });
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
            

            var dsr = $( "#distance-slider-range" ).slider({
                value: 200,
                max: 20000,
                range: "max",
                animate: true,
                slide: function( event, ui ) {
                    $( "#distance-amount-top" ).val(ui.value);
                }
            });

            $("input.amount").change(function() {
                var $this = $(this);
                
                console.log("1" + " " + $this);
                console.log("2 " + $this.data("index"));
                console.log("3 " + $this.val());
                console.log("4 " + $("#amount-bottom").val());
                console.log("5 " + $("#amount-top").val());
                
                if (($this.data("index") == 0 && ($this.val() < $("#amount-bottom").val())) || ($this.data("index") == 1 && ($this.val() > $("#amount-top").val())))
                {
                    $("#slider-range").slider("values", $this.data("index"), $this.val());
                }
                
            });

            $( "#distance-amount-top" ).on( "keyup", function() {
                console.log(this.value);
                dsr.slider( "value", this.value );
            });
            
            //$( "#amount-top" ).val($( "#slider-range" ).slider( "values", 0 ));
            //$( "#amount-bottom" ).val($( "#slider-range" ).slider( "values", 1 ));
            //$( "#distance-amount-top" ).val($( "#distance-slider-range" ).slider( "value" ));
        }
    });
});
