package v1

import (
	"encoding/json"
	"forklift-bat-backend/middleware"
	"forklift-bat-backend/model"
	"forklift-bat-backend/util"
	"forklift-bat-backend/validate"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const UserAuthURL = "https://oa.app.swirecocacola.com/OAuth/Work/GetUserId"

type UserID struct {
	UserID string `json:"UserId"`
}

type returnData struct {
	Token string `json:"token"`
	ID    string `json:"id"`
}

func UserAuth(c *gin.Context) {
	code := c.Query("code")
	authbody := validate.UserAuth{
		Code: code,
	}
	err := validate.UserAuthValidator(&authbody)
	GinUtil := util.Gin{Context: c}
	if err != nil {
		GinUtil.Response(400, 0, err.Error(), "")
		return
	}
	url := UserAuthURL + "?code=" + code
	resp, err := http.Get(url)
	if err != nil {
		GinUtil.Response(500, 0, err.Error(), "")
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		GinUtil.Response(500, 0, err.Error(), "")
		return
	}
	restring := string(body)
	bodyjson := &UserID{}
	json.Unmarshal([]byte(restring), &bodyjson)
	if bodyjson.UserID == "" {
		GinUtil.Response(400, 0, "User not found", "")
		return
	}
	user := model.User{}.GetUser(bodyjson.UserID)
	if len(user) == 1 {
		token, err := middleware.GenToken(bodyjson.UserID)
		if err != nil {
			GinUtil.Response(500, 0, "Internal Error", "")
			return
		}
		data := &returnData{
			Token: token,
			ID:    bodyjson.UserID,
		}
		GinUtil.Response(200, 1, "success", data)
	} else {
		GinUtil.Response(403, 0, "access denied", "")
	}
}

func CheckJWT(c *gin.Context) {
	token := c.Query("token")
	id := c.Query("id")
	claims := &middleware.CustomClaims{}
	claims, err := middleware.ParseToken(token)
	GinUtil := util.Gin{Context: c}
	if err != nil {
		GinUtil.Response(400, 0, err.Error(), "")
		return
	}
	if id == claims.UserID {
		GinUtil.Response(200, 1, "Success", "")
	} else {
		GinUtil.Response(400, 0, "ID doesn`t match", "")
	}
}

func GetForklift(c *gin.Context) {
	warehouse := c.Query("dc")
	forkliftcat := c.Query("cat")
	forksbody := validate.GetForksBody{
		Warehouse: warehouse,
		Category:  forkliftcat,
	}
	err := validate.GetForksBodyValidator(&forksbody)
	GinUtil := util.Gin{Context: c}
	if err == nil {
		forklifts := model.Forklift{}.GetForkliftNo(warehouse, forkliftcat)
		fullnumbers := []map[string]string{} // blank slice
		for _, forklift := range forklifts {
			warehouseget := model.Warehouse{}.GetWarehouse(forklift.Warehouse)
			if len(warehouseget) == 0 {
				GinUtil.Response(500, 0, "Internal Error", "")
				return
			}
			forkcat := model.ForkCat{}.GetForkCat(forklift.Category)
			if len(forkcat) == 0 {
				GinUtil.Response(500, 0, "Internal Error", "")
				return
			}
			resmap := make(map[string]string)
			resmap["warehouse"] = forklift.Warehouse
			resmap["number"] = forklift.Number
			resmap["category"] = forklift.Category
			resmap["dcname"] = warehouseget[0].Name
			resmap["forkcat"] = forkcat[0].Name
			resmap["forkcatno"] = forkcat[0].Number
			fullnumbers = append(fullnumbers, resmap)
		}
		GinUtil.Response(200, 1, "", fullnumbers)
	} else {
		GinUtil.Response(400, 0, err.Error(), "")
	}
}

func AddForklift(c *gin.Context) {
	warehouse := c.PostForm("dc")
	category := c.PostForm("cat")
	no := c.PostForm("no")
	forkbody := validate.InsertForksBody{
		Warehouse: warehouse,
		Category:  category,
		No:        no,
	}
	err := validate.InsertForksBodyValidator(&forkbody)
	GinUtil := util.Gin{Context: c}
	if err == nil {
		warehouseget := model.Warehouse{}.GetWarehouse(warehouse)
		if len(warehouseget) == 0 {
			GinUtil.Response(400, 0, "No such warehouse", "")
			panic("No such warehouse")
		}
		forkbody := model.Forklift{
			Warehouse: warehouse,
			Category:  category,
			Number:    no,
			FullNo:    warehouse + category + no,
		}
		result := model.InsertForks(&forkbody)
		if result == nil {
			GinUtil.Response(200, 1, warehouse, category)
		} else {
			GinUtil.Response(500, 0, result.Error(), "")
		}
	} else {
		GinUtil.Response(400, 0, err.Error(), "")
	}
}

func AddWarehouse(c *gin.Context) {
	name := c.PostForm("name")
	no := c.PostForm("no")
	dcbody := validate.InsertWarehouseBody{
		Name:   name,
		Number: no,
	}
	err := validate.InsertWarehouseValidator(&dcbody)
	GinUtil := util.Gin{Context: c}
	if err == nil {
		dcinsert := model.Warehouse{
			Number: no,
			Name:   name,
		}
		result := model.InsertWarehouse(&dcinsert)
		if result == nil {
			GinUtil.Response(200, 1, "success", "")
		} else {
			GinUtil.Response(500, 0, result.Error(), "")
		}
	} else {
		GinUtil.Response(400, 0, err.Error(), "")
	}
}

func AddForkCat(c *gin.Context) {
	name := c.PostForm("name")
	no := c.PostForm("no")
	forkbody := validate.InsertForkCatBody{
		Name:   name,
		Number: no,
	}
	err := validate.InsertForkCatValidator(&forkbody)
	GinUtil := util.Gin{Context: c}
	if err == nil {
		forkinsert := model.ForkCat{
			Number: no,
			Name:   name,
		}
		result := model.InsertForkCat(&forkinsert)
		if result == nil {
			GinUtil.Response(200, 1, "success", "")
		} else {
			GinUtil.Response(500, 0, result.Error(), "")
		}
	} else {
		GinUtil.Response(400, 0, err.Error(), "")
	}
}

func SwitchBattery(c *gin.Context) {
	bat := c.PostForm("bat")
	fork := c.PostForm("fork") //forklift number instead of type
	dc := c.PostForm("dc")
	inout := c.PostForm("inout")
	user := c.PostForm("user")
	forkcat := c.PostForm("forkcat")
	GinUtil := util.Gin{Context: c}
	validBody := validate.SwitchBatBody{
		No:        bat,
		ForkNo:    fork,
		Warehouse: dc,
		Switch:    inout,
		UserID:    user,
	}
	validResult := validate.SwichBatValidator(&validBody)
	if validResult != nil {
		GinUtil.Response(400, 0, validResult.Error(), "")
		return
	}
	status := map[string]string{
		"on":  "已装入",
		"off": "已卸下",
	}
	batGet := model.Battery{}.GetBatStatus(bat, forkcat, dc)
	if len(batGet) == 0 {
		GinUtil.Response(400, 0, "没有这个电瓶", "")
		return
	}
	if batGet[0].Status == inout {
		GinUtil.Response(400, 0, "错误！电瓶"+status[inout], "")
		return
	}
	if inout == "off" && batGet[0].LastSeenAt != fork && batGet[0].LastSeenAt != "" {
		GinUtil.Response(400, 0, "错误！此电瓶上次未装在此叉车上", "")
		return
	}
	// for batteryswitch table, full number is warehouse + forklift number + battery number
	fullno := dc + fork + bat
	// for battery table, full number is warehouse + forklift type number + battery number
	fullnocat := dc + forkcat + bat
	lastBatSwitch := model.SwitchBattery{}.GetBatSwitch(fullno)
	updateTime := time.Unix(0, 0)
	if len(lastBatSwitch) != 0 {
		updateTime = lastBatSwitch[0].UpdatedAt
	}
	if time.Now().Sub(updateTime).Seconds() < float64(util.TimeInterval) {
		GinUtil.Response(400, 0, "两次提交间隔时间太短", "")
		return
	}
	addBatSwichBody := model.SwitchBattery{
		FullNo:     fullno,
		BatNumber:  bat,
		ForkNumber: fork,
		UserID:     user,
		Operation:  inout,
	}
	result := model.UpdateBatStatus(&model.Battery{Status: inout, LastSeenAt: fork}, &addBatSwichBody, fullno, fullnocat)
	if result != nil {
		GinUtil.Response(400, 0, result.Error(), "")
		return
	}
	GinUtil.Response(200, 1, "Success", "")
}

func AddBattery(c *gin.Context) {
	no := c.PostForm("no")
	fork := c.PostForm("fork")
	dc := c.PostForm("dc")
	batteryValidBody := validate.InsertBatteryBody{
		Number:      no,
		Forkliftcat: fork,
		Warehouse:   dc,
	}
	valid := validate.InsertBatteryValidator(&batteryValidBody)
	GinUtil := util.Gin{Context: c}
	if valid != nil {
		GinUtil.Response(400, 0, valid.Error(), "")
		return
	}
	addBatteryBody := model.Battery{
		Number:    no,
		Forklift:  fork, // Type of forklift, not specific forklift SN
		Warehouse: dc,
		Status:    "off",
		FullNo:    dc + fork + no,
	}
	insertResult := model.InsertBattery(&addBatteryBody)
	if insertResult != nil {
		GinUtil.Response(400, 0, insertResult.Error(), "")
		return
	} else {
		GinUtil.Response(200, 0, "Success", "")
	}
}
