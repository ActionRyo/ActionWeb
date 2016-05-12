package main

import (
	"ActionWeb/Dal"
	"ActionWeb/Entity"
	//"fmt"
	"github.com/astaxie/beego/session"
	//_ "github.com/astaxie/session/providers/memory"
	"html/template"
	//"io/ioutil"
	"net/http"
	"strconv"
	//"time"
	//"fmt"
)

// 登陆处
func loginHandler(w http.ResponseWriter, r *http.Request) {

	// 判断是否点击登陆
	if r.Method == "POST" {

		// 获取表单数据
		r.ParseForm()
		var regAcc, regPwd string
		regAcc = r.FormValue("account")
		regPwd = r.FormValue("password")

		if regAcc != "" && regPwd != "" {
			var oneUser *Entity.UserInfo

			oneUser = Dal.SelectUser(regAcc, regPwd)

			if oneUser.UserID == 0 {
				w.Write([]byte("不存在该用户或发生未知错误"))
				return
			}

			if oneUser.UserID > 0 {
				sess, err := globalSessions.SessionStart(w, r)
				if err != nil {
					return
				}
				w.Header().Set("Content-Type", "text/html")
				sess.Set("UserID", oneUser.UserID)

				//expiration := time.Now()
				//expiration = expiration.AddDate(0, 0, 1)
				//cookie := http.Cookie{Name: "UserID", Value: "asasasa", Path: "/", Expires: expiration}
				//http.SetCookie(w, &cookie)

				http.Redirect(w, r, "/home/main", http.StatusFound)
			}

		} else {
			w.Write([]byte("请输入账户或密码"))
			return
		}
	}

	renderTemplate(w, "Login")
}

// 注销
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := globalSessions.SessionStart(w, r)
	if err != nil {
		http.Redirect(w, r, "/login/login", http.StatusFound)
	}

	uacc := sess.Get("UserID")
	if uacc == nil {
		http.Redirect(w, r, "/login/login", http.StatusFound)
	}

	sess.Delete("UserID")
	http.Redirect(w, r, "/login/login", http.StatusFound)

	renderTemplate(w, "LogOut")
}

// 注册
func registHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		// 获取表单数据
		r.ParseForm()
		var regAcc, regPwd string
		regAcc = r.FormValue("account")
		regPwd = r.FormValue("password")

		if regAcc != "" && regPwd != "" {

			var addInt int
			addInt = Dal.AddUser(regAcc, regPwd)
			if addInt == 1 {
				http.Redirect(w, r, "/login/login", http.StatusFound)
			}
			if addInt == 0 {
				w.Write([]byte("已经存在相同的账户名"))
				return
			} else {
				w.Write([]byte("发生未知错误，注册失败"))
				return
			}
		} else {
			w.Write([]byte("请输入账户名或者密码"))
			return
		}

	}

	renderTemplate(w, "regist")
}

// 主页
func homeHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := globalSessions.SessionStart(w, r)
	if err != nil {
		http.Redirect(w, r, "/login/login", http.StatusFound)
	}

	uacc := sess.Get("UserID")
	if uacc == nil {
		http.Redirect(w, r, "/login/login", http.StatusFound)
	}

	var lstBooks []*Entity.BookInfo
	lstBooks = Dal.SelectAllBook()

	// 加载页面
	t, err := template.ParseFiles("View/main.html")
	if err != nil {
		panic(err)
	}

	if lstBooks == nil {
		t.Execute(w, nil)
	} else {
		t.Execute(w, &Entity.ListBooks{ArrBooks: lstBooks})
	}
}

// 添加书本信息
func addBookHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := globalSessions.SessionStart(w, r)
	if err != nil {
		http.Redirect(w, r, "/login/login", http.StatusFound)
	}

	uacc := sess.Get("UserID")
	if uacc == nil {
		http.Redirect(w, r, "/login/login", http.StatusFound)
	}

	if r.Method == "POST" {
		// 获取表单数据
		r.ParseForm()

		var book Entity.BookInfo
		book.BookCode = r.FormValue("code")
		book.BookName = r.FormValue("name")
		book.UserID, _ = uacc.(int)
		if book.BookCode != "" && book.BookName != "" {

			var addInt int
			addInt = Dal.AddBook(book)
			if addInt > 0 {
				http.Redirect(w, r, "/home/main", http.StatusFound)
			}

			if addInt == 0 {
				w.Write([]byte("已经存在相同的图书编码"))
				return
			} else {
				w.Write([]byte("发生未知错误，添加失败"))
				return
			}

		} else {
			w.Write([]byte("请输入书本信息"))
			return
		}
	}

	renderTemplate(w, "addBook")
}

