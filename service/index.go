package service

import(
	"github.com/labstack/echo/v4"
	"net/http"
	"echolabstack/model"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron"
	 "database/sql"
	 "time"
      "fmt"
	  "log"

)

func NewAPI(c echo.Context) error {
	return c.String(http.StatusOK, "New API created")
}

func PdfAPI(c echo.Context) error{
	pdfPath := "static/Echo_static.pdf"
	return c.File(pdfPath)

}

var db *sql.DB
var err error
var lastBatchTime time.Time

func EmailIDAPI(c echo.Context) error {
	// Parse request body
	var email model.Email
	if err := c.Bind(&email); err != nil {
		return err
	}

	// Insert email into database
	_, err:= db.Exec("INSERT INTO subscribers (email,posted) VALUES (?, ?)", email.Email, email.Posted)

	if err != nil {
		fmt.Println("Error inserting into database:", err)
		return c.String(http.StatusConflict, "Email already subscribed")
	}

	StartCron()
	return c.String(http.StatusCreated, "Subscribed successfully")
}
func StartCron() {
	c := cron.New()
	_ = c.AddFunc("1 * * ? * *", post)
	c.Start()
}
// @cron(run every 1 min)
func post() {


	// Create a new batch
	var newBatchID int64
	
	newBatchID, err = insertBatch()
	if err != nil {
		log.Println("Error creating batch:", err)
		return
	}
	log.Printf("New batch created: %d", newBatchID)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Println(err)
	}
	_, err = db.Exec("UPDATE subscribers SET batch_id = ? WHERE posted = false", newBatchID)
	stmt, err := db.Prepare("SELECT email FROM subscribers WHERE posted = false AND batch_id = ?")
	if err != nil {
		log.Println("error*****")
		panic(err)
	}
	mail, err := stmt.Query(newBatchID)
	if err != nil {
		log.Println("error")
		panic(err)
	}
	defer mail.Close()
	var emailIds []string

	// Iterate through the mail
	for mail.Next() {
		var email string
		if err := mail.Scan(&email); err != nil {
			log.Fatal(err)
		}
		emailIds = append(emailIds, email)
	}

	if len(emailIds) > 0 {
		sendMail(emailIds)
	}

	_, err = db.Exec("UPDATE subscribers SET posted = true WHERE posted = false AND batch_id = ?", newBatchID)

	if err != nil {
		log.Fatal(err)
	}
}

func sendMail(whoSubscribed []string) {
	//send email(admin@netxd.com,whoSubscribed)
	log.Println(whoSubscribed)
}

func insertBatch() (int64, error) {
	result, err := db.Exec("INSERT INTO batch DEFAULT VALUES")
	if err != nil {
		return -1, err
	}
	batchID, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return batchID, nil
}
