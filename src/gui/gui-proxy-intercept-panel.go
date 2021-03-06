package gui

import (
	"gitee.com/snxamdf/golcl/lcl"
	"gitee.com/snxamdf/golcl/lcl/types"
	"gitee.com/snxamdf/golcl/lcl/types/colors"
	"gitee.com/snxamdf/http-server/src/common"
	"gitee.com/snxamdf/http-server/src/consts"
	"gitee.com/snxamdf/http-server/src/entity"
	"path/filepath"
	"sync"
)

//代理拦截Panel
type ProxyInterceptPanel struct {
	mutex                       sync.Mutex //锁
	TPanel                      *lcl.TPanel
	State                       int32            //当前状态
	UrlAddrEdit                 *lcl.TEdit       //拦截地址
	StateLabel                  *lcl.TStaticText //拦截状态
	StateOkBtn                  *lcl.TButton     //拦截状态确认按钮
	interceptPageControl        *lcl.TPageControl
	ProxyInterceptRequestPanel  *ProxyInterceptRequestPanel  //代理拦截请求Panel
	ProxyInterceptResponsePanel *ProxyInterceptResponsePanel //代理拦截响应Panel
	ProxyInterceptSettingPanel  *ProxyInterceptSettingPanel  //代理拦截配置Panel
	InterceptQueue              *common.Queue                //拦截队列
	InterceptQueueProxyDetail   *entity.ProxyDetail          //当前拦截队列处理的代理
	IsUseInterceptQueue         bool                         //是否正在使用拦截队列
}

//代理拦截请求Panel
type ProxyInterceptRequestPanel struct {
	TPanel              *lcl.TPanel
	ParamsGrid          *lcl.TStringGrid
	ParamsGridRowCount  int32
	HeadersGrid         *lcl.TStringGrid
	HeadersGridRowCount int32
	TBodyPanel          *ProxyInterceptRequestBodyPanel
}

//代理拦截请求Body Panel
type ProxyInterceptRequestBodyPanel struct {
	TPanel               *lcl.TPanel
	RawPanel             *lcl.TPanel
	RawMemo              *lcl.TMemo
	FormDataGridPanel    *lcl.TPanel
	FormDataGrid         *lcl.TStringGrid
	FormDataGridOpenFile *lcl.TOpenDialog
	FormDataGridList     []*entity.FormDataGridList
	FormDataGridRowCount int32
}

//代理拦截响应Panel
type ProxyInterceptResponsePanel struct {
	TPanel *lcl.TPanel
}

//代理拦截配置Panel
type ProxyInterceptSettingPanel struct {
	TPanel        *lcl.TPanel
	OnOffBtn      *lcl.TImageButton
	InterceptGrid *lcl.TStringGrid
	//InterceptGridRowCount   int32
	InterceptGridConfigData *[]*entity.ProxyInterceptConfig
}

//proxy intercept request UI
func (m *ProxyInterceptRequestPanel) initUI() {
	//Tabs 的控制标签
	resetPVars()
	pLeft = 0
	pTop = 0
	pHeight = m.TPanel.Height()/2 - 50
	pWidth = m.TPanel.Width()
	reqPageControl := lcl.NewPageControl(m.TPanel)
	reqPageControl.SetParent(m.TPanel)
	reqPageControl.SetBounds(pLeft, pTop, pWidth, pHeight)
	reqPageControl.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkRight))

	//--- begin --- Request Query Params
	var paramsSheet = lcl.NewTabSheet(reqPageControl) //标签页
	paramsSheet.SetPageControl(reqPageControl)
	paramsSheet.SetCaption("　Request Query Params　")
	paramsSheet.SetAlign(types.AlClient)
	paramsPanel := lcl.NewPanel(m.TPanel) // 标签页
	paramsPanel.SetParent(paramsSheet)
	paramsPanel.SetAlign(types.AlClient)
	//按钮
	var reqQueryParamAddBtn = lcl.NewButton(m.TPanel)
	reqQueryParamAddBtn.SetParent(paramsSheet)
	reqQueryParamAddBtn.SetCaption("　添加参数　")
	reqQueryParamAddBtn.SetBounds(460, 1, 60, 30)
	reqQueryParamAddBtn.SetOnClick(func(sender lcl.IObject) {
		m.QueryParamsGridAdd("", "")
	})

	//ParamsGrid
	m.ParamsGrid = lcl.NewStringGrid(paramsPanel)
	m.ParamsGrid.SetParent(paramsPanel)
	m.ParamsGrid.SetFixedCols(0)
	m.ParamsGrid.SetFixedColor(colors.ClGreen)
	m.ParamsGrid.SetAlign(types.AlClient)
	m.ParamsGrid.SetBorderStyle(types.BsNone)
	m.ParamsGrid.SetFlat(true)
	m.ParamsGrid.SetOptions(m.ParamsGrid.Options().Include(types.GoAlwaysShowEditor, types.GoCellHints, types.GoEditing, types.GoTabs))
	m.ParamsGrid.SetOnSetEditText(func(sender lcl.IObject, aCol, aRow int32, value string) {
		if aCol == 1 || aCol == 2 {
			if aRow == m.ParamsGridRowCount-1 && value != "" {
				m.QueryParamsGridAdd("", "")
			}
		}
	})
	m.ParamsGrid.SetOnButtonClick(func(sender lcl.IObject, aCol, aRow int32) {
		if aCol == 3 {
			if m.ParamsGridRowCount > 2 {
				m.ParamsGrid.DeleteRow(aRow)
				m.ParamsGridRowCount--
			}
		}
	})
	m.RequestQueryParamsGridHead() //请求拦截参数表格头
	m.ParamsGrid.SetRow(m.ParamsGridRowCount)
	m.ParamsGrid.SetRowCount(m.ParamsGridRowCount)
	//m.QueryParamsGridAdd("", "") //默认添加一条
	//--- end --- Request Query Params

	//--- begin --- Request Headers
	var headersSheet = lcl.NewTabSheet(reqPageControl) //标签页
	headersSheet.SetPageControl(reqPageControl)
	headersSheet.SetCaption("　Request Headers　")
	headersSheet.SetAlign(types.AlClient)
	headersPanel := lcl.NewPanel(m.TPanel) // 标签页
	headersPanel.SetParent(headersSheet)
	headersPanel.SetAlign(types.AlClient)
	//按钮
	var reqHeaderAddBtn = lcl.NewButton(m.TPanel)
	reqHeaderAddBtn.SetParent(headersSheet)
	reqHeaderAddBtn.SetCaption("　添加请求头　")
	reqHeaderAddBtn.SetBounds(460, 1, 80, 30)
	reqHeaderAddBtn.SetOnClick(func(sender lcl.IObject) {
		m.HeaderGridAdd("", "")
	})
	//HeadersGrid
	m.HeadersGrid = lcl.NewStringGrid(headersPanel)
	m.HeadersGrid.SetParent(headersPanel)
	m.HeadersGrid.SetFixedCols(0)
	m.HeadersGrid.SetFixedColor(colors.ClGreen)
	m.HeadersGrid.SetAlign(types.AlClient)
	m.HeadersGrid.SetBorderStyle(types.BsNone)
	m.HeadersGrid.SetFlat(true)
	m.HeadersGrid.SetOptions(m.HeadersGrid.Options().Include(types.GoAlwaysShowEditor, types.GoCellHints, types.GoEditing, types.GoTabs))
	m.HeadersGrid.SetOnSetEditText(func(sender lcl.IObject, aCol, aRow int32, value string) {
		if aCol == 1 || aCol == 2 {
			if aRow == m.HeadersGridRowCount-1 && value != "" {
				m.HeaderGridAdd("", "")
			}
		}
	})
	m.HeadersGrid.SetOnButtonClick(func(sender lcl.IObject, aCol, aRow int32) {
		if aCol == 3 {
			if m.HeadersGridRowCount > 2 {
				m.HeadersGrid.DeleteRow(aRow)
				m.HeadersGridRowCount--
			}
		}
	})
	m.HeaderGridHead()
	m.HeadersGrid.SetRow(m.HeadersGridRowCount)
	m.HeadersGrid.SetRowCount(m.ParamsGridRowCount)
	//m.HeaderGridAdd("", "")
	//--- end --- Request Headers

	//--- begin --- Request Body
	resetPVars()
	pLeft = 0
	pTop = reqPageControl.Height()
	pHeight = m.TPanel.Height()/2 + 50
	pWidth = m.TPanel.Width()
	reqPageControl = lcl.NewPageControl(m.TPanel)
	reqPageControl.SetParent(m.TPanel)
	reqPageControl.SetBounds(pLeft, pTop, pWidth, pHeight)
	reqPageControl.SetAnchors(types.NewSet(types.AkLeft, types.AkBottom, types.AkTop, types.AkRight))
	var bodySheet = lcl.NewTabSheet(reqPageControl) //标签页
	bodySheet.SetPageControl(reqPageControl)
	bodySheet.SetCaption("　Request Body　")
	bodySheet.SetAlign(types.AlClient)
	m.TBodyPanel.TPanel = lcl.NewPanel(m.TPanel) // 标签页
	m.TBodyPanel.TPanel.SetParent(bodySheet)
	m.TBodyPanel.TPanel.SetAlign(types.AlClient)
	resetPVars()
	pLeft = 30
	pTop = 5
	var rdoRaw = lcl.NewRadioButton(m.TBodyPanel.TPanel)
	rdoRaw.SetParent(m.TBodyPanel.TPanel)
	rdoRaw.SetCaption("raw/json")
	rdoRaw.SetLeft(pLeft)
	rdoRaw.SetTop(pTop)
	rdoRaw.SetOnClick(func(sender lcl.IObject) {
		m.TBodyPanel.bodyRdoCheckClick(0)
	})
	var rdoFormData = lcl.NewRadioButton(m.TBodyPanel.TPanel)
	rdoFormData.SetParent(m.TBodyPanel.TPanel)
	rdoFormData.SetCaption("form-data/x-www-form-urlencoded/binary")
	rdoFormData.SetLeft(rdoRaw.Left() + 120)
	rdoFormData.SetTop(pTop)
	rdoFormData.SetOnClick(func(sender lcl.IObject) {
		m.TBodyPanel.bodyRdoCheckClick(1)
	})

	m.TBodyPanel.initUI()
	//--- end --- Request Body

	//最后初始body选中
	rdoRaw.SetChecked(true)
	m.TBodyPanel.bodyRdoCheckClick(0)
}

