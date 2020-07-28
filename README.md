dapi
--------------------------------------------------

`dapi` is a Go `sql.Driver` for the [AWS RDS Data API](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/data-api.html)

### Motivation

I wanted to the productivity of using [gorm](https://gorm.io) combined with the utility of
convenience of the RDS Data API.  Looking around, I couldn't find anything that fit
the bill and hence, `dapi`

### QuickStart

`dapi` is intended to work as a standard golang `sql.Driver` and specifically as a driver
usable by `gorm`

```go
func main() {
    var (
        s           = session.Must(session.NewSession(aws.NewConfig()))
        api         = rdsdataservice.New(s)
        driver      = dapi.New(api)
        database    = "the database name"
        secretARN   = "secret arn holding database credentials"
        resourceARN = "resource arn"
        dsn         = fmt.Sprintf("secret=%v resource=%v database=%v", secretARN, resourceARN, database)
        dialect     = "mysql"
    )

	sql.Register(dialect, driver)
	db, err := gorm.Open(dialect, dsn)
    // at this point you can use gorm as you normally would
}
```


### Maturity

This project is very new and should not be used for production. Your mileage may vary.
