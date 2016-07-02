module.exports = function (app) {
	require('./stats')(app);
	require('./tags')(app);
};
