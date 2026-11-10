package queue
import("context";"database/sql";"encoding/json";"errors";"log";"time")
func RunConsumer(ctx context.Context,db *sql.DB,url string)error{
 conn,err:=Connect(url);if err!=nil{return err};defer conn.Close()
 ch,err:=conn.Channel();if err!=nil{return err};defer ch.Close()
 if err=Declare(ch);err!=nil{return err};if err=ch.Qos(5,0,false);err!=nil{return err}
 deliveries,err:=ch.Consume(Name,"",false,false,false,false,nil);if err!=nil{return err}
 for{select{case <-ctx.Done():return;case d,ok:= <-deliveries:
  if !ok{return errors.New("rabbitmq delivery channel closed")}
  var job Job;if err=json.Unmarshal(d.Body,&job);err!=nil||job.MessageID<=0{log.Printf("invalid job: %v",err);_ =d.Reject(false);continue}
  jobCtx,cancel:=context.WithTimeout(ctx,10*time.Second);err=simulate(jobCtx,db,job.MessageID);cancel()
  if err!=nil{log.Printf("job %d: %v",job.MessageID,err);_ =d.Nack(false,true);time.Sleep(time.Second);continue}
  _ =d.Ack(false)
 }}
}
func simulate(ctx context.Context,db *sql.DB,id int64)error{result,err:=db.ExecContext(ctx,"UPDATE messages SET status='SENT',provider_message_id=CONCAT('mock-',id) WHERE id=? AND status='QUEUED'",id);if err!=nil{return err};n,err:=result.RowsAffected();if err!=nil{return err};if n>0{log.Printf("MOCK WhatsApp message sent: id=%d",id)};return nil}
