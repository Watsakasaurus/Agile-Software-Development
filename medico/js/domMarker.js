// Add DOM Marker to the map, showing the prices of the procedures.

// Add DOM Marker to the map.
function addDomMarker(lati, lon, price, name) {
    var outerElement = document.createElement('div'),
        innerElement = document.createElement('div');

    // Style the DOM Marker
    outerElement.style.userSelect = 'none';
    outerElement.style.webkitUserSelect = 'none';
    outerElement.style.msUserSelect = 'none';
    outerElement.style.mozUserSelect = 'none';
    outerElement.style.cursor = 'default';

    // More Style
    innerElement.style.color = 'black';
    innerElement.style.backgroundColor = 'white';
    innerElement.style.border = '2px solid black';
    innerElement.style.font = 'normal 12px arial';
    innerElement.style.lineHeight = '12px'
    innerElement.style.borderRadius = '15%';

    // More Style
    innerElement.style.paddingTop = '8px';
    innerElement.style.paddingBottom = '8px';
    innerElement.style.paddingLeft = '8px';
    innerElement.style.paddingRight = '8px';
    innerElement.style.width = '100%';
    innerElement.style.height = '100%';

    // Negative Margin for the purpose of centering the text.
    innerElement.style.marginTop = '-10px';
    innerElement.style.marginBottom = '-10px';

    // Join the two elements together.
    outerElement.appendChild(innerElement);

    // Add text to the DOM element
    innerElement.innerHTML = name + "<br>$" + price;

    // Change opacity on mouseover.
    function changeOpacity(evt) {
        evt.target.style.opacity = 0.6;
    };

    // Change opacity to normal when mouse leaves.
    function changeOpacityToOne(evt) {
        evt.target.style.opacity = 1;
    };

    // Create DOM Icon and add / remove the opacity listeners.
    var domIcon = new H.map.DomIcon(outerElement, {
        // Everytime Marker enters view, run the function,
        onAttach: function (clonedElement, domIcon, domMarker) {
            clonedElement.addEventListener('mouseover', changeOpacity);
            clonedElement.addEventListener('mouseout', changeOpacityToOne);
        },
        // When marker leaves view, run the function.
        onDetach: function (clonedElement, domIcon, domMarker) {
            clonedElement.removeEventListener('mouseover', changeOpacity);
            clonedElement.removeEventListener('mouseout', changeOpacityToOne);
        }
    });

    // Set Marker location. 
    var DOM_Marker = new H.map.DomMarker({ lat: lati, lng: lon }, {
        icon: domIcon
    });
    map.addObject(DOM_Marker);
}