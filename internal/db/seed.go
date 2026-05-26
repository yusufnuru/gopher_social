package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	"github.com/yusufnuru/gopher_social/internal/store"
)

var usernames = []string{
	"shadowfox01", "pixelnova22", "silenttiger", "bluecomet77", "neonfalcon", "rapidwolf99", "stormbyte", "darkorbitx", "lunarvibe", "crypticzen", "ironpanda", "frostflare", "echohunter", "quantumjay", "wildmatrix", "solarblitz", "velvetstorm", "ghostlynx", "hypernovae", "ninjacobra", "turboraven", "voidrunner", "alphadrift", "cosmicpixel", "silverfang", "nightpulse",
	"emberstrike", "glitchwizard", "thunderroot", "vortexspark", "atomiccrane", "stealthotter", "mechashade", "frozenorbit", "skyphantom", "mysticquark", "nebulashift", "phantomhawk", "cyberbadger", "stellarforge", "blazekitten", "gravitylock", "electricmoth", "midnightbyte", "rocketlizard", "crimsonwave", "diamondvolt", "ultracipher", "zenithrider", "omegaecho",
}

var titles = []string{
	"The Rise of Quantum Computing", "Why Sleep Is Your Superpower", "Go vs Rust in 2025", "Urban Farming Is Taking Off", "The Loneliness Epidemic", "Dark Mode Is Overrated", "Learning a Language After 30", "The Death of the Open Office", "Why Film Photography Survived", "The Hidden Cost of Fast Fashion", "AI Is Writing Your News", "Cold Plunges Aren't Magic", "How Rust Is Entering the Linux Kernel", "The Quiet Luxury Trend", "Remote Work Changed Cities Forever", "The Comeback of Board Games", "Electric Bikes Are Replacing Cars", "Microplastics Are Everywhere Now", "The 4-Day Work Week Experiment", "Why Handwriting Still Matters",
}

var contents = []string{
	"Quantum computers are no longer just theoretical. Engineers are now solving problems that classical machines can't touch.",
	"Skipping sleep to grind harder is counterproductive. Your brain consolidates memory and repairs itself only during deep sleep.",
	"Both languages dominate systems programming, but Go wins on simplicity while Rust wins on memory safety guarantees.",
	"Rooftop gardens and vertical farms are reshaping how cities think about food supply and sustainability.",
	"Despite being more connected than ever online, people report feeling more isolated. Community spaces are making a comeback.",
	"Studies show dark mode doesn't reduce eye strain as much as we think. Ambient lighting matters far more.",
	"Adults actually have advantages over kids when learning languages — stronger grammar intuition and better study habits.",
	"Startups are ditching open floor plans after research showed they kill deep work and increase distractions.",
	"Digital was supposed to kill film. Instead, a whole new generation is shooting on 35mm for the texture and intentionality.",
	"A single pair of jeans requires 1,800 gallons of water to produce. Slow fashion is no longer a niche movement.",
	"More outlets are using AI to draft breaking news. The debate isn't if it's happening — it's how to label it.",
	"The viral cold plunge trend has benefits, but most claims are exaggerated. Consistency in sleep and exercise matters more.",
	"For the first time in decades, Linux is welcoming a second language. Rust modules are shipping in production kernels.",
	"Logos are out. Understated, high-quality basics are in. Fashion is shifting away from loud branding.",
	"Mid-size cities saw population booms as remote workers left expensive metros. That shift isn't reversing.",
	"Tabletop gaming is a billion-dollar industry again. People are craving screen-free social experiences.",
	"In dense cities, e-bikes are outselling electric cars. They're cheaper, faster in traffic, and easier to park.",
	"Researchers have found microplastics in human blood, lungs, and even placentas. The long-term effects are still unknown.",
	"Companies that tried a 4-day week reported higher productivity and lower turnover. Most refused to go back.",
	"Taking notes by hand improves retention compared to typing. The friction forces your brain to process and summarize.",
}

var tags = []string{
	"quantum", "health", "go", "sustainability", "society", "design", "language", "productivity", "sustainability", "journalism", "health", "linux", "minimalism", "remote-work", "tabletop",
	"ebike", "transport", "environment", "work", "productivity", "writing",
}

var comments = []string{
	"Quantum supremacy is still overhyped.",
	"8 hours of sleep changed my life.",
	"Rust is worth the learning curve.",
	"My city just approved a rooftop farm bill!",
	"Third places really need to come back.",
	"Dark mode looks cooler though.",
	"Started Spanish at 32, can confirm.",
	"My focus went up 10x working from home.",
	"The grain of film is just unmatched.",
	"Thrifting is the move.",
	"As long as AI news is labeled I am fine with it.",
	"The science on cold plunges is pretty weak.",
	"About time Linux modernized.",
	"Finally fashion that actually lasts.",
	"Moved from SF to Boise, never looking back.",
	"Nothing beats Catan with friends.",
	"Sold my car, bought an e-bike, no regrets.",
	"Plastic was a mistake.",
	"4 days should be the standard by now.",
	"Switched back to pen and paper for notes.",
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)

	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			log.Println("Error creating user:", err)
			return
		}
	}

	tx.Commit()

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post:", err)
			return
		}
	}

	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment:", err)
			return
		}
	}

	log.Println("Seeding complete")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := range num {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
		}
	}

	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	for i := range num {
		posts[i] = &store.Post{
			UserID:  users[rand.Intn(len(users))].ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}
	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	cmts := make([]*store.Comment, num)
	for i := range num {
		cmts[i] = &store.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}
	return cmts
}
