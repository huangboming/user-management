// docker run --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=password -d mysql
// MYSQL_URI = root:password@tcp(127.0.0.1:3306)/users
// 需要先在MySQL中创建一个叫user的数据库

package main

func main() {
	Server, _ := InitializeServer()
	Server.LoginDB()
	Server.SetupRoute()
	Server.RunServer()
}
