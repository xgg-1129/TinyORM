##挖新坑，实现持久层框架

ORM全称为Object Relational Mapping，对象关系映射，是使用对象将数据持久化到关系型数据库中

###  表     -->  类
### 记 录   --> 对象
### 字 段   --> 属性

实现效果

    type User struct{

    Name string，
    Age  string
    }   


    orm.createTable(&User)
    orm.Save(&User{"Tom", 18})
    var users []User
    orm.Find(&users)


go提供的sql标准库

   

    func main() {
	db, _ := sql.Open("数据库类型名字", "库名字")  //返回个指针，通过这个指针对后续操作
	defer func() { _ = db.Close() }()
	_, _ = db.Exec("DROP TABLE IF EXISTS User;")
	_, _ = db.Exec("CREATE TABLE User(Name text);")
	result, err := db.Exec("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam")
	if err == nil {
		affected, _ := result.RowsAffected()
		log.Println(affected)
	}
	row := db.QueryRow("SELECT Name FROM User LIMIT 1")
	var name string
	if err := row.Scan(&name); err == nil {
		log.Println(name)
	}
}
    