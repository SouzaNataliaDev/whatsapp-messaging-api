package contact
import("context";"database/sql";"errors";"net/http";"regexp";"strings";"github.com/gin-gonic/gin")
var e164=regexp.MustCompile(`^\+[1-9][0-9]{7,14}$`)
type CreateRequest struct{CompanyID int64 `json:"company_id" binding:"required"`;Name string `json:"name" binding:"required"`;Phone string `json:"phone" binding:"required"`;OptedIn bool `json:"opted_in"`}
type Repository struct{DB *sql.DB}
type Service struct{Repo *Repository}
type Handler struct{Service *Service}
func(s *Service)Create(ctx context.Context,r CreateRequest)(int64,error){r.Name=strings.TrimSpace(r.Name);if r.CompanyID<=0||r.Name==""||!e164.MatchString(r.Phone){return 0,errors.New("invalid company_id, name or phone; phone must use E.164 format")};result,err:=s.Repo.DB.ExecContext(ctx,"INSERT INTO contacts(company_id,name,phone,opted_in) VALUES(?,?,?,?)",r.CompanyID,r.Name,r.Phone,r.OptedIn);if err!=nil{return 0,err};return result.LastInsertId()}
func(h *Handler)Create(c *gin.Context){var r CreateRequest;if err:=c.ShouldBindJSON(&r);err!=nil{c.JSON(http.StatusBadRequest,gin.H{"error":"invalid request"});return};id,err:=h.Service.Create(c.Request.Context(),r);if err!=nil{c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()});return};c.JSON(http.StatusCreated,gin.H{"id":id,"company_id":r.CompanyID,"name":r.Name,"phone":r.Phone,"opted_in":r.OptedIn})}
