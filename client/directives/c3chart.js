var c3 = require('c3');

// adapted from github.com/wasilak/angular-c3-simple

module.exports = function (app) {
	app
		.service('c3ChartService', function () {
			return {};
		})

		.directive('c3Chart', function (c3ChartService) {
			return {
				restrict: 'EA',
				scope: {
					config: '='
				},
				template: '<div></div>',
				replace: true,
				controller: function ($scope, $element) {
					// Wait until id is set before binding chart to this id
					$scope.$watch($element, function () {
						console.log('chart!');

						if ('' === $element[0].id) {
							return;
						}

						$scope.config.bindto = '#' + $element[0].id;

						$scope.$watch('config', function (newConfig, oldConfig) {
							console.log('newconfig', newConfig);

							if (!(newConfig.data.columns || newConfig.data.json || newConfig.data.rows)) {
								return;
							}

							c3ChartService[$scope.config.bindto] = c3.generate(newConfig);

							if (!newConfig.size) {
								c3ChartService[$scope.config.bindto].resize();
							}

							$scope.$watch('config.data', function (newData, oldData) {
								if ($scope.config.bindto) {
									c3ChartService[$scope.config.bindto].load(newData);
								}
							}, true);
						}, true);
					});
				}
			};
		});
};
