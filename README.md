# Gin Gorm
An gin middleware to provide every gin action an gorm transaction

# Install
`go get github.com/wppurking/gin-gorm`


# Usage
1. First `gin.Use` this middleware
2. In the `gin Action` method use `gin_gorm.Tx(c)` to get the gorm.DB transactino to user.
```
func main() {
	router := gin.Default()

	router.POST("/form_post", func(c *gin.Context) {
		message := c.PostForm("message")
		nick := c.DefaultPostForm("nick", "anonymous")

		tx := gin_gorm.Tx(c)
		var user model.User
		tx.Where("nick = ?", nick).First(&user)

		c.JSON(200, gin.H{
			"status":  "posted",
			"message": message,
			"nick":    nick,
			"user":    user,
		})
	})
	router.Run(":8080")
}

``` 
