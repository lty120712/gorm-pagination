# gorm-pagination

`gorm-pagination` 是一个自定义的 GORM 分页插件，提供了简单的分页功能，适用于 GORM 的数据库查询。该插件能够自动处理分页查询的 SQL 操作，支持灵活的分页配置，并且可以与 GORM 的其他操作配合使用。

## 安装

首先，确保你已经安装了 Go 并且已经初始化了 Go Modules。

然后，使用 `go get` 命令将该分页插件添加到你的项目中：

```bash
go get github.com/lty120712/gorm-pagination
```
使用方法
1. 引入包
在你的 Go 文件中，首先引入该分页插件：
```go
import "github.com/lty120712/gorm-pagination"
```
2. 定义分页请求和返回结构体
创建分页请求和返回结果的结构体。请求结构体用于接收分页参数，返回结构体用于存储分页结果。
```go
// 分页请求结构体
type PageRequest struct {
	Page     int `json:"page"`     // 当前页
	PageSize int `json:"pageSize"` // 每页数据条数
}

// 分页结果结构体
type PageResult[T any] struct {
	Records  []T   `json:"records"`  // 当前页数据
	Total    int64 `json:"total"`    // 总记录数
	Page     int   `json:"page"`     // 当前页
	PageSize int   `json:"pageSize"` // 每页数据条数
}
```
3. 执行分页查询
在你的控制器或者业务逻辑中，调用 gorm-pagination 插件的 Paginate 函数来执行分页查询。
```go
package controller

import (
	"github.com/lty120712/gorm-pagination"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"your_project/models"
)

// 查询角色并分页
func (con RoleController) FindHandler(c *gin.Context) {
	var req models.PageRequest
	var roleList []models.Role

	// 解析请求中的分页参数
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		con.Error(c, err.Error())
		return
	}

	// 调用分页函数
	var result models.PageResult[models.Role]
	tx := db.Db.Model(&models.Role{})
	if err := gormpagination.Paginate(tx, req.Page, req.PageSize, &result); err != nil {
		con.Error(c, err.Error())
		return
	}

	// 返回分页数据
	con.Success(c, result)
}
```
4. 原理
该分页插件的原理基于以下步骤：
```text
1. 分页参数校验：首先确保传入的 page 和 pageSize 参数是合法的。page 必须大于 0，如果传入的 page 小于等于 0，插件会将其纠正为 1。pageSize 必须大于 0，如果传入的 pageSize 小于等于 0，插件会将其纠正为 10。

2. 获取总记录数：插件会执行一次 COUNT 查询，以获取符合条件的总记录数（total）。该步骤用于计算总页数以及分页的边界。

3. 分页查询：插件使用 LIMIT 和 OFFSET 来进行分页查询。具体的 LIMIT 为每页数据条数，OFFSET 根据当前页和每页条数来计算。

4. 返回结果：查询到的数据会被存储在 result 中，其中包括：
5. 返回:
    Records：当前页的数据列表。
    Total：总记录数。
    Page：当前页码。
    PageSize：每页条数。
```

5. 配置与自定义

你可以根据自己的需求进一步配置插件或调整分页逻辑。以下是一些常见的自定义需求和解决方案。

5.1. 自定义默认值

在插件内部，我们已经对 page 和 pageSize 做了默认值的处理。如果你想要修改默认值，可以自行修改插件代码中的校验逻辑。

5.2. 分页后返回 tx 对象

如果你需要在分页操作后继续使用 tx 对象进行其他查询，可以将 tx 作为返回值获取：
```go
tx, err := gormpagination.Paginate(tx, req.Page, req.PageSize, &result)
if err != nil {
    // 处理错误
}
// 可以继续使用 tx 进行其他操作
```
5.3 自定义查询条件

你可以通过在传入的 tx 上继续链式调用 Where、Order 等方法，添加自定义的查询条件。例如：
```go
tx := db.Db.Model(&models.Role{}).Where("status = ?", "active")
tx, err := gormpagination.Paginate(tx, req.Page, req.PageSize, &result)
```
示例
以下是一个完整的分页查询例子，展示了如何使用该插件来分页查询角色数据。
``` go
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lty120712/gorm-pagination"
	"gorm.io/gorm"
	"your_project/models"
)

func main() {
	// 初始化数据库连接（假设已配置）
	db, err := gorm.Open(/* 数据库配置 */)
	if err != nil {
		panic("failed to connect database")
	}

	// 初始化 Gin 路由
	r := gin.Default()

	// 定义查询角色的分页接口
	r.POST("/roles", func(c *gin.Context) {
		var req models.PageRequest
		var result models.PageResult[models.Role]

		// 获取分页请求参数
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// 执行分页查询
		tx := db.Model(&models.Role{})
		if err := gormpagination.Paginate(tx, req.Page, req.PageSize, &result); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// 返回分页结果
		c.JSON(200, result)
	})

	// 启动服务器
	r.Run(":8080")
}
```
常见问题
1. 分页查询返回空数据怎么办？

    如果查询没有符合条件的数据，result.Records 将为空数组。如果总记录数为零，你的结果会显示空数据。这是正常现象，意味着没有符合条件的数据。

2. 如何处理分页查询时的错误？

    在分页查询时，我们使用了 tx.Error 来捕获数据库操作错误。你可以在调用 Paginate 函数时处理这些错误，或者直接将错误传递给前端。

3. 如何修改分页查询的默认参数？

    默认情况下，分页插件会将 page 和 pageSize 设置为 1 和 10。如果你想修改这些默认值，可以直接在插件代码中调整它们的初始值。