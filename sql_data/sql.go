package sql_data

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"log"
)

type Item struct {
	Id                 int
	Docker_name        sql.NullString
	Docker_id          sql.NullString
	Docker_create_user sql.NullString
	Docker_password    sql.NullString
	Host               sql.NullString
	Bind               sql.NullString
	Docker_volume      sql.NullString
	Docker_info        sql.NullString
	Create_time        sql.NullString
	Update_time        sql.NullString
	Delete_flag        int
}

func db_insert(docker_name string, docker_id string, docker_create_user string, docker_password string, host string, bind string, docker_volume string, docker_info string, delete_flag int) {

	delete_flag = 0;
	db, err := sqlx.Connect("mysql", "root:root@(localhost:3306)/go_docker")
	if err != nil {
		log.Fatalln(err)
	}
	value := map[string]interface{}{
		"docker_name":        docker_name,
		"docker_id":          docker_id,
		"docker_create_user": docker_create_user,
		"docker_password":    docker_password,
		"host":               host,
		"bind":               bind,
		"docker_volume":      docker_volume,
		"docker_info":        docker_info,
		"delete_flag":        0,
	}

	_, err = db.NamedExec(`INSERT INTO tbl_docker(docker_name,docker_id,docker_create_user,docker_password,host,bind,docker_volume,docker_info,delete_flag)VALUES (:docker_name,:docker_id,:docker_create_user,:docker_password,:host,:bind,:docker_volume,:docker_info,:delete_flag)`, value)
}

func Db() []Item {
	db, err := sqlx.Connect("mysql", "root:root@(localhost:3306)/go_docker")
	if err != nil {
		log.Fatalln(err)
	}
	//fmt.Println(db)
	item := []Item{}
	err = db.Select(&item, "SELECT * FROM tbl_docker")
	//fmt.Println(reflect.TypeOf(item))
	return item
}
