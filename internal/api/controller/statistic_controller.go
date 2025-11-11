package controller

import (
	"github.com/CRS-Project/crs-backend/internal/api/service"
	"github.com/CRS-Project/crs-backend/internal/pkg/meta"
	"github.com/CRS-Project/crs-backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type (
	StatisticController interface {
		GetAOCAndCommentChart(ctx *gin.Context)
		GetCommentCard(ctx *gin.Context)
		GetCommentUserChart(ctx *gin.Context)
		GetCommentUserData(ctx *gin.Context)
	}

	statisticController struct {
		statisticService service.StatisticService
	}
)

func NewStatistic(statisticService service.StatisticService) StatisticController {
	return &statisticController{
		statisticService: statisticService,
	}
}

func (c *statisticController) GetAOCAndCommentChart(ctx *gin.Context) {
	packageId := ctx.Param("package_id")
	res, err := c.statisticService.GetAOCAndCommentChart(ctx.Request.Context(), packageId)
	if err != nil {
		response.NewFailed("failed to get statistic", err).Send(ctx)
		return
	}

	response.NewSuccess("success get statistic", res).Send(ctx)
}

func (c *statisticController) GetCommentCard(ctx *gin.Context) {
	packageId := ctx.Param("package_id")
	res, err := c.statisticService.GetCommentCard(ctx.Request.Context(), packageId)
	if err != nil {
		response.NewFailed("failed to get statistic", err).Send(ctx)
		return
	}

	response.NewSuccess("success get statistic", res).Send(ctx)
}

func (c *statisticController) GetCommentUserChart(ctx *gin.Context) {
	packageId := ctx.Param("package_id")
	res, err := c.statisticService.GetCommentUserChart(ctx.Request.Context(), packageId)
	if err != nil {
		response.NewFailed("failed to get statistic", err).Send(ctx)
		return
	}

	response.NewSuccess("success get statistic", res).Send(ctx)
}

func (c *statisticController) GetCommentUserData(ctx *gin.Context) {
	packageId := ctx.Param("package_id")
	res, metares, err := c.statisticService.GetCommentUserData(ctx.Request.Context(), packageId, meta.New(ctx))
	if err != nil {
		response.NewFailed("failed to get statistic", err).Send(ctx)
		return
	}

	response.NewSuccess("success get statistic", res, metares).Send(ctx)
}
