package main

import (
	"database/sql"
	"fmt"
	"html"
	"time"

	// _ "github.com/go-sql-driver/mysql"

	_ "github.com/lib/pq"
	"github.com/mmcdole/gofeed"
	"github.com/nlopes/slack"
)

type Feed struct {
	Url     string
	Channel string
}

// rss_url = [
//     {url: 'https://www.upwork.com/ab/feed/jobs/rss?q=vue&sort=recency&job_type=hourly%2Cfixed&contractor_tier=1%2C2%2C3&proposals=0-4%2C5-9%2C10-14%2C15-19%2C20-49&budget=100-499%2C500-999%2C1000-4999%2C5000-&workload=as_needed%2Cpart_time%2Cfull_time&duration_v3=week%2Cmonth%2Csemester%2Congoing&verified_payment_only=1&connect_price=0-2%2C4%2C6&paging=0%3B50&api_params=1&securityToken=47afe9a9c98905215eff0ebb88e56c6569dcdc4d488154ffd5919b79fe7dea7de636ce5f4f17916ec086e90afb58c177a39d0220953e320ac563245149f0f2a3&userUid=700193349553074176&orgUid=700193349557268481', channel: '#upwork_feed_vue'},
//     {url: 'https://www.upwork.com/ab/feed/jobs/rss?q=node&sort=recency&job_type=hourly%2Cfixed&contractor_tier=1%2C2%2C3&proposals=0-4%2C5-9%2C10-14%2C15-19%2C20-49&budget=100-499%2C500-999%2C1000-4999%2C5000-&workload=as_needed%2Cpart_time%2Cfull_time&duration_v3=week%2Cmonth%2Csemester%2Congoing&verified_payment_only=1&connect_price=0-2%2C4%2C6&paging=0%3B50&api_params=1&securityToken=47afe9a9c98905215eff0ebb88e56c6569dcdc4d488154ffd5919b79fe7dea7de636ce5f4f17916ec086e90afb58c177a39d0220953e320ac563245149f0f2a3&userUid=700193349553074176&orgUid=700193349557268481', channel: '#upwork_feed_node'},
//     {url: 'https://www.upwork.com/ab/feed/jobs/rss?q=python&sort=recency&job_type=hourly%2Cfixed&contractor_tier=1%2C2%2C3&proposals=0-4%2C5-9%2C10-14%2C15-19%2C20-49&budget=100-499%2C500-999%2C1000-4999%2C5000-&workload=as_needed%2Cpart_time%2Cfull_time&duration_v3=week%2Cmonth%2Csemester%2Congoing&verified_payment_only=1&connect_price=0-2%2C4%2C6&paging=0%3B50&api_params=1&securityToken=47afe9a9c98905215eff0ebb88e56c6569dcdc4d488154ffd5919b79fe7dea7de636ce5f4f17916ec086e90afb58c177a39d0220953e320ac563245149f0f2a3&userUid=700193349553074176&orgUid=700193349557268481', channel: '#upwork_feed_python'},
//     {url: 'https://www.upwork.com/ab/feed/jobs/rss?q=django&sort=recency&job_type=hourly%2Cfixed&contractor_tier=1%2C2%2C3&proposals=0-4%2C5-9%2C10-14%2C15-19%2C20-49&budget=100-499%2C500-999%2C1000-4999%2C5000-&workload=as_needed%2Cpart_time%2Cfull_time&duration_v3=week%2Cmonth%2Csemester%2Congoing&verified_payment_only=1&connect_price=0-2%2C4%2C6&paging=0%3B50&api_params=1&securityToken=47afe9a9c98905215eff0ebb88e56c6569dcdc4d488154ffd5919b79fe7dea7de636ce5f4f17916ec086e90afb58c177a39d0220953e320ac563245149f0f2a3&userUid=700193349553074176&orgUid=700193349557268481', channel: '#upwork_feed_python'},
//     {url: 'https://www.upwork.com/ab/feed/jobs/rss?q=quasar&sort=recency&job_type=hourly%2Cfixed&contractor_tier=1%2C2%2C3&proposals=0-4%2C5-9%2C10-14%2C15-19%2C20-49&budget=100-499%2C500-999%2C1000-4999%2C5000-&workload=as_needed%2Cpart_time%2Cfull_time&duration_v3=week%2Cmonth%2Csemester%2Congoing&verified_payment_only=1&connect_price=0-2%2C4%2C6&paging=0%3B50&api_params=1&securityToken=47afe9a9c98905215eff0ebb88e56c6569dcdc4d488154ffd5919b79fe7dea7de636ce5f4f17916ec086e90afb58c177a39d0220953e320ac563245149f0f2a3&userUid=700193349553074176&orgUid=700193349557268481', channel: '#upwork_feed_vue'},
//     {url: 'https://www.upwork.com/ab/feed/jobs/rss?q=react&sort=recency&job_type=hourly%2Cfixed&contractor_tier=1%2C2%2C3&proposals=0-4%2C5-9%2C10-14%2C15-19%2C20-49&budget=100-499%2C500-999%2C1000-4999%2C5000-&workload=as_needed%2Cpart_time%2Cfull_time&duration_v3=week%2Cmonth%2Csemester%2Congoing&verified_payment_only=1&connect_price=0-2%2C4%2C6&paging=0%3B50&api_params=1&securityToken=47afe9a9c98905215eff0ebb88e56c6569dcdc4d488154ffd5919b79fe7dea7de636ce5f4f17916ec086e90afb58c177a39d0220953e320ac563245149f0f2a3&userUid=700193349553074176&orgUid=700193349557268481', channel: '#upwork_feed_react'},
//     {url: 'https://www.upwork.com/ab/feed/jobs/rss?q=magento&sort=recency&job_type=hourly%2Cfixed&contractor_tier=1%2C2%2C3&proposals=0-4%2C5-9%2C10-14%2C15-19%2C20-49&budget=100-499%2C500-999%2C1000-4999%2C5000-&workload=as_needed%2Cpart_time%2Cfull_time&duration_v3=week%2Cmonth%2Csemester%2Congoing&verified_payment_only=1&connect_price=0-2%2C4%2C6&paging=0%3B50&api_params=1&securityToken=47afe9a9c98905215eff0ebb88e56c6569dcdc4d488154ffd5919b79fe7dea7de636ce5f4f17916ec086e90afb58c177a39d0220953e320ac563245149f0f2a3&userUid=700193349553074176&orgUid=700193349557268481', channel: '#upwork_feed_magento'},
//     {url: 'https://www.upwork.com/ab/feed/jobs/rss?q=laravel&sort=recency&job_type=hourly%2Cfixed&contractor_tier=1%2C2%2C3&proposals=0-4%2C5-9%2C10-14%2C15-19%2C20-49&budget=100-499%2C500-999%2C1000-4999%2C5000-&workload=as_needed%2Cpart_time%2Cfull_time&duration_v3=week%2Cmonth%2Csemester%2Congoing&verified_payment_only=1&connect_price=0-2%2C4%2C6&paging=0%3B50&api_params=1&securityToken=47afe9a9c98905215eff0ebb88e56c6569dcdc4d488154ffd5919b79fe7dea7de636ce5f4f17916ec086e90afb58c177a39d0220953e320ac563245149f0f2a3&userUid=700193349553074176&orgUid=700193349557268481', channel: '#upwork_feed_laravel'},
//     {url: 'https://www.upwork.com/ab/feed/jobs/rss?q=php&sort=recency&job_type=hourly%2Cfixed&contractor_tier=1%2C2%2C3&proposals=0-4%2C5-9%2C10-14%2C15-19%2C20-49&budget=100-499%2C500-999%2C1000-4999%2C5000-&workload=as_needed%2Cpart_time%2Cfull_time&duration_v3=week%2Cmonth%2Csemester%2Congoing&verified_payment_only=1&connect_price=0-2%2C4%2C6&paging=0%3B50&api_params=1&securityToken=47afe9a9c98905215eff0ebb88e56c6569dcdc4d488154ffd5919b79fe7dea7de636ce5f4f17916ec086e90afb58c177a39d0220953e320ac563245149f0f2a3&userUid=700193349553074176&orgUid=700193349557268481', channel: '#upwork_feed_php'},
//     {url: 'https://www.upwork.com/ab/feed/jobs/rss?q=shopify&sort=recency&job_type=hourly%2Cfixed&contractor_tier=1%2C2%2C3&proposals=0-4%2C5-9%2C10-14%2C15-19%2C20-49&budget=100-499%2C500-999%2C1000-4999%2C5000-&workload=as_needed%2Cpart_time%2Cfull_time&duration_v3=week%2Cmonth%2Csemester%2Congoing&verified_payment_only=1&connect_price=0-2%2C4%2C6&paging=0%3B50&api_params=1&securityToken=47afe9a9c98905215eff0ebb88e56c6569dcdc4d488154ffd5919b79fe7dea7de636ce5f4f17916ec086e90afb58c177a39d0220953e320ac563245149f0f2a3&userUid=700193349553074176&orgUid=700193349557268481', channel: '#upwork_feed_shopify'},
//     {url: 'https://www.upwork.com/ab/feed/jobs/rss?q=golang&sort=recency&job_type=hourly%2Cfixed&contractor_tier=1%2C2%2C3&proposals=0-4%2C5-9%2C10-14%2C15-19%2C20-49&budget=100-499%2C500-999%2C1000-4999%2C5000-&workload=as_needed%2Cpart_time%2Cfull_time&duration_v3=week%2Cmonth%2Csemester%2Congoing&verified_payment_only=1&connect_price=0-2%2C4%2C6&paging=0%3B50&api_params=1&securityToken=47afe9a9c98905215eff0ebb88e56c6569dcdc4d488154ffd5919b79fe7dea7de636ce5f4f17916ec086e90afb58c177a39d0220953e320ac563245149f0f2a3&userUid=700193349553074176&orgUid=700193349557268481', channel: '#upwork_feed_go},
//   ]

