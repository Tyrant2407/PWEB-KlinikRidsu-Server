package login

import (
	"net/http"
	"time"
	"KlinikRidsu/databases"
	"KlinikRidsu/session"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func loginHandler(c *gin.Context, db *gorm.DB, user databases.User) {
	expirationTime := time.Now().Add(time.Hour)
	expiresIn := time.Until(expirationTime)
	err := session.SetSession(c.Writer, c.Request, "user", map[string]interface{}{
		"nors":          user.NoRS,
		"nik":           user.NIK,
		"nama":          user.Nama,
		"tgl_lahir":     user.TglLahir,
		"jenis_kelamin": user.JenisKelamin,
		"expires_in": expiresIn.String(),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan sesi!"})
		return
	}

	token, err := session.GenerateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghasilkan token!"})
		return
	}

	tokenString, ok := token["token"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan token!"})
		return
	}
	c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)
	// c.SetCookie("token", tokenString, 3600, "/", "xhaeffer.me:11092", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Login Berhasil!",
		"token":   token,
		"session" : convertMap(session.GetSession(c.Request)),
	})
}