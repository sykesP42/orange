package handlers

import (
	"database/sql"
	"net/http"
	"server-go/internal/config"
	"server-go/internal/database"
	"server-go/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB  *sql.DB
	Cfg *config.Config
}

func NewAuthHandler(db *sql.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{DB: db, Cfg: cfg}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		RealName string `json:"real_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "请输入用户名和密码"})
		return
	}

	var exists int
	h.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", req.Username).Scan(&exists)
	if exists > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "用户名已存在"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "密码加密失败"})
		return
	}

	realName := req.RealName
	if realName == "" {
		realName = req.Username
	}

	result, err := h.DB.Exec(`INSERT INTO users (username, password, real_name, role, status) 
		VALUES (?, ?, ?, ?, ?)`,
		req.Username, string(hashedPassword), realName, "user", "pending")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": err.Error()})
		return
	}

	id, _ := result.LastInsertId()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       id,
		"username": req.Username,
		"role":     "user",
		"realName": realName,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.Cfg.JWTSecret))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "生成token失败"})
		return
	}

	database.AddLog("注册", realName, realName, "新用户【"+realName+"】注册")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "注册成功，等待管理员审核",
		"data": gin.H{
			"token":    tokenString,
			"id":       id,
			"username": req.Username,
			"real_name": realName,
			"role":     "user",
			"status":   "pending",
		},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "请输入用户名和密码"})
		return
	}

	var user models.User
	var realName, avatar, email, phone, bio, status sql.NullString
	var teamID sql.NullInt64
	var createTime, updateTime string

	err := h.DB.QueryRow(`
		SELECT id, username, password, real_name, role, avatar, email, phone, bio, status, team_id, create_time, update_time
		FROM users WHERE username = ?`, req.Username).Scan(
		&user.ID, &user.Username, &user.Password, &realName, &user.Role,
		&avatar, &email, &phone, &bio, &status, &teamID, &createTime, &updateTime,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, gin.H{"code": 401, "message": "用户名或密码错误"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": err.Error()})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "用户名或密码错误"})
		return
	}

	if status.Valid && status.String == "pending" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "账号待审核，请联系管理员"})
		return
	}

	if status.Valid && status.String != "active" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "账号已禁用"})
		return
	}

	user.RealName = realName.String
	user.Avatar = avatar.String
	user.Email = email.String
	user.Phone = phone.String
	user.Bio = bio.String
	user.Status = status.String
	user.CreateTime = createTime
	user.UpdateTime = updateTime
	if teamID.Valid {
		user.TeamID = &teamID.Int64
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
		"realName": user.RealName,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.Cfg.JWTSecret))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "生成token失败"})
		return
	}

	var teamName, teamColor sql.NullString
	if user.TeamID != nil {
		h.DB.QueryRow("SELECT name, color FROM teams WHERE id = ?", *user.TeamID).Scan(&teamName, &teamColor)
	}

	responseData := gin.H{
		"token":      tokenString,
		"id":         user.ID,
		"username":   user.Username,
		"real_name":  user.RealName,
		"role":       user.Role,
		"avatar":     user.Avatar,
		"team_id":    user.TeamID,
		"team_name":  teamName.String,
		"team_color": teamColor.String,
		"status":     user.Status,
	}

	database.AddLog("登录", user.RealName, user.RealName, "用户登录")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data":    responseData,
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	var user models.User
	var realName, avatar, email, phone, bio, status sql.NullString
	var teamID sql.NullInt64
	var createTime, updateTime string

	err := h.DB.QueryRow(`
		SELECT id, username, password, real_name, role, avatar, email, phone, bio, status, team_id, create_time, update_time
		FROM users WHERE id = ?`, userID).Scan(
		&user.ID, &user.Username, &user.Password, &realName, &user.Role,
		&avatar, &email, &phone, &bio, &status, &teamID, &createTime, &updateTime,
	)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "用户不存在"})
		return
	}

	user.RealName = realName.String
	user.Avatar = avatar.String
	user.Email = email.String
	user.Phone = phone.String
	user.Bio = bio.String
	user.Status = status.String
	user.CreateTime = createTime
	user.UpdateTime = updateTime
	if teamID.Valid {
		user.TeamID = &teamID.Int64
	}

	var teamName, teamColor sql.NullString
	if user.TeamID != nil {
		h.DB.QueryRow("SELECT name, color FROM teams WHERE id = ?", *user.TeamID).Scan(&teamName, &teamColor)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"real_name":  user.RealName,
			"role":       user.Role,
			"avatar":     user.Avatar,
			"email":      user.Email,
			"phone":      user.Phone,
			"bio":        user.Bio,
			"team_id":    user.TeamID,
			"team_name":  teamName.String,
			"team_color": teamColor.String,
			"status":     user.Status,
		},
	})
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	realName, _ := c.Get("realName")

	var req struct {
		RealName string `json:"real_name"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Bio      string `json:"bio"`
		Avatar   string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": err.Error()})
		return
	}

	updates := []string{}
	args := []interface{}{}

	if req.RealName != "" {
		updates = append(updates, "real_name = ?")
		args = append(args, req.RealName)
	}
	if req.Email != "" {
		updates = append(updates, "email = ?")
		args = append(args, req.Email)
	}
	if req.Phone != "" {
		updates = append(updates, "phone = ?")
		args = append(args, req.Phone)
	}
	if req.Bio != "" {
		updates = append(updates, "bio = ?")
		args = append(args, req.Bio)
	}
	if req.Avatar != "" {
		updates = append(updates, "avatar = ?")
		args = append(args, req.Avatar)
	}

	if len(updates) > 0 {
		updates = append(updates, "update_time = CURRENT_TIMESTAMP")
		args = append(args, userID)
		query := "UPDATE users SET " + string(join(updates, ", ")) + " WHERE id = ?"
		h.DB.Exec(query, args...)
	}

	database.AddLog("更新资料", realName.(string), realName.(string), "更新个人资料")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("userID")
	realName, _ := c.Get("realName")

	var req struct {
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "请输入原密码和新密码"})
		return
	}

	var currentPassword string
	h.DB.QueryRow("SELECT password FROM users WHERE id = ?", userID).Scan(&currentPassword)

	if err := bcrypt.CompareHashAndPassword([]byte(currentPassword), []byte(req.OldPassword)); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "原密码错误"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "密码加密失败"})
		return
	}

	h.DB.Exec("UPDATE users SET password = ?, update_time = CURRENT_TIMESTAMP WHERE id = ?", hashedPassword, userID)

	database.AddLog("修改密码", realName.(string), realName.(string), "修改登录密码")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "密码修改成功"})
}

