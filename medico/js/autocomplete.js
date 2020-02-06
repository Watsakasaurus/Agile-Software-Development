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
                step: 250,
                min: 0,
                max: response.price_max,
                values: [ 0, response.price_max ],
                slide: function( event, ui ) {
                    for (var i=0; i<ui.values.length; i++)
                    {
                        $("input.amount[data-index=" + i + "]").val(ui.values[i]);
                    }
                }
            });

            var dsr = $( "#distance-slider-range" ).slider({
                value: 50,
                max: 10000,
                step: 10,
                range: "max",
                animate: true,
                slide: function( event, ui ) {
                    $( "#distance-amount-top" ).val(ui.value);
                }
            });

            $("input.amount").change(function() {
                var $this = $(this);
                
                if (($this.data("index") == 0 && ($this.val() < $("#amount-bottom").val())) || ($this.data("index") == 1 && ($this.val() > $("#amount-top").val())))
                {
                    $("#slider-range").slider("values", $this.data("index"), $this.val());
                }
                else if (($this.data("index") == 0 && ($this.val() > $("#amount-bottom").val()))) 
                {
                    $("#slider-range").slider("values", $this.data("index"), $("#amount-bottom").val());
                }
                else if (($this.data("index") == 1 && ($this.val() < $("#amount-top").val()))) 
                {
                    $("#slider-range").slider("values", $this.data("index"), $("#amount-top").val());
                }
            });

            $( "#distance-amount-top" ).on( "keyup", function() {
                console.log(this.value);
                dsr.slider( "value", this.value );
            });

            let procedure = localStorage.getItem("procedure")
            let minPrice = parseInt(localStorage.getItem("minPrice"))
            let maxPrice = parseInt(localStorage.getItem("maxPrice"))
            let proximity = parseInt(localStorage.getItem("proximity"))
            let zipCode = parseInt(localStorage.getItem("zipCode"))

            isNaN(zipCode) == false ? $("#zip-code").val(zipCode) : null
            procedure != null ? $("#injury-type").val(procedure) : null
            isNaN(minPrice) == false ? $("#amount-bottom").val(minPrice) : $( "#amount-bottom" ).val($( "#slider-range" ).slider( "values", 0 ));
            isNaN(maxPrice) == false ? $("#amount-top").val(maxPrice) : $( "#amount-top" ).val($( "#slider-range" ).slider( "values", 1 ));
            isNaN(proximity) == false ? $("#distance-amount-top").val(proximity) : $( "#distance-amount-top" ).val($( "#distance-slider-range" ).slider( "value" ));


            $("#slider-range").slider("values", 0, $("#amount-bottom").val());
            $("#slider-range").slider("values", 1, $("#amount-top").val());
            $("#distance-slider-range").slider("value", $("#distance-amount-top").val());
        }
    });
});
