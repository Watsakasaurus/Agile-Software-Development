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
                values: [ response.price_min, response.price_max ],
                slide: function( event, ui ) {
                    for (var i=0; i<ui.values.length; i++)
                    {
                        $("input.amount[data-index=" + i + "]").val(ui.values[i]);
                    }
                }
            });

            var dsr = $( "#distance-slider-range" ).slider({
                value: 200,
                max: 20000,
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
            
            $( "#amount-bottom" ).val($( "#slider-range" ).slider( "values", 0 ));
            $( "#amount-top" ).val($( "#slider-range" ).slider( "values", 1 ));
            $( "#distance-amount-top" ).val($( "#distance-slider-range" ).slider( "value" ));
        }
    });
});
