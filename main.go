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
	"time"
	"context"
	"io"
	"net/http"
	//"reflect"
	"log"
	"database/sql"
)

/*dokcer関連*/
func docker() {
	ctx := context.Background()
	cl, err := client.NewClient("http://192.168.111.146:2375", client.DefaultVersion, &http.Client{Transport: http.DefaultTransport}, map[string]string{})
	//cl, err := client.NewEnvClient()

	fmt.Println(err)
	fmt.Println(cl.ClientVersion())
	list, err := cl.ImageList(ctx, types.ImageListOptions{All: true})
	if err != nil {
		panic(err)
	}
	for _, image := range list {
		fmt.Println(image.RepoTags)
	}

	export := nat.PortSet{"25565/tcp": struct{}{}}
	portbind := nat.PortMap{
		"25565/tcp": []nat.PortBinding{
			{
				HostPort: "25565",
			},
		},
	}

	conf := &container.Config{
		Image:        "sasenomura/test",
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

	time.Sleep(60 * time.Second)
	cl.ContainerStop(ctx, resp.ID, nil)
	cl.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{
		Force: true,
	})

}

func docker_list() {
	ctx := context.Background()

	//server gave HTTP response to HTTPS client goroutineへの対処
	tr := &http.Transport{}
	cl, err := client.NewClient("http://192.168.111.146:2375", client.DefaultVersion, &http.Client{Transport: tr}, map[string]string{})

	//fmt.Println(err)
	//fmt.Println(cl.ClientVersion())

	list, err := cl.ImageList(ctx, types.ImageListOptions{All: true})
	if err != nil {
		panic(err)
	}

	fmt.Println(list)
	for _, image := range list {
		fmt.Println(image.RepoTags)
	}

	cl.Close()
}

func docker_run(host nat.Port, bind nat.Port) (string) {
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
				HostPort: "8080",
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

	db, err := sqlx.Connect("mysql", "root:root@(localhost:3306)/go_docker")
	if err != nil {
		log.Fatalln(err)
	}
	value := map[string]interface{}{
		"docker_name":        docker_name,
		"docker_id":          "Smuth",
		"docker_create_user": "bensmith@allblacks.nz",
		"docker_password":    "bensmith@allblacks.nz",
		"host":               "bensmith@allblacks.nz",
		"bind":               "bensmith@allblacks.nz",
		"docker_volume":      "bensmith@allblacks.nz",
		"docker_info":        "bensmith@allblacks.nz",
		"delete_flag":        1,
	}

	_, err = db.NamedExec(`INSERT INTO tbl_docker(docker_name,docker_id,docker_create_user,docker_password,host,bind,docker_volume,docker_info,delete_flag)VALUES (:docker_name,:docker_id,:docker_create_user,:docker_password,:host,:bind,:docker_volume,:docker_info,:delete_flag)`, value)
}

func main() {

	
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("template/*.html")),
	}

	fmt.Println(docker_run(nat.Port("80/tcp"),nat.Port("80")))

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

	e.POST("/new", func(i echo.Context) error {
		//docker_run("80/tcp", "")
		fmt.Println(i.Request())
		return i.JSON(http.StatusOK, "ok")
	}).Name = "info"

	// サーバー起動
	e.Start(":80")

}
