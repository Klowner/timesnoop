module.exports = function (app) {
	app.factory('Matcher', function ($resource) {
		return $resource('/matchers/:id');
	});
};
