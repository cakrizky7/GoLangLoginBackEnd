package controllers

import (
	"GoLangLoginBackEnd/db"
	"GoLangLoginBackEnd/models"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
func UsersLoginCheck(c echo.Context) error {
	// Get JWT
	jwt_val := c.Get("user").(*jwt.Token)
	claims := jwt_val.Claims.(jwt.MapClaims)
	username := claims["username"].(string)

	return c.JSON(http.StatusOK, gin.H{
		"data": username,
	})
}
func UsersReLogin(c echo.Context) error {
	con := db.CreateCon()
	// filter := new(models.Filter)
	users := new(models.Users)
	hash_password := ""
	//Get The Payloads
	if err := c.Bind(users); err != nil {
		return c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Payload Error",
		})
	}
	query := "SELECT password,nama_lengkap FROM users WHERE username='" + users.Username + "' AND token='" + users.Token + "' AND token_expired >= UTC_TIMESTAMP()"
	rows, err := con.Query(query)
	defer con.Close()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"in":      "user not found",
		})
	}
	for rows.Next() {
		if err := rows.Scan(&hash_password, &users.Nama_lengkap); err != nil {
			return c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})

		} else {
			// Create token
			token := jwt.New(jwt.SigningMethodHS256)

			// Set claims
			claims := token.Claims.(jwt.MapClaims)
			claims["username"] = users.Username
			claims["exp"] = time.Now().Add(time.Minute * 1).Unix()

			// Generate encoded token and send it as response.
			t, err := token.SignedString([]byte("secret"))
			if err != nil {
				return c.JSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
				})
			}
			claims["username"] = users.Username
			token_expired := time.Now().Add(time.Minute * 30)
			users.Token_expired = token_expired
			claims["exp"] = token_expired.Unix()
			token_refresh, err := token.SignedString([]byte("refresh"))
			if err != nil {
				return c.JSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
				})
			}
			//Save Refresh Token
			stmt, err := con.Prepare("update " + users.TableName() + " set token=?, token_expired=? where username=?")
			if err != nil {
				return c.JSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
				})
			}
			_, err = stmt.Exec(token_refresh, users.Token_expired, users.Username)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
				})
			}
			return c.JSON(http.StatusOK, gin.H{
				"message":      "done",
				"nama":         users.Nama_lengkap,
				"username":     users.Username,
				"token":        t,
				"refreshToken": token_refresh,
			})
		}
	}
	defer con.Close()

	return c.JSON(http.StatusUnauthorized, gin.H{
		"message": "error",
	})
}
func UsersLogin(c echo.Context) error {
	con := db.CreateCon()      // Prepare the DB connection
	users := new(models.Users) // Prepare variable for payload
	hash_password := ""        // Prepare variable for hash the payload password

	//Get The Payloads
	if err := c.Bind(users); err != nil {
		return c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Payload Error",
		})
	}
	query := "SELECT password,nama_lengkap FROM  " + users.TableName() + "  WHERE username='" + users.Username + "' LIMIT 1"
	rows, err := con.Query(query)
	// defer con.Close()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	for rows.Next() {

		if err := rows.Scan(&hash_password, &users.Nama_lengkap); err != nil {
			return c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})

		} else {
			if comparePasswords(hash_password, []byte(users.Password)) {
				// Create tokens
				token := jwt.New(jwt.SigningMethodHS256)

				// Set claims
				claims := token.Claims.(jwt.MapClaims)
				claims["username"] = users.Username
				claims["exp"] = time.Now().Add(time.Minute * 1).Unix()

				// Generate encoded token and send it as response.
				t, err := token.SignedString([]byte("secret"))
				if err != nil {
					return c.JSON(http.StatusInternalServerError, gin.H{
						"message": err.Error(),
					})
				}

				// Generate refresh Token
				claims["username"] = users.Username
				token_expired := time.Now().Add(time.Minute * 30)
				users.Token_expired = token_expired
				claims["exp"] = token_expired.Unix()
				token_refresh, err := token.SignedString([]byte("refresh"))
				if err != nil {
					return c.JSON(http.StatusInternalServerError, gin.H{
						"message": err.Error(),
					})
				}

				//Save Refresh Token to DB
				fmt.Println(token_expired)
				stmt, err := con.Prepare("update " + users.TableName() + " set token=?, token_expired=? where username=?")
				if err != nil {
					return c.JSON(http.StatusInternalServerError, gin.H{
						"message": err.Error(),
					})
				}
				_, err = stmt.Exec(token_refresh, token_expired, users.Username)
				if err != nil {
					return c.JSON(http.StatusInternalServerError, gin.H{
						"message": err.Error(),
					})
				}

				//Return the Data
				return c.JSON(http.StatusOK, gin.H{
					"message":      "done",
					"nama":         users.Nama_lengkap,
					"username":     users.Username,
					"token":        t,
					"refreshToken": token_refresh,
				})
			}
		}
	}
	defer con.Close()

	return c.JSON(http.StatusUnauthorized, gin.H{
		"message": "username or password error",
	})
}

func UsersRegister(c echo.Context) error {
	con := db.CreateCon()
	users := new(models.Users)
	//Get The Payloads
	if err := c.Bind(users); err != nil {
		return err
	}
	// Modif the Model that has been filled with the payloads
	users.Id = uuid.Must(uuid.NewRandom())
	users.Created_at = time.Now()
	users.Updated_at = time.Now()
	users.Password = hashAndSalt([]byte(users.Password))
	if err := c.Validate(users); err != nil {
		return err
	}
	insForm, err := con.Prepare("INSERT INTO " + users.TableName() + " (id,username,password,nama_lengkap,created_at,updated_at) VALUES(?,?,?,?,?,?)")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	insForm.Exec(users.Id, users.Username, users.Password, users.Nama_lengkap, users.Created_at, users.Updated_at)
	defer con.Close()
	return c.JSON(http.StatusOK, gin.H{
		"message":  "done",
		"nama":     users.Nama_lengkap,
		"username": users.Username,
	})
}