//代理拦截请求Body Panel UI
func (m *ProxyInterceptRequestBodyPanel) initUI() {
	//raw
	resetPVars()
	pLeft = 10
	pTop = 31
	m.RawPanel = lcl.NewPanel(m.TPanel)
	m.RawPanel.SetParent(m.TPanel)
	m.RawPanel.SetBounds(pLeft, pTop, m.TPanel.Width()-20, m.TPanel.Height()-41)
	m.RawPanel.SetAnchors(types.NewSet(types.AkLeft, types.AkBottom, types.AkTop, types.AkRight))
	m.RawPanel.SetVisible(false)
	m.RawMemo = lcl.NewMemo(m.RawPanel)
	m.RawMemo.SetParent(m.RawPanel)
	m.RawMemo.SetAlign(types.AlClient)

	//form-data & x-www-form-urlencoded
	m.FormDataGridOpenFile = lcl.NewOpenDialog(m.TPanel)
	m.FormDataGridOpenFile.SetTitle("选择上传文件")
	m.FormDataGridPanel = lcl.NewPanel(m.TPanel)
	m.FormDataGridPanel.SetParent(m.TPanel)
	m.FormDataGridPanel.SetBounds(pLeft, pTop, m.TPanel.Width()-20, m.TPanel.Height()-41)
	m.FormDataGridPanel.SetAnchors(types.NewSet(types.AkLeft, types.AkBottom, types.AkTop, types.AkRight))
	m.FormDataGridPanel.SetVisible(false)
	m.FormDataGrid = lcl.NewStringGrid(m.FormDataGridPanel)
	m.FormDataGrid.SetParent(m.FormDataGridPanel)
	m.FormDataGrid.SetFixedCols(0)
	m.FormDataGrid.SetFixedColor(colors.ClGreen)
	m.FormDataGrid.SetBorderStyle(types.BsNone)
	m.FormDataGrid.SetFlat(true)
	m.FormDataGrid.SetOptions(m.FormDataGrid.Options().Include(types.GoAlwaysShowEditor, types.GoCellHints, types.GoEditing, types.GoTabs, types.GoRowHighlight))
	m.FormDataGrid.SetAlign(types.AlClient)
	m.FormDataGrid.SetOnSetEditText(func(sender lcl.IObject, aCol, aRow int32, value string) {
		var idxRow = aRow - 1
		row := m.FormDataGridList[idxRow]
		if aCol == 0 { //列 类型
			m.FormDataGrid.SetOptions(m.FormDataGrid.Options().Exclude(types.GoEditing))
			//解决同步到列表问题
			common.NewDebounce(10).Start(func() { //是个线程操作
				lcl.ThreadSync(func() {
					if value == "Text" {
						m.FormDataGrid.SetCells(3, aRow, "---")
						row.FileValue = ""
					} else if value == "File" {
						if row.FileValue == "" {
							m.FormDataGrid.SetCells(3, aRow, "选择文件")
						} else {
							_, fileName := filepath.Split(row.FileValue)
							m.FormDataGrid.SetCells(3, aRow, fileName)
						}
					} else if value != "Text" && value != "File" {
						value = "Text"
						m.FormDataGrid.SetCells(0, aRow, value)
					}
					row.Type = value
				})
			})
			m.FormDataGrid.CanFocus()
			m.FormDataGrid.SetOptions(m.FormDataGrid.Options().Include(types.GoEditing))
		} else if aCol == 1 {
			row.Key = value
		} else if aCol == 2 {
			row.TextValue = value
		} else if aCol == 3 {
			//row.FileValue = value
		}
		if aCol == 1 || aCol == 2 {
			if aRow == m.FormDataGridRowCount-1 && (row.Key != "" || row.TextValue != "") {
				m.FormDataGridAdd("Text", "", "")
			}
		}
	})
	m.FormDataGrid.SetOnButtonClick(func(sender lcl.IObject, aCol, aRow int32) {
		var idxRow = aRow - 1
		//按钮触发
		if aCol == 4 { //删除行
			if m.FormDataGridRowCount > 2 {
				m.FormDataGrid.DeleteRow(aRow)
				var before = m.FormDataGridList[:idxRow]
				var after = m.FormDataGridList[idxRow+1:]
				m.FormDataGridList = append(before, after...)
				m.FormDataGridRowCount--
			}
		} else if aCol == 3 { //选择文件
			row := m.FormDataGridList[idxRow]
			if row.Type == "File" {
				//解决同步到列表问题
				m.FormDataGrid.SetOptions(m.FormDataGrid.Options().Exclude(types.GoEditing))
				//解决同步到列表问题
				common.NewDebounce(1).Start(func() { //是个线程操作
					lcl.ThreadSync(func() { //需要主线程同步
						if m.FormDataGridOpenFile.Execute() {
							var filePath = m.FormDataGridOpenFile.FileName()
							row.FileValue = filePath
							_, fileName := filepath.Split(filePath)
							m.FormDataGrid.SetCells(3, aRow, fileName)
						}
					})
				})
				m.FormDataGrid.SetOptions(m.FormDataGrid.Options().Include(types.GoEditing))
			}
		}
	})
	m.FormDataGridHead()
	m.FormDataGrid.SetRow(m.FormDataGridRowCount)
	m.FormDataGrid.SetRowCount(m.FormDataGridRowCount)
	//m.FormDataGridAdd("", "")
	//按钮
	var reqFormAddBtn = lcl.NewButton(m.FormDataGridPanel)
	reqFormAddBtn.SetParent(m.FormDataGridPanel)
	reqFormAddBtn.SetCaption("　添加参数　")
	reqFormAddBtn.SetBounds(520, 2, 80, 30)
	reqFormAddBtn.SetOnClick(func(sender lcl.IObject) {
		m.FormDataGridAdd("Text", "", "")
		//for k, v := range m.FormDataGridList {
		//	fmt.Printf("%v,%+v\n", k, v)
		//}
	})
}

