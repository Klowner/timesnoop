var path = require('path'),
	webpack = require('webpack'),
	SplitByPathPlugin = require('webpack-split-by-path');

var plugins = [
	new webpack.optimize.DedupePlugin(),
	new SplitByPathPlugin([{
		name: 'libs',
		path: path.join(__dirname, 'node_modules'),
	}]),
	//new webpack.optimize.UglifyJsPlugin(),
];

module.exports = {
	context: __dirname,
	entry: {
		main: path.join(__dirname, 'client/index.js')
	},

	output: {
		path: path.join(__dirname, 'src/app/static/'),
		publicPath: '/static/',
		filename: '[name].js'
	},

	plugins: plugins,

	module: {
		loaders: [
			{ test: /\.css$/, loader: 'style!css' },
			{ test: /\.html$/, loader: 'html' },
			{ test: /\.jsx?$/, loader: 'babel?presets[]=es2015', exclude: /node_modules/ },
			{ test: /\.scss$/, loader: 'style!css!sass' }
		]
	},

	resolve: {
		moduleDirectories: ['node_modules', 'bower_components', 'client'],
		extensions: ['', '.js']
	},

	sassLoader: {
		includePaths: [path.resolve(__dirname, './client/style')]
	}
};
