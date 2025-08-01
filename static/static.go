package static

import (
	"context"
	"embed"
	"io/fs"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

//go:embed dist/*
var StaticFiles embed.FS

// SetupStaticRoutes 设置静态文件路由
func SetupStaticRoutes(h *server.Hertz) error {
	// 获取嵌入的文件系统，去掉 dist 前缀
	staticFS, err := fs.Sub(StaticFiles, "dist")
	if err != nil {
		return err
	}

	// 使用 Hertz 官方推荐的 StaticFS 方法
	h.StaticFS("/", &app.FS{
		Root:               "./static/dist",         // 使用文件路径
		PathRewrite:        app.NewPathSlashesStripper(0), // 不移除任何路径前缀
		IndexNames:         []string{"index.html"},        // 默认索引文件
		CacheDuration:      time.Hour * 24,                // 缓存24小时
		Compress:           true,                           // 启用压缩
		AcceptByteRange:    true,                          // 支持范围请求
		GenerateIndexPages: false,                         // 不生成目录列表
		PathNotFound: func(ctx context.Context, c *app.RequestContext) {
			// SPA 应用支持：所有非文件请求都返回 index.html
			path := string(c.Path())
			
			// 如果是 API 路径，返回 404
			if len(path) >= 4 && path[:4] == "/api" {
				c.String(404, "API endpoint not found")
				return
			}
			
			// 如果请求路径不包含文件扩展名，返回 index.html 支持前端路由
			if !containsFileExtension(path) {
				indexData, err := fs.ReadFile(staticFS, "index.html")
				if err != nil {
					c.String(404, "File not found")
					return
				}
				c.Header("Content-Type", "text/html; charset=utf-8")
				c.Data(200, "text/html; charset=utf-8", indexData)
				return
			}
			
			// 其他情况返回 404
			c.String(404, "File not found")
		},
	})

	return nil
}

// containsFileExtension 检查路径是否包含文件扩展名
func containsFileExtension(path string) bool {
	commonExts := []string{".html", ".css", ".js", ".json", ".png", ".jpg", ".jpeg", ".gif", ".svg", ".ico", ".woff", ".woff2", ".ttf", ".eot", ".map", ".txt", ".xml"}
	
	for _, ext := range commonExts {
		if len(path) >= len(ext) && path[len(path)-len(ext):] == ext {
			return true
		}
	}
	
	return false
}