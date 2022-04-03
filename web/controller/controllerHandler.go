/**
  @Prject: goProjects
  @Dev_Software: GoLand
  @File : controllerHandler
  @Time : 2018/10/18 14:31 
  @Author : AllenIverson
*/

package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fabric-identity/service"
	"github.com/fabric-identity/web/utils"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var cuser User

func (app *Application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	//将请求body解析到dst中
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		//如果在解析的过程中发生错误，开始分类...
		var syntaxError *json.SyntaxError
		var unmarshalTypeError * json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		//使用errors.As()函数检查错误类型是否为*json.SyntaxError。如果是，返回错误
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contain badly-formed JSON (at charcter %d)")
		//某些情况下，因为语法错误Decode()函数可能返回io.ErrUnexpectedEOF错误。
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contain badly-formed JSON")
		//同样捕获*json.UnmarshalTypeError错误，这些错误是因为JSON值和接收对象不匹配。如果错误对应到特定到字段，
		//我们可以指出哪个字段错误方便客户端debug
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		//如果解码到内容是空对象，会返回io.EOF。
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		//如果decode()函数传入一个非空的指针，将返回json.InvalidUnmarshalError。
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}
	return nil
}

//档案管理
func (app *Application) QueryManagement (w http.ResponseWriter, r *http.Request)  {
	data := &struct{
		Infos []*ArchivesInfo
	}{
		Infos:ArchivesInfos,
	}
	ShowView(w, r, "danganmain.html", data)
}
//添加
func (app *Application)DanganAdd(w http.ResponseWriter, r *http.Request){
	ShowView(w, r, "Dangan-add.html", nil)
}

func (app *Application)OperateManagement(w http.ResponseWriter, r *http.Request){
	data := &struct{
		Infos []*ArchivesInfo
	}{
		Infos:ArchivesInfos,
	}
	ShowView(w, r, "danganmain2.html", data)
}
func (app *Application)AddArchives(w http.ResponseWriter, r *http.Request)()  {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	index := strings.Index(string(bodyBytes), "=")
	requestBody := string(bodyBytes)[index+1:]

	var formData service.Archives
	err = json.Unmarshal([]byte(requestBody), &formData)
	//err = json.NewDecoder(r.Body).Decode(&formData)
	if err != nil{
		fmt.Println("AdSet json转结构体出错！err>>> ",err)
	}
	stuArchives[formData.UserName] = append(stuArchives[formData.UserName], &formData)

	timeStr:=time.Now().Format("2006-01-02 15:04:05")

	ArchivesID := utils.GetRandomString(8)

	archiveInfo := ArchivesInfo{
		ArchivesID: ArchivesID,
		Operator:cuser.LoginName,
		CurrentUser:formData.UserName,
		CreateTime:timeStr,
		UpdateTime:timeStr,
		InfoMsg:"",
	}
	ArchivesInfos = append(ArchivesInfos,&archiveInfo)

	data := &struct{
		Infos []*ArchivesInfo
	}{
		Infos:ArchivesInfos,
	}
	ShowView(w, r, "Danganmain2.html", data)
}
//编辑
func (app *Application)DanganModifyShow(w http.ResponseWriter, r *http.Request){
	// 根据证书编号与姓名查询信息
	//name := r.FormValue("name")
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	indexOne := strings.Index(string(bodyBytes), "=")
	indexTwo := strings.Index(string(bodyBytes), "!")
	requestBody := string(bodyBytes)[indexOne+1:indexTwo]
	if _,ok := stuArchives[requestBody];!ok{
		fmt.Println("该用户不存在")
	}

	data := &struct {
		Archives service.Archives
		CurrentUser User
	}{
		Archives:*stuArchives[requestBody][0],
		CurrentUser:cuser,
	}

	ShowView(w, r, "Dangan-edit.html", data)
}

func (app *Application)DanganModify(w http.ResponseWriter, r *http.Request){
	modifyArchives := service.Archives{
		UserName: r.FormValue("name"),
		Sex: r.FormValue("sex"),
		BirthPlace: r.FormValue("birthplace"),
		BirthDay: r.FormValue("birthday"),
		Age:r.FormValue("age"),
		Religion:r.FormValue("religion"),
		PostCodes:r.FormValue("postcodes"),
		TypeofCompany:r.FormValue("typeofcompany"),
		TypeofWork:r.FormValue("typeofwork"),
		Address:r.FormValue("address"),
		Telephone:r.FormValue("telephone"),
		Role:r.FormValue("role"),
		Married:r.FormValue("married"),
		Educated:r.FormValue("educated"),
	}
	stuArchives[modifyArchives.UserName] = append(stuArchives[modifyArchives.UserName], &modifyArchives)
	data := &struct {
		ModifyInfo []*service.Archives
	}{
		ModifyInfo:stuArchives[modifyArchives.UserName],
	}
	ShowView(w, r, "Dangan-history.html",data)
}


