module.exports = function (app) {
	require('./matchers')(app);
	require('./stats')(app);
	require('./tags')(app);
};
