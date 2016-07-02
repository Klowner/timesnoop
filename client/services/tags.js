module.exports = function (app) {
	app.factory('TagsService', function ($resource) {
		function all() {
			return $resource('api/tags', {}, {
				getData: {
					method: 'GET',
					isArray: true
				}
			});
		}

		return {
			all: all
		};
	});
};
