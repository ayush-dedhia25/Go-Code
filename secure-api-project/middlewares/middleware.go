package middleware

import (
   "log"
   "time"
   "strings"
   "net/http"
   "github.com/justinas/alice"
   "secure-api-project/myJWT"
   "secure-api-project/db"
   "secure-api-project/templates"
)

func recoverHandler(next http.Handler) http.Handler {
   fn := func(res http.ResponseWriter, req *http.Request) {
      defer func() {
         if err := recover(); err != nil {
            log.Panic("Recovered Panic: %+v", err)
            http.Error(res, http.StatusText(500), 500)
         }
      }()
      next.ServeHTTP(res, req)
   }
   return http.HandlerFunc(fn)
}

func authHandler(next http.Handler) http.Handler {
   fn := func(res http.ResponseWriter, req *http.Request) {
      switch req.URL.Path {
         case "/restricted", "/logout", "/deleteUser":
            log.Println("In auth restricted section!")
            authCookie, authErr := req.Cookie("AuthToken")
            if authErr == http.ErrNoCookie {
               log.Println("Unauthorized attempt! No auth cookie was found.")
               nullifyTokenCookies(&res, req)
               http.Error(res, http.StatusText(401), 401)
               return
            }
            else if authErr != nil {
               log.Panic("Panic: %+v", authErr)
               nullifyTokenCookies(&res, req)
               http.Error(res, http.StatusText(500), 500)
               return
            }
            
            refreshCookie, refreshCookieErr := req.Cookie("RefreshToken")
            if refreshCookieErr == http.ErrNoCookie {
               log.Println("Unauthorized attempt! No refresh cookie was found.")
               nullifyTokenCookies(&res, req)
               http.Redirect(res, req, "/login", 302)
               return
            }
            else if refreshCookieErr != nil {
               log.Panic("Panic: %+v", refreshCookieErr)
               nullifyTokenCookies(&res, req)
               http.Error(res, http.StatusText(500), 500)
               return
            }
            
            requestCsrfToken := grabCsrfFromRequest(req)
            log.Println(requestCsrfToken)
            
            authTokenString, refreshTokenString, csrfSecret, err := myJWT.CheckAndRefreshTokens(authCookie.Value, refreshCookie.Value, requestCsrfToken)
            if err != nil {
               if err.Error() == "unauthorized" {
                  log.Println("Unauthorized attempt! No JWT's valid.")
                  http.Error(res, http.StatusText(401), 401)
                  return
               }
               else {
                  log.Panic("Error not nil")
                  log.Panic("Panic: %+v", err)
                  http.Error(res, http.StatusText(500), 500)
               }
            }
            log.Println("Successfully recreated JWT")
            res.Header().Set("Access-Control-Allow-Origin", "*")
            setAuthAndRefreshCookies(&res, authTokenString, refreshTokenString)
            res.Header().Set("X-CSRF-Token", csrfSecret)
         default:
            // No checks necessary :)
      }
      next.ServeHTTP(res, req)
   }
}

func logicHandler(res http.ResponseWriter, req *http.Request) {
   switch req.URL.Path {
      case "/restricted":
         csrfSecret := grabCsrfFromRequest(req)
         templates.RenderTemplate(res, "restricted", &templates.DashboardPage{csrfSecret, "Hello, Ayush!"})
      case "/login":
         switch req.Method {
            case "GET":
            case "POST":
            default:
         }
      case "/register":
         switch req.Method {
            case "GET":
               templates.RenderTemplate(res, "register", &templates.RegisterPage{false, ""})
            case "POST":
               req.ParseForm()
               log.Println(req.Form)
               _, uuid, err := db.FetchUserByUsername(strings.Join(req.Form["username"], ""))
               if err == nil {
                  res.WriteHeader(http.StatusUnauthorized)
               } else {
                  role := 
                  uuid, err := db.StoreUser(strings.Join(req.Form["username"], ""), strings.Join(req.Form["password"], ""), role)
                  if err != nil {
                     http.Error(res, http.StatusText(500), 500)
                  }
                  log.Println("uuid: " + uuid)
                  
                  authToken, refreshToken, csrfSecret, err := myJWT.CreateNewTokens(uuid, role)
                  if err != nil {
                     http.Error(res, http.StatusText(500), 500)
                  }
                  setAuthAndRefreshCookies(&res, authToken, refreshToken)
                  res.Header().Set("X-CSRF-Token", csrfSecret)
                  res.WriteHeader(http.StatusOK)
               }
            default:
               res.WriteHeader(http.StatusMethodNotAllowed)
         }
      case "/logout":
         nullifyTokenCookies(&res, req)
         http.Redirect(res, req, "/login", 302)
      case "/deleteUser":
      default: 
   }
}

func nullifyTokenCookies(res *http.ResponseWriter, req *http.Request) {
   authCookie := http.Cookie{
      Name: "AuthToken",
      Value: "",
      Expires: time.Now().Add(-1000 * time.Hour),
      HttpOnly: true,
   }
   http.SetCookie(*res, &authCookie)
   
   refreshCookie := http.Cookie{
      Name: "RefreshToken",
      Value: "",
      Expires: time.Now().Add(-1000 * time.Hour),
      HttpOnly: true,
   }
   http.SetCookie(*res, &refreshCookie)
   
   refreshCookie, refreshErr := req.Cookie("RefreshToken")
   if refreshErr == http.ErrNoCookie {
      return
   } else if refreshErr != nil {
      log.Panic("Error %+v", refreshErr)
      http.Error(*res, http.StatusText(500), 500)
   }
   
   myJWT.RevokeRefreshToken(refreshCookie.Value)
}

func setAuthAndRefreshCookies(res *http.ResponseWriter, authTokenString string, refreshTokenString string) {
   authCookie := http.Cookie{
      Name: "AuthToken",
      Value: authTokenString,
      HttpOnly: true,
   }
   http.SetCookie(*res, &authCookie)
   
   refreshCookie := http.Cookie{
      Name: "RefreshToken",
      Value: refreshTokenString,
      HttpOnly: true,
   }
   http.SetCookie(*res, &refreshCookie)
}

func grabCsrfFromRequest(req *http.Request) string {
   csrfFromForm := req.FormValue("X-CSRF-Token")
   if csrfFromForm != "" {
      return csrfFromForm
   } else {
      return req.Header.Get("X-CSRF-Token")
   }
}

func NewHandler() http.Handler {
   return alice.New(recoverHandler, authHandler).ThenFunc(logicHandler)
}