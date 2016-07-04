var _ = require('lodash');

module.exports = function (app) {
	app.directive('matchTagsWidget', function () {
		return {
			restrict: 'EA',
			scope: {
				matcher: '=',
				tags: '='
			},
			transclude: true,
			template: require('../templates/directives/matchtagswidget.html'),
			controller: function ($scope, $element, Tag, Matchers2Tags) {


				$scope.enableTagMenu = function () {
					$scope.availableTags = _.filter($scope.tags, function (tag) {
						return $scope.matcher.tag_ids.indexOf(tag.id) === -1;
					});
				};

				$scope.addTag = function (tag) {
					Matchers2Tags
						.link($scope.matcher.id, tag.id)
						.then(function () {
							$scope.availableTags = [];
							$scope.matcher.tag_ids.push(tag.id);
						});
				};

				$scope.getTags = function () {
					var out = [];
					angular.forEach($scope.matcher.tag_ids, function (tid) {
						var match = _.find($scope.tags, {id: tid});
						if (match) {
							out.push(match);
						}
					});
					return out;
				};

				$scope.removeTag = function (tag) {
					Matchers2Tags
						.unlink($scope.matcher.id, tag.id)
						.then(function () {
							$scope.matcher.tag_ids.splice($scope.matcher.tag_ids.indexOf(tag.id), 1);
						});
				};
			}
		};
	});
};
