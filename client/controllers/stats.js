var _ = require('lodash');

module.exports = function (app) {
	app.controller('StatsController', function ($scope, $stateParams, $state, StatsService, Stats, Tag) {

		/*
		$scope.chart = {
			data: {
				type: 'pie',
			},
			legend: {
				show: false
			}
		};
		*/


		$scope.tags = Tag.query();

		$scope.tagsChart = {
			data: {
				type: 'pie',
				onclick: function (d, i) {
					var tag = _.find($scope.tags, {name: d.name});

					if (tag) {
						$state.go('stats', {parentTagId: tag.id});
					}
				}
			},
			legend: { show: true }
		};

		Stats.getTagTotals($stateParams.parentTagId).then(function (result) {
			console.log(result);
			$scope.tagsChart.data.columns = _.map(result.data, function (record) {
				return [record.name, record.duration];
			});

			$scope.options = {
				chart: {
					type: 'sunburstChart',
					height: 450
				}
			};

			$scope.data = [{
				name: "all",
				children: [
					{name: "one", size: 200 },
					{name: "two",
						children: [
							{name: 'two sub', size: 200},
							{name: 'twosub2', size: 200}
						]
					}
				]
			}];
		});

		/*
		StatsService.getData(function (data) {
			$scope.collection = data;

			$scope.chart.data.columns = _.map(_.take(data, 10), function (record) {
				return [record.title, record.duration];
			});
			});
		*/
	});
};
