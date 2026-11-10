package company
import("context";"database/sql";"errors";"net/http";"strings";"github.com/gin-gonic/gin")
type Repository struct{DB *sql.DB}
type Service struct{Repo *Repository}
type Handler struct{Service *Service}
type CreateRequest struct{Name string `json:"name" binding:"required"`}
func(s *Service)Create(ctx context.Context,name string)(int64,error){name=strings.TrimSpace(name);if name==""{return 0,errors.New("name is required")};result,err:=s.Repo.DB.ExecContext(ctx,"INSERT INTO companies(name) VALUES(?)",name);if err!=nil{return 0,err};return result.LastInsertId()}
func(h *Handler)Create(c *gin.Context){var req CreateRequest;if err:=c.ShouldBindJSON(&req);err!=nil{c.JSON(http.StatusBadRequest,gin.H{"error":"name is required"});return};id,err:=h.Service.Create(c.Request.Context(),req.Name);if err!=nil{c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()});return};c.JSON(http.StatusCreated,gin.H{"id":id,"name":strings.TrimSpace(req.Name)})}
