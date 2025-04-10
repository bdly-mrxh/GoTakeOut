package admin

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"takeout/common/constant"
	"takeout/common/logger"
	"takeout/common/response"
	"takeout/common/utils"
)

// CommonController 通用控制器
type CommonController struct{}

// NewCommonController 创建通用控制器
func NewCommonController() *CommonController {
	return &CommonController{}
}

// RealUpload 文件上传
func (c *CommonController) RealUpload(ctx *gin.Context) {
	// 获取上传的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		logger.Error("获取上传文件失败", zap.Error(err))
		response.BadRequest(ctx, "请选择要上传的文件")
		return
	}

	// 创建OSS上传工具
	ossUploader, err := utils.NewOSSUploader()
	if err != nil {
		logger.Error("创建OSS上传工具失败", zap.Error(err))
		response.ServerError(ctx, "文件上传服务初始化失败")
		return
	}

	// 上传文件到OSS
	fileUrl, err := ossUploader.UploadFile(file)
	if err != nil {
		logger.Error("上传文件到OSS失败", zap.Error(err))
		response.ServerError(ctx, "文件上传失败")
		return
	}

	// 返回上传成功的文件信息
	logger.Info("文件上传成功", zap.String("url", fileUrl))
	response.Success(ctx, constant.MsgSuccess, fileUrl)
}

// Upload 跳过上传图片逻辑
func (c *CommonController) Upload(ctx *gin.Context) {
	uploadFile, err := ctx.FormFile("file")
	if err != nil || uploadFile == nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}
	path := "https://photo.jpg"
	response.Success(ctx, constant.MsgSuccess, path)
}
