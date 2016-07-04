var _ = require('lodash');

module.exports = function (app) {
	app.controller('StatsController', function ($scope, StatsService, Stats) {

		$scope.chart = {
			data: {
				type: 'pie',
			},
			legend: {
				show: false
			}
		};

		$scope.tagsChart = {
			data: { type: 'pie' },
			legend: { show: false }
		};

		Stats.getTagTotals().then(function (result) {
			console.log(result);
			$scope.tagsChart.data.columns = _.map(result.data, function (record) {
				console.log('record', record);
				return [record.name, record.duration];
			});
			console.log($scope.tagsChart.data.columns);
		});

		StatsService.getData(function (data) {
			$scope.collection = data;

			$scope.chart.data.columns = _.map(_.take(data, 10), function (record) {
				return [record.title, record.duration];
			});

		});
	});
};
