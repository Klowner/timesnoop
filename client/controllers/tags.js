module.exports = function (app) {
	app.controller('TagsController', function ($scope, Tag) {

		$scope.tags = Tag.query();

		$scope.tag = {};

		$scope.addTag = function () {
			Tag.save($scope.tag);
			$scope.tags.push($scope.tag);
		};
	});
};
