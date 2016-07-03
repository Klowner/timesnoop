module.exports = function (app) {
	require('./main')(app);
	require('./matchers')(app);
	require('./stats')(app);
	require('./tags')(app);
};