//清空Body表格
func (m *ProxyInterceptRequestBodyPanel) ClearFormDataGrid() {
	lcl.ThreadSync(func() {
		m.FormDataGrid.Clear()
		m.FormDataGridRowCount = 1
		m.FormDataGrid.SetRow(m.FormDataGridRowCount)
		m.FormDataGrid.SetRowCount(m.FormDataGridRowCount)
	})
}

//请求Body表格添加
func (m *ProxyInterceptRequestBodyPanel) FormDataGridAdd(t, key, value string) {
	lcl.ThreadSync(func() {
		var n = &entity.FormDataGridList{
			Type: t,
			Key:  key,
		}
		if t == "Text" {
			n.TextValue = value
		} else if t == "File" {
			n.FileValue = value
		}
		m.FormDataGridList = append(m.FormDataGridList, n)
		//在指定位置播放一行
		m.FormDataGrid.InsertColRow(false, m.FormDataGridRowCount)
		m.FormDataGrid.SetCells(0, m.FormDataGridRowCount, "Text")
		m.FormDataGrid.SetCells(1, m.FormDataGridRowCount, key)
		m.FormDataGrid.SetCells(2, m.FormDataGridRowCount, value)
		m.FormDataGrid.SetCells(3, m.FormDataGridRowCount, "---")
		m.FormDataGrid.SetCells(4, m.FormDataGridRowCount, "删除")
		m.FormDataGridRowCount++
		m.FormDataGrid.SetRowCount(m.FormDataGridRowCount)
	})
}

//请求Body表格头
func (m *ProxyInterceptRequestBodyPanel) FormDataGridHead() {
	var colType = m.FormDataGrid.Columns().Add()
	colType.SetWidth(50)
	colType.Title().SetCaption("TYPE")
	colType.SetButtonStyle(types.CbsPickList)
	colType.Title().SetAlignment(types.TaCenter)
	colType.SetAlignment(types.TaCenter)
	colType.PickList().Add("Text")
	colType.PickList().Add("File")

	var colKey = m.FormDataGrid.Columns().Add()
	colKey.SetWidth(150)
	colKey.Title().SetCaption("Key")

	var colTextValue = m.FormDataGrid.Columns().Add()
	colTextValue.SetWidth(150)
	colTextValue.Title().SetCaption("Text Value")

	var colFileValue = m.FormDataGrid.Columns().Add()
	colFileValue.SetWidth(100)
	colFileValue.Title().SetCaption("Select File")
	colFileValue.Title().SetAlignment(types.TaCenter)
	colFileValue.SetButtonStyle(types.CbsButtonColumn)
	colFileValue.SetAlignment(types.TaRightJustify)

	var delBtn = m.FormDataGrid.Columns().Add()
	delBtn.SetWidth(60)
	delBtn.Title().SetCaption("操作")
	delBtn.Title().SetAlignment(types.TaCenter)
	delBtn.SetButtonStyle(types.CbsButtonColumn)
	delBtn.SetAlignment(types.TaCenter)
}

//body radio 按钮点击
func (m *ProxyInterceptRequestBodyPanel) bodyRdoCheckClick(t int) {
	m.RawPanel.SetVisible(t == 0)
	m.FormDataGridPanel.SetVisible(t == 1)
}

//请求拦截头添加
func (m *ProxyInterceptRequestPanel) HeaderGridAdd(key, value string) {
	lcl.ThreadSync(func() {
		//在指定位置播放一行
		m.HeadersGrid.InsertColRow(false, m.HeadersGridRowCount)
		m.HeadersGrid.SetCells(0, m.HeadersGridRowCount, "1")
		m.HeadersGrid.SetCells(1, m.HeadersGridRowCount, key)
		m.HeadersGrid.SetCells(2, m.HeadersGridRowCount, value)
		m.HeadersGrid.SetCells(3, m.HeadersGridRowCount, "删除")
		m.HeadersGridRowCount++
		m.HeadersGrid.SetRowCount(m.HeadersGridRowCount)
	})
}

//清空头列表
func (m *ProxyInterceptRequestPanel) ClearHeaderGrid() {
	lcl.ThreadSync(func() {
		m.HeadersGrid.Clear()
		m.HeadersGridRowCount = 1
		m.HeadersGrid.SetRow(m.HeadersGridRowCount)
		m.HeadersGrid.SetRowCount(m.HeadersGridRowCount)
	})
}

//请求拦截参数表格头
func (m *ProxyInterceptRequestPanel) HeaderGridHead() {
	var chkBox = m.HeadersGrid.Columns().Add()
	chkBox.SetWidth(30)
	chkBox.SetButtonStyle(types.CbsCheckboxColumn)
	chkBox.Title().SetCaption("启用")

	var colNo = m.HeadersGrid.Columns().Add()
	colNo.SetWidth(180)
	colNo.Title().SetCaption("Key")
	colNo.Title().SetAlignment(types.TaCenter)
	colNo.SetAlignment(types.TaLeftJustify)

	var colAddr = m.HeadersGrid.Columns().Add()
	colAddr.SetWidth(180)
	colAddr.Title().SetCaption("Value")
	colAddr.Title().SetAlignment(types.TaCenter)
	colAddr.SetAlignment(types.TaLeftJustify)

	var delBtn = m.HeadersGrid.Columns().Add()
	delBtn.SetWidth(60)
	delBtn.Title().SetCaption("操作")
	delBtn.Title().SetAlignment(types.TaCenter)
	delBtn.SetButtonStyle(types.CbsButtonColumn)
	delBtn.SetAlignment(types.TaCenter)
}

//清空参数表格
func (m *ProxyInterceptRequestPanel) ClearQueryParamsGrid() {
	lcl.ThreadSync(func() {
		m.ParamsGrid.Clear()
		m.ParamsGridRowCount = 1
		m.ParamsGrid.SetRow(m.ParamsGridRowCount)
		m.ParamsGrid.SetRowCount(m.ParamsGridRowCount)
	})
}

//请求拦截参数列表添加
func (m *ProxyInterceptRequestPanel) QueryParamsGridAdd(key, value string) {
	lcl.ThreadSync(func() {
		//在指定位置播放一行
		m.ParamsGrid.InsertColRow(false, m.ParamsGridRowCount)
		m.ParamsGrid.SetCells(0, m.ParamsGridRowCount, "1")
		m.ParamsGrid.SetCells(1, m.ParamsGridRowCount, key)
		m.ParamsGrid.SetCells(2, m.ParamsGridRowCount, value)
		m.ParamsGrid.SetCells(3, m.ParamsGridRowCount, "删除")
		m.ParamsGridRowCount++
		m.ParamsGrid.SetRowCount(m.ParamsGridRowCount)
	})
}

//请求拦截参数表格头
func (m *ProxyInterceptRequestPanel) RequestQueryParamsGridHead() {
	var chkBox = m.ParamsGrid.Columns().Add()
	chkBox.SetWidth(30)
	chkBox.SetButtonStyle(types.CbsCheckboxColumn)
	chkBox.Title().SetCaption("启用")

	var colNo = m.ParamsGrid.Columns().Add()
	colNo.SetWidth(180)
	colNo.Title().SetCaption("Key")
	colNo.Title().SetAlignment(types.TaCenter)
	colNo.SetAlignment(types.TaLeftJustify)

	var colAddr = m.ParamsGrid.Columns().Add()
	colAddr.SetWidth(180)
	colAddr.Title().SetCaption("Value")
	colAddr.Title().SetAlignment(types.TaCenter)
	colAddr.SetAlignment(types.TaLeftJustify)

	var delBtn = m.ParamsGrid.Columns().Add()
	delBtn.SetWidth(60)
	delBtn.Title().SetCaption("操作")
	delBtn.Title().SetAlignment(types.TaCenter)
	delBtn.SetButtonStyle(types.CbsButtonColumn)
	delBtn.SetAlignment(types.TaCenter)
}

//response
func (m *ProxyInterceptResponsePanel) initUI() {
	resetPVars()
	pLeft = 0
	pTop = 25
	pWidth = m.TPanel.Width()
	pHeight = m.TPanel.Height()

	reqPageControl := lcl.NewPageControl(m.TPanel) //Tabs 的控制标签
	reqPageControl.SetParent(m.TPanel)
	reqPageControl.SetBounds(pLeft, pTop, pWidth, pHeight)
	reqPageControl.SetAlign(types.AlClient)

	sheet := lcl.NewTabSheet(reqPageControl) //标签页
	sheet.SetPageControl(reqPageControl)
	sheet.SetCaption("　Response Headers　")
	sheet.SetAlign(types.AlClient)
	headersPanel := lcl.NewPanel(m.TPanel) // 标签页
	headersPanel.SetParent(sheet)
	headersPanel.SetBounds(0, 0, pWidth, pHeight)
	headersPanel.SetAlign(types.AlClient)

	sheet = lcl.NewTabSheet(reqPageControl) //标签页
	sheet.SetPageControl(reqPageControl)
	sheet.SetCaption("　Response Body　")
	sheet.SetAlign(types.AlClient)
	bodyPanel := lcl.NewPanel(m.TPanel) // 标签页
	bodyPanel.SetParent(sheet)
	bodyPanel.SetBounds(0, 0, pWidth, pHeight)
	bodyPanel.SetAlign(types.AlClient)
}

//代理拦截配置Panel
func (m *ProxyInterceptSettingPanel) initUI() {
	m.InterceptGridConfigData = &[]*entity.ProxyInterceptConfig{}
	//开关按钮
	m.OnOffBtn = lcl.NewImageButton(m.TPanel)
	m.OnOffBtn.SetParent(m.TPanel)
	m.OnOffBtn.SetImageCount(1)
	m.OnOffBtn.SetBounds(m.TPanel.Width()-80, 10, 68, 32)
	m.OnOffBtn.SetCursor(types.CrHandPoint)
	//列表
	m.InterceptGrid = lcl.NewStringGrid(m.TPanel)
	m.InterceptGrid.SetParent(m.TPanel)
	m.InterceptGrid.SetFixedCols(0)
	m.InterceptGrid.SetFixedColor(colors.ClGreen)
	m.InterceptGrid.SetBorderStyle(types.BsNone)
	m.InterceptGrid.SetFlat(true)
	m.InterceptGrid.SetOptions(m.InterceptGrid.Options().Include(types.GoAlwaysShowEditor, types.GoCellHints, types.GoEditing, types.GoTabs, types.GoRowHighlight))
	m.InterceptGrid.SetBounds(0, 50, m.TPanel.Width(), m.TPanel.Height()-50)
	m.InterceptGrid.SetAnchors(types.NewSet(types.AkLeft, types.AkBottom, types.AkTop, types.AkRight))
	//编辑事件
	m.InterceptGrid.SetOnSetEditText(func(sender lcl.IObject, aCol, aRow int32, value string) {
		if aCol == 1 && aRow > 0 { //URL地址列
			configData := (*m.InterceptGridConfigData)[aRow-1]
			configData.SetInterceptUrl(value)
			configData.Option = consts.PIOption3
			configData.Index = aRow - 1
			entity.ProxyInterceptConfigChan <- m.InterceptGridConfigData //发送到通道
			if aRow == m.InterceptGrid.RowCount()-1 && value != "" {
				m.InterceptGridAdd("")
			}
		}
	})
	m.OnOffBtn.SetOnClick(func(sender lcl.IObject) {
		entity.ProxyInterceptConfigEnable = !entity.ProxyInterceptConfigEnable
		if entity.ProxyInterceptConfigEnable {
			m.OnOffBtn.Picture().LoadFromFSFile("resources/btn-on.png")
		} else {
			m.OnOffBtn.Picture().LoadFromFSFile("resources/btn-off.png")
		}
	})
	m.OnOffBtn.Click() //执行一次 把图片加载进来
	//列表中的按钮点击事件
	m.InterceptGrid.SetOnButtonClick(func(sender lcl.IObject, aCol, aRow int32) {
		if aCol == 2 && aRow > 0 { //删除行
			if m.InterceptGrid.RowCount() > 2 {
				var before = (*m.InterceptGridConfigData)[:aRow-1]
				var after = (*m.InterceptGridConfigData)[aRow:]
				//取出删除数据
				*m.InterceptGridConfigData = append(before, after...)
				entity.ProxyInterceptConfigChan <- m.InterceptGridConfigData //发送到通道
				m.InterceptGrid.DeleteRow(aRow)
			}
		}
	})
	//checkbox 事件
	m.InterceptGrid.SetOnCheckboxToggled(func(sender lcl.IObject, aCol, aRow int32, aState types.TCheckBoxState) {
		if aCol == 0 && aRow > 0 {
			var checked = aState == types.CbChecked
			configData := (*m.InterceptGridConfigData)[aRow-1]
			configData.Option = consts.PIOption3
			configData.SetEnable(checked)
			entity.ProxyInterceptConfigChan <- m.InterceptGridConfigData //发送到通道
		}
	})
	m.InterceptGridHead()
	m.InterceptGrid.SetRowCount(1)
	m.InterceptGridAdd("")
}

//请求拦截参数列表添加
func (m *ProxyInterceptSettingPanel) InterceptGridAdd(URL string) {
	if URL != "" {
		configData := m.InterceptGridConfigData
		for _, igcd := range *configData {
			if igcd.InterceptUrl() == URL {
				return
			}
		}
	}
	lcl.ThreadSync(func() {
		m.InterceptGridConfigDataAdd(URL)
		var count = m.InterceptGrid.RowCount()
		//在指定位置播放一行
		m.InterceptGrid.InsertColRow(false, count)
		m.InterceptGrid.SetCells(0, count, "1")
		m.InterceptGrid.SetCells(1, count, URL)
		m.InterceptGrid.SetCells(2, count, "删除")
		m.InterceptGrid.SetRowCount(count + 1)
	})
}

//拦截配置列表添加配置数据
func (m *ProxyInterceptSettingPanel) InterceptGridConfigDataAdd(URL string) {
	configData := &entity.ProxyInterceptConfig{Index: -1}
	configData.SetEnable(true)
	configData.SetInterceptUrl(URL)
	configData.Option = consts.PIOption1
	*m.InterceptGridConfigData = append(*m.InterceptGridConfigData, configData)
	entity.ProxyInterceptConfigChan <- m.InterceptGridConfigData //发送到通道
}

//请求拦截参数表格头
func (m *ProxyInterceptSettingPanel) InterceptGridHead() {
	var chkBox = m.InterceptGrid.Columns().Add()
	chkBox.SetWidth(30)
	chkBox.SetButtonStyle(types.CbsCheckboxColumn)
	chkBox.Title().SetCaption("启用")

	var colNo = m.InterceptGrid.Columns().Add()
	colNo.SetWidth(m.TPanel.Width() - 100)
	colNo.Title().SetCaption("拦截地址-URL")
	colNo.Title().SetAlignment(types.TaCenter)
	colNo.SetAlignment(types.TaLeftJustify)

	var delBtn = m.InterceptGrid.Columns().Add()
	delBtn.SetWidth(60)
	delBtn.Title().SetCaption("操作")
	delBtn.Title().SetAlignment(types.TaCenter)
	delBtn.SetButtonStyle(types.CbsButtonColumn)
	delBtn.SetAlignment(types.TaCenter)
}

