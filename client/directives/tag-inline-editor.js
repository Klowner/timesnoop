module.exports = function (app) {
	app.directive('tagInlineEditor', function () {
		return {
			restrict: 'EA',
			scope: {
				tag: '='
			},
			template: require('../templates/directives/tag-inline-editor.html'),
			replace: false,
			controller: function ($scope, $element, $timeout, Tag) {
				var saveTimeout;

				if (!($scope.tag instanceof Tag)) {
					$scope.tag = new Tag($scope.tag);
				}

				function save() {
					new $scope.tag.$save();
				}

				function debouncedSave(newValue, oldValue) {
					if (newValue != oldValue) {
						$timeout.cancel(saveTimeout);
						saveTimeout = $timeout(save, 400);
					}
				}

				$element.addClass('tag-inline-editor');

				$scope.removeTag = function (tag) {
					console.log('remove!', tag);


					tag.$remove().then(function () {
						console.log('removed tag');
					});

				};

				$scope.addSubTag = function (tag) {

					var child = new Tag({
						parent_id: tag.id,
						color: tag.color
					});

					child.$save().then(function () {
						tag.children = tag.children || [];
						tag.children.push(child);
					});

				};

				$scope.$watch('tag.name', debouncedSave);
				$scope.$watch('tag.parent_id', debouncedSave);
				$scope.$watch('tag.color', debouncedSave);
			}
		};
	});
};
