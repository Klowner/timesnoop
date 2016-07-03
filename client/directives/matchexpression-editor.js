module.exports = function (app) {

	app.directive('matchTagsWidget', function () {
		return {
			restrict: 'EA',
			scope: {
				matcher: '='
			},
			template: require('../templates/directives/matchtagswidget.html'),
			controller: function ($scope, $element, Tag, Matchers2Tags) {


				$scope.enableTagMenu = function () {
					$scope.availableTags = Tag.query();
				};

				$scope.addTag = function (tag) {
					Matchers2Tags
						.link($scope.matcher.id, tag.id)
						.then(function () {
							$scope.availableTags = [];
						});
				};
			}
		};
	});

	app.directive('matchExpressionEditor', function () {
		return {
			restrict: 'EA',
			scope: {
				matchers: '='
			},
			template: require('../templates/directives/matchexpression-editor.html'),
			replace: false,
			controller: function ($scope, $element) {


				$scope.matcher = {}; // for appending new records

				$scope.deleteMatcher = function (matcher) {
					matcher.$delete().then(function () {

						var i = $scope.matchers.indexOf(matcher);
						$scope.matchers.splice(i, 1);
					});
				};

				$scope.addMatcher = function () {
					console.log($scope.matcher);
				};
			}
		};
	});
};
