package routers

import (
	"github.com/astaxie/beego"
	"github.com/horse_race/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/get_status", &controllers.MainController{}, "get:GetStatus")
	beego.Router("/update_status", &controllers.MainController{}, "post:UpdateStatus")
	beego.Router("/get_prize", &controllers.MainController{}, "get:GetPrize")
	beego.Router("/update_prize", &controllers.MainController{}, "post:UpdatePrize")
	beego.Router("/get_prizeurl", &controllers.MainController{}, "get:GetPrizeUrl")
	beego.Router("/update_prizeurl", &controllers.MainController{}, "post:UpdatePrizeUrl")

	beego.Router("/get_one_horse/:id([0-9]+)", &controllers.MainController{}, "get:GetOneHorse")
	beego.Router("/get_all_horse", &controllers.MainController{}, "get:GetAllHorse")

	beego.Router("/update_winnerid", &controllers.MainController{}, "post:UpdateRuslt")
	beego.Router("/update_bet", &controllers.MainController{}, "post:UpdateBet")

	beego.Router("/get_all_bet", &controllers.MainController{}, "get:GetAllBet")
	beego.Router("/get_all_winlog", &controllers.MainController{}, "get:GetAllWinlog")
	beego.Router("/get_Winlog/:all", &controllers.MainController{}, "get:GetWinlog")
}
