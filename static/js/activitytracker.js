angular.module('activityTracker', ['ngRoute', 'ngResource'])

	.factory('StatsService', function($resource) {
		return $resource('stats/day/2016-07-01', {}, {
			getData: {method:'GET', isArray: true}
		});
	})

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

	.controller('StatsController', function ($scope, $route, $routeParams, $location, StatsService) {
		StatsService.getData(function (data) {
			$scope.collection = data;
		});
	})


	.directive('prime', function ($state) {
		return {
			restrict: 'EA',
			scope: {
				collection: '='
			},
			controller: function ($scope) {

			},
			link: function (scope, elem, attrs) {

			},
			templateUrl: 'tmpl/list.html'
		};
	})

	.config(function ($routeProvider, $locationProvider) {
		$routeProvider
			.when('/', {
				controller: 'MainController'
			})

			.when('/tags', {
				templateUrl: 'tmpl/tags.html',
				controller: 'TagsController'
			})

			.when('/statistics', {
				templateUrl: 'tmpl/list.html',
				controller: 'StatsController'
			})
		;

		$locationProvider.html5Mode(false);
	})

	.run(function () {
		console.log('starting app');
	});



