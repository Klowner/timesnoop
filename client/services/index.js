module.exports = function (app) {
	require('./matchers')(app);
	require('./me2tags')(app);
	require('./stats')(app);
	require('./tags')(app);
};
