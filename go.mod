module Crawler-go

go 1.15

require (
	github.com/PuerkitoBio/goquery v1.6.0
	github.com/go-ini/ini v1.62.0
	gorm.io/driver/mysql v1.0.3
	gorm.io/gorm v1.20.8
)

replace (
	Crawler-go/models => ./models
	Crawler-go/setting => ./setting
)
