import { h, Component } from "preact";
import { DateTime } from "luxon";

import humanizer from "../duration";
import rpc from "../rpc";
import css from "./mta.css";

// Times Square - 42 St.
const defaultLatLon = { Lat: 40.7589545, Lon: -73.9849801 };

class MTA extends Component {
  constructor() {
    super();
    this.state = { stations: [] };
  }

  refreshFeed(coordinates) {
    rpc.request(
      "GetClosest",
      coordinates || defaultLatLon,
      (err, error, response) => {
        if (response && response.Stations) {
          this.setState({ stations: response.Stations });
        }
      }
    );
  }

  componentDidMount() {
    this.refreshFeed(this.props.coordinates);
    this.feedInterval = setInterval(() => {
      this.refreshFeed(this.props.coordinates);
    }, 30000);
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.coordinates.Lat != this.props.coordinates.Lat) {
      this.refreshFeed(nextProps.coordinates);
    }
  }

  componentWillUnmount() {
    clearInterval(this.feedInterval);
  }

  renderArrival(header, trips = []) {
    let timings = trips
      .map(v => {
        const arrival = DateTime.fromISO(v.Arrival, { setZone: true });
        return (
          <span>
            {v.RouteID}:{" "}
            {humanizer(this.props.now.diff(arrival).toObject().milliseconds)}{" "}
            <br />
          </span>
        );
      })
      .slice(0, 10) || [
      <span>
        -<br />
      </span>
    ];

    while (timings.length < 5) {
      timings.push(<br />);
    }
    return (
      <div>
        <h4>{header}</h4>
        {(timings.length && timings) || "-"}
      </div>
    );
  }

  renderStation(station) {
    if (!station) return;
    let { Schedules } = station;
    const schedules = Schedules || {};
    const updated = Math.round(
      DateTime.fromISO(station.Updated, { setZone: true })
        .diff(this.props.now, "minutes")
        .toObject().minutes
    );

    return (
      <pre className={css.station}>
        <p>
          <strong>{station.Name}</strong>
        </p>
        {this.renderArrival("Uptown / Manhattan", schedules.N)}
        {this.renderArrival("Downtown / Brooklyn", schedules.S)}
        <p className={css.updated}>
          <small>
            {(updated && <span>~{Math.abs(updated)} min ago</span>) || (
              <span>recently</span>
            )}
          </small>
        </p>
      </pre>
    );
  }

  render(props, state) {
    let stations = state.stations;
    return (
      <div>
        <pre className={css.station}>{props.now.toLocaleString(DateTime.DATETIME_MED_WITH_SECONDS)}</pre>
        <div className={css.container}>
          {stations.map(v => this.renderStation(v))}
        </div>
      </div>
    );
  }
}

export default MTA;
