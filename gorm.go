package gin_gorm

const GormTx = "gorm_tx"

// GormMiddleware 为 gin 提供 gorm 的中间件.
// 1. 内部 panic 事务退回, 并且继续向上 panic
// 2. *gorm.DB 有错误, 事务退回
// 3. gin.Context 有错误, 事务退回
// 4. 否则, 事务正常提交
func GormMiddleware(tx *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		itx := tx.New().Begin()
		if itx.Error != nil {
			// 事务开启失败, 交给 gin 的最外层 Recover 去处理
			panic("gorm.Tx fails")
		}
		// Set example variable
		c.Set(GormTx, itx)

		defer func() {
			// 如果内部有 painc, 回退事务, 继续向上 panic
			if p := recover(); p != nil {
				itx.Rollback()
				panic(p)
			} else if itx.Error != nil {
				itx.Rollback()
			} else if len(c.Errors) > 0 {
				itx.Rollback()
			} else if c.Writer.Status() < 200 || c.Writer.Status() >= 300 {
				itx.Rollback()
			} else {
				itx.Commit()
			}
		}()

		c.Next()
	}
}
