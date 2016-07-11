module.exports = function (app) {
	app.controller('TagsController', function ($scope, Tag) {

		$scope.tags = Tag.query();

		$scope.tag = {};

		$scope.color = '#ff00ff';

		$scope.addTag = function () {
			Tag.save($scope.tag);
			$scope.tags.push($scope.tag);
		};

		$scope.tagtree = {children: Tag.tree()};

		$scope.treeOptions = {
			dropped: function (event) {
				// Assign the new parent id
				var droppedTag = event.source.nodeScope.tag,
					newParentTag = event.dest.nodesScope.tag;

				droppedTag.parent_id = newParentTag.id || 0;

				droppedTag.$save().then(function () {
					// reload the tree
					$scope.tagtree = {children: Tag.tree()};
				});
			}
		};
	});
};
