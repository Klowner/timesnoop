module.exports = function (app) {
	app.directive('tagtreeEditor', function () {
		return {
			restrict: 'EA',
			scope: {
				tags: '='
			},
			template: require('../templates/directives/tagtree-editor.html'),
			replace: false,
			controller: function ($scope, $element, Tag) {

				$scope.reloadTags = function () {
					$scope.tags = {children: Tag.tree()};
				};

				$scope.options = {
					dropped: function (event) {
						// Assign the new parent id
						var droppedTag = event.source.nodeScope.tag,
							newParentTag = event.dest.nodesScope.tag;

						droppedTag.parent_id = newParentTag.id || 0;

						droppedTag.$save().then(function () {
							// reload the tree
							$scope.reloadTags();
						});
					}
				};
			}
		};
	});
};
