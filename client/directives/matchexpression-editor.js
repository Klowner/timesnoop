var _ = require('lodash');

module.exports = function (app) {
	app.directive('matchExpressionEditor', function () {
		return {
			restrict: 'EA',
			scope: {
				matchers: '='
			},
			template: require('../templates/directives/matchexpression-editor.html'),
			replace: false,
			controller: function ($scope, $element, Tag) {


				$scope.matcher = {}; // for appending new records
				$scope.tags = Tag.query();

				$scope.deleteMatcher = function (matcher) {
					matcher.$delete().then(function () {

						var i = $scope.matchers.indexOf(matcher);
						$scope.matchers.splice(i, 1);
					});
				};

				$scope.addMatcher = function () {
					console.log($scope.matcher);
				};

				$scope.getTag = function (id) {
					return _.find($scope.tags, {id: id});
				};
			}
		};
	});
};
