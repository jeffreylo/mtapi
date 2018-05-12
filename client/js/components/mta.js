import { h, Component } from "preact";
import { DateTime } from "luxon";

import rpc from "../rpc";
import css from "./mta.css";

// Times Square - 42 St.
const defaultLatLon = { Lat: 40.7589545, Lon: -73.9849801 };

class MTA extends Component {
  constructor() {
    super();
    this.state = { stations: [], now: DateTime.utc() };
  }

  refreshFeed(coordinates) {
    rpc.request(
      "GetClosest",
      coordinates || defaultLatLon,
      (err, error, response) => {
        if (response && response.Stations) {
          this.setState({ stations: response.Stations, now: DateTime.utc() });
        }
      }
    );
  }

  componentDidMount() {
    this.refreshFeed(this.props.coordinates);
    this.timer = setInterval(() => {
      this.refreshFeed(this.props.coordinates);
    }, 30000);
  }

  componentWillReceiveProps(nextProps) {
    this.refreshFeed(nextProps.coordinates);
  }

  componentWillUnmount() {
    clearInterval(this.timer);
  }

  renderArrival(header, trips = []) {
    const t = trips.filter(v => {
      const arrival = DateTime.fromISO(v.Arrival, { setZone: true });
      return (
        Math.round(arrival.diff(this.state.now, "minutes").toObject().minutes) >
        0
      );
    });

    let timings = t
      .map(v => {
        const arrival = DateTime.fromISO(v.Arrival, { setZone: true });
        const time = Math.round(
          arrival.diff(this.state.now, "minutes").toObject().minutes
        );
        return (
          <span>
            {v.RouteID}: {time} min<br />
          </span>
        );
      })
      .slice(0, 5) || [<span>-<br /></span>];
    if (timings.length < 5) {
      while (timings.length < 5) {
        timings.push(<br />);
      }
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
    const updated = Math.abs(DateTime.fromISO(station.Updated, { setZone: true }).diff(this.state.now, "seconds").toObject().seconds);

    return (
      <pre className={css.station}>
        <p>
          <strong>{station.Name}</strong>
        </p>
        {updated && <p>{updated}s ago</p>}
        {this.renderArrival("N", schedules.N)}
        {this.renderArrival("S", schedules.S)}
      </pre>
    );
  }

  render(props, state) {
    let stations = state.stations;
    return <div className={css.container}>{stations.map(v => this.renderStation(v))}</div>;
  }
}

export default MTA;
