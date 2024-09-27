package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"sysbitBroker/config"
	"sysbitBroker/middleware"
	"sysbitBroker/pkg"
	"sysbitBroker/route"

	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	//Initial Configuration
	cfg := config.NewConfig()

	gin.SetMode(cfg.Server.RunMode)
	// r = router
	router := gin.New()
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!!!")
	})

	// OpenID Initialized
	openIDProvider := pkg.NewOpenIDProvider(cfg)
	openIDCfg := pkg.NewOpenIDConfig(cfg, openIDProvider)

	//Postgre Initialize
	pg, err := pkg.NewPgx(cfg)
	if err != nil {
		fmt.Println(err.Error())
		// log.Fatal(logging.Postgres, logging.Startup, err.Error(), nil)
	}

	//Redis Initialize
	rdb, err := pkg.NewRedis(cfg)
	if err != nil {
		fmt.Println(err.Error())
		// log.Fatal(logging.Postgres, logging.Startup, err.Error(), nil)
	}
	api := router.Group("/api")
	//Add check usertoken
	m := api.Group("/m")
	{
		//Handle for User
		auth := m.Group("/auth")
		route.User(auth, cfg, pg, openIDProvider, openIDCfg)

		//Handle for Lesson
		lesson := m.Group("/lesson", middleware.JwtAuthMiddleware(cfg.JWT.RefreshTokenSecret))
		route.Lesson(lesson, cfg, pg, rdb)

	}

	// defer conn.Close(context.Background())

	// Create a Router without the Token Authenitcation
	// router := mux.NewRouter()
	// router.HandleFunc("/api/AppToken", getAppToken).Methods("PUT")
	// router.HandleFunc("/api/AdminToken", getAdminToken).Methods("PUT")
	// router.HandleFunc("/api/chkApi", chkApi).Methods("GET")
	// router.HandleFunc("/api/regApp", regApp).Methods("POST")

	// Admin
	// tokenRouter := router.PathPrefix("/").Subrouter()
	// tokenRouter.Use(chkToken)

	// Admin
	// tokenRouter.HandleFunc("/api/getLessonHeaders/{modCode}", getLessonHeaders).Methods("GET")
	// tokenRouter.HandleFunc("/api/updLessonHeader/{lessonCode}", updLessonHeader).Methods("POST")
	// tokenRouter.HandleFunc("/api/updQuiz/{quizCode}/{lingoCode}", updQuiz).Methods("POST")
	// tokenRouter.HandleFunc("/api/updLessonStep/{lessonCode}/{lingoCode}", updLessonStep).Methods("POST")
	// tokenRouter.HandleFunc("/api/getLessonDetail/{lessonCode}/{lingoCode}", getLessonDetail).Methods("GET")
	// tokenRouter.HandleFunc("/api/getQuizDetail/{quizCode}/{lingoCode}", getQuizDetail).Methods("GET")

	// App
	// tokenRouter.HandleFunc("/api/syncApp", syncApp).Methods("GET")
	// tokenRouter.HandleFunc("/api/updProgress/{lessonCode}/{result}", updProgress).Methods("PUT")
	// tokenRouter.HandleFunc("/api/getProgress", getProgress).Methods("GET")

	srv := &http.Server{
		Addr:    cfg.Server.ExternalPort,
		Handler: router,
	}
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")

}

// Tokens

// func getAppToken(w http.ResponseWriter, r *http.Request) {
// 	auth.GetAppToken(w, r, conn)
// }

// func getAdminToken(w http.ResponseWriter, r *http.Request) {
// 	auth.GetAdminToken(w, r, conn)
// }

// // APIs

// func chkApi(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Welcome to the InglesGuru API")
// 	fmt.Println("Endpoint Hit: InglesGuru API")
// }

// func syncApp(w http.ResponseWriter, r *http.Request) {

// 	data.SyncApp(w, r, conn)
// }

// func getQuizDetail(w http.ResponseWriter, r *http.Request) {

// 	w.Header().Set("Content-Type", "application/json")

// 	vars := mux.Vars(r)
// 	quizCode := vars["quizCode"]
// 	lingoCode := vars["lingoCode"]
// 	data.GetQuizDetail(w, r, conn, quizCode, lingoCode)
// }

// func getLessonDetail(w http.ResponseWriter, r *http.Request) {

// 	w.Header().Set("Content-Type", "application/json")