var (
	feeds = []Feed{
		{"https://www.upwork.com/ab/feed/jobs/rss?q=node&sort=recency&job_type=hourly%2Cfixed&contractor_tier=1%2C2%2C3&proposals=0-4%2C5-9%2C10-14%2C15-19%2C20-49&budget=100-499%2C500-999%2C1000-4999%2C5000-&workload=as_needed%2Cpart_time%2Cfull_time&duration_v3=week%2Cmonth%2Csemester%2Congoing&verified_payment_only=1&connect_price=0-2%2C4%2C6&paging=0%3B50&api_params=1&securityToken=47afe9a9c98905215eff0ebb88e56c6569dcdc4d488154ffd5919b79fe7dea7de636ce5f4f17916ec086e90afb58c177a39d0220953e320ac563245149f0f2a3&userUid=700193349553074176&orgUid=700193349557268481", "#upwork_feed_node"},
		{"https://www.upwork.com/ab/feed/jobs/rss?q=golang&sort=recency&job_type=hourly%2Cfixed&contractor_tier=1%2C2%2C3&proposals=0-4%2C5-9%2C10-14%2C15-19%2C20-49&budget=100-499%2C500-999%2C1000-4999%2C5000-&workload=as_needed%2Cpart_time%2Cfull_time&duration_v3=week%2Cmonth%2Csemester%2Congoing&verified_payment_only=1&connect_price=0-2%2C4%2C6&paging=0%3B50&api_params=1&securityToken=47afe9a9c98905215eff0ebb88e56c6569dcdc4d488154ffd5919b79fe7dea7de636ce5f4f17916ec086e90afb58c177a39d0220953e320ac563245149f0f2a3&userUid=700193349553074176&orgUid=700193349557268481", "#upwork_feed_golang"},
		{"https://www.upwork.com/ab/feed/jobs/rss?q=react&sort=recency&job_type=hourly%2Cfixed&contractor_tier=1%2C2%2C3&proposals=0-4%2C5-9%2C10-14%2C15-19%2C20-49&budget=100-499%2C500-999%2C1000-4999%2C5000-&workload=as_needed%2Cpart_time%2Cfull_time&duration_v3=week%2Cmonth%2Csemester%2Congoing&verified_payment_only=1&connect_price=0-2%2C4%2C6&paging=0%3B50&api_params=1&securityToken=47afe9a9c98905215eff0ebb88e56c6569dcdc4d488154ffd5919b79fe7dea7de636ce5f4f17916ec086e90afb58c177a39d0220953e320ac563245149f0f2a3&userUid=700193349553074176&orgUid=700193349557268481", "#upwork_feed_react"},
		{"https://www.upwork.com/ab/feed/jobs/rss?q=node&sort=recency&job_type=hourly%2Cfixed&contractor_tier=1%2C2%2C3&proposals=0-4%2C5-9%2C10-14%2C15-19%2C20-49&budget=100-499%2C500-999%2C1000-4999%2C5000-&workload=as_needed%2Cpart_time%2Cfull_time&duration_v3=week%2Cmonth%2Csemester%2Congoing&verified_payment_only=1&connect_price=0-2%2C4%2C6&paging=0%3B50&api_params=1&securityToken=47afe9a9c98905215eff0ebb88e56c6569dcdc4d488154ffd5919b79fe7dea7de636ce5f4f17916ec086e90afb58c177a39d0220953e320ac563245149f0f2a3&userUid=700193349553074176&orgUid=700193349557268481", "#upwork_feed_node"},
		{"https://www.upwork.com/ab/feed/jobs/rss?q=python&sort=recency&job_type=hourly%2Cfixed&contractor_tier=1%2C2%2C3&proposals=0-4%2C5-9%2C10-14%2C15-19%2C20-49&budget=100-499%2C500-999%2C1000-4999%2C5000-&workload=as_needed%2Cpart_time%2Cfull_time&duration_v3=week%2Cmonth%2Csemester%2Congoing&verified_payment_only=1&connect_price=0-2%2C4%2C6&paging=0%3B50&api_params=1&securityToken=47afe9a9c98905215eff0ebb88e56c6569dcdc4d488154ffd5919b79fe7dea7de636ce5f4f17916ec086e90afb58c177a39d0220953e320ac563245149f0f2a3&userUid=700193349553074176&orgUid=700193349557268481", "#upwork_feed_python"},
		{"https://www.upwork.com/ab/feed/jobs/rss?q=django&sort=recency&job_type=hourly%2Cfixed&contractor_tier=1%2C2%2C3&proposals=0-4%2C5-9%2C10-14%2C15-19%2C20-49&budget=100-499%2C500-999%2C1000-4999%2C5000-&workload=as_needed%2Cpart_time%2Cfull_time&duration_v3=week%2Cmonth%2Csemester%2Congoing&verified_payment_only=1&connect_price=0-2%2C4%2C6&paging=0%3B50&api_params=1&securityToken=47afe9a9c98905215eff0ebb88e56c6569dcdc4d488154ffd5919b79fe7dea7de636ce5f4f17916ec086e90afb58c177a39d0220953e320ac563245149f0f2a3&userUid=700193349553074176&orgUid=700193349557268481", "#upwork_feed_python"},
		{"https://www.upwork.com/ab/feed/jobs/rss?q=ror&sort=recency&job_type=hourly%2Cfixed&contractor_tier=1%2C2%2C3&proposals=0-4%2C5-9%2C10-14%2C15-19%2C20-49&budget=100-499%2C500-999%2C1000-4999%2C5000-&workload=as_needed%2Cpart_time%2Cfull_time&duration_v3=week%2Cmonth%2Csemester%2Congoing&verified_payment_only=1&connect_price=0-2%2C4%2C6&paging=0%3B50&api_params=1&securityToken=47afe9a9c98905215eff0ebb88e56c6569dcdc4d488154ffd5919b79fe7dea7de636ce5f4f17916ec086e90afb58c177a39d0220953e320ac563245149f0f2a3&userUid=700193349553074176&orgUid=700193349557268481", "#upwork_feed_ror"},
		{"https://www.upwork.com/ab/feed/jobs/rss?q=java&sort=recency&job_type=hourly%2Cfixed&contractor_tier=1%2C2%2C3&proposals=0-4%2C5-9%2C10-14%2C15-19%2C20-49&budget=100-499%2C500-999%2C1000-4999%2C5000-&workload=as_needed%2Cpart_time%2Cfull_time&duration_v3=week%2Cmonth%2Csemester%2Congoing&verified_payment_only=1&connect_price=0-2%2C4%2C6&paging=0%3B50&api_params=1&securityToken=47afe9a9c98905215eff0ebb88e56c6569dcdc4d488154ffd5919b79fe7dea7de636ce5f4f17916ec086e90afb58c177a39d0220953e320ac563245149f0f2a3&userUid=700193349553074176&orgUid=700193349557268481", "#upwork_feed_java"},

		// {"https://www.upwork.com/ab/feed/jobs/rss?q=Nodejs&sort=recency&t=0,1&contractor_tier=1,2,3&client_hires=0&proposals=0-4,5-9,10-14&amount=500-50000&payment_verified=1&hourly_rate=8-50&api_params=1&securityToken=47afe9a9c98905215eff0ebb88e56c6569dcdc4d488154ffd5919b79fe7dea7de636ce5f4f17916ec086e90afb58c177a39d0220953e320ac563245149f0f2a3&userUid=700193349553074176&orgUid=700193349557268481", "#upwork_feed_node"},
		// {"https://www.upwork.com/ab/feed/jobs/rss?q=golang&sort=recency&t=0,1&contractor_tier=1,2,3&client_hires=0&proposals=0-4,5-9,10-14&amount=500-50000&payment_verified=1&hourly_rate=8-50&api_params=1&securityToken=47afe9a9c98905215eff0ebb88e56c6569dcdc4d488154ffd5919b79fe7dea7de636ce5f4f17916ec086e90afb58c177a39d0220953e320ac563245149f0f2a3&userUid=700193349553074176&orgUid=700193349557268481", "#upwork_feed_golang"},
		// {"https://www.upwork.com/ab/feed/jobs/rss?q=reactjs&sort=recency&t=0,1&contractor_tier=1,2,3&client_hires=0&proposals=0-4,5-9,10-14&amount=500-50000&payment_verified=1&hourly_rate=8-50&api_params=1&securityToken=47afe9a9c98905215eff0ebb88e56c6569dcdc4d488154ffd5919b79fe7dea7de636ce5f4f17916ec086e90afb58c177a39d0220953e320ac563245149f0f2a3&userUid=700193349553074176&orgUid=700193349557268481", "#upwork_feed_react"},
		// {"https://www.upwork.com/ab/feed/jobs/rss?q=Express&sort=recency&t=0,1&contractor_tier=1,2,3&client_hires=0&proposals=0-4,5-9,10-14&amount=500-50000&payment_verified=1&hourly_rate=8-50&api_params=1&securityToken=47afe9a9c98905215eff0ebb88e56c6569dcdc4d488154ffd5919b79fe7dea7de636ce5f4f17916ec086e90afb58c177a39d0220953e320ac563245149f0f2a3&userUid=700193349553074176&orgUid=700193349557268481", "#upwork_feed_node"},
		// {"https://www.upwork.com/ab/feed/jobs/rss?q=python&sort=recency&t=0,1&contractor_tier=1,2,3&client_hires=0&proposals=0-4,5-9,10-14&amount=500-50000&payment_verified=1&hourly_rate=8-50&api_params=1&securityToken=47afe9a9c98905215eff0ebb88e56c6569dcdc4d488154ffd5919b79fe7dea7de636ce5f4f17916ec086e90afb58c177a39d0220953e320ac563245149f0f2a3&userUid=700193349553074176&orgUid=700193349557268481", "#upwork_feed_python"},
		// {"https://www.upwork.com/ab/feed/jobs/rss?q=django&sort=recency&t=0,1&contractor_tier=1,2,3&client_hires=0&proposals=0-4,5-9,10-14&amount=500-50000&payment_verified=1&hourly_rate=8-50&api_params=1&securityToken=47afe9a9c98905215eff0ebb88e56c6569dcdc4d488154ffd5919b79fe7dea7de636ce5f4f17916ec086e90afb58c177a39d0220953e320ac563245149f0f2a3&userUid=700193349553074176&orgUid=700193349557268481", "#upwork_feed_python"},
	}

	db *sql.DB
	sc *slack.Client
)

