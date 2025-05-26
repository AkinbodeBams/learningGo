package db

import (
	"context"
	"database/sql"
	"math/rand"

	"fmt"
	"log"

	"github.com/akinbodeBams/social/internal/store"
)
var usernames = []string{
    "alex_brown", "jessica_lee", "michael_smith", "sophia_jones", "daniel_kim",
    "emily_clark", "ryan_taylor", "olivia_white", "james_moore", "mia_davis",
    "ethan_hall", "ava_martin", "noah_allen", "isabella_thomas", "liam_scott",
    "grace_harris", "lucas_miller", "chloe_adams", "jack_turner", "ella_wright",
    "benjamin_cooper", "amelia_evans", "henry_wilson", "scarlett_morris", "sebastian_hughes",
    "lily_bennett", "logan_ross", "zoe_morgan", "jayden_hill", "nora_price",
    "leo_green", "hannah_watson", "nathan_wood", "ariana_reed", "caleb_bailey",
    "eva_cook", "christian_foster", "leah_kelly", "dylan_henderson", "audrey_simpson",
    "isaac_barnes", "samantha_patterson", "owen_long", "layla_wells", "julian_richards",
    "ruby_perry", "aaron_ward", "hazel_carter", "miles_powell", "piper_hamilton",
}

var titles = []string{
    "10 Habits of Highly Productive People",
    "A Beginner's Guide to Mindful Living",
    "Why Learning to Code Is Easier Than You Think",
    "The Power of Small Daily Wins",
    "Mastering Time Management in a Digital World",
    "Top 5 Books That Changed My Perspective on Life",
    "From Burnout to Balance: My Work-Life Reset",
    "How to Start a Side Hustle That Actually Works",
    "Exploring the Future of Artificial Intelligence",
    "Minimalism: Decluttering Your Life, Not Just Your Room",
    "Building Mental Resilience in Tough Times",
    "What I Learned After Failing My First Startup",
    "The Art of Saying No Without Guilt",
    "How to Create a Morning Routine That Sticks",
    "Remote Work: Tips for Staying Focused and Motivated",
    "The Science Behind Better Sleep",
    "5 Ways to Boost Your Creativity Today",
    "Why Journaling Could Be the Habit That Changes Your Life",
    "Redefining Success in the Age of Hustle",
    "Simple Nutrition Tips That Actually Work",
}

var contents = []string{
    "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
    "Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
    "Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris.",
    "Duis aute irure dolor in reprehenderit in voluptate velit esse.",
    "Excepteur sint occaecat cupidatat non proident, sunt in culpa.",
    "Curabitur pretium tincidunt lacus. Nulla gravida orci a odio.",
    "Donec sit amet eros. Lorem ipsum dolor sit amet, consectetur.",
    "Etiam vel tortor sodales tellus ultricies commodo.",
    "Aliquam erat volutpat. In congue. Etiam justo.",
    "Aliquam lorem ante, dapibus in, viverra quis, feugiat a, tellus.",
    "Vestibulum ante ipsum primis in faucibus orci luctus.",
    "Aenean leo ligula, porttitor eu, consequat vitae, eleifend ac, enim.",
    "Integer tincidunt. Cras dapibus. Vivamus elementum semper nisi.",
    "Aenean vulputate eleifend tellus. Aenean leo ligula.",
    "Phasellus viverra nulla ut metus varius laoreet.",
    "Quisque rutrum. Aenean imperdiet. Etiam ultricies nisi vel augue.",
    "Curabitur ullamcorper ultricies nisi. Nam eget dui.",
    "Maecenas tempus, tellus eget condimentum rhoncus.",
    "Fusce vulputate eleifend sapien. Vestibulum purus quam.",
    "Morbi in sem quis dui placerat ornare. Pellentesque odio nisi.",
}

var tags = []string{
    "tech","health","travel",
    "finance",
    "coding",
    "music",
    "books",
    "food",
    "lifestyle",
    "fitness",
    "education",
    "science",
    "gaming",
    "history",
    "design",
    "art",
    "startup",
    "nature",
    "culture",
    "news",
}

var comments = []string{
    "Great post! Thanks for sharing.",
    "I completely agree with your points.",
    "This was really helpful, appreciate it.",
    "Interesting perspective!",
    "Could you explain more about this topic?",
    "Loved this article.",
    "This cleared up a lot of confusion.",
    "Thanks for the insights!",
    "I learned something new today.",
    "Not sure I agree with this.",
    "This is spot on!",
    "Nice work, keep it up!",
    "You make a strong argument.",
    "This was a fun read.",
    "I’ll be sharing this with friends.",
    "Helpful for beginners, thanks!",
    "This could use more details.",
    "Looking forward to more like this.",
    "Awesome content!",
    "Can you do a follow-up post on this?",
}



func Seed(store store.Storage, db *sql.DB) {
ctx:= context.Background()
users := generateUsers(100)
tx,_ := db.BeginTx(ctx,nil)
for _, user := range users {
	if err:= store.Users.Create(ctx,tx,user); err!=nil {
		log.Println("Error creating user", err)
		return 
	}
}
tx.Commit()
posts:= generatePost(200,users)
for _, post := range posts {
	if err:= store.Posts.Create(ctx,post); err!=nil {
		log.Println("Error creating post", err)
		return 
	}
}

comments := generateComments(500, users, posts)
for _, cms := range comments {
	if err:= store.Comments.Create(ctx,cms); err!=nil {
		log.Println("Error creating comment", err)
		return 
	}
}

log.Println("seeding completed")
return 
}

func generateUsers(num int) []*store.User{ 
users:= make([]*store.User,num)
pwd := store.Password{}
_ = pwd.Set("test123") 
for i := 0; i< num; i++ {
	username:= usernames[i%len(usernames)] + fmt.Sprintf("%d",i)
	users[i] = &store.User{
        
		Username:username,
		Email: username+ "@example.com",
		 Password:pwd,
       
	}
}
return users
}

func generatePost(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]
		posts[i] = &store.Post{
			UserId:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Tags:    []string{tags[rand.Intn(len(tags))]},
		}
	}
	return posts
}


func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	cms := make([]*store.Comment, num)
	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]
		post := posts[rand.Intn(len(posts))]
		comment := comments[rand.Intn(len(comments))]
		cms[i] = &store.Comment{
			UserId:  user.ID,
			PostId:  post.ID,
			Content: comment,
		}
	}
	return cms
}
