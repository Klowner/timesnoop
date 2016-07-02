module.exports = function (app) {
	app.controller('TagsController', function ($scope, TagsService) {

		$scope.tags = TagsService.query();

		console.log(TagsService);
	});
};
