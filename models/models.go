package models

import (
	"github.com/Unknwon/com"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path"
	"strconv"
	"time"
)

const (
	// 设置数据库路径
	_DB_NAME = "data/db.db"
	// 设置数据库名称
	_SQLITE3_DRIVER = "sqlite3"
)

// 参数设置
type Configure struct {
	Id     int
	CName  string `orm:"unique"`
	CValue string
}

// 赛马
type Horse struct {
	Id   int
	Name string `orm:"unique"`
	Win  int
	Lose int
}

type Bet struct {
	Id      int
	Uid     int `orm:"unique"`
	Uname   string
	Uavtar  string
	Horseid int
	Money   int
}

type Winlog struct {
	Id      int
	Uid     int
	Uname   string
	Uavtar  string
	Horseid int
	Money   int
	Time    string
}

func RegisterDB() {
	// 检查数据库文件
	if !com.IsExist(_DB_NAME) {
		os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm)
		os.Create(_DB_NAME)
	}

	// 注册模型
	orm.RegisterModel(new(Configure), new(Horse), new(Bet), new(Winlog))
	// 注册驱动（“sqlite3” 属于默认注册，此处代码可省略）
	orm.RegisterDriver(_SQLITE3_DRIVER, orm.DR_Sqlite)
	// 注册默认数据库
	orm.RegisterDataBase("default", _SQLITE3_DRIVER, _DB_NAME, 10)

	// 自动建表
	orm.RunSyncdb("default", false, true)

	DefaultData()
}

func DefaultData() {
	o := orm.NewOrm()
	var configure Configure
	configure.CName = "status"
	configure.CValue = "0" //0未启动，1启动，2已押注

	_, err := o.Insert(&configure)
	if err != nil {
		beego.Error(err)
	} else {
		beego.Trace("初始化status……")
	}

	configure.CName = "prize"
	configure.CValue = ""

	_, err = o.Insert(&configure)
	if err != nil {
		beego.Error(err)
	} else {
		beego.Trace("初始化prize……")
	}

	configure.CName = "prize_url"
	configure.CValue = "http://yun.baozoumanhua.com/Project/RageMaker0/images/0/5.png"

	_, err = o.Insert(&configure)
	if err != nil {
		beego.Error(err)
	} else {
		beego.Trace("初始化prize_url……")
	}

	// 初始化赛马信息
	horse := []Horse{
		{Name: "Oe", Win: 0, Lose: 0},
		{Name: "To", Win: 0, Lose: 0},
		{Name: "Tr", Win: 0, Lose: 0},
		{Name: "Fr", Win: 0, Lose: 0},
		{Name: "Fe", Win: 0, Lose: 0},
		{Name: "Sx", Win: 0, Lose: 0},
	}

	successNums, err := o.InsertMulti(1, horse)
	if err != nil {
		beego.Error(err)
	} else {
		beego.Trace("初始化" + strconv.FormatInt(successNums, 10) + "匹马……")
	}
}

func getValue(CName string) (*Configure, error) {
	o := orm.NewOrm()
	configure := new(Configure)
	qs := o.QueryTable("configure")
	err := qs.Filter("cname", CName).One(configure)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	return configure, nil
}

func updateValue(Name, Value string) int64 {
	o := orm.NewOrm()
	configure := Configure{CName: Name}
	if err := o.Read(&configure, "CName"); err == nil {
		configure.CValue = Value
		if num, err := o.Update(&configure, "CValue"); err != nil {
			beego.Error(err)
			return 0
		} else {
			return num
		}
	} else {
		beego.Error(err)
		return -1
	}
}

func GetStatus() *Configure {
	configure := new(Configure)
	configure, _ = getValue("status")
	return configure
}

func UpdateStatus(CValue string) int64 {
	i, err := strconv.Atoi(CValue)
	if err != nil {
		beego.Error(err)
		return -1001
	}
	if i >= 0 && i <= 2 {

		configure := new(Configure)
		configure = GetStatus()

		if CValue == configure.CValue {
			return -1010
		}

		if i == 0 {
			clearBet()
		}
		num := updateValue("status", CValue)
		return num

	} else {
		return -1011
	}
}

func GetPrize() *Configure {
	configure := new(Configure)
	configure, _ = getValue("prize")
	return configure
}

func UpdatePrize(CValue string) int64 {
	num := updateValue("prize", CValue)
	return num
}

func GetPrizeUrl() *Configure {
	configure := new(Configure)
	configure, _ = getValue("prize_url")
	return configure
}

func UpdatePrizeUrl(CValue string) int64 {
	num := updateValue("prize_url", CValue)
	return num
}

