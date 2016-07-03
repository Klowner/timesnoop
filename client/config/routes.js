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
			})

			.state('tags', {
				url: '/tags',
				template: require('../templates/tags.html'),
				controller: 'TagsController'
			})

			.state('matchers', {
				url: '/matchers',
				template: require('../templates/matchers.html'),
				controller: 'MatchersController'
			})
		;

		$urlRouterProvider.otherwise('/');

	});
};