//代理拦截配置Panel
func (m *ProxyInterceptPanel) initUI() {
	resetPVars()
	pLeft = 0
	pTop = 30
	pWidth = m.TPanel.Width()
	pHeight = m.TPanel.Height() - pTop

	m.interceptPageControl = lcl.NewPageControl(m.TPanel) //Tabs 的控制标签
	m.interceptPageControl.SetParent(m.TPanel)
	m.interceptPageControl.SetBounds(pLeft, pTop, pWidth, pHeight)
	m.interceptPageControl.SetAnchors(types.NewSet(types.AkLeft, types.AkBottom, types.AkTop, types.AkRight))

	sheetInterReq := lcl.NewTabSheet(m.interceptPageControl) //标签页
	sheetInterReq.SetPageControl(m.interceptPageControl)
	sheetInterReq.SetCaption("　拦截请求　")
	sheetInterReq.SetAlign(types.AlClient)
	m.ProxyInterceptRequestPanel.TPanel = lcl.NewPanel(m.TPanel) //ProxyInterceptRequestPanel 标签页
	m.ProxyInterceptRequestPanel.TPanel.SetParent(sheetInterReq)
	m.ProxyInterceptRequestPanel.TPanel.SetBounds(0, 0, pWidth, pHeight)
	m.ProxyInterceptRequestPanel.TPanel.SetAlign(types.AlClient)

	sheetInterRes := lcl.NewTabSheet(m.interceptPageControl) //标签页
	sheetInterRes.SetPageControl(m.interceptPageControl)
	sheetInterRes.SetCaption("　拦截响应　")
	sheetInterRes.SetAlign(types.AlClient)
	m.ProxyInterceptResponsePanel.TPanel = lcl.NewPanel(m.TPanel) //responsePanel 标签页
	m.ProxyInterceptResponsePanel.TPanel.SetParent(sheetInterRes)
	m.ProxyInterceptResponsePanel.TPanel.SetBounds(0, 0, pWidth, pHeight)
	m.ProxyInterceptResponsePanel.TPanel.SetAlign(types.AlClient)

	sheetInterSet := lcl.NewTabSheet(m.interceptPageControl) //标签页
	sheetInterSet.SetPageControl(m.interceptPageControl)
	sheetInterSet.SetCaption("　拦截配置　")
	sheetInterSet.SetAlign(types.AlClient)
	m.ProxyInterceptSettingPanel.TPanel = lcl.NewPanel(m.TPanel) //responsePanel 标签页
	m.ProxyInterceptSettingPanel.TPanel.SetParent(sheetInterSet)
	m.ProxyInterceptSettingPanel.TPanel.SetBounds(0, 0, pWidth, pHeight)
	m.ProxyInterceptSettingPanel.TPanel.SetAlign(types.AlClient)

	//拦截地址
	urlAddrLabel := lcl.NewLabel(m.TPanel)
	urlAddrLabel.SetParent(m.TPanel)
	urlAddrLabel.SetBounds(5, 6, 0, 0)
	urlAddrLabel.SetCaption("被拦截地址:")
	m.UrlAddrEdit = lcl.NewEdit(m.TPanel)
	m.UrlAddrEdit.SetParent(m.TPanel)
	m.UrlAddrEdit.SetReadOnly(true)
	m.UrlAddrEdit.SetBounds(75, 2, m.TPanel.Width()-80, 30)

	//状态栏标签
	state := lcl.NewStaticText(m.TPanel)
	state.SetParent(m.TPanel)
	state.SetBounds(300, pTop, 40, 20)
	state.Font().SetSize(13)
	state.Font().SetStyle(types.NewSet(types.FsBold))
	state.Font().SetColor(colors.ClBlue)
	state.SetCaption("状态: ")
	m.StateLabel = lcl.NewStaticText(m.TPanel)
	m.StateLabel.SetParent(m.TPanel)
	m.StateLabel.SetBounds(342, pTop, 180, 20)
	m.StateLabel.Font().SetSize(13)
	m.StateLabel.Font().SetStyle(types.NewSet(types.FsBold))
	m.StateLabel.Font().SetColor(0x46D12E) //绿0x46D12E 红0x8000FF
	m.StateLabel.SetCaption("--")

	//状态拦截，等待确认按钮
	m.StateOkBtn = lcl.NewButton(m.TPanel)
	m.StateOkBtn.SetParent(m.TPanel)
	m.StateOkBtn.SetCaption(" 确 认 ")
	m.StateOkBtn.Font().SetSize(12)
	m.StateOkBtn.SetBounds(m.TPanel.Width()-80, pTop-1, 70, 25)
	m.StateOkBtn.SetOnClick(func(sender lcl.IObject) {
		m.StateOkBtn.SetVisible(false)
		var state = m.State
		m.stateReset()
		//发送处理信号，InterceptQueueProxyDetail 当前正在处理的代理请求
		if state == consts.SIGNAL10 {
			m.InterceptQueueProxyDetail.ProxyInterceptSignal <- consts.SIGNAL11
			m.StateLabel.SetCaption("请求发送中...")
		} else if state == consts.SIGNAL20 {
			m.InterceptQueueProxyDetail.ProxyInterceptSignal <- consts.SIGNAL21
			m.StateLabel.SetCaption("请求响应中...")
		}
	})
	m.StateOkBtn.SetVisible(false)

	//初始化子组件
	m.ProxyInterceptRequestPanel.initUI()
	m.ProxyInterceptResponsePanel.initUI()
	m.ProxyInterceptSettingPanel.initUI()
}

//将要拦截的代理请求添加到队列 一次只处理一个
func (m *ProxyInterceptPanel) interceptQueue(proxyDetail *entity.ProxyDetail) {
	if proxyDetail != nil && proxyDetail.IsAddTaskQueue {
		//先添加到队列
		//fmt.Println("向队列添加")
		m.InterceptQueue.Push(proxyDetail)
	}
	m.mutex.Lock()
	if !m.IsUseInterceptQueue {
		//未被占用-设置成已被占用
		m.IsUseInterceptQueue = true
		m.mutex.Unlock()
		m.handlerInterceptQueue()
	} else {
		m.mutex.Unlock()
	}
}

