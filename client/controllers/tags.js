module.exports = function (app) {
	app.controller('TagsController', function ($scope, Tag) {

		$scope.tags = Tag.query();

		$scope.tag = {};

		$scope.addTag = function () {
			console.log('add tag', $scope.tag);

		};
	});
};
