package campaign
import("context";"database/sql";"errors";"net/http";"strconv";"strings";"github.com/gin-gonic/gin")
type CreateRequest struct{CompanyID int64 `json:"company_id" binding:"required"`;Name string `json:"name" binding:"required"`;TemplateName string `json:"template_name" binding:"required"`}
type Campaign struct{ID int64 `json:"id"`;CompanyID int64 `json:"company_id"`;Name string `json:"name"`;TemplateName string `json:"template_name"`;Status string `json:"status"`}
type Repository struct{DB *sql.DB}
type Service struct{Repo *Repository}
type Handler struct{Service *Service}
func(s *Service)Create(ctx context.Context,r CreateRequest)(int64,error){r.Name=strings.TrimSpace(r.Name);r.TemplateName=strings.TrimSpace(r.TemplateName);if r.CompanyID<=0||r.Name==""||r.TemplateName==""{return 0,errors.New("invalid campaign fields")};result,err:=s.Repo.DB.ExecContext(ctx,"INSERT INTO campaigns(company_id,name,template_name) VALUES(?,?,?)",r.CompanyID,r.Name,r.TemplateName);if err!=nil{return 0,err};return result.LastInsertId()}
func(s *Service)Get(ctx context.Context,id int64)(Campaign,error){var r Campaign;err:=s.Repo.DB.QueryRowContext(ctx,"SELECT id,company_id,name,template_name,status FROM campaigns WHERE id=?",id).Scan(&r.ID,&r.CompanyID,&r.Name,&r.TemplateName,&r.Status);return r,err}
func(h *Handler)Create(c *gin.Context){var r CreateRequest;if err:=c.ShouldBindJSON(&r);err!=nil{c.JSON(http.StatusBadRequest,gin.H{"error":"invalid request"});return};id,err:=h.Service.Create(c.Request.Context(),r);if err!=nil{c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()});return};c.JSON(http.StatusCreated,gin.H{"id":id,"status":"DRAFT"})}
func(h *Handler)Get(c *gin.Context){id,err:=strconv.ParseInt(c.Param("id"),10,64);if err!=nil||id<=0{c.JSON(http.StatusBadRequest,gin.H{"error":"invalid id"});return};result,err:=h.Service.Get(c.Request.Context(),id);if errors.Is(err,sql.ErrNoRows){c.JSON(http.StatusNotFound,gin.H{"error":"campaign not found"});return};if err!=nil{c.JSON(http.StatusInternalServerError,gin.H{"error":"database error"});return};c.JSON(http.StatusOK,result)}
