package handlers

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"strconv"

	"github.com/devhalfdog/dynocord-api/database"
	e "github.com/devhalfdog/dynocord-api/errors"
	"github.com/devhalfdog/dynocord-api/twitch"
	"github.com/devhalfdog/dynocord-api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/nfnt/resize"
)

var (
	model = database.Model{}
)

func Initial() error {
	model = *database.New(
		utils.Environment("CHAT_DB_IP"),
		utils.Environment("CHAT_DB_PORT"),
		utils.Environment("CHAT_DB_USER"),
		utils.Environment("CHAT_DB_PASSWORD"),
		utils.Environment("CHAT_DB_DB"),
	)

	err := model.Connect()
	if err != nil {
		return err
	}

	return nil
}

// 방송 스크린 캡쳐
func GetStreamCapture(ctx *fiber.Ctx) error {
	streamer := ctx.Params("streamer")
	headers := ctx.GetReqHeaders()

	if headers["Client-Id"] != utils.Environment("STREAM_CLIENT") {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Does not match the token",
		})
	}

	file, err := twitch.GetStreamScreenShot(streamer)
	if err != nil {
		return ctx.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	// imageName := strings.Split(file, string(os.PathSeparator))
	imgUrl, err := UploadImage(file, streamer)
	if err != nil {
		return ctx.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "failed image upload",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": imgUrl,
	})
}

// 이미지 리사이즈
func ResizeImage(ctx *fiber.Ctx) error {
	streamer := ctx.Params("streamer")
	widthStr := ctx.Params("width")
	heightStr := ctx.Params("height")
	headers := ctx.GetReqHeaders()

	if headers["Client-Id"] != utils.Environment("STREAM_CLIENT") {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Does not match the token",
		})
	}

	if streamer == "" || widthStr == "" || heightStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "does not input parameters",
		})
	}

	width, err := strconv.Atoi(widthStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "invaild width",
		})
	}

	height, err := strconv.Atoi(heightStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "invaild height",
		})
	}

	file, err := os.Open(fmt.Sprintf("./static/%s.png", streamer))
	if err != nil {
		return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{
			"error":   true,
			"message": "does not file exists",
		})
	}

	img, err := png.Decode(file)
	if err != nil {
		fmt.Println(err)
		return ctx.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "does not image decode",
		})
	}

	resizeFilePath := fmt.Sprintf("./static/resize/%s.png", streamer)
	resizeM := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
	out, err := os.Create(resizeFilePath)
	if err != nil {
		return ctx.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "failed resizing to image",
		})
	}
	defer out.Close()

	err = png.Encode(out, resizeM)
	if err != nil {
		return ctx.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "failed encoding to image",
		})
	}

	err = os.Remove(fmt.Sprintf("./static/%s.png", streamer))
	if err != nil {
		log.Println(err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": fmt.Sprintf("/static/resize/%s.png", streamer),
	})
}

// 채팅 불러오기
func GetChat(ctx *fiber.Ctx) error {
	streamer := ctx.Params("streamer")
	before := ctx.Params("before")

	var bInt int64 = 0
	var err error

	if before != "" {
		bInt, err = strconv.ParseInt(before, 10, 32)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "cannot parse int",
			})
		}
	}

	chat, err := model.GetChat(streamer, bInt)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": e.ErrGetChat,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": chat,
	})
}

// 채팅 저장
func SaveChat(ctx *fiber.Ctx) error {
	c := new(database.Chat)

	if err := ctx.BodyParser(&c); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": e.ErrBindJSON,
		})
	}

	err := model.CreateChat(c)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": e.ErrCreateChat,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"message": c,
	})
}
