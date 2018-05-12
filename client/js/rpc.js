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

export default jayson(callServer);
