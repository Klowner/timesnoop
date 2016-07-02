module.exports = function (app) {
	app.factory('Tag', function ($resource) {
		return $resource('api/tags/:id');
	});
};