//查看


func (app *Application)RootDanganCheck(w http.ResponseWriter, r *http.Request){
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	indexOne := strings.Index(string(bodyBytes), "=")
	indexTwo := strings.Index(string(bodyBytes), "!")
	requestBody := string(bodyBytes)[indexOne+1:indexTwo]
	if _,ok := stuArchives[requestBody];!ok{
		fmt.Println("该用户不存在")
	}
	checkArchives := stuArchives[requestBody][0]

	data := &struct{
		Archives *service.Archives
	}{
		Archives:checkArchives,
	}
	ShowView(w, r, "Dangan-look2.html",data)
}
func (app *Application)DanganCheck(w http.ResponseWriter, r *http.Request){

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	indexOne := strings.Index(string(bodyBytes), "=")
	indexTwo := strings.Index(string(bodyBytes), "!")
	requestBody := string(bodyBytes)[indexOne+1:indexTwo]

	if _,ok := stuArchives[requestBody];!ok{
		fmt.Println("该用户不存在")
	}

	checkArchives := stuArchives[requestBody][0]

	data := &struct{
		Archives *service.Archives
	}{
		Archives:checkArchives,
	}
	ShowView(w, r, "Dangan-look.html",data)
}

func (app *Application)LookArchives(w http.ResponseWriter, r *http.Request){
	ShowView(w, r, "Dangan-edit.html", nil)
}


func (app *Application) LoginView(w http.ResponseWriter, r *http.Request)  {

	ShowView(w, r, "login.html", nil)
}

func (app *Application) RegisterView(w http.ResponseWriter, r *http.Request)  {

	ShowView(w, r, "register.html", nil)
}

func (app *Application) Index(w http.ResponseWriter, r *http.Request)  {
	ShowView(w, r, "index.html", nil)
}

func (app *Application) Help(w http.ResponseWriter, r *http.Request)  {
	data := &struct {
		CurrentUser User
	}{
		CurrentUser:cuser,
	}
	ShowView(w, r, "help.html", data)
}
//用户注册
func (app *Application) Register(w http.ResponseWriter, r *http.Request) {
	loginName := r.FormValue("loginName")
	password := r.FormValue("password")
	rule := r.FormValue("rule")

	newUser := &User{
		LoginName:		loginName,
		Password:		password,
		IsAdmin: 		rule,
	}
	for _,k := range users {
		if k.LoginName == newUser.LoginName {
			fmt.Println("已经存在")
			//w.Write([]byte("a user with this loginName already exists"))
		}
	}
	users = append(users, *newUser)
	ShowView(w, r, "login.html", nil)

}

// 用户登录
func (app *Application) Login(w http.ResponseWriter, r *http.Request) {
	loginName := r.FormValue("loginName")
	password := r.FormValue("password")
	var flag bool
	for _, user := range users {
		if user.LoginName == loginName && user.Password == password {
			cuser = user
			flag = true
			break
		}
	}

	data := &struct {
		CurrentUser User
		Flag bool
	}{
		CurrentUser:cuser,
		Flag:false,
	}

	if flag {
		// 登录成功
		ShowView(w, r, "help.html", data)
	}else{
		// 登录失败
		data.Flag = true
		data.CurrentUser.LoginName = loginName
		ShowView(w, r, "login.html", data)
	}
}

// 用户登出
func (app *Application) LoginOut(w http.ResponseWriter, r *http.Request)  {
	cuser = User{}
	ShowView(w, r, "login.html", nil)
}

// 显示添加信息页面
func (app *Application) AddEduShow(w http.ResponseWriter, r *http.Request)  {
	data := &struct {
		CurrentUser User
		Msg string
		Flag bool
	}{
		CurrentUser:cuser,
		Msg:"",
		Flag:false,
	}
	ShowView(w, r, "addEdu.html", data)
}

func (app *Application) AddScoreShow (w http.ResponseWriter, r *http.Request){
	data := &struct {
		CurrentUser User
		Msg string
		Flag bool
	}{
		CurrentUser:cuser,
		Msg:"",
		Flag:false,
	}
	ShowView(w, r, "uploadScore.html", data)
}

