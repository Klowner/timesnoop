module.exports = function (app) {
	app.factory('StatsService', function ($resource) {
		return $resource('stats/day/2016-07-03', {}, {
			getData: {method:'GET', isArray: true},
		});
	});

	app.factory('Stats', function ($http) {
		var url = '/stats';

		return {
			getTagTotals: function () {
				return $http({
					url: url + '/tags',
					method: 'GET'
				});
			}
		};
	});
};