type FeedDetails struct {
	Title string
	Link  string
}

func init() {
	// create a MySQL database instance-
	// print("KUGIUGIUGIGIUGIUGUGIUGIUJ")

	var err error
	// db, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/rssfeeds")
	db, err = sql.Open("postgres", "postgres://superdbuser:Password1234@localhost/rssfeeds")
	if err != nil {
		print("KUGIUGIUGIGIUGIUGUGIUGIUJ", err)
		panic(err)
	}

	// defer db.Close()
	// db.SetMaxOpenConns(250) // Adjust the value as per your requirements

	// db.SetMaxIdleConns(250) // Adjust the value as per your requirements

	// db.SetConnMaxLifetime(time.Minute * 5) // Adjust the value as per your requirements

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS feed_details (title VARCHAR(255), link VARCHAR(255), PRIMARY KEY(title))")
	if err != nil {
		print("xxbcvhbx")
		panic(err)
	}
	// println("JYGYUGYGYYH")

	// Replace with your Slack bot token
	// slackBotToken := "xoxb-5024190216176-4993950875446-iwQ7hlX1MoE2A6kwXpa263DM"
	slackBotToken := "xoxb-4994422805927-5042532385296-4HDreDBsmgQTpOVL5GjObnGG"

	sc = slack.New(slackBotToken)
	// print("\nJHGFYG========")
}