//添加成绩信息
func (app *Application) AddScore(w http.ResponseWriter, r *http.Request) {
	var formDataArray []*service.Score
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	indexLeft := []int{}
	indexRight := []int{}
	for i:=0;i<len(string(bodyBytes));i++ {
		if string(bodyBytes)[i] == 123{

			indexLeft = append(indexLeft, i)
		}
		if string(bodyBytes)[i] == 125{
			indexRight = append(indexRight, i)
		}
	}
	requestBody := []string{}
	for i:=0;i<len(indexLeft);i++{
		str := string(bodyBytes)[indexLeft[i]:indexRight[i]+1]
		requestBody = append(requestBody, str)
	}
	for i:=0;i<len(requestBody);i++{
		formData := service.Score{}
		err = json.Unmarshal([]byte(requestBody[i]), &formData)
		if err != nil{
			fmt.Println("AdSet json转结构体出错！err>>> ",err)
		}
		formDataArray = append(formDataArray, &formData)
	}


	data := &struct {
		Items []*service.Score
		StuNum string
		StuName string
	}{
		Items:formDataArray,

		StuNum:formDataArray[0].StuNum,
		StuName:"",
	}
	addStu := StuScore{data.Items,data.StuNum,data.StuName}
	stuScores = append(stuScores, addStu)

	ShowView(w,r,"queryScoreResult.html",data)
}
func (app *Application) ShowScore(w http.ResponseWriter, r *http.Request)  {
	ShowView(w, r, "queryScore.html", nil)
}

// 添加信息
func (app *Application) AddEdu(w http.ResponseWriter, r *http.Request)  {

	edu := service.Education{
		Name:r.FormValue("name"),
		Gender:r.FormValue("gender"),
		Nation:r.FormValue("nation"),
		EntityID:r.FormValue("entityID"),
		Place:r.FormValue("place"),
		BirthDay:r.FormValue("birthDay"),
		EnrollDate:r.FormValue("enrollDate"),
		GraduationDate:r.FormValue("graduationDate"),
		SchoolName:r.FormValue("schoolName"),
		Major:r.FormValue("major"),
		QuaType:r.FormValue("quaType"),
		Mode:r.FormValue("mode"),
		Graduation:r.FormValue("graduation"),
		CertNo:r.FormValue("certNo"),
		Photo:r.FormValue("photo"),
	}

	app.Setup.SaveEdu(edu)
	/*transactionID, err := app.Setup.SaveEdu(edu)

	data := &struct {
		CurrentUser User
		Msg string
		Flag bool
	}{
		CurrentUser:cuser,
		Flag:true,
		Msg:"",
	}

	if err != nil {
		data.Msg = err.Error()
	}else{
		data.Msg = "信息添加成功:" + transactionID
	}*/

	//ShowView(w, r, "addEdu.html", data)
	r.Form.Set("certNo", edu.CertNo)
	r.Form.Set("name", edu.Name)
	app.FindCertByNoAndName(w, r)
}

func (app *Application) QueryPage(w http.ResponseWriter, r *http.Request)  {
	data := &struct {
		CurrentUser User
		Msg string
		Flag bool
	}{
		CurrentUser:cuser,
		Msg:"",
		Flag:false,
	}
	ShowView(w, r, "query.html", data)
}

// 根据证书编号与姓名查询信息
func (app *Application) FindCertByNoAndName(w http.ResponseWriter, r *http.Request)  {
	certNo := r.FormValue("certNo")
	name := r.FormValue("name")
	result, err := app.Setup.FindEduByCertNoAndName(certNo, name)
	var edu = service.Education{}
	json.Unmarshal(result, &edu)

	//fmt.Println("根据证书编号与姓名查询信息成功：")
	//fmt.Println(edu)

	data := &struct {
		Edu service.Education
		CurrentUser User
		Msg string
		Flag bool
		History bool
	}{
		Edu:edu,
		CurrentUser:cuser,
		Msg:"",
		Flag:false,
		History:false,
	}

	if err != nil {
		data.Msg = err.Error()
		data.Flag = true
	}

	ShowView(w, r, "queryResult.html", data)
}

func (app *Application) QueryPage2(w http.ResponseWriter, r *http.Request)  {
	data := &struct {
		CurrentUser User
		Msg string
		Flag bool
	}{
		CurrentUser:cuser,
		Msg:"",
		Flag:false,
	}
	ShowView(w, r, "query2.html", data)
}

