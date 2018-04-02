package main

import (
	"github.com/moby/moby/client"
	"github.com/docker/go-connections/nat"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"

	"html/template"
	"fmt"
	//"time"
	"context"
	"io"
	"net/http"
	//"reflect"
	"log"
	"database/sql"
	"time"
)

/*dokcer関連*/
//func docker() {
//	ctx := context.Background()
//	cl, err := client.NewClient("http://192.168.111.146:2375", client.DefaultVersion, &http.Client{Transport: http.DefaultTransport}, map[string]string{})
//	//cl, err := client.NewEnvClient()
//
//	fmt.Println(err)
//	fmt.Println(cl.ClientVersion())
//	list, err := cl.ImageList(ctx, types.ImageListOptions{All: true})
//	if err != nil {
//		panic(err)
//	}
//	for _, image := range list {
//		fmt.Println(image.RepoTags)
//	}
//
//	export := nat.PortSet{"25565/tcp": struct{}{}}
//	portbind := nat.PortMap{
//		"25565/tcp": []nat.PortBinding{
//			{
//				HostPort: "25565",
//			},
//		},
//	}
//
//	conf := &container.Config{
//		Image:        "sasenomura/test",
//		ExposedPorts: export,
//	}
//
//	host_conf := &container.HostConfig{
//		AutoRemove:   true,
//		PortBindings: portbind,
//	}
//
//	net_conf := &network.NetworkingConfig{
//	}
//
//	resp, err := cl.ContainerCreate(ctx, conf, host_conf, net_conf, "")
//	if err := cl.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
//		panic(err)
//	}
//
//	time.Sleep(60 * time.Second)
//	cl.ContainerStop(ctx, resp.ID, nil)
//	cl.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{
//		Force: true,
//	})
//
//}

//func docker_list() {
//	ctx := context.Background()
//
//	//server gave HTTP response to HTTPS client goroutineへの対処
//	tr := &http.Transport{}
//	cl, err := client.NewClient("http://192.168.111.146:2375", client.DefaultVersion, &http.Client{Transport: tr}, map[string]string{})
//
//	//fmt.Println(err)
//	//fmt.Println(cl.ClientVersion())
//
//	list, err := cl.ImageList(ctx, types.ImageListOptions{All: true})
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Println(list)
//	for _, image := range list {
//		fmt.Println(image.RepoTags)
//	}
//
//	cl.Close()
//}

func docker_run(host nat.Port, bind string) (string) {
	ctx := context.Background()

	//server gave HTTP response to HTTPS client goroutineへの対処
	tr := &http.Transport{}
	cl, err := client.NewClient("http://192.168.111.146:2375", client.DefaultVersion, &http.Client{Transport: tr}, map[string]string{})

	list, err := cl.ImageList(ctx, types.ImageListOptions{All: true})
	if err != nil {
		panic(err)
	}

	fmt.Println(list)
	for _, image := range list {
		fmt.Println(image.RepoTags)
	}

	export := nat.PortSet{host: struct{}{}}
	portbind := nat.PortMap{
		host: []nat.PortBinding{
			{
				HostPort: bind,
			},
		},
	}

	conf := &container.Config{
		Image:        "httpd:latest",
		ExposedPorts: export,
	}

	host_conf := &container.HostConfig{

		AutoRemove:   true,
		PortBindings: portbind,
	}

	net_conf := &network.NetworkingConfig{
	}

	resp, err := cl.ContainerCreate(ctx, conf, host_conf, net_conf, "")
	if err := cl.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	cl.Close()
	return resp.ID
}
func docker_del(docker_id string) {
	sttime, _ := time.ParseDuration("60s")
	ctx := context.Background()

	//server gave HTTP response to HTTPS client goroutineへの対処
	tr := &http.Transport{}
	cl, _ := client.NewClient("http://192.168.111.146:2375", client.DefaultVersion, &http.Client{Transport: tr}, map[string]string{})
	fmt.Println(cl.ContainerStop(ctx, docker_id, &sttime))

	remove_conf := types.ContainerRemoveOptions{
		Force: true,
	}
	cl.ContainerRemove(ctx, docker_id, remove_conf)

}

/*テンプレートエンジン関連*/
type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

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

type Set struct {
	Docler_id       string
	Docker_name     string
	Docker_password string
	Host            string
	Bind            string
	//Docker_volume   string
	Docker_info string
}

func validation(set Set) bool {
	if set.Docker_name == "" || set.Host == "" || set.Bind == "" {
		return false
	} else {
		return true
	}
}

//Sql関連

func db() []Item {
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
func db_del(docker_id string) {
	db, err := sqlx.Connect("mysql", "root:root@(localhost:3306)/go_docker")
	if err != nil {
		log.Fatalln(err)
	}
	id := map[string]interface{}{"docker_id": docker_id}
	_, err = db.NamedExec("DELETE FROM tbl_docker WHERE docker_id = :docker_id", id)
}

func main() {

	//docker_stop("94fceeaaaa3d")
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("template/*.html")),
	}

	//fmt.Println(docker_run(nat.Port("80/tcp"),"8080"))

	e := echo.New()
	e.Static("/static", "assets")

	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Renderer = renderer

	e.GET("/index", func(i echo.Context) error {
		//var item []Item = db()
		//docker_list()
		return i.Render(http.StatusOK, "index.html", map[string]interface{}{
			"name": "Dolly!",
			//"docker": item,
		})
	}).Name = "foobar"

	e.GET("/list", func(i echo.Context) error {
		var item []Item = db()
		return i.JSON(http.StatusOK, item)

	}).Name = "list"

	//テンプレートを返すだけ
	e.GET("/get_new", func(i echo.Context) error {
		return i.Render(http.StatusOK, "new.html", map[string]interface{}{
			"name": "Dolly!",
			//"docker": item,
		})
	}).Name = "get_new"

	//新しく作るとき
	e.POST("/post_new", func(i echo.Context) error {
		//docker_run("80/tcp", "")
		var docker_name string = i.Request().FormValue("name")
		var docker_password string = i.Request().FormValue("password")
		//var ip = i.Request().FormValue("ip")
		var host string = i.Request().FormValue("host")
		var bind string = i.Request().FormValue("bind")
		var info string = i.Request().FormValue("info")

		//構造体のセット
		set := Set{}
		set.Docker_name = docker_name
		set.Docker_password = docker_password
		set.Host = host
		set.Bind = bind
		set.Docker_info = info
		//確認用の関数
		if validation(set) {
			set.Docler_id = docker_run(nat.Port(set.Host), set.Bind)
			db_insert(
				set.Docker_name,
				set.Docler_id,
				"",
				set.Docker_password,
				set.Host,
				set.Bind,
				"",
				set.Docker_info,
				0,
			)
		}

		return i.Redirect(http.StatusMovedPermanently, "/get_new")
	}).Name = "post_new"

	e.GET("/delete/:docker_id", func(i echo.Context) error {
		var id string = i.Param("docker_id")
		docker_del(id)
		db_del(id)
		fmt.Println(id)

		return i.Redirect(http.StatusOK, "/index")
	}).Name = "delete"

	e.GET("/teapot", func(i echo.Context) error {
		return i.String(http.StatusTeapot, "418. I’m a teapot.")
	}).Name = "teapot"

	// サーバー起動
	e.Start(":80")

}
