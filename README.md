# config

## 初始化

```
type Config struct {
	Wechat struct {
		Appid  string `yaml:"appid"`
		Secret string `yaml:"secret"`
	} `yaml:"wechat"`
}

var conf = new(Config)

func init() {
	goo.LoadConfig(".yaml", conf)
}
```

## 调用示例

```
func main() {
	fmt.Println(conf.Wechat.Appid, conf.Wechat.Secret)
}
```

# web server

```
func main() {
	s := goo.NewServer()
	s.GET("/", goo.Handler(UserList{}))
	s.Run(":18000")
}

type UserList struct {
	Name string `form:"name"`
}

func (ul UserList) DoHandle(c *gin.Context) *goo.Response {
	if err := c.ShouldBind(&ul); err != nil {
		return goo.Error(40010, "参数错误", err.Error())
	}
	return goo.Success(gin.H{"name": ul.Name})
}
```

# db

## 初始化

```
func init() {
    conf := goo.DBConfig{
        Driver: "mysql",
        Master: "test:123456789@tcp(192.168.1.100:3316)/test",
        Slaves: []string{
            "test:123456789@tcp(192.168.1.100:3326)/test",
        },
        LogModel: true,
        MaxIdle:  5,
        MaxOpen:  10,
    }
    goo.DBInit(conf)
}
```

## 调用示例

```
type User struct {
	Id       int64  `xorm:"id"`
	Nickname string `xorm:"nickname"`
}

func (*User) TableName() string {
	return "u_user"
}

func main() {
	u := &User{}
	has, err := goo.DB().Where("id = ?", 1).Get(u)
	if err != nil {
		log.Println(err.Error())
		return
	}
	if !has {
		log.Println("用户信息不存在")
	}
	log.Println(u.Id, u.Nickname)
}
```

# redis

## 初始化

```
func init() {
    conf := goo.RedisConfig{
        Addr: "192.168.1.100:6379", 
        Password: "123456789", 
        DB: 99, 
        Prefix: "tt:"
    }
    goo.RedisInit(conf)
}
```

## 调用示例

```
func main() {
    err := goo.Redis().Set("name", "lqt", 0).Err()
    fmt.Println(err)
    
    name := goo.Redis().Get("name").Val()
    fmt.Println(name)
}
```