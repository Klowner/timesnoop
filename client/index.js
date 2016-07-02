//@require "../templates/**.html"
var app = require('./app'),
	controllers = require('./controllers');

// Style
require('./style/main.scss');
require('../node_modules/c3/c3.min.css');

// Components
require('./controllers')(app);
require('./services')(app);
require('./filters')(app);
require('./directives')(app);
require('./config')(app);

app.run(function ($templateCache) {

	$templateCache.put('templates/forms/tag.html', require('./templates/forms/tag.html'));

	console.log('templtecache');
	console.log('running app');
});


