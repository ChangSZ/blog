package service

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/ChangSZ/golib/log"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/ChangSZ/blog/conf"
)

type MinioService struct{}

func NewMinio() *MinioService {
	return &MinioService{}
}

func (s *MinioService) Upload(ctx context.Context, file multipart.File, filename string) (string, error) {
	// 创建Minio客户端对象
	minioClient, err := minio.New(conf.Cnf.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.Cnf.MinioAccessKey, conf.Cnf.MinioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return "", fmt.Errorf("创建客户端失败: %w", err)
	}
	found, err := minioClient.BucketExists(ctx, conf.Cnf.MinioBucketName)
	if err != nil {
		return "", err
	}
	if !found {
		// 创建存储桶并设置只读权限
		if err := minioClient.MakeBucket(ctx, conf.Cnf.MinioBucketName, minio.MakeBucketOptions{}); err != nil {
			return "", fmt.Errorf("新建桶失败: %w", err)
		}
		// 设置存储桶策略
		policy, err := s.createBucketPolicy(conf.Cnf.MinioBucketName)
		if err != nil {
			return "", fmt.Errorf("创建桶策略失败: %w", err)
		}
		if err := minioClient.SetBucketPolicy(ctx, conf.Cnf.MinioBucketName, policy); err != nil {
			return "", fmt.Errorf("设置桶策略失败: %w", err)
		}
	} else {
		log.WithTrace(ctx).Info("存储桶已经存在！")
	}

	// 生成存储对象的名称
	if filename == "" {
		filename = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	objectName := fmt.Sprintf("%s/%s", time.Now().Format("20060102"), filename)
	// 上传文件到Minio
	_, err = minioClient.PutObject(ctx, conf.Cnf.MinioBucketName, objectName, file, -1, minio.PutObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("PUT文件失败: %w", err)
	}
	log.WithTrace(ctx).Infof("文件上传成功: %v", objectName)
	// 构建返回结果
	url := fmt.Sprintf("%s/console/minio/file?bucket=%s&objectName=%s",
		conf.Cnf.AppUrl, conf.Cnf.MinioBucketName, objectName)
	return url, nil
}

/**
 * 创建存储桶的访问策略，设置为只读权限
 */
func (s *MinioService) createBucketPolicy(bucketName string) (string, error) {
	policy := map[string]interface{}{
		"version": time.Now().Format("2006-01-02"),
		"statement": []map[string]interface{}{
			{
				"effect":    "Allow",
				"principal": "*",
				"action":    "s3:GetObject",
				"resource":  fmt.Sprintf("arn:aws:s3:::%s/*.**", bucketName),
			},
		},
	}
	policyJSON, err := json.Marshal(policy)
	return string(policyJSON), err
}

func (s *MinioService) Delete(ctx context.Context, objectName string) error {
	// 创建Minio客户端对象
	minioClient, err := minio.New(conf.Cnf.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.Cnf.MinioAccessKey, conf.Cnf.MinioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return fmt.Errorf("创建客户端失败: %w", err)
	}

	// 删除文件
	return minioClient.RemoveObject(ctx, conf.Cnf.MinioBucketName, objectName, minio.RemoveObjectOptions{})
}

func (s *MinioService) PresignedURL(ctx context.Context, bucketName, objectName string) (string, error) {
	// 创建Minio客户端对象
	minioClient, err := minio.New(conf.Cnf.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.Cnf.MinioAccessKey, conf.Cnf.MinioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return "", fmt.Errorf("创建客户端失败: %w", err)
	}

	// 生成预签名URL
	presignedURL, err := minioClient.PresignedGetObject(ctx, bucketName, objectName, 24*time.Hour, nil)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}
