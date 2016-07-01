module.exports = function (app) {
	require('./main')(app);
	require('./stats')(app);
};
