module.exports = function (app) {
	app.factory('Tag', function ($resource) {
		return $resource('/tags/:id', {id: '@id'}, {
			tree: {
				method: 'GET',
				url: '/tags/tree',
				isArray: true
			}
		});
	});
};
