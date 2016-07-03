module.exports = function (app) {
	app.factory('Matchers2Tags', function ($http) {
		var url = '/me2tags';

		return {
			link: function (matcherId, tagId) {
				return $http({
					url: url,
					method: 'POST',
					data: {
						mId: matcherId,
						tagId: tagId
					}
				});
			},

			unlink: function (matcherId, tagId) {
				return $http({
					url: url,
					method: 'DELETE',
					data: {
						mId: matcherId,
						tagId: tagId
					}
				});
			}
		};
	});
};
