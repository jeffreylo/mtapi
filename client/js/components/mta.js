import { h, Component } from "preact";
import { DateTime } from "luxon";

import rpc from "../rpc";

// Times Square - 42 St.
const defaultLatLon = { Lat: 40.7589545, Lon: -73.9849801 };

class MTA extends Component {
  constructor() {
    super();
    this.state = { stations: [], now: DateTime.utc() };
  }

  refreshFeed() {
    rpc.request("GetClosest", (this.props || {}).coordinates || defaultLatLon, (err, error, response) => {
      if (response && response.Stations) {
        this.setState({ stations: response.Stations, now: DateTime.utc() });
      }
    });
  }

  componentDidMount() {
    this.refreshFeed();
    this.timer = setInterval(this.refreshFeed, 30000);
  }

  componentWillUnmount() {
    clearInterval(this.timer);
  }

  renderArrival(trips) {
    if (!trips) return;

    const t = trips.filter((v) => {
      const arrival = DateTime.fromISO(v.Arrival, { setZone: true })
      return Math.round(arrival.diff(this.state.now, 'minutes').toObject().minutes) > 0
    });

    return t.map((v) => {
      const arrival = DateTime.fromISO(v.Arrival, { setZone: true })
      const time = Math.round(arrival.diff(this.state.now, 'minutes').toObject().minutes);
      return <span>{v.RouteID}: {time} min<br /></span>
    }).slice(0, 5);
  }

  renderStation(station) {
    if (!station) return;

    let schedules = (station.Schedules || {});
    return (
      <div>
        <h3>{station.Name}</h3>
        <h4>N</h4>
        {(schedules.N || []).length > 0 &&
         this.renderArrival(schedules.N)}
        <h4>S</h4>
        {(schedules.S || []).length > 0 &&
         this.renderArrival(schedules.S)}
      </div>
    );
  }

  render(props, state) {
    let stations = state.stations;
    return <div>{stations.map(v => this.renderStation(v))}</div>;
  }
}

export default MTA;
