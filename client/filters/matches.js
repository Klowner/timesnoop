module.exports = function (app) {
	app.filter('matchesExpression', function () {
		return function (input, scope) {
			console.log(input, scope);
			return input;
		};
	});
};