func (app *Application) QueryPage3(w http.ResponseWriter, r *http.Request)  {
	data := &struct {
		CurrentUser User
		Msg string
		Flag bool
	}{
		CurrentUser:cuser,
		Msg:"",
		Flag:false,
	}
	ShowView(w, r, "query3.html", data)
}
//根据身份证号查询成绩信息

func (app *Application) ByNameFindScore (w http.ResponseWriter, r *http.Request)  {
	certNo := r.FormValue("stuNo")

	//name := r.FormValue("name")
	result := StuScore{}
	for i:=0;i<len(stuScores);i++{

		if stuScores[i].StuNum == certNo{
			result.StuClss = stuScores[i].StuClss
			result.StuName = stuScores[i].StuName
			result.StuNum = certNo
		}
	}
	data := &struct {
		Items []*service.Score
		StuNum string
		StuName string
	}{
		Items:result.StuClss,
		StuNum:result.StuNum,
		StuName:result.StuName,
	}
	ShowView(w,r,"queryScoreResult.html",data)
}


// 根据身份证号码查询信息
func (app *Application) FindByID(w http.ResponseWriter, r *http.Request)  {
	entityID := r.FormValue("entityID")
	result, err := app.Setup.FindEduInfoByEntityID(entityID)
	var edu = service.Education{}
	json.Unmarshal(result, &edu)

	data := &struct {
		Edu service.Education
		CurrentUser User
		Msg string
		Flag bool
		History bool
	}{
		Edu:edu,
		CurrentUser:cuser,
		Msg:"",
		Flag:false,
		History:true,
	}

	if err != nil {
		data.Msg = err.Error()
		data.Flag = true
	}

	ShowView(w, r, "queryResult.html", data)
}

func (app *Application) FindByIDHistory(w http.ResponseWriter, r *http.Request)  {
	entityID := r.FormValue("entityID")
	result, err := app.Setup.FindEduInfoByEntityID(entityID)
	var edu = service.Education{}
	json.Unmarshal(result, &edu)

	data := &struct {
		Edu service.Education
		CurrentUser User
		Msg string
		Flag bool
		History bool
	}{
		Edu:edu,
		CurrentUser:cuser,
		Msg:"",
		Flag:false,
		History:true,
	}

	if err != nil {
		data.Msg = err.Error()
		data.Flag = true
	}

	ShowView(w, r, "queryHistory.html", data)
}

// 修改/添加新信息
func (app *Application) ModifyShow(w http.ResponseWriter, r *http.Request)  {
	// 根据证书编号与姓名查询信息
	certNo := r.FormValue("certNo")
	name := r.FormValue("name")
	result, err := app.Setup.FindEduByCertNoAndName(certNo, name)

	var edu = service.Education{}
	json.Unmarshal(result, &edu)

	data := &struct {
		Edu service.Education
		CurrentUser User
		Msg string
		Flag bool
	}{
		Edu:edu,
		CurrentUser:cuser,
		Flag:true,
		Msg:"",
	}

	if err != nil {
		data.Msg = err.Error()
		data.Flag = true
	}

	ShowView(w, r, "modify.html", data)
}

// 修改/添加新信息
func (app *Application) Modify(w http.ResponseWriter, r *http.Request) {
	edu := service.Education{
		Name:r.FormValue("name"),
		Gender:r.FormValue("gender"),
		Nation:r.FormValue("nation"),
		EntityID:r.FormValue("entityID"),
		Place:r.FormValue("place"),
		BirthDay:r.FormValue("birthDay"),
		EnrollDate:r.FormValue("enrollDate"),
		GraduationDate:r.FormValue("graduationDate"),
		SchoolName:r.FormValue("schoolName"),
		Major:r.FormValue("major"),
		QuaType:r.FormValue("quaType"),
		Mode:r.FormValue("mode"),
		Graduation:r.FormValue("graduation"),
		CertNo:r.FormValue("certNo"),
		Photo:r.FormValue("photo"),
	}

	//transactionID, err := app.Setup.ModifyEdu(edu)
	app.Setup.ModifyEdu(edu)

	/*data := &struct {
		Edu service.Education
		CurrentUser User
		Msg string
		Flag bool
	}{
		CurrentUser:cuser,
		Flag:true,
		Msg:"",
	}

	if err != nil {
		data.Msg = err.Error()
	}else{
		data.Msg = "新信息添加成功:" + transactionID
	}

	ShowView(w, r, "modify.html", data)
	*/

	r.Form.Set("certNo", edu.CertNo)
	r.Form.Set("name", edu.Name)
	app.FindCertByNoAndName(w, r)
}