func GetOneHorse(id int) (*Horse, error) {
	o := orm.NewOrm()
	horse := new(Horse)
	qs := o.QueryTable("horse")
	err := qs.Filter("id", id).One(horse)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	return horse, nil
}

func GetAllHorse() ([]*Horse, error) {
	o := orm.NewOrm()

	horses := make([]*Horse, 0)

	qs := o.QueryTable("horse")
	_, err := qs.All(&horses)
	return horses, err
}

func UpdateRaceRuselt(WinnerId int) int {
	configure := new(Configure)
	configure = GetStatus()
	if configure.CValue != "2" {
		return -1101
	}

	o := orm.NewOrm()

	// 事务开始
	err := o.Begin()

	winNum, err := o.QueryTable("horse").Filter("id", WinnerId).Update(orm.Params{
		"win": orm.ColValue(orm.Col_Add, 1),
	})
	if err != nil || winNum != 1 {
		beego.Error(err)
		return -1
	}

	loseNum, err := o.QueryTable("horse").Exclude("id", WinnerId).Update(orm.Params{
		"lose": orm.ColValue(orm.Col_Add, 1),
	})
	if err != nil || loseNum != 5 {
		beego.Error(err)
		return -2
	}

	if winNum == 1 && loseNum == 5 {
		err = o.Commit()
	} else {
		err = o.Rollback()
		return -3
	}

	if err != nil {
		beego.Error(err)
		return 0
	}

	// 将中奖者的记录从bet表移动到winlog表
	bets, err := GetBet(WinnerId)
	// winlogs := make([]*Winlog,0)
	if err != nil {
		beego.Error(err)
		if err != orm.ErrNoRows {
			return -5
		}
	} else {
		var winlog Winlog
		timeStr := time.Now().Format("2006-01-02 15:04:05")
		qs := o.QueryTable("winlog")
		i, _ := qs.PrepareInsert()
		for _, val := range bets {
			winlog.Uid = val.Uid
			winlog.Uname = val.Uname
			winlog.Uavtar = val.Uavtar
			winlog.Horseid = val.Horseid
			winlog.Money = val.Money
			winlog.Time = timeStr
			_, err := i.Insert(&winlog)
			if err != nil {
				beego.Error(err)
				return -4
			}
		}
		i.Close()
	}

	UpdateStatus("0")
	return 1
}

func UpdateBet(uId int, uName string, uAvtar string, horseId int, money int) int {
	configure := new(Configure)
	configure = GetStatus()
	if configure.CValue != "1" {
		return -1100
	}

	o := orm.NewOrm()
	var bet Bet
	bet.Uid = uId

	err := o.Read(&bet, "Uid")
	if err == orm.ErrNoRows {
		bet.Uname = uName
		bet.Uavtar = uAvtar
		bet.Horseid = horseId
		bet.Money = money

		_, err = o.Insert(&bet)
		if err != nil {
			beego.Error(err)
			return -1
		}
	} else if err == nil {
		bet.Money = money
		_, err = o.QueryTable("bet").Filter("Uid", uId).Update(orm.Params{
			"money": orm.ColValue(orm.Col_Add, money),
		})
		if err != nil {
			beego.Error(err)
			return -2
		}
	} else {
		beego.Error(err)
		return -3
	}
	return 1
}

func GetAllBet() ([]*Bet, error) {
	o := orm.NewOrm()

	bets := make([]*Bet, 0)

	qs := o.QueryTable("bet")
	_, err := qs.All(&bets)
	return bets, err
}

func GetBet(horseId int) ([]*Bet, error) {
	o := orm.NewOrm()
	Bets := make([]*Bet, 0)
	qs := o.QueryTable("Bet")
	_, err := qs.Filter("horseId", horseId).All(&Bets)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	return Bets, nil
}

func clearBet() int64 {
	o := orm.NewOrm()
	num, err := o.QueryTable("bet").Exclude("id", "slene").Delete()
	if err != nil {
		beego.Error(err)
		return num
	}
	return -1
}

func GetAllWinlog() ([]*Winlog, error) {
	o := orm.NewOrm()

	Winlogs := make([]*Winlog, 0)

	qs := o.QueryTable("Winlog")
	_, err := qs.All(&Winlogs)
	return Winlogs, err
}

func GetWinlog(uId int) (*Winlog, error) {
	o := orm.NewOrm()
	Winlog := new(Winlog)
	qs := o.QueryTable("Winlog")
	err := qs.Filter("uid", uId).One(Winlog)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	return Winlog, nil
}
