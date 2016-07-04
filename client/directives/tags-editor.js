module.exports = function (app) {
	app.directive('tagsEditor', function () {
		return {
			restrict: 'EA',
			scope: {
				tags: '='
			},
			template: require('../templates/directives/tags-editor.html'),
			replace: false,
			controller: function ($scope, $element, Tag) {

				$scope.tag = {};

				$scope.deleteTag = function (tag) {
					tag.$delete().then(function () {
						var i = $scope.tags.indexOf(tag);
						$scope.tags.splice(i, 1);
					});
				};

				$scope.addTag = function () {
					console.log($scope.tag);

					var t = new Tag($scope.tag);
					if (t.parent) {
						t.parent_id = t.parent.id;
						delete t.parent;
					}

					t.$save().then(function () {
						$scope.tags.push(t);
						$scope.tag = {};
					});
				};
			}
		};
	});
};
