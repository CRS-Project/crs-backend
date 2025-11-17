package config

import (
	"fmt"

	"log"
	"os"

	"github.com/CRS-Project/crs-backend/db"
	"github.com/CRS-Project/crs-backend/internal/api/controller"
	"github.com/CRS-Project/crs-backend/internal/api/repository"
	"github.com/CRS-Project/crs-backend/internal/api/routes"
	"github.com/CRS-Project/crs-backend/internal/api/service"
	"github.com/CRS-Project/crs-backend/internal/middleware"
	mailer "github.com/CRS-Project/crs-backend/internal/pkg/email"
	"github.com/CRS-Project/crs-backend/internal/pkg/google/oauth"
	"github.com/gin-gonic/gin"
)

type RestConfig struct {
	server *gin.Engine
}

func NewRest() RestConfig {
	db := db.New()
	app := gin.Default()
	server := NewRouter(app)
	middleware := middleware.New(db)

	var (
		//=========== (PACKAGE) ===========//
		mailerService mailer.Mailer = mailer.New()
		oauthService  oauth.Oauth   = oauth.New()
		// awsS3Service  storage.AwsS3 = storage.NewAwsS3()

		//=========== (REPOSITORY) ===========//
		userRepository                      repository.UserRepository                      = repository.NewUser(db)
		packageRepository                   repository.PackageRepository                   = repository.NewPackage(db)
		userDisciplineRepository            repository.UserDisciplineRepository            = repository.NewUserDiscipline(db)
		commentRepository                   repository.CommentRepository                   = repository.NewComment(db)
		documentRepository                  repository.DocumentRepository                  = repository.NewDocument(db)
		areaOfConcernGroupRepository        repository.AreaOfConcernGroupRepository        = repository.NewAreaOfConcernGroup(db)
		areaOfConcernRepository             repository.AreaOfConcernRepository             = repository.NewAreaOfConcern(db)
		areaOfConcernConsolidatorRepository repository.AreaOfConcernConsolidatorRepository = repository.NewAreaOfConcernConsolidator(db)
		statisticRepository                 repository.StatisticRepository                 = repository.NewStatistic(db)

		//=========== (SERVICE) ===========//
		authService               service.AuthService               = service.NewAuth(userRepository, mailerService, oauthService, db)
		userService               service.UserService               = service.NewUser(userRepository, userDisciplineRepository, areaOfConcernConsolidatorRepository, packageRepository, db)
		packageService            service.PackageService            = service.NewPackage(packageRepository, userRepository, db)
		userDisciplineService     service.UserDisciplineService     = service.NewUserDiscipline(userDisciplineRepository, db)
		documentService           service.DocumentService           = service.NewDocument(documentRepository, packageRepository, userRepository, db)
		commentService            service.CommentService            = service.NewComment(commentRepository, documentRepository, areaOfConcernRepository, userRepository, db)
		areaOfConcernGroupService service.AreaOfConcernGroupService = service.NewAreaOfConcernGroup(areaOfConcernGroupRepository, areaOfConcernRepository, packageRepository, commentRepository, userRepository, userDisciplineRepository, db)
		areaOfConcernService      service.AreaOfConcernService      = service.NewAreaOfConcern(areaOfConcernRepository, areaOfConcernGroupRepository, areaOfConcernConsolidatorRepository, commentRepository, packageRepository, userRepository, userDisciplineRepository, db)
		statisticService          service.StatisticService          = service.NewStatistic(statisticRepository, commentRepository, documentRepository, areaOfConcernRepository, userRepository, db)

		//=========== (CONTROLLER) ===========//
		authController               controller.AuthController               = controller.NewAuth(authService)
		packageController            controller.PackageController            = controller.NewPackage(packageService)
		userController               controller.UserController               = controller.NewUser(userService)
		userDisciplineController     controller.UserDisciplineController     = controller.NewUserDiscipline(userDisciplineService)
		documentController           controller.DocumentController           = controller.NewDocument(documentService)
		commentController            controller.CommentController            = controller.NewComment(commentService)
		areaOfConcernGroupController controller.AreaOfConcernGroupController = controller.NewAreaOfConcernGroup(areaOfConcernGroupService)
		areaOfConcernController      controller.AreaOfConcernController      = controller.NewAreaOfConcern(areaOfConcernService)
		statisticController          controller.StatisticController          = controller.NewStatistic(statisticService)
	)

	// Register all routes
	routes.Auth(server, authController, middleware)
	routes.User(server, userController, middleware)
	routes.Package(server, packageController, middleware)
	routes.UserDiscipline(server, userDisciplineController, middleware)
	routes.Document(server, documentController, middleware)
	routes.AreaOfConcernGroup(server, areaOfConcernGroupController, middleware)
	routes.AreaOfConcern(server, areaOfConcernController, middleware)
	routes.Comment(server, commentController, middleware)
	routes.Statistic(server, statisticController, middleware)

	return RestConfig{
		server: server,
	}
}

func (ap *RestConfig) Start() {
	port := os.Getenv("APP_PORT")
	host := os.Getenv("APP_HOST")
	if port == "" {
		port = "8080"
	}

	serve := fmt.Sprintf("%s:%s", host, port)
	if err := ap.server.Run(serve); err != nil {
		log.Panicf("failed to start server: %s", err)
	}
	log.Println("server start on port ", serve)
}
