module.exports = function (app) {
	app.factory('Tag', function ($resource) {
		return $resource('/tags/:id', {id: '@id'});
	});
};