func SendSlackNotification(channel, title, link, description string) error {
	attachment := slack.Attachment{
		Fallback:   title,
		Color:      "#36a64f",
		Title:      title,
		Text:       html.UnescapeString(description),
		TitleLink:  link,
		MarkdownIn: []string{"text"},
	}
	message := slack.MsgOptionAttachments(attachment)
	text := slack.MsgOptionText("New RSS Feed Item :mag_right:", false)
	_, _, err := sc.PostMessage(channel, text, message)
	// print("--------098765354909ijhbnn0i-0m-0u9nyhun-mum-mu------------")
	return err
}

// func ProcessRssFeed() {
// 	for _, feed := range feeds {
// 		fp := gofeed.NewParser()
// 		feedItems, _ := fp.ParseURL(feed.Url)

// 		for _, item := range feedItems.Items {
// 			title := item.Title
// 			link := item.Link

// 			// check if the feed item has already been processed
// 			rows, err := db.Query("SELECT title FROM feed_details WHERE title = $1", title)
// 			if err != nil {
// 				// fmt.Println(err, "IUB*&T&*87v8787tv6786786v87")
// 				continue
// 			}
// 			// time.Sleep(10 * time.Second)
// 			defer rows.Close()

// 			if rows.Next() {
// 				// feed item already exists in database, skip it
// 				// print("JJJJJJJJJJ\n")
// 				continue
// 			}
// 			// print("KKKKKKKKKKK")

