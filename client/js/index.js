import { h, render } from "preact";
import queryString from "query-string";

import { getLocation } from "./location";
import MTA from "./components/mta";

// Acquire coordinates from query-string or location storage.
let coordinates;
const latitude = parseFloat(localStorage.getItem("latitude"));
const longitude = parseFloat(localStorage.getItem("longitude"));

if (location.search) {
  const parsed = queryString.parse(location.search);
  if (parsed.lat && parsed.lon) {
    coordinates = { Lat: parseFloat(parsed.lat), Lon: parseFloat(parsed.lon) };
  }
} else if (latitude && longitude) {
  coordinates = { Lat: latitude, Lon: longitude };
}

const renderApp = coordinates => {
  const rootNode = document.getElementById("root");
  render(<MTA coordinates={coordinates} />, rootNode, rootNode.lastChild);
};

// Render app with existing or specified knowledge.
renderApp(coordinates);

// If we don't have navigator permissions or a specified lat/lon, prompt for
// permissions.
if (!location.search) {
  getLocation().then(location => {
    const coordinates = {
      Lat: location.coords.latitude,
      Lon: location.coords.longitude
    };
    localStorage.setItem("latitude", location.coords.latitude);
    localStorage.setItem("longitude", location.coords.longitude);
    renderApp(coordinates);
  });
}
