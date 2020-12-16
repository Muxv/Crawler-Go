package models

import (
	"strconv"
	"strings"
	"time"
)

type BookInfo struct {
	Topk      int     `gorm:"primaryKey" json:"top"`
	ChName    string  `json:"ch_name"`
	EnName    string  `json:"en_name"`
	BasicInfo string  `json:"basic_info"`
	Rank      float64 `json:"rank"`
	RankNum   int     `json:"rank_num"`
	Comment   string  `json:"comment"`
}

type Top struct {
	BookInfo
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (b BookInfo) String() string {
	var builder strings.Builder
	builder.WriteString("排名: " + strconv.Itoa(b.Topk) + "\n")
	builder.WriteString("书名: " + b.ChName + "\n")
	if b.EnName != "" {
		builder.WriteString("英文名: " + b.EnName + "\n")
	}
	builder.WriteString("书本基本信息: " + b.BasicInfo + "\n")
	rStr := strconv.FormatFloat(b.Rank, 'f', 1, 64)
	builder.WriteString("评价: " + rStr + "(" + strconv.Itoa(b.RankNum) + "人评价)\n")
	builder.WriteString("评语: " + b.Comment + "\n")

	return builder.String()
}
