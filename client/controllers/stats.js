var _ = require('lodash'),
	d3 = require('d3');

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

		Stats.getTagTotalsTree().then(function (result) {
			console.log('resultdata', result.data);

			//$scope.options = {
				//chart: {
					//type: 'sunburstChart',
					//height: 700,
					//color: d3.scale.category20c(),
					//duration: 250,

					//sunburst: {
						//mode: 'size',
						//key: function (d) { return d.name; }
					//}
				//},
			//};

			//$scope.data = [{
				//name: 'all',
				//children: result.data
			//}];

			$scope.sunoptions = {
				width: 200,
				height: 200
			};

			$scope.sundata = {
				name: 'all',
				duration: 0,
				children: result.data
			};

			console.log($scope.sundata);
		});

		Stats.getTagTotals($stateParams.parentTagId).then(function (result) {
			console.log(result);
			$scope.tagsChart.data.columns = _.map(result.data, function (record) {
				return [record.name, record.duration];
			});

			//$scope.options = {
				//chart: {
					//type: 'sunburstChart',
					//height: 700,
					//color: d3.scale.category20c(),
					//duration: 250,

					//sunburst: {
						//mode: 'duration',
						//key: function (d) { return d.name; }
					//}
				//},

			//};

			//$scope.data = [
				//{
					//name: 'root',
					//children: [
						//{
							//name: 'a',
							//children: [
								//{ name: "one", size: 10, duration: 20 },
								//{ name: "two", size: 20, duration: 40},
							//]
						//},

						//{
							//name: 'b',
							//children: [
								//{ name: "b/one", size: 10, duration: 2},
								//{ name: "b/two", size: 20, duration: 2},
							//]
						//},
					//]
				//}
			//];

				//name: "flare",
				//children: [
					//{
						//name: "one",
						//size: 128
					//},
					//{
						//name: "two",
						//size: 256
					//}
				//]
			//}];

		});
	});
};
