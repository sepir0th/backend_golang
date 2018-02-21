package main

import (
	"log"
	"net/http"
	"net/url"
	"os"

	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
	"gopkg.in/session.v1"
	"github.com/kataras/iris"
)

var (
	globalSessions *session.Manager
)

func init() {
	globalSessions, _ = session.NewManager("memory", `{"cookieName":"gosessionid","gclifetime":3600}`)
	go globalSessions.GC()
}

func MainOauth(application *iris.Application) {
	manager := manage.NewDefaultManager()
	// token store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	clientStore := store.NewClientStore()
	clientStore.Set("222222", &models.Client{
		ID:     "222222",
		Secret: "22222222",
		Domain: "http://localhost:9095",
	})
	manager.MapClientStorage(clientStore)

	srv := server.NewServer(server.NewConfig(), manager)
	srv.SetUserAuthorizationHandler(userAuthorizeHandler)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	//http.HandleFunc("/login", loginHandler)
	application.Handle("", "/login", func(ctx iris.Context) {
		loginHandler(ctx.ResponseWriter(), ctx.Request())
	})

	//http.HandleFunc("/auth", authHandler)
	application.Handle("", "/auth", func(ctx iris.Context) {
		authHandler(ctx.ResponseWriter(), ctx.Request())
	})

	//http.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
	//	err := srv.HandleAuthorizeRequest(w, r)
	//	if err != nil {
	//		http.Error(w, err.Error(), http.StatusBadRequest)
	//	}
	//})
	application.Handle("", "/authorize", func(ctx iris.Context) {
		err := srv.HandleAuthorizeRequest(ctx.ResponseWriter(), ctx.Request())
		if err != nil {
			http.Error(ctx.ResponseWriter(), err.Error(), http.StatusBadRequest)
		}
	})

	//http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
	//	err := srv.HandleTokenRequest(w, r)
	//	if err != nil {
	//		http.Error(w, err.Error(), http.StatusInternalServerError)
	//	}
	//})
	application.Handle("", "/token", func(ctx iris.Context) {
		err := srv.HandleTokenRequest(ctx.ResponseWriter(), ctx.Request())
		if err != nil {
			http.Error(ctx.ResponseWriter(), err.Error(), http.StatusInternalServerError)
		}
	})
}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	us, err := globalSessions.SessionStart(w, r)
	uid := us.Get("UserID")
	if uid == nil {
		if r.Form == nil {
			r.ParseForm()
		}
		us.Set("Form", r.Form)
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}
	userID = uid.(string)
	us.Delete("UserID")
	return
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		us, err := globalSessions.SessionStart(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		us.Set("LoggedInUserID", "000000")
		w.Header().Set("Location", "/auth")
		w.WriteHeader(http.StatusFound)
		return
	}
	outputHTML(w, r, "./static/login.html")
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	us, err := globalSessions.SessionStart(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if us.Get("LoggedInUserID") == nil {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}
	if r.Method == "POST" {
		form := us.Get("Form").(url.Values)
		u := new(url.URL)
		u.Path = "/authorize"
		u.RawQuery = form.Encode()
		w.Header().Set("Location", u.String())
		w.WriteHeader(http.StatusFound)
		us.Delete("Form")
		us.Set("UserID", us.Get("LoggedInUserID"))
		return
	}
	outputHTML(w, r, "./static/auth.html")
}

func outputHTML(w http.ResponseWriter, req *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
}