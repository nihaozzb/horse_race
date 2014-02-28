package controllers

import (
	"github.com/astaxie/beego"
	"github.com/horse_race/models"
	"strconv"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	models.DefaultData()
	this.Ctx.WriteString("Hi")
}

func (this *MainController) ErrorCode() {
	this.Data["json"] = map[string]string{
		"1001": "传参错误",
		"1010": "提交的值与数据库中的相同",
		"1011": "传参范围有误",
		"1100": "当前状态为不可押注",
	}
	this.ServeJson()
}

func (this *MainController) GetStatus() {
	config := models.GetStatus()
	if config != nil {
		this.Data["json"] = &config
	} else {
		this.Data["json"] = map[string]string{"error": "Hehe,Jiu Bu Gao Su Ni Shi Na Li Chu Cuo Le"}
	}
	this.ServeJson()
}

func (this *MainController) UpdateStatus() {
	CValue := this.Input().Get("status")
	num := models.UpdateStatus(CValue)
	this.Data["json"] = map[string]int64{"num": num}
	this.ServeJson()
}

func (this *MainController) GetPrize() {
	config := models.GetPrize()
	if config != nil {
		this.Data["json"] = &config
	} else {
		this.Data["json"] = map[string]string{"error": "Hehe,Jiu Bu Gao Su Ni Shi Na Li Chu Cuo Le"}
	}
	this.ServeJson()
}

func (this *MainController) UpdatePrize() {
	CValue := this.Input().Get("prize")
	num := models.UpdatePrize(CValue)
	this.Data["json"] = map[string]int64{"num": num}
	this.ServeJson()
}

func (this *MainController) GetPrizeUrl() {
	config := models.GetPrizeUrl()
	if config != nil {
		this.Data["json"] = &config
	} else {
		this.Data["json"] = map[string]string{"error": "Hehe,Jiu Bu Gao Su Ni Shi Na Li Chu Cuo Le"}
	}
	this.ServeJson()
}

func (this *MainController) UpdatePrizeUrl() {
	CValue := this.Input().Get("prizeurl")
	res := models.UpdatePrizeUrl(CValue)
	this.Data["json"] = map[string]int64{"res": res}
	this.ServeJson()
}

func (this *MainController) GetOneHorse() {
	idStr := this.Ctx.Input.Param(":id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		beego.Error(err)
		return
	}
	horse, _ := models.GetOneHorse(idInt)
	if horse != nil {
		this.Data["json"] = &horse
		this.ServeJson()
	}
	return
}

func (this *MainController) GetAllHorse() {
	this.Data["json"], _ = models.GetAllHorse()
	this.ServeJson()
}

func (this *MainController) UpdateRuslt() {
	winnerId_str := this.Input().Get("winnerid")
	winnerId_int, err := strconv.Atoi(winnerId_str)

	res := 0
	if err != nil {
		beego.Error(nil)
		res = -3
	} else {
		res = models.UpdateRaceRuselt(winnerId_int)
	}

	this.Data["json"] = map[string]int{"res": res}
	this.ServeJson()
}

func (this *MainController) UpdateBet() {
	res := 0

	uId_str := this.Input().Get("uid")
	uId_int, err := strconv.Atoi(uId_str)
	if err != nil {
		beego.Error(err)
		res = -10
	}

	uName := this.Input().Get("uname")
	uAvtar := this.Input().Get("uavtar")

	horseId_str := this.Input().Get("horseid")
	horseId_int, err := strconv.Atoi(horseId_str)
	if err != nil {
		beego.Error(err)
		res = -11
	}

	money_str := this.Input().Get("money")
	money_int, err := strconv.Atoi(money_str)
	if err != nil {
		beego.Error(err)
		res = -12
	}

	if res == 0 {
		res = models.UpdateBet(uId_int, uName, uAvtar, horseId_int, money_int)
	}

	this.Data["json"] = map[string]int{"res": res}
	this.ServeJson()
}

func (this *MainController) GetAllBet() {
	this.Data["json"], _ = models.GetAllBet()
	this.ServeJson()
}

func (this *MainController) GetAllWinlog() {
	var err error
	this.Data["json"], err = models.GetAllWinlog()
	if err != nil {
		beego.Error(err)
		return
	}
	this.ServeJson()
}

func (this *MainController) GetWinlog() {
	uId_str := this.Input().Get("uid")
	uId_int, err := strconv.Atoi(uId_str)
	if err != nil {
		beego.Error(err)
		return
	} else {
		this.Data["json"], _ = models.GetWinlog(uId_int)
	}
	this.ServeJson()
}
