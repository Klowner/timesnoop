module.exports = function (app) {
	app.controller('MatchersController', function ($scope, Matcher) {

		$scope.matcher = {};

		$scope.$watch('matcher.expression', function (newValue, oldValue) {
			console.log(newValue);
		});

		$scope.addMatcher = function () {
			console.log($scope.matcher);
		};
	});
};
