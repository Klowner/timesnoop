var _ = require('lodash');

module.exports = function (app) {
	app.controller('StatsController', function ($scope, StatsService) {

		$scope.chart = {
			data: {
				type: 'pie',
			},
			legend: {
				show: false
			}
		};

		StatsService.getData(function (data) {
			$scope.collection = data;

			$scope.chart.data.columns = _.map(_.take(data, 10), function (record) {
				return [record.title, record.duration];
			});
		});
	});
};
