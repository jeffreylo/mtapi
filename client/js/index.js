import { h, render } from "preact";
import queryString from "query-string";

import MTA from "./components/mta";

const parsed = queryString.parse(location.search);
let coordinates;
if (parsed.lat && parsed.lon) {
  coordinates = {Lat: parseFloat(parsed.lat), Lon: parseFloat(parsed.lon)}
}
render(<MTA coordinates={coordinates} />, document.getElementById("root"));