// 			fmt.Println(feed.Channel)

// 			description := item.Description
// 			channel := feed.Channel

// 			err = SendSlackNotification(channel, title, link, description)
// 			if err != nil {
// 				fmt.Println(err, "IUGUIGh67tg67t67t67t67tg67gt")
// 				continue
// 			}

// 			// insert the feed item into the database
// 			_, err = db.Exec("INSERT INTO feed_details (title, link) VALUES (?, ?)", title, link)
// 			if err != nil {
// 				// fmt.Println(err, "KGJUHJGHHGHGHG987897907")
// 				continue
// 			}
// 		}
// 	}
// }

func ProcessRssFeed() {
	print("JJJJJJJJJJJJJJJJJJJJ\n")
	for _, feed := range feeds {
		fp := gofeed.NewParser()
		feedItems, err := fp.ParseURL(feed.Url)
		if err != nil {
			fmt.Println("Error parsing RSS feed:", err)
			continue
		}

		for _, item := range feedItems.Items {
			title := item.Title
			link := item.Link

			// Check if the feed item already exists in the database
			var existingTitle string
			err := db.QueryRow("SELECT title FROM feed_details WHERE title = $1", title).Scan(&existingTitle)
			switch {
			case err == sql.ErrNoRows:
				// Feed item does not exist, proceed with processing
			case err != nil:
				// Error occurred while checking for duplicate, log and continue
				fmt.Println("Error checking for duplicate:", err)
				continue
			default:
				// Feed item already exists in the database, skip it
				// fmt.Printf("Duplicate feed item found: %s\n", title)
				continue
			}
			// If the control reaches here, it means the feed item is not a duplicate
			description := item.Description
			channel := feed.Channel

			err = SendSlackNotification(channel, title, link, description)
			if err != nil {
				fmt.Println("Error sending Slack notification:", err)
				continue
			}

			// Insert the feed item into the database
			_, err = db.Exec("INSERT INTO feed_details (title, link) VALUES ($1, $2)", title, link)
			if err != nil {
				fmt.Println("Error inserting into database:", err)
				continue
			}
		}
		time.Sleep(1 * time.Second)
	}
}

// func main() {
// 	router := mux.NewRouter()

// 	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		// print("cxvkjcxbkjv")
// 		fmt.Fprintf(w, "Welcome to my RSS feed app!")
// 	})

// 	// run the RSS feed processing in the background
// 	go func() {
// 		for {
// 			ProcessRssFeed()
// 			time.Sleep(120 * time.Second)
// 		}
// 	}()

// 	http.ListenAndServe(":5050", router)
// }

func main() {
	// Initialize your database and Slack client
	// init()

	// Run ProcessRssFeed() periodically using a ticker
	ticker := time.Tick(60 * time.Second) // Adjust the interval here

	for range ticker {
		ProcessRssFeed()
	}
}
