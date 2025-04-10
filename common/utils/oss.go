package utils

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	iUtils "github.com/iWyh2/go-myUtils/utils"
	"go.uber.org/zap"
	"mime/multipart"
	"path/filepath"
	"takeout/common/constant"
	"takeout/common/errs"
	"takeout/common/global"
	"takeout/common/logger"
)

// OSSUploader 阿里云OSS上传工具
type OSSUploader struct {
	client      *oss.Client
	bucket      *oss.Bucket
	ossConfig   *global.OSSConfig
	initialized bool
}

// NewOSSUploader 创建OSS上传工具
func NewOSSUploader() (*OSSUploader, error) {
	ossConfig := &global.Config.OSS

	// 创建OSSClient实例
	client, err := oss.New(ossConfig.Endpoint, ossConfig.AccessKeyID, ossConfig.AccessKeySecret)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeInternalError, "创建OSS客户端失败")
	}

	// 获取存储空间
	bucket, err := client.Bucket(ossConfig.BucketName)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeInternalError, "获取OSS Bucket失败")
	}

	return &OSSUploader{
		client:      client,
		bucket:      bucket,
		ossConfig:   ossConfig,
		initialized: true,
	}, nil
}

// UploadFile 上传文件到OSS
func (u *OSSUploader) UploadFile(file *multipart.FileHeader) (string, error) {
	if !u.initialized {
		return "", errs.New(constant.CodeInternalError, "OSS上传工具未初始化")
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		return "", errs.Wrap(err, constant.CodeInternalError, "打开文件失败")
	}
	defer func(src multipart.File) {
		err = src.Close()
		if err != nil {
			logger.Error("文件关闭失败", zap.Error(err))
		}
	}(src)

	// 获取原文件的扩展名
	extension := filepath.Ext(file.Filename)
	// 创建新文件名称
	objectName := iUtils.UUID() + extension

	// 上传文件到OSS
	err = u.bucket.PutObject(objectName, src)
	if err != nil {
		return "", errs.Wrap(err, constant.CodeInternalError, "上传文件到OSS失败")
	}

	// 返回文件URL
	return fmt.Sprintf("https://%s.%s/%s", u.ossConfig.BucketName, u.ossConfig.Endpoint, objectName), nil
}
