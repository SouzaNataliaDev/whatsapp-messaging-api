package campaign
import("context";"database/sql";"errors";"net/http";"strconv";"github.com/gin-gonic/gin")
func(s *Service)Send(ctx context.Context,id int64)(int64,error){
 tx,err:=s.Repo.DB.BeginTx(ctx,nil);if err!=nil{return 0,err};defer tx.Rollback()
 var companyID int64;var status string
 err=tx.QueryRowContext(ctx,"SELECT company_id,status FROM campaigns WHERE id=? FOR UPDATE",id).Scan(&companyID,&status);if err!=nil{return 0,err}
 if status!="DRAFT"{return 0,errors.New("campaign already started")}
 _,err=tx.ExecContext(ctx,"INSERT INTO messages(campaign_id,contact_id,status) SELECT ?,id,'QUEUED' FROM contacts WHERE company_id=? AND opted_in=1",id,companyID);if err!=nil{return 0,err}
 _,err=tx.ExecContext(ctx,"INSERT INTO message_outbox(message_id) SELECT id FROM messages WHERE campaign_id=?",id);if err!=nil{return 0,err}
 var count int64;if err=tx.QueryRowContext(ctx,"SELECT COUNT(*) FROM messages WHERE campaign_id=?",id).Scan(&count);err!=nil{return 0,err}
 state:="COMPLETED";if count>0{state="QUEUED"}
 if _,err=tx.ExecContext(ctx,"UPDATE campaigns SET status=? WHERE id=?",state,id);err!=nil{return 0,err}
 return count,tx.Commit()
}
func(h *Handler)Send(c *gin.Context){id,err:=strconv.ParseInt(c.Param("id"),10,64);if err!=nil||id<=0{c.JSON(400,gin.H{"error":"invalid id"});return};count,err:=h.Service.Send(c.Request.Context(),id);if errors.Is(err,sql.ErrNoRows){c.JSON(404,gin.H{"error":"campaign not found"});return};if err!=nil{if err.Error()=="campaign already started"{c.JSON(409,gin.H{"error":err.Error()})}else{c.JSON(500,gin.H{"error":"could not queue campaign"})};return};c.JSON(http.StatusAccepted,gin.H{"campaign_id":id,"queued":count,"mode":"mock"})}
func(h *Handler)Stats(c *gin.Context){id,err:=strconv.ParseInt(c.Param("id"),10,64);if err!=nil||id<=0{c.JSON(400,gin.H{"error":"invalid id"});return};var total,sent,queued int64;err=h.Service.Repo.DB.QueryRowContext(c.Request.Context(),"SELECT COUNT(*),COALESCE(SUM(status='SENT'),0),COALESCE(SUM(status='QUEUED'),0) FROM messages WHERE campaign_id=?",id).Scan(&total,&sent,&queued);if err!=nil{c.JSON(500,gin.H{"error":"database error"});return};c.JSON(200,gin.H{"campaign_id":id,"total":total,"sent_mock":sent,"queued":queued})}
