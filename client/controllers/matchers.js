module.exports = function (app) {
	app.controller('MatchersController', function ($scope, Matcher, Tag) {

		$scope.matcher = {};
		$scope.matchers = Matcher.query();
		$scope.tags = Tag.query();

		$scope.$watch('matcher.expression', function (newValue, oldValue) {
			console.log(newValue);
		});

		$scope.addMatcher = function () {
			console.log($scope.matcher);
		};
	});


	app.controller('CategorizeController', function ($scope, $http, Matcher) {
		var expression;

		$scope.entries = [];
		$scope.matcher = {};
		$scope.valid = true;

		$scope.updateRecords = function () {
			return $http({
				method: 'GET',
				url: '/stats/unmatched'
			}).then(function (result) {
				$scope.entries = result.data;
			});

		};

		$scope.matchesMatcher = function (scope) {
			return !expression || expression.test(scope.entry.title);
		};

		$scope.$watch('matcher.expression', function (newValue, oldValue) {
			try {
				expression = new RegExp(newValue);
				$scope.valid = true;
			} catch (e) {
				$scope.valid = false;
				expression = null;
			}
		});

		$scope.addMatcher = function () {
			console.log('add matcher');
			var m = new Matcher($scope.matcher);
			m.$save()
				.then($scope.updateRecords);
		};

		$scope.updateRecords();

	});
};
