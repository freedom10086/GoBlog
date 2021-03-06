package router

import (
	"database/sql"
	"fmt"
	"goBlog/logger"
	"goBlog/repository"
	"io"
	"net/http"
	"sync"
	"time"
)

type RegisterHandler struct {
	BaseHandler
}

func (h *RegisterHandler) DoAuth(method int, r *http.Request) error {
	if method == MethodGet || method == MethodPost {
		return nil
	}
	return h.BaseHandler.DoAuth(method, r)
}

type CompeteRegData struct {
	BasePageData
	PostUrl  string
	Token    string
	Username string
	Email    string
}

//token null /regiest ->登陆页面
//登陆页面 ->dopost -> 发邮件 -> 点击连接 -> user.doPost 插入数据库
func (h *RegisterHandler) DoGet(w http.ResponseWriter, r *http.Request) {
	mode := r.FormValue("mod")
	switch mode {
	case "done":
		token := r.FormValue("token")
		if t, err := repository.ValidRegToken(token, config.SecretKey); err == nil {
			// check if already register
			if user, err := repository.GetUserByEmail(t.Email); err == sql.ErrNoRows {
				Template(w, &TemplateData{
					Title: "完成注册",
					Css:   []string{"style.css"},
					Js:    []string{"base.js", "particles.js"},
					Data: &CompeteRegData{
						BasePageData: BasePageData{0},
						PostUrl:      "/users",
						Token:        token,
						Email:        t.Email,
						Username:     t.Username},
				},
					"page.gohtml", "register-done.gohtml")
			} else {
				if err != nil {
					// error happens
					InternalError(w, r, err)
				} else {
					logger.I("already reg %v %v", user, err)
					// user already exist
					io.WriteString(w, fmt.Sprintf("already reg %s %s", user.Username, user.Password))
				}
			}
		} else {
			logger.E("invalid token failed %s %v", token, err)
			Unauthorized(w, r, err.Error())
		}
		return
	case "checkUsername":
		if u := r.FormValue("username"); !repository.CheckUsername(u) {
			Error(w, u+"用户名不可用", 400)
		} else {
			io.WriteString(w, "ok")
		}
		return
	case "checkEmail":
		if e := r.FormValue("email"); !repository.CheckEmail(e) {
			Error(w, e+"邮箱不可用", 400)
		} else {
			io.WriteString(w, "ok")
		}
		return
	default:
		Template(w, &TemplateData{
			Css:  []string{"style.css"},
			Js:   []string{"base.js", "particles.js"},
			Data: &BasePageData{TabIndex: 0},
		}, "page.gohtml", "register.gohtml")
	}
}

//填好用户名 邮箱假注册
//真注册链接在邮件
func (h *RegisterHandler) DoPost(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	email := r.PostFormValue("email")
	if username == "" || email == "" {
		BadParameter(w, r, "用户名或者邮箱不能为空")
		return
	}

	if !repository.CheckUsername(username) {
		Error(w, username+"用户名不可用", 400)
		return
	}

	if !repository.CheckEmail(email) {
		Error(w, email+"邮件不可用", 400)
		return
	}

	token := repository.GenRegToken(username, email, config.SecretKey, time.Minute*30)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		content := fmt.Sprintf("你好：%s!<br><b>请点击以下链接激活你的%s账号，完成注册</b>"+
			"<br>验证邮箱:<a href=\"%s\">%s</a><br><b>注意请在%d分钟内完成操作</b>",
			username,
			config.SiteName,
			"https://"+config.SiteIpAddr+config.SitePortSSL+"/register?mod=done&token="+token,
			email,
			30,
		)
		err := repository.SendHtmlMail(email, "验证你的注册邮件-"+config.SiteName, content)
		if err != nil {
			logger.E("error send email to %s %v", email, err)
		}
	}()
	wg.Wait()
	logger.I("register %s %s", username, email)
	s := fmt.Sprintf("注册确认链接已经发送到你的邮箱:%s,请在%d分钟内完成验证", email, 30)
	io.WriteString(w, s)
}
