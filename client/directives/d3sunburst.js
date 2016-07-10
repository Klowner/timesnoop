var d3 = require('d3'),
	d3tip = require('d3-tip');

module.exports = function (app) {
	app.directive('d3Sunburst', function () {
		return {
			restrict: 'EA',
			scope: {
				options: '=',
				data: '='
			},
			link: function (scope, element, attrs) {
				var width = 960,
					height = 700,
					radius = Math.min(width, height) / 2 - 10;

				var vis = d3.select(element[0]).append('svg')
					.attr('width', width)
					.attr('height', height)
					.append('g')
					.attr('transform', 'translate(' + width / 2 + ',' + (height/2+10) + ')'),

				colors = d3.scale.category20c(),
				color = function (d) {
					console.log(d.color);
					return d.color || (d.parent && d.parent.color) || '#ffffff';
				},


				partition = d3.layout.partition()
					.sort(null)
					.value(function (d) { return 1; }),

				x = d3.scale.linear().range([0, 2 * Math.PI]),
				y = d3.scale.linear().range([0, radius]),

				arc = d3.svg.arc()
					.startAngle(function (d) { return Math.max(0, Math.min(2 * Math.PI, x(d.x))); })
					.endAngle(function (d) { return Math.max(0, Math.min(2 * Math.PI, x(d.x + d.dx))); })
					.innerRadius(function (d) { return Math.max(0, y(d.y)); })
					.outerRadius(function (d) { return Math.max(0, y(d.y + d.dy)); }),

				tip = d3tip()
					.offset([-10, 0])
					.html(function (d) {
						return d.name + ' ' + d.duration;
					}),

				node;

				vis.call(tip);

				scope.$watch('data', function (newData, oldData) {
					if (!newData) {
						return;
					}

					node = newData;


					var path = vis.datum(newData).selectAll('path')
						.data(partition.value(function (d) { return d.duration; }).nodes)
						.enter().append('path')
						.attr('d', arc)
						.style('fill', color) //function (d) { return color((d.children ? d : d.parent).name); })
						.on('click', click)
						.on('mouseover', tip.show)
						.on('mouseout', tip.hide)
						.each(stash);


					function click(d) {
						node = d;
						path.transition()
							.duration(1000)
							.attrTween('d', arcTweenZoom(d));
					}

				});

				function stash(d) {
					d.x0 = d.x;
					d.dx0 = d.dx;
				}

				function arcTweenZoom(d) {
					var xd = d3.interpolate(x.domain(), [d.x, d.x + d.dx]),
						yd = d3.interpolate(y.domain(), [d.y, 1]),
						yr = d3.interpolate(y.range(), [d.y ? 20 : 0, radius]);

					return function (d, i) {
						return i ?
							function (t) { return arc(d); } :
							function (t) { x.domain(xd(t)); y.domain(yd(t)).range(yr(t)); return arc(d); };
					};
				}
			}
		};
	});
};
