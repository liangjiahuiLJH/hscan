package server

import (
	"hscan/schema"
	"net/http"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gin-gonic/gin"
)

func (s *Server) format(txs []*schema.Transaction) {

	for i := range txs {
		var logs sdk.ABCIMessageLogs
		s.cdc.UnmarshalJSON([]byte(txs[i].RawMessages), &logs)
		s.l.Printf("log is %+v", logs)
		txs[i].Messages = logs
	}

}

func (s *Server) txs(c *gin.Context) {
	limit := c.DefaultQuery("limit", "5")
	iLimit, _ := strconv.ParseInt(limit, 10, 64)
	if iLimit <= 0 {
		iLimit = 5
	}

	var txs []*schema.Transaction

	if err := s.db.Order("id DESC").Limit(iLimit).Find(&txs).Error; err != nil {
		s.l.Printf("query blocks from db failed")
	}

	s.format(txs)

	c.JSON(http.StatusOK, gin.H{
		"paging": map[string]interface{}{
			"total":  1,
			"before": 2,
			"after":  3,
		},
		"data": txs,
	})

}
