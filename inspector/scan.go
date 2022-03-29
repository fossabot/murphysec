package inspector

import (
	"fmt"
	"io/fs"
	"murphysec-cli-simple/api"
	"murphysec-cli-simple/display"
	"murphysec-cli-simple/logger"
	"path/filepath"
	"strings"
)

func Scan(ctx *ScanContext) (interface{}, error) {
	ui := ctx.UI()
	if e := ctx.FillProjectInfo(); e != nil {
		logger.Debug.Printf("v+\n", e)
		logger.Err.Println(e)
		return nil, ErrGetProjectInfo
	}
	ui.Display(display.MsgInfo, fmt.Sprint("项目名称：", ctx.ProjectName))
	ui.UpdateStatus(display.StatusRunning, "正在创建扫描任务，请稍候······")
	if e := createTask(ctx); e != nil {
		ui.Display(display.MsgError, fmt.Sprint("项目创建失败：", e.Error()))
		logger.Err.Println("Create task failed.", e.Error())
		logger.Debug.Printf("%+v", e)
		return nil, e
	}
	ui.Display(display.MsgInfo, fmt.Sprint("项目创建成功，项目唯一标识：", ctx.TaskId))
	ui.UpdateStatus(display.StatusRunning, "正在进行扫描...")
	if e := managedInspectScan(ctx); e != nil {
		logger.Debug.Println("Managed inspect failed.", e.Error())
		logger.Debug.Printf("%v", e)
	}

	{
		// todo: refactor
		enableCxx := false
		filepath.Walk(ctx.ProjectDir, func(path string, info fs.FileInfo, err error) error {
			if enableCxx {
				return filepath.SkipDir
			}
			if strings.HasPrefix(info.Name(), ".") {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			enableCxx = CxxExtSet[filepath.Ext(info.Name())]
			return nil
		})
		if enableCxx {
			FileHashScan(ctx)
		}
	}

	ui.UpdateStatus(display.StatusRunning, "项目扫描结束，正在提交信息...")
	if e := submitModuleInfo(ctx); e != nil {
		ui.Display(display.MsgError, fmt.Sprint("信息提交失败：", e.Error()))
		logger.Debug.Printf("%+v", e)
		logger.Err.Println(e.Error())
	}
	if ctx.EnableDeepScan && shouldUploadFile(ctx) {
		logger.Info.Printf("deep scan enabled, upload source code")
		ui.UpdateStatus(display.StatusRunning, "正在上传文件到服务端以进行深度检测")
		if e := UploadCodeFile(ctx); e != nil {
			ui.Display(display.MsgError, "上传文件失败："+e.Error())
		} else {
			ui.Display(display.MsgInfo, "上传文件成功")
		}
		ui.ClearStatus()
	}
	ui.UpdateStatus(display.StatusRunning, "检测中，等待返回结果...")

	if e := api.StartCheck(ctx.TaskId); e != nil {
		ui.Display(display.MsgError, "启动检测失败："+e.Error())
		logger.Err.Println("send start check command failed.", e.Error())
		return nil, e
	}
	ui.ClearStatus()
	resp, e := api.QueryResult(ctx.TaskId)
	ui.ClearStatus()
	if e != nil {
		ui.Display(display.MsgError, "获取检测结果失败："+e.Error())
		logger.Err.Println("query result failed.", e.Error())
		return nil, e
	}
	ctx.ScanResult = resp
	return nil, nil
}
