//## Here Maps example interactive map with markers.
/**  @param  {H.Map} map      A HERE Map instance within the application */

// Add a single marker to the map with the provided lat longs.
function addMarkerToMap(lati, long) {
    var newMarker = new H.map.Marker({ lat: lati, lng: long });
    map.addObject(newMarker);
}

// // Add Markers to the map based on lat-lon coordinates. 
// function addMarkersToMap(map) {
//     // Get lats and Longs from the API.
//     getLatLongs(42, -71, 200, 247, 1000000, 1000)
//         .then((latLongs) => {
//             console.log("LATLONGS _ _ _ ");
//             console.log(latLongs);
//             var lats = latLongs.lats;
//             var longs = latLongs.longs;

//             for (var i = 0; i < lats.length; i++) {
//                 var newMarker = new H.map.Marker({ lat: lats[i], lng: longs[i] });
//                 map.addObject(newMarker);
//             }
//         })
// }

// Initialize communication with Here Maps platform. 
// API Key declaration.
var platform = new H.service.Platform({
    'apikey': 'NZC4iABZusQ2Pr6JcnvLLSFn6vr8iytfJzet-EjuFpg'
});
var defaultLayers = platform.createDefaultLayers();

// Initialize a map.
// Centered on Philadelphia, PA, USA for the sake of seeing example markers.
var map = new H.Map(document.getElementById('map'),
    defaultLayers.vector.normal.map, {
    center: { lat: 42, lng: -71 },
    zoom: 10,
    pixelRatio: window.devicePixelRatio || 1
});

// Resize the map to occupy the entire div or element it is embedded in.
window.addEventListener('resize', () => map.getViewPort().resize());

// Make the map interactive. (Can move around it. Click on things, etc.)
// MapEvents enables the event system. 
// Behavior enables interactions for pan/zoom (mouse wheel, touch screen pinch, etc)
var behavior = new H.mapevents.Behavior(new H.mapevents.MapEvents(map))

// Create the zoom in/out buttons / bar and menu components.
var ui = H.ui.UI.createDefault(map, defaultLayers);

// After map is loaded, add the markers to the embedded instance.
// window.onload = function () {
//     addMarkersToMap(map);
// }

// Remove ALL markers from the map.
function removeAllMarkersFromMap() {
    map.removeObjects(map.getObjects())
}
