package Dal

import (
	"ActionWeb/Entity"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// 新增用户
// 返回值：0：用户已经存在；1：添加成功；909：添加失败
func AddUser(userAcc, pwd string) int {
	var sqlCmd string = "insert into userinfo(useraccount,userpwd) values(?,?)"

	db, err := sql.Open("mysql", "root:1@tcp(127.0.0.1:3306)/ryo?charset=utf8")
	if err != nil {
		fmt.Println(err.Error())
		return 909
	} else {

		var isEx int

		// 验证用户名是否存在
		isEx = IsExistUser(db, userAcc)
		if isEx == 1 {
			fmt.Print("该用户名已经存在\n")
			return 0
		}

		if isEx == 909 {
			fmt.Print("发生未知错误\n")
			return 909
		}

		result, err1 := db.Prepare(sqlCmd)

		defer result.Close()

		if err1 != nil {
			fmt.Println(err1.Error())
			return 909
		}

		rows, errsRows := result.Exec(userAcc, pwd)

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

// 修改用户信息
// 返回值：0：用户不存在；1：更新成功；909：更新失败
func UpdateUser(acc, newPwd string) int {
	db, err := sql.Open("mysql", "root:1@tcp(127.0.0.1:3306)/ryo?charset=utf8")
	if err != nil {
		fmt.Println(err.Error())
		return 909
	}

	var isEx int

	// 验证用户名是否存在
	isEx = IsExistUser(db, acc)
	if isEx == 0 {
		fmt.Print("该用户不存在\n")
		return 0
	}

	if isEx == 909 {
		fmt.Print("发生未知错误\n")
		return 909
	}

	var sqlCmd = "update userinfo set userpwd = ? where userAccount = ?"
	result, errOper := db.Prepare(sqlCmd)
	if errOper != nil {
		fmt.Println(errOper.Error())
		return 909
	}

	rows, errEffect := result.Exec(newPwd, acc)
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

// 删除用户信息
// 返回值：0：用户不存在；1：删除成功；909：删除失败
func DeleteUser(acc string) int {
	db, err := sql.Open("mysql", "root:1@tcp(127.0.0.1:3306)/ryo?charset=utf8")
	if err != nil {
		fmt.Println(err.Error())
		return 909
	}

	var isEx int

	// 验证用户名是否存在
	isEx = IsExistUser(db, acc)
	if isEx == 0 {
		fmt.Print("该用户不存在\n")
		return 0
	}

	if isEx == 909 {
		fmt.Print("发生未知错误\n")
		return 909
	}

	var sqlCmd = "delete from userinfo where userAccount = ?"
	result, errOper := db.Prepare(sqlCmd)
	if errOper != nil {
		fmt.Println(errOper.Error())
		return 909
	}

	rows, errEffect := result.Exec(acc)
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

// 查询用户
// 返回值：nil：用户不存在；否则返回用户信息
func SelectUser(acc, pwd string) *Entity.UserInfo {
	db, err := sql.Open("mysql", "root:1@tcp(127.0.0.1:3306)/ryo?charset=utf8")
	var oneUser Entity.UserInfo
	if err != nil {
		fmt.Printf("connect err")
		oneUser.UserID = 0
		return &oneUser
	}

	var sqlCmd string

	sqlCmd = "select userid,useraccount,userpwd from userinfo where useraccount=? and userpwd =?"

	rows, err1 := db.Query(sqlCmd, acc, pwd)
	if err1 != nil {
		fmt.Println(err1.Error())
		oneUser.UserID = 0
		return &oneUser
	}
	defer rows.Close()

	//cols, err2 := rows.Columns()
	//if err2 != nil {
	//	fmt.Println(err2.Error())
	//	return &oneUser
	//}
	//for i := range cols {
	//	fmt.Print(cols[i])
	//	fmt.Print("\t")
	//}

	var userid int
	var uacc, uapwd string
	for rows.Next() {
		if err3 := rows.Scan(&userid, &uacc, &uapwd); err3 == nil {
			oneUser.UserID = userid
			oneUser.UserAccount = uacc
			oneUser.UserPwd = uapwd
		}
	}

	return &oneUser
}

// 验证用户是否存在
// 返回值：0：用户不存在；1：用户存在；909：发生错误
func IsExistUser(db *sql.DB, acc string) int {
	var sqlCmd string = "select userid from userinfo where useraccount = ?"
	row, err := db.Query(sqlCmd, acc)

	if err != nil {
		fmt.Println(err.Error())

		return 909
	}

	defer row.Close()

	var uid int

	for row.Next() {
		errs := row.Scan(&uid)
		if errs != nil {
			fmt.Println(errs.Error())
			return 909
		}

		if uid > 0 {
			return 1
		}
	}

	return 0
}