// 	vars := mux.Vars(r)
// 	lessonCode := vars["lessonCode"]
// 	lingoCode := vars["lingoCode"]
// 	data.GetLessonDetail(w, r, conn, lessonCode, lingoCode)
// }

// func getLessonHeaders(w http.ResponseWriter, r *http.Request) {

// 	w.Header().Set("Content-Type", "application/json")

// 	vars := mux.Vars(r)
// 	modCode := vars["modCode"]

// 	data.GetLessonHeaders(w, r, conn, modCode)
// }

// func updQuiz(w http.ResponseWriter, r *http.Request) {

// 	w.Header().Set("Content-Type", "application/json")

// 	vars := mux.Vars(r)
// 	quizCode := vars["quizCode"]
// 	lingoCode := vars["lingoCode"]

// 	data.UpdQuiz(w, r, conn, quizCode, lingoCode)
// }

// func updLessonStep(w http.ResponseWriter, r *http.Request) {

// 	w.Header().Set("Content-Type", "application/json")

// 	vars := mux.Vars(r)
// 	lessonCode := vars["lessonCode"]
// 	lingoCode := vars["lingoCode"]

// 	data.UpdLessonStep(w, r, conn, lessonCode, lingoCode)
// }

// func updProgress(w http.ResponseWriter, r *http.Request) {

// 	w.Header().Set("Content-Type", "application/json")

// 	vars := mux.Vars(r)
// 	lessonCode := vars["lessonCode"]
// 	result := vars["result"]

// 	data.UpdProgress(w, r, conn, lessonCode, result)
// }

// func updLessonHeader(w http.ResponseWriter, r *http.Request) {

// 	w.Header().Set("Content-Type", "application/json")

// 	vars := mux.Vars(r)
// 	lessonCode := vars["lessonCode"]
// 	modCode := vars["modCode"]

// 	data.UpdLessonHeader(w, r, conn, modCode, lessonCode)
// }

// func getProgress(w http.ResponseWriter, r *http.Request) {

// 	w.Header().Set("Content-Type", "application/json")

// 	data.GetProgress(w, r, conn)
// }

// func regApp(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	data.RegisterApp(w, r, conn)
// }

// func chkToken(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		tokenString := r.Header.Get("Authorization")
// 		if tokenString == "" {
// 			errHandler.ErrMsg(w, errors.New("no token"), http.StatusUnauthorized)
// 			return
// 		}

// 		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			return secretKey, nil
// 		})

// 		if err != nil {

// 			if err.Error() == "Token is expired" {
// 				errHandler.ErrMsg(w, err, http.StatusForbidden)
// 			} else {
// 				errHandler.ErrMsg(w, err, http.StatusUnauthorized)
// 			}

// 			return
// 		}

// 		if !token.Valid {
// 			errHandler.ErrMsg(w, errors.New("token not valid"), http.StatusUnauthorized)
// 			return
// 		}

// 		claims, ok := token.Claims.(jwt.MapClaims)
// 		if !ok {
// 			errHandler.ErrMsg(w, errors.New("claims not valid"), http.StatusUnauthorized)
// 			return
// 		}

// 		tokenType := claims["tokenType"].(string)

// 		if tokenType == "admin" {
// 			name, ok1 := claims["name"].(string)
// 			roleCode, ok2 := claims["roleCode"].(string)
// 			email, ok3 := claims["email"].(string)

// 			if !ok1 {
// 				errHandler.ErrMsg(w, errors.New("name in claims not valid"), http.StatusUnauthorized)
// 				return
// 			}

// 			if !ok2 {
// 				errHandler.ErrMsg(w, errors.New("rolecode in claims not valid"), http.StatusUnauthorized)
// 				return
// 			}

// 			if !ok3 {
// 				errHandler.ErrMsg(w, errors.New("email in claims not valid"), http.StatusUnauthorized)
// 				return
// 			}

// 			ctx := context.WithValue(r.Context(), "name", name)
// 			ctx = context.WithValue(ctx, "roleCode", roleCode)
// 			ctx = context.WithValue(ctx, "email", email)

// 			next.ServeHTTP(w, r.WithContext(ctx))

// 		} else {

// 			appId, ok := claims["appId"].(string)

// 			if !ok {
// 				errHandler.ErrMsg(w, errors.New("appId in claims not valid"), http.StatusUnauthorized)
// 				return
// 			}

// 			ctx := context.WithValue(r.Context(), "appId", appId)
// 			next.ServeHTTP(w, r.WithContext(ctx))
// 		}

// 	})
// }
