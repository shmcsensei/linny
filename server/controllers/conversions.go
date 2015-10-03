package controllers

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/simonchong/linny/common"
	"github.com/simonchong/linny/constants"
	"github.com/simonchong/linny/server/controllers/conversions"
	"github.com/simonchong/linny/server/wrappers"
	"github.com/zenazn/goji/web"
)

func ConversionsJS(ac *wrappers.AppContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {

	code, err := ioutil.ReadFile("./resources/conversion.js")
	if err != nil {
		panic(err)
	}
	unix := strconv.FormatInt(time.Now().Unix(), 10)
	tag := r.FormValue("t")

	body := "(function(h, v, t, g) {" + string(code) + "})('" + r.Host + "','" + constants.MeasureDir + "','" + tag + "'," + unix + ");"

	w.Header().Set(
		"Content-Type",
		"application/javascript",
	)
	fmt.Fprint(w, body)

	return 200, nil
}

func Conversions(ac *wrappers.AppContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {

	adID, errA := conversions.GetCookie(r)
	if errA != nil {
		return http.StatusOK, errA
	}

	originIP, _, errIP := net.SplitHostPort(r.RemoteAddr)
	if errIP != nil {
		return http.StatusOK, errIP
	}

	timeGen, errT := common.FormTime("g", r)
	if errT != nil {
		return http.StatusOK, errT
	}

	sessionID, errS := wrappers.GetSessionCookie(r)
	if errS != nil {
		return http.StatusOK, errS
	}

	conversionTag := r.FormValue("t")
	//TODO limit to 255 characters

	referer := r.Header.Get("referer")

	fmt.Println("Conversion Origin IP", originIP)
	fmt.Println("Conversion Referer", referer)
	fmt.Println("Conversion Gen Time", timeGen)
	fmt.Println("Conversion adID", adID)
	fmt.Println("Conversion Conversion Tag", conversionTag)
	fmt.Println("Conversion Session", sessionID)

	//Add to DB
	ac.Data.AdConversions.Insert(adID, referer, originIP, conversionTag, sessionID)

	//Send GIF response
	gif, _ := base64.StdEncoding.DecodeString("R0lGODlhAQABAIABAP///wAAACwAAAAAAQABAAACAkQBADs=")
	w.Header().Set("Content-Type", "image/gif")
	io.WriteString(w, string(gif))

	return 200, nil
}