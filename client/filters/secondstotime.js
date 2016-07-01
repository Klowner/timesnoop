module.exports = function (app) {
	app.filter('secondsToTime', function () {
		return function (_seconds) {
			var hours = Math.floor(_seconds / 3600),
				minutes = Math.floor((_seconds % 3600) / 60),
				seconds = Math.floor(_seconds % 60),
				out = [];

			if (hours) {
				out.push(hours + ' hours,');
			}
			if (minutes) {
				out.push(minutes + ' minutes,');
			}
			if (seconds) {
				out.push(seconds + ' seconds');
			}
			return out.join(' ');
		};
	});
};
