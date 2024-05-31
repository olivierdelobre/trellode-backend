package api

// Set the routes for API
func (s *server) Routes() {
	v1 := s.Router.Group("/trellode-api/v1")

	s.Router.GET("/healthcheck", s.getHealthcheck)
	s.Router.GET("/liveness", s.getLiveness)

	v1.POST("/users/register", s.registerUser)
	v1.POST("/users/authenticate", s.authenticate)

	v1.GET("/boards/:id", s.getBoard)
	v1.GET("/boards", s.getBoards)
	v1.POST("/boards", s.createBoard)
	v1.PUT("/boards/:id", s.updateBoard)
	v1.DELETE("/boards/:id", s.deleteBoard)

	v1.GET("/lists/:id", s.getList)
	v1.POST("/lists", s.createList)
	v1.PUT("/lists/:id", s.updateList)
	v1.DELETE("/lists/:id", s.deleteList)
	//v1.PUT("/lists/:id/order/:idsordered", s.updateCardsOrder)
	v1.PUT("/lists/:id/order", s.updateCardsOrder)

	v1.GET("/cards/:id", s.getCard)
	v1.POST("/cards", s.createCard)
	v1.PUT("/cards/:id", s.updateCard)
	v1.DELETE("/cards/:id", s.deleteCard)

	v1.GET("/comments/:id", s.getComment)
	v1.GET("/cards/:id/comments", s.getComments)
	v1.POST("/comments", s.createComment)
	v1.PUT("/comments/:id", s.updateComment)
	v1.DELETE("/comments/:id", s.deleteComment)

	v1.GET("/backgrounds/:id", s.getBackground)
	v1.GET("/backgrounds", s.getBackgrounds)
	v1.POST("/backgrounds", s.createBackground)
	v1.DELETE("/backgrounds/:id", s.deleteBackground)

	v1.GET("/checklists/:id", s.getChecklist)
	v1.POST("/checklists", s.createChecklist)
	v1.PUT("/checklists/:id", s.updateChecklist)
	v1.DELETE("/checklists/:id", s.deleteChecklist)
	v1.GET("/checklistitems/:id", s.getChecklistItem)
	v1.POST("/checklistitems", s.createChecklistItem)
	v1.PUT("/checklistitems/:id", s.updateChecklistItem)
	v1.DELETE("/checklistitems/:id", s.deleteChecklistItem)

	v1.GET("/logs", s.getLogs)

	v1.OPTIONS("/users/register", s.options)
	v1.OPTIONS("/users/authenticate", s.options)
	v1.OPTIONS("/boards", s.options)
	v1.OPTIONS("/boards/:id", s.options)
	v1.OPTIONS("/boards/:id/lists", s.options)
	v1.OPTIONS("/lists", s.options)
	v1.OPTIONS("/lists/:id", s.options)
	v1.OPTIONS("/lists/:id/cards", s.options)
	v1.OPTIONS("/cards", s.options)
	v1.OPTIONS("/cards/:id", s.options)
	v1.OPTIONS("/cards/:id/comments", s.options)
	v1.OPTIONS("/comments", s.options)
	v1.OPTIONS("/comments/:id", s.options)
	v1.OPTIONS("/backgrounds", s.options)
	v1.OPTIONS("/backgrounds/:id", s.options)
	v1.OPTIONS("/logs", s.options)
	v1.OPTIONS("/lists/:id/order", s.options)
	v1.OPTIONS("/checklists", s.options)
	v1.OPTIONS("/checklists/:id", s.options)
	v1.OPTIONS("/checklistitems", s.options)
	v1.OPTIONS("/checklistitems/:id", s.options)

	//v1.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