func (h *AuthHandler) GetUsers(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	offset := (page - 1) * pageSize
	statusFilter := c.Query("status")
	roleFilter := c.Query("role")
	teamIDFilter := c.Query("team_id")

	query := `SELECT id, username, real_name, role, avatar, bio, email, phone, status, team_id, create_time, approved_by, approved_time 
		FROM users WHERE is_deleted = 0`
	countQuery := `SELECT COUNT(*) FROM users WHERE is_deleted = 0`
	args := []interface{}{}
	countArgs := []interface{}{}

	if statusFilter != "" {
		query += " AND status = ?"
		countQuery += " AND status = ?"
		args = append(args, statusFilter)
		countArgs = append(countArgs, statusFilter)
	}
	if roleFilter != "" {
		query += " AND role = ?"
		countQuery += " AND role = ?"
		args = append(args, roleFilter)
		countArgs = append(countArgs, roleFilter)
	}
	if teamIDFilter != "" {
		query += " AND team_id = ?"
		countQuery += " AND team_id = ?"
		teamIDInt, _ := strconv.Atoi(teamIDFilter)
		args = append(args, teamIDInt)
		countArgs = append(countArgs, teamIDInt)
	}

	query += " ORDER BY create_time DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	var total int
	h.DB.QueryRow(countQuery, countArgs...).Scan(&total)

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": err.Error()})
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var username, realName, avatar, bio, email, phone, status sql.NullString
		var role string
		var teamID sql.NullInt64
		var createTime, approvedBy, approvedTime sql.NullString
		rows.Scan(&id, &username, &realName, &role, &avatar, &bio, &email, &phone, &status, &teamID, &createTime, &approvedBy, &approvedTime)

		var teamName, teamColor sql.NullString
		if teamID.Valid {
			h.DB.QueryRow("SELECT name, color FROM teams WHERE id = ?", teamID.Int64).Scan(&teamName, &teamColor)
		}

		user := map[string]interface{}{
			"id":            id,
			"username":      username.String,
			"real_name":     realName.String,
			"role":          role,
			"avatar":        avatar.String,
			"bio":           bio.String,
			"email":         email.String,
			"phone":         phone.String,
			"status":        status.String,
			"team_id":       nil,
			"team_name":     teamName.String,
			"team_color":    teamColor.String,
			"create_time":   createTime.String,
			"approved_by":   approvedBy.String,
			"approved_time": approvedTime.String,
		}
		if teamID.Valid {
			user["team_id"] = teamID.Int64
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":     users,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

func (h *AuthHandler) UpdateUser(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		RealName string `json:"real_name"`
		Role     string `json:"role"`
		TeamID   *int64 `json:"team_id"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Status   string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": err.Error()})
		return
	}

	var oldRealName, oldStatus, oldRole sql.NullString
	var oldTeamID sql.NullInt64
	h.DB.QueryRow("SELECT real_name, status, role, team_id FROM users WHERE id = ?", id).Scan(&oldRealName, &oldStatus, &oldRole, &oldTeamID)

	updates := []string{}
	args := []interface{}{}

	if req.RealName != "" {
		updates = append(updates, "real_name = ?")
		args = append(args, req.RealName)
	}
	if req.Role != "" {
		updates = append(updates, "role = ?")
		args = append(args, req.Role)
	}
	if req.TeamID != nil {
		updates = append(updates, "team_id = ?")
		args = append(args, req.TeamID)
	}
	if req.Email != "" {
		updates = append(updates, "email = ?")
		args = append(args, req.Email)
	}
	if req.Phone != "" {
		updates = append(updates, "phone = ?")
		args = append(args, req.Phone)
	}
	if req.Status != "" {
		updates = append(updates, "status = ?")
		args = append(args, req.Status)
	}

	if len(updates) > 0 {
		updates = append(updates, "update_time = CURRENT_TIMESTAMP")
		args = append(args, id)
		query := "UPDATE users SET " + string(join(updates, ", ")) + " WHERE id = ?"
		h.DB.Exec(query, args...)
	}

	adminRealName, _ := c.Get("realName")
	database.AddLog("更新用户", oldRealName.String, adminRealName.(string), "更新用户【"+oldRealName.String+"】信息")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

func (h *AuthHandler) ApproveUser(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	var realName, username string
	h.DB.QueryRow("SELECT real_name, username FROM users WHERE id = ?", id).Scan(&realName, &username)

	if realName == "" {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "用户不存在"})
		return
	}

	adminRealName, _ := c.Get("realName")
	now := time.Now().Format(time.RFC3339)

	h.DB.Exec(`UPDATE users SET status = ?, approved_by = ?, approved_time = ?, update_time = ? WHERE id = ?`,
		"active", adminRealName, now, now, id)

	h.DB.Exec(`INSERT INTO notifications (user_id, type, title, content, priority, from_user, create_time)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id, "account_approved", "🎉 账号审核通过",
		"恭喜！您的账号【"+username+"】已通过管理员审核，现在可以正常使用了。",
		"important", adminRealName, now)

	database.AddLog("审核用户", realName, adminRealName.(string), "审核通过用户【"+realName+"】")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "审核通过"})
}

