module.exports = function (app) {
	app.factory('TagsService', function ($resource) {
		return $resource('api/tags/:id');
	});
};
