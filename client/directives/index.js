module.exports = function (app) {
	//require('./c3chart')(app);
	require('./matchexpression-editor')(app);
	require('./matchtagswidget')(app);
	require('./tags-editor')(app);
};
