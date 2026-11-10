package main
import (
 "context"
 "log"
 "os"
 "os/signal"
 "syscall"
 "time"
 "github.com/SouzaNataliaDev/whatsapp-messaging-api/internal/database"
 "github.com/SouzaNataliaDev/whatsapp-messaging-api/internal/queue"
)
func main(){
 db,err:=database.Connect(os.Getenv("MYSQL_DSN"));if err!=nil{log.Fatal(err)};defer db.Close()
 url:=os.Getenv("RABBITMQ_URL");if url==""{log.Fatal("RABBITMQ_URL required")}
 ctx,stop:=signal.NotifyContext(context.Background(),os.Interrupt,syscall.SIGTERM);defer stop()
 go queue.RunOutbox(ctx,db,url)
 for ctx.Err()==nil{
  if err=queue.RunConsumer(ctx,db,url);err!=nil{log.Printf("consumer reconnect: %v",err)}
  select{case <-ctx.Done():return;case <-time.After(3*time.Second):}
 }
}