func (h *AuthHandler) DeleteUser(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	var realName string
	h.DB.QueryRow("SELECT real_name FROM users WHERE id = ?", id).Scan(&realName)

	if realName == "" {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "用户不存在"})
		return
	}

	adminRealName, _ := c.Get("realName")
	h.DB.Exec("UPDATE users SET is_deleted = 1, update_time = CURRENT_TIMESTAMP WHERE id = ?", id)
	database.AddLog("删除用户", realName, adminRealName.(string), "删除用户【"+realName+"】")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	var realName, username string
	h.DB.QueryRow("SELECT real_name, username FROM users WHERE id = ?", id).Scan(&realName, &username)

	if realName == "" {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "用户不存在"})
		return
	}

	defaultPassword := "123456"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "密码加密失败"})
		return
	}

	h.DB.Exec("UPDATE users SET password = ?, update_time = CURRENT_TIMESTAMP WHERE id = ?", hashedPassword, id)

	adminRealName, _ := c.Get("realName")
	database.AddLog("重置密码", realName, adminRealName.(string), "重置用户【"+realName+"】密码")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "密码已重置为: 123456"})
}

func join(slice []string, sep string) string {
	switch len(slice) {
	case 0:
		return ""
	case 1:
		return slice[0]
	}
	result := slice[0]
	for i := 1; i < len(slice); i++ {
		result += sep + slice[i]
	}
	return result
}