//处理拦截队列内容 一次只处理一个
func (m *ProxyInterceptPanel) handlerInterceptQueue() {
	//取出一个
	if value, err := m.InterceptQueue.Pop(); err == nil {
		m.InterceptQueueProxyDetail = value.(*entity.ProxyDetail)
		//fmt.Println("队列剩余", m.InterceptQueue.Len(), "当前处理", m.InterceptQueueProxyDetail.TargetUrl)
		//自动处理拦截队列内容 一次只处理一个
		if m.InterceptQueueProxyDetail.State == consts.P2 { //请求处理
			m.updateRequestUI(m.InterceptQueueProxyDetail)
		} else if m.InterceptQueueProxyDetail.State == consts.P4 { //响应处理
			m.updateResponseUI(m.InterceptQueueProxyDetail)
		}
		//监听当前处理的
		go func() {
		EXIT:
			for {
				select {
				case signal, ok := <-m.InterceptQueueProxyDetail.ProxyInterceptSignal:
					if ok {
						//10:开始请求拦截 11:结束请求拦截， 20:开始响应拦截 21:结束响应拦截
						m.State = signal
						//fmt.Println("for select signal", signal)
						if signal == consts.SIGNAL10 { //10:开始请求拦截 - 阻塞请求
							m.stateOkBtnSetVisible(true)
							m.switchRequestPage()
							m.updateStateUI(0x8000FF, "请求拦截，待确认")
						} else if signal == consts.SIGNAL20 { //20:开始响应拦截 - 阻塞响应
							m.stateOkBtnSetVisible(true)
							m.switchResponsePage()
							m.updateStateUI(0x8000FF, "响应拦截，待确认")
						} else if signal == consts.SIGNAL22 { //请求超时-请求响应失败
							m.switchResponsePage()
							m.updateStateUI(0x8000FF, "请求响应失败-超时")
						} else if signal == consts.SIGNAL23 { //请求超时-响应成功
							m.switchResponsePage()
							m.updateStateUI(0x8000FF, "请求响应成功")
						} else if signal == consts.SIGNAL24 { //请求超时-响应客户端失败
							m.switchResponsePage()
							m.updateStateUI(0x8000FF, "响应客户端失败")
						} else if signal == consts.SIGNAL30 {
							break EXIT
						}
					}
				}
			}
			//关闭这个通道
			close(m.InterceptQueueProxyDetail.ProxyInterceptSignal)
			//处理完立即置空
			m.InterceptQueueProxyDetail = nil
			m.State = -1
			m.mutex.Lock()
			m.IsUseInterceptQueue = false
			m.mutex.Unlock()
			//结束 - 调出队列启用下一个拦截
			m.interceptQueue(nil)
		}()
		//通知一下proxy-server，已准备好，可以继续处理下面流程
		if m.InterceptQueueProxyDetail != nil {
			m.InterceptQueueProxyDetail.ProxyInterceptSignal <- consts.SIGNAL01
		}
	} else {
		m.mutex.Lock()
		//处理完立即置空
		m.InterceptQueueProxyDetail = nil
		m.State = -1
		m.IsUseInterceptQueue = false
		m.mutex.Unlock()
	}
}

//更新状态UI
func (m *ProxyInterceptPanel) stateOkBtnSetVisible(b bool) {
	lcl.ThreadSync(func() {
		m.StateOkBtn.SetVisible(b)
	})
}

//更新状态UI
func (m *ProxyInterceptPanel) updateStateUI(color types.TColor, caption string) {
	lcl.ThreadSync(func() {
		m.StateLabel.Font().SetColor(color)
		m.StateLabel.SetCaption(caption)
	})
}

//更新拦截到的RequestUI
func (m *ProxyInterceptPanel) updateRequestUI(proxyDetail *entity.ProxyDetail) {
	m.ProxyInterceptRequestPanel.ClearHeaderGrid()
	m.ProxyInterceptRequestPanel.ClearQueryParamsGrid()
	m.ProxyInterceptRequestPanel.TBodyPanel.ClearFormDataGrid()
	m.UrlAddrEdit.SetText(proxyDetail.TargetUrl)
	//query params
	for key, param := range proxyDetail.Request.URLParams {
		for _, p := range param {
			m.ProxyInterceptRequestPanel.QueryParamsGridAdd(key, p)
		}
	}
	//header
	for key, header := range proxyDetail.Request.Header {
		for _, p := range header {
			m.ProxyInterceptRequestPanel.HeaderGridAdd(key, p)
		}
	}
	//body - raw/json
	m.ProxyInterceptRequestPanel.TBodyPanel.RawMemo.SetText(proxyDetail.Request.Body)
	//form-data/x-www-form-urlencoded/binary
	for key, form := range proxyDetail.Request.PostForm {
		for _, p := range form {
			m.ProxyInterceptRequestPanel.TBodyPanel.FormDataGridAdd("", key, p)
		}
	}

}

//更新拦截到的ResponseUI
func (m *ProxyInterceptPanel) updateResponseUI(proxyDetail *entity.ProxyDetail) {
}

//状态相关的值和属性重置
func (m *ProxyInterceptPanel) stateReset() {
	m.State = -1
	m.StateLabel.SetCaption("--")
	m.StateLabel.Font().SetColor(0x46D12E)
}

//切换至 request sheet
func (m *ProxyInterceptPanel) switchRequestPage() {
	m.interceptPageControl.SetActivePageIndex(0)
}

//切换至 response sheet
func (m *ProxyInterceptPanel) switchResponsePage() {
	m.interceptPageControl.SetActivePageIndex(1)
}

//切换至 config sheet
func (m *ProxyInterceptPanel) switchConfigPage() {
	m.interceptPageControl.SetActivePageIndex(2)
}
