import humanizeDuration from "humanize-duration";

const humanizer = humanizeDuration.humanizer({
  language: "shortEn",
  round: true,
  units: ["m"],
  languages: {
    shortEn: {
      y: function() {
        return "y";
      },
      mo: function() {
        return "mo";
      },
      w: function() {
        return "w";
      },
      d: function() {
        return "d";
      },
      h: function() {
        return "h";
      },
      m: function() {
        return "m";
      },
      s: function() {
        return "s";
      },
      ms: function() {
        return "ms";
      }
    }
  }
});

export default humanizer;
