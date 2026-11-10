package main

import (
 "log"
 "net/http"
 "os"
 "github.com/SouzaNataliaDev/whatsapp-messaging-api/internal/campaign"
 "github.com/SouzaNataliaDev/whatsapp-messaging-api/internal/company"
 "github.com/SouzaNataliaDev/whatsapp-messaging-api/internal/contact"
 "github.com/SouzaNataliaDev/whatsapp-messaging-api/internal/database"
 "github.com/gin-gonic/gin"
)
func main() {
 dsn:=os.Getenv("MYSQL_DSN")
 if dsn=="" {log.Fatal("MYSQL_DSN is required")}
 db,err:=database.Connect(dsn);if err!=nil {log.Fatal(err)};defer db.Close()
 r:=gin.Default()
 r.GET("/health",func(c *gin.Context){if err:=db.PingContext(c.Request.Context());err!=nil {c.JSON(http.StatusServiceUnavailable,gin.H{"status":"unhealthy"});return};c.JSON(http.StatusOK,gin.H{"status":"ok","database":"connected"})})
 companies:=&company.Handler{Service:&company.Service{Repo:&company.Repository{DB:db}}}
 contacts:=&contact.Handler{Service:&contact.Service{Repo:&contact.Repository{DB:db}}}
 campaigns:=&campaign.Handler{Service:&campaign.Service{Repo:&campaign.Repository{DB:db}}}
 v1:=r.Group("/v1")
 v1.POST("/companies",companies.Create)
 v1.POST("/contacts",contacts.Create)
 v1.POST("/campaigns",campaigns.Create)
 v1.GET("/campaigns/:id",campaigns.Get)
 v1.POST("/campaigns/:id/send",campaigns.Send)
 v1.GET("/campaigns/:id/stats",campaigns.Stats)
 port:=os.Getenv("PORT");if port=="" {port="8080"}
 log.Fatal(r.Run(":"+port))
}
