package main

import (
    "github.com/codegangsta/martini-contrib/render"
    "github.com/codegangsta/martini"
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
    "github.com/coopernurse/gorp"
    "fmt"
    "log"
    "time"
    "os"
)

func main() {
    db, err := sql.Open("mysql", "root@/promdash?parseTime=true")
    if err != nil {
        fmt.Println("Open Error:\n", err.Error())
        return
    }
    defer db.Close()

    dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
    dbmap.AddTableWithName(Dashboard{}, "dashboards").SetKeys(true, "id")
    dbmap.TraceOn("[gorp]", log.New(os.Stdout, "myapp:", log.Lmicroseconds))

    // dash := &Dashboard{Id: 0, Name:"GoDash", Slug:"go-dash", CreatedAt:time.Now(), UpdatedAt:time.Now(), DashboardJSON:"{}"}

    // err = dbmap.Insert(dash)
    // if err != nil {
    //     fmt.Println(err)
    //     return
    // }

    m := martini.Classic()
    m.Use(render.Renderer(render.Options{
        Layout:     "layout",
        Delims: render.Delims{"{[{", "}]}"},
        Extensions: []string{".html"}}))

    m.Get("/", func (r render.Render) {
        var dashboards []Dashboard
        _, err := dbmap.Select(&dashboards, "select * from dashboards order by id")
        if err != nil {
            panic(err)
        }
        fmt.Println(len(dashboards))
        r.HTML(200, "dashboards/index", &dashboards)
    })

    m.Get("/:slug", func (params martini.Params, r render.Render) {
        var dashboard Dashboard
        err := dbmap.SelectOne(&dashboard, "select * from dashboards where slug=?", params["slug"])
        if err != nil {
            // render the to be created 404 page
            r.HTML(404, "dashboards/index", nil)
            return
        }
        r.HTML(200, "dashboards/index", &dashboard)
    })

    m.Run()
}

type Dashboard struct {
    Id             int64  `db:"id"`
    Name           string `db:"name"`
    DashboardJSON  []byte `db:"dashboard_json"`
    Slug           string `db:"slug"`
    CreatedAt      time.Time `db:"created_at"`
    UpdatedAt      time.Time `db:"updated_at"`
}
