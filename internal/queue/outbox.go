package queue
import("context";"database/sql";"log";"time")
func RunOutbox(ctx context.Context,db *sql.DB,url string){ticker:=time.NewTicker(2*time.Second);defer ticker.Stop();for{select{case <-ctx.Done():return;case <-ticker.C:if err:=publishBatch(ctx,db,url);err!=nil{log.Printf("outbox: %v",err)}}}}
func publishBatch(ctx context.Context,db *sql.DB,url string)error{
 conn,err:=Connect(url);if err!=nil{return err};defer conn.Close()
 ch,err:=conn.Channel();if err!=nil{return err};defer ch.Close()
 if err=Declare(ch);err!=nil{return err};if err=ch.Confirm(false);err!=nil{return err}
 rows,err:=db.QueryContext(ctx,"SELECT id,message_id FROM message_outbox WHERE published_at IS NULL ORDER BY id LIMIT 100");if err!=nil{return err}
 type entry struct{id,msg int64};var items []entry
 for rows.Next(){var x entry;if err=rows.Scan(&x.id,&x.msg);err!=nil{rows.Close();return err};items=append(items,x)}
 err=rows.Err();rows.Close();if err!=nil{return err}
 for _,x:=range items{if err=Publish(ctx,ch,x.msg);err!=nil{return err};if !ch.WaitForConfirms(){return context.DeadlineExceeded};if _,err=db.ExecContext(ctx,"UPDATE message_outbox SET published_at=NOW() WHERE id=? AND published_at IS NULL",x.id);err!=nil{return err}}
 return nil
}
