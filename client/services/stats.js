module.exports = function (app) {
	app.factory('StatsService', function ($resource) {
		return $resource('stats/day/2016-07-03', {}, {
			getData: {method:'GET', isArray: true}
		});
	});
};
