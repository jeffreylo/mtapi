import jayson from "jayson/lib/client/browser";
import fetch from "node-fetch";

const callServer = function(request, callback) {
  const options = {
    method: "POST",
    body: request, // request is a string
    headers: {
      "Content-Type": "application/json"
    }
  };
  const location = window.location;
  fetch(`${location.protocol}//${location.host}/rpc`, options)
    .then(function(res) {
      return res.text();
    })
    .then(function(text) {
      callback(null, text);
    })
    .catch(function(err) {
      callback(err);
    });
};

const rpc = jayson(callServer);

export default rpc;

export const GetClosestStations = coordinates => {
  if (isEmptyObject(coordinates)) {
    throw new Error("coordinates cannot be empty");
  }
  return new Promise((resolve, reject) => {
    rpc.request("GetClosestStations", {NumStations: 5, ...coordinates}, (err, error, result) => {
      if (error) reject(error);
      if (result && result.Stations) {
        resolve(result.Stations);
      }
    });
  });
};

const isEmptyObject = obj =>
  !!obj && Object.keys(obj).length === 0 && obj.constructor === Object;