// 修改图书信息
func editBookHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := globalSessions.SessionStart(w, r)
	if err != nil {
		http.Redirect(w, r, "/login/login", http.StatusFound)
	}

	uacc := sess.Get("UserID")
	if uacc == nil {
		http.Redirect(w, r, "/login/login", http.StatusFound)
	}

	r.ParseForm()

	var book *Entity.BookInfo
	if r.Method != "POST" {
		// 获取url的参数值
		bid, _ := strconv.Atoi(r.Form.Get("bid"))

		book = Dal.SelectBookByID(bid)

		// 验证图书是否
		var errMsg string
		errMsg = "编号为: [" + r.Form.Get("bid") + "] 的书本不存在！"
		if book.TableID == 0 {
			w.Write([]byte(errMsg))
			return
		}
	}

	// 判断是否修改后提交
	if r.Method == "POST" {

		// 获取表单数据
		var books Entity.BookInfo
		books.TableID, _ = strconv.Atoi(r.FormValue("tableid"))
		books.BookCode = r.FormValue("code")
		books.BookName = r.FormValue("name")
		books.UserID, _ = uacc.(int)
		if books.BookCode != "" && books.BookName != "" {

			var addInt int
			addInt = Dal.UpdateBook(books)
			if addInt == 1 {
				http.Redirect(w, r, "/home/main", http.StatusFound)
			}
			if addInt == 0 {
				w.Write([]byte("不存在书本信息"))
				return
			}

			if addInt == 101 {
				w.Write([]byte("书本编码已经存在"))
				return
			} else {
				w.Write([]byte("更新书本信息失败"))
				return
			}
		} else {
			w.Write([]byte("请输入书本信息"))
			return
		}
	}

	t, err := template.ParseFiles("View/EditBook.html")
	if err != nil {
		panic(err)
	}

	if book == nil {
		t.Execute(w, nil)
	} else {
		t.Execute(w, book)
	}
}

// 删除图书信息
func delBookHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := globalSessions.SessionStart(w, r)
	if err != nil {
		http.Redirect(w, r, "/login/login", http.StatusFound)
	}

	uacc := sess.Get("UserID")
	if uacc == nil {
		http.Redirect(w, r, "/login/login", http.StatusFound)
		return
	}

	r.ParseForm()
	var bid int

	// 获取url的参数值
	bid, _ = strconv.Atoi(r.Form.Get("bid"))

	var delInt = Dal.DeleteBook(bid)

	if delInt == 0 {
		w.Write([]byte("不存在该书本信息"))
		return
	}

	if delInt == 1 {
		http.Redirect(w, r, "/home/main", http.StatusFound)
	} else {
		w.Write([]byte("删除失败"))
		return
	}

	t, err := template.ParseFiles("view/DelBook.html")
	if err != nil {
		panic(err)
	}

	t.Execute(w, nil)
}

// 定义全局Session
var globalSessions *session.Manager

//然后在init函数中初始化Session
func init() {
	globalSessions, _ = session.NewManager("memory", `{"cookieName":"gosessionid","gclifetime":3600}`)
	go globalSessions.GC()
}

// 加载页面(无需又初始数据)
func renderTemplate(w http.ResponseWriter, tmpl string) {
	t, err := template.ParseFiles("View/" + tmpl + ".html")
	if err != nil {
		panic(err)
	}

	t.Execute(w, nil)
}

func main() {

	http.HandleFunc("/login/", loginHandler)
	http.HandleFunc("/regist/", registHandler)
	http.HandleFunc("/home/", homeHandler)
	http.HandleFunc("/add/", addBookHandler)
	http.HandleFunc("/edit/", editBookHandler)
	http.HandleFunc("/del/", delBookHandler)
	http.HandleFunc("/out/", logoutHandler)
	//http.HandleFunc("/edit/", editHandler)
	//http.HandleFunc("/login/", saveHandler)
	http.ListenAndServe(":8080", nil)
}
