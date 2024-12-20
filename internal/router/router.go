package router

import (
	"borrow_book/internal/handler"

	"github.com/gin-gonic/gin"
)

type AppRouter struct {
	bookController   *handler.BookHandler
	authorController *handler.AuthorHandler
	borrowController *handler.BorrowHandler
	swaggerRouter    *SwaggerRouter
}

func NewAppRouter(
	bookController *handler.BookHandler,
	authorController *handler.AuthorHandler,
	borrowController *handler.BorrowHandler,
	swaggerRouter *SwaggerRouter,
) *AppRouter {
	return &AppRouter{
		bookController:   bookController,
		authorController: authorController,
		borrowController: borrowController,
		swaggerRouter:    swaggerRouter,
	}
}

func (a *AppRouter) RegisterBookRoutes(r *gin.RouterGroup) {
	public := r.Group("/books")
	{
		public.GET("", a.bookController.ListBooks)
		public.GET("/:id", a.bookController.GetBook)
		public.POST("", a.bookController.CreateBook)
		public.PUT("/:id", a.bookController.UpdateBook)
		public.DELETE("/:id", a.bookController.DeleteBook)
	}
}

func (a *AppRouter) RegisterAuthorRoutes(r *gin.RouterGroup) {
	public := r.Group("/authors")
	{
		public.GET("", a.authorController.ListAuthors)
		public.GET("/:id", a.authorController.GetAuthor)
		public.POST("", a.authorController.CreateAuthor)
		public.PUT("/:id", a.authorController.UpdateAuthor)
		public.DELETE("/:id", a.authorController.DeleteAuthor)
	}
}

func (a *AppRouter) RegisterBorrowRoutes(r *gin.RouterGroup) {
	public := r.Group("/borrows")
	{
		public.GET("", a.borrowController.ListBorrows)
		public.GET("/:id", a.borrowController.GetBorrow)
		public.POST("", a.borrowController.CreateBorrow)
		public.PUT("/:id", a.borrowController.UpdateBorrow)
		public.DELETE("/:id", a.borrowController.DeleteBorrow)
	}
}

// RegisterSwaggerRoutes sets up the route for Swagger API documentation
func (a *AppRouter) RegisterSwaggerRoutes(r *gin.RouterGroup) {
	// Check if SwaggerRouter is initialized before registering
	if a.swaggerRouter != nil {
		a.swaggerRouter.Register(r)
	}
}
