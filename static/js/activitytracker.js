


angular.module('activityTracker', ['ngRoute'])

	.controller('MainController', function ($scope, $route, $routeParams, $location) {
		$scope.$route = $route;
		$scope.$location = $location;
		$scope.$routeParams = $routeParams;
	})

	.controller('TagsController', function ($scope, $route, $routeParams, $location) {
		$scope.$route = $route;
		$scope.$location = $location;
		$scope.$routeParams = $routeParams;
	})


	.config(function ($routeProvider, $locationProvider) {
		$routeProvider
			.when('/', {
				controller: 'MainController'
			})

			.when('/tags', {
				templateUrl: 'tmpl/tags.html',
				controller: 'TagsController'
			});

		$locationProvider.html5Mode(false);
	})

	.run(function () {
		console.log('starting app');
	});



