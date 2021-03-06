module.exports = function (app) {
	app.factory('StatsService', function ($resource) {
		return $resource('stats/day/2016-07-03', {}, {
			getData: {method:'GET', isArray: true},
		});
	});

	app.factory('Stats', function ($http) {
		var url = '/stats';

		return {
			getTagTotals: function (parentId) {
				return $http({
					url: url + '/tags' + (parentId ? '/' + parentId : ''),
					method: 'GET'
				});
			},

			getTagTotalsTree: function () {
				return $http({
					url: url + '/tags/tree',
					method: 'GET'
				});
			}


		};
	});
};
