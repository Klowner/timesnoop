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
			link: function (scope, el, attrs) {
				var w, h,
					vis = d3.select(el[0]).append('svg'),
					sunburst = vis.append('g').attr('class', 'sunburst'),
					radius = 100,
					x = d3.scale.linear().range([0, 2*Math.PI]),
					y = d3.scale.linear().range([0, radius]),
					colors = d3.scale.category20c(),
					color = function (d) {
							return d.color || (d.parent && d.parent.color) || '#ffffff';
						},
					partition = d3.layout.partition()
						.sort(null)
						.value(function (d) { return 1; }),

					arc = d3.svg.arc()
						.startAngle(function (d) { return Math.max(0, Math.min(2 * Math.PI, x(d.x))); })
						.endAngle(function (d) { return Math.max(0, Math.min(2 * Math.PI, x(d.x + d.dx))); })
						.innerRadius(function (d) { return Math.max(0, y(d.y)); })
						.outerRadius(function (d) { return Math.max(0, y(d.y + d.dy)); }),

					tip = d3tip()
						.attr('class', 'd3-tip')
						.offset(function () {
							console.log(this);
							return [0, 0];
						})
						.direction('s')
						.html(function (d) {
							return d.name + ' ' + (d.duration / 60).toFixed(2) + 'hrs';
						}),

					sunburstPath,
					node;


				vis.call(tip);

				el = el[0];

				angular.element(el).addClass('d3-sunburst');

				scope.$watch(function () {
					w = el.clientWidth;
					h = el.clientHeight;
					return w + h;
				}, resize);

				function resize() {
					radius = Math.min(w, h) / 2 - 5;
					vis.attr({width: w, height: h});
					y.range([0, radius]);
					sunburst.attr('transform', 'translate(' + (w / 2) + ',' + (h / 2 + 5) + ')');
					update();
					rescale();
				}

				scope.$watch('data', update);

				function update() {
					if (!scope.data) { return; }

					sunburst.selectAll('*').remove();

					var path = sunburst.datum(scope.data).selectAll('path')
						.data(partition.value(function (d) { return d.duration; }).nodes)
						.enter().append('path')
						.attr('d', arc)
						.style('fill', color)
						.on('click', click)
						.on('mouseover', tip.show)
						.on('mouseout', tip.hide)
						.each(stash);

					function click(d) {
						if (d !== node) {
							node = d;
							path.transition()
								.duration(1000)
								.attrTween('d', arcTweenZoom(d));
						}
					}
				}

				function stash(d) {
					d.x0 = d.x;
					d.dx0 = d.dx;
				}

				function rescale() {

					sunburst.selectAll('path')
						.transition()
						.duration(1000)
						.call(y);
				}

				function arcTweenZoom(d) {
					var xd = d3.interpolate(x.domain(), [d.x, d.x + d.dx]),
						yd = d3.interpolate(y.domain(), [d.y, 1]),
						yr = d3.interpolate(y.range(), [d.y ? 30 : 0, radius]);

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
