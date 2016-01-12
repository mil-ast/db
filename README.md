# Golang database package

### Пример использования
```go
package main

import(
    "fmt"
    "github.com/mil-ast/db"
)

func main(){
    defer func() {
        if err := recover(); err != nil {
            fmt.Println(err)
        }
    }()
	
	// параметры соединения с базой
	options := db.Options{
		DriverName     : "mysql",
		DataSourceName : "login:password@tcp(127.0.0.1:3306)",
		DbName         : "database_name",
	}

	// устанавливаем соединение
	conn, err := db.CreateConnection(options)
	if err != nil {
		panic(err)
	}
	
	if err = conn.Ping(); err != nil {
		panic(err)
	}
	
	row := conn.QueryRow("SELECT `field`,`field_2` FROM `table` WHERE `id`=?", 1)
	
	var value, value_2 string
	row.Scan(&value, &value_2)
	
	fmt.Println(value, value_2)
	
	// использование соединения в других местах
	my_func()
}

func my_func() {
	// в других пакетах или функциях получаем только инстанс, соединение устанавливать не нужно
	conn, err := db.GetConnection()
	if err != nil {
		panic(err)
	}
	
	row := conn.QueryRow("SELECT `field` FROM `table` WHERE `id`=?", 2)
	var value string
	row.Scan(&value)
	
	fmt.Println(value)
}
```