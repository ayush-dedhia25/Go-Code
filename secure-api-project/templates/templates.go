package templates

import (
   "log"
   "net/http"
   "html/template"
)

type LoginPage struct {
   BAlertUser bool
   AlertMsg   string
}

type RegisterPage struct {
   BAlertUser bool
   AlertMsg   string
}

type DashboardPage struct {
   BAlertUser bool
   AlertMsg   string
}

var template = template.Must(template.ParseFiles(
   "./server/templates/templateFiles/login.tmpl",
   "./server/templates/templateFiles/register.tmpl",
   "./server/templates/templateFiles/restricted.tmpl"
))

func RenderTemplate(res http.ResponseWriter, tmpl string, p interface{}) {
   err := template.ExecuteTemplate(res, tmpl + ".tmpl", p)
   if err != nil {
      log.Printf("You've got template err %v", err)
      http.Error(res, err.Error(), http.StatusInternalServerError)
   }
}