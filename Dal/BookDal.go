package Dal

import (
	"ActionWeb/Entity"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// 查询所有的书本信息
// 返回值：nil：没有相关数据或者出错，否则返回全部数据
func SelectAllBook() []*Entity.BookInfo {

	var lstBook []*Entity.BookInfo
	db, err := sql.Open("mysql", "root:1@tcp(127.0.0.1:3306)/ryo?charset=utf8")

	if err != nil {
		fmt.Printf("connect err")
		return lstBook
	}

	var sqlCmd string

	sqlCmd = "select tableid,bookcode,bookname,userid from bookinfo"

	rows, err1 := db.Query(sqlCmd)
	if err1 != nil {
		fmt.Println(err1.Error())
		return lstBook
	}

	defer rows.Close()

	var bid int
	var code string
	var bname string
	var uid int

	for rows.Next() {
		if err := rows.Scan(&bid, &code, &bname, &uid); err == nil {
			var book Entity.BookInfo
			book.TableID = bid
			book.BookCode = code
			book.BookName = bname
			book.UserID = uid

			lstBook = append(lstBook, &book)
		}
	}

	return lstBook
}

// 根据ID查询书本信息
// 返回值：nil：没有相关数据或者出错，否则匹配的数据信息
func SelectBookByID(bid int) *Entity.BookInfo {
	var book Entity.BookInfo
	db, err := sql.Open("mysql", "root:1@tcp(127.0.0.1:3306)/ryo?charset=utf8")
	if err != nil {
		fmt.Printf("connect err")
		book.TableID = 0
		return &book
	}

	var sqlCmd string
	sqlCmd = "select a.* from bookinfo a where a.tableid=?"
	rows, err1 := db.Query(sqlCmd, bid)
	if err1 != nil {
		fmt.Println(err1.Error())
		book.TableID = 0
		return &book
	}

	defer rows.Close()

	var sbid int
	var sbname string
	var userid int
	var code string

	for rows.Next() {
		if err := rows.Scan(&sbid, &code, &sbname, &userid); err == nil {

			book.TableID = sbid
			book.BookCode = code
			book.BookName = sbname
			book.UserID = userid

			return &book
		}
	}

	return &book
}

// 添加书本信息
// 返回值：0：存在相同的编码的书本；1：添加成功；909：添加失败
func AddBook(book Entity.BookInfo) int {

	var sqlCmd string = "insert into bookinfo(bookcode,bookname,userid) values(?,?,?)"

	db, err := sql.Open("mysql", "root:1@tcp(127.0.0.1:3306)/ryo?charset=utf8")
	if err != nil {
		fmt.Println(err.Error())
		return 909
	} else {

		// 验证编码是否存在
		var isCode = CodeIsExist(db, book.BookCode)

		// 存在相同的书本编码
		if isCode == 1 {
			return 0
		}

		// 发生错误
		if isCode == 909 {
			return 909
		}

		result, err1 := db.Prepare(sqlCmd)

		defer result.Close()

		if err1 != nil {
			fmt.Println(err1.Error())
			return 909
		}

		rows, errsRows := result.Exec(book.BookCode, book.BookName, book.UserID)

		if errsRows != nil {
			fmt.Printf(errsRows.Error())
			return 909
		}

		count, errsCount := rows.RowsAffected()

		if errsCount != nil {
			fmt.Printf(errsRows.Error())
			return 909
		}

		if count > 0 {
			fmt.Print("添加成功\n")
			return 1
		} else {
			return 909
		}
	}
}

// 删除书本信息
// 返回值：0：不存在书本信息；1：删除成功；909：删除失败
func DeleteBook(bid int) int {
	db, err := sql.Open("mysql", "root:1@tcp(127.0.0.1:3306)/ryo?charset=utf8")
	if err != nil {
		fmt.Println(err.Error())
		return 909
	}

	// 验证书本信息是否存在
	var bookIsExist = IsExist(db, bid)
	if bookIsExist == 0 {
		return 0
	}

	if bookIsExist == 909 {
		return 909
	}

	var sqlCmd = "delete from bookinfo where tableid = ?"
	result, errOper := db.Prepare(sqlCmd)
	if errOper != nil {
		fmt.Println(errOper.Error())
		return 909
	}

	rows, errEffect := result.Exec(bid)
	if errEffect != nil {
		fmt.Println(errOper.Error())
		return 909
	}

	count, errCount := rows.RowsAffected()
	if errCount != nil {
		fmt.Println(errCount.Error())
		return 909
	}

	if count == 1 {
		fmt.Print("删除信息成功\n")
		return 1
	} else {
		return 909
	}
}

// 修改书本信息
// 返回值：0：不存在书本信息；1：更新成功；101：书本编码已经存在；909：更新失败
func UpdateBook(book Entity.BookInfo) int {
	db, err := sql.Open("mysql", "root:1@tcp(127.0.0.1:3306)/ryo?charset=utf8")
	if err != nil {
		fmt.Println(err.Error())
		return 909
	}

	// 验证书本信息是否存在
	var bookIsExit int
	bookIsExit = IsExist(db, book.TableID)

	if bookIsExit == 0 {
		return 0
	}

	if bookIsExit == 909 {
		return 909
	}

	var bCodeIsExist int
	bCodeIsExist = CodeEditIsExist(db, book.BookCode, book.TableID)
	if bCodeIsExist == 1 {
		return 101
	}

	if bookIsExit == 909 {
		return 909
	}

	var sqlCmd = "update bookinfo a set a.bookcode=?,a.bookname=?,a.userid=? where a.tableid = ?"
	result, errOper := db.Prepare(sqlCmd)
	if errOper != nil {
		fmt.Println(errOper.Error())
		return 909
	}

	rows, errEffect := result.Exec(book.BookCode, book.BookName, book.UserID, book.TableID)
	if errEffect != nil {
		fmt.Println(errOper.Error())
		return 909
	}

	count, errCount := rows.RowsAffected()
	if errCount != nil {
		fmt.Println(errCount.Error())
		return 909
	}

	if count == 1 {
		fmt.Print("修改信息成功\n")
		return 1
	} else {
		return 909
	}
}

// 验证书本是否存在
// 返回值：0：不存在；1：存在；909：发生错误
func IsExist(db *sql.DB, tid int) int {
	var sqlCmd string = "select a.tableid from BookInfo a where a.tableid = ?"
	row, err := db.Query(sqlCmd, tid)

	if err != nil {
		fmt.Print("1")
		fmt.Println(err.Error())
		return 909
	}

	defer row.Close()

	var stid int

	for row.Next() {
		errs := row.Scan(&stid)
		if errs != nil {
			fmt.Println(errs.Error())
			fmt.Print("2")
			return 909
		}

		if stid > 0 {
			fmt.Print("3")
			return 1
		}
	}

	return 0
}

// 验证书本编码是否存在
// 返回值：0：不存在；1：存在；909：发生错误
func CodeIsExist(db *sql.DB, code string) int {
	var sqlCmd string = "select a.tableid from BookInfo a where a.bookcode = ?"
	row, err := db.Query(sqlCmd, code)

	if err != nil {
		fmt.Println(err.Error())
		return 909
	}

	defer row.Close()

	var stid int

	for row.Next() {
		errs := row.Scan(&stid)
		if errs != nil {
			fmt.Println(errs.Error())
			return 909
		}

		if stid > 0 {
			return 1
		}
	}

	return 0
}

// 验证书本编码是否存在(更新用)
// 返回值：0：不存在；1：存在；909：发生错误
func CodeEditIsExist(db *sql.DB, code string, tid int) int {
	var sqlCmd string = "select a.tableid from BookInfo a where a.bookcode = ? and a.tableid<>?"
	row, err := db.Query(sqlCmd, code, tid)

	if err != nil {
		fmt.Println(err.Error())
		return 909
	}

	defer row.Close()

	var stid int

	for row.Next() {
		errs := row.Scan(&stid)
		if errs != nil {
			fmt.Println(errs.Error())
			return 909
		}

		if stid > 0 {
			return 1
		}
	}

	return 0
}
