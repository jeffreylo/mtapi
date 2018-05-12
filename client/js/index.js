import { h, render } from "preact";
import queryString from "query-string";

import MTA from "./components/mta";

const latitude = parseFloat(localStorage.getItem("latitude"));
const longitude = parseFloat(localStorage.getItem("longitude"));

let coordinates;
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

renderApp(coordinates);

const getLocation = () => {
  const geolocation = navigator.geolocation;

  const location = new Promise((resolve, reject) => {
    if (!geolocation) {
      reject(new Error("navigator not supported"));
    }
    geolocation.getCurrentPosition(
      position => {
        resolve(position);
      },
      () => {
        reject(new Error("navigator permission denied"));
      }
    );
  });

  return location;
};

getLocation().then(location => {
  const coordinates = {
    Lat: location.coords.latitude,
    Lon: location.coords.longitude
  };
  localStorage.setItem("latitude", location.coords.latitude);
  localStorage.setItem("longitude", location.coords.longitude);
  renderApp(coordinates);
});
