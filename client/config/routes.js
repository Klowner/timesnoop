var mainTemplate = require('!ngtemplate!html!../templates/main.html');


module.exports = function (app) {
	app.config(function ($stateProvider, $urlRouterProvider) {
		$stateProvider
			.state('main', {
				url: '/',
				template: require('../templates/main.html'),
				controller: 'MainController'
			})

			.state('stats', {
				url: '/stats',
				template: require('../templates/stats.html'),
				controller: 'StatsController'
			});

		$urlRouterProvider.otherwise('/');

	});
};