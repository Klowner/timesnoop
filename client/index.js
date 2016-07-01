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

app.run(function () {
	console.log('running app');
});


