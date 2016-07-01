module.exports = function (app) {
	app.factory('StatsService', function ($resource) {
		return $resource('stats/day/2016-07-01', {}, {
			getData: {method:'GET', isArray: true}
		});
	});
};
