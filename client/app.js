var angular = require('angular'),
	nv = require('nvd3');

module.exports = angular.module('activityTracker', [
	require('angular-resource'),
	require('angular-ui-router'),
	require('angular-nvd3'),
	require('angular-ui-tree'),
]);
