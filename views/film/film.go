package film

import (
	"github.com/labstack/echo"
	"net/http"
	"github.com/lempiy/echo_api/types"
	"github.com/lempiy/echo_api/models"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

func Add(c echo.Context) error {
	f := new(types.PostFilm)
	if err := c.Bind(f); err != nil {
		return err
	}

	err := models.Film.Create(f.Film, f.Genres)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Something went wrong",
		})
	}
	return c.JSON(http.StatusOK, map[string]bool{
		"success": true,
	})
}

func Get(c echo.Context) error {
	g := new(types.GetFilmParams)
	if err := c.Bind(g); err != nil {
		return err
	}
	films, err := models.Film.Read(g)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Something went wrong",
		})
	}
	return c.JSON(http.StatusOK, films)
}

type rentData struct {
	FilmID int `json:"film_id"`
}

func Rent(c echo.Context) error {
	r := new(rentData)
	if err := c.Bind(r); err != nil {
		return err
	}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID, isOK := claims["user_id"]
	if !isOK {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "Wrong JWT claims",
		})
	}

	notExist, err := models.Film.Rent(r.FilmID, int(userID.(float64)))

	if notExist {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "Film with such ID doesn't exist.",
		})
	}

	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Something went wrong",
		})
	}
	return c.JSON(http.StatusOK, map[string]bool{
		"success": true,
	})
}

func FinishRent(c echo.Context) error {
	r := new(rentData)
	if err := c.Bind(r); err != nil {
		return err
	}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID, isOK := claims["user_id"]
	if !isOK {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "Wrong JWT claims",
		})
	}

	notExist, err := models.Film.FinishRent(r.FilmID, int(userID.(float64)))

	if notExist {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "Film with such ID doesn't exist or had never been rented by user.",
		})
	}

	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Something went wrong",
		})
	}
	return c.JSON(http.StatusOK, map[string]bool{
		"success": true,
	})
}

func GetRentedFilms(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID, isOK := claims["user_id"]
	if !isOK {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "Wrong JWT claims",
		})
	}
	g := new(types.GetFilmParams)
	if err := c.Bind(g); err != nil {
		return err
	}
	films, err := models.Film.ReadRentedFilms(int(userID.(float64)), g)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Something went wrong",
		})
	}
	return c.JSON(http.StatusOK, films)
}
