package queue
import("context";"encoding/json";"time";amqp "github.com/rabbitmq/amqp091-go")
const Name="whatsapp.messages"
type Job struct{MessageID int64 `json:"message_id"`}
func Connect(url string)(*amqp.Connection,error){return amqp.Dial(url)}
func Declare(ch *amqp.Channel)error{_,err:=ch.QueueDeclare(Name,true,false,false,false,nil);return err}
func Publish(ctx context.Context,ch *amqp.Channel,id int64)error{b,_:=json.Marshal(Job{MessageID:id});return ch.PublishWithContext(ctx,"",Name,false,false,amqp.Publishing{ContentType:"application/json",DeliveryMode:amqp.Persistent,Body:b,MessageId:time.Now().Format(time.RFC3339Nano)})}
