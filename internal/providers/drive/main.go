package drive

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"passman/internal/interfaces"
	"passman/internal/providers"
)

type googleDrive struct {
	srv *drive.Service
}

func (g *googleDrive) GetFile(filename string) ([]byte, error) {
	if g.srv == nil {
		g.connect()
	}

	listRes, err := g.srv.Files.List().Do(googleapi.QueryParameter("q", fmt.Sprintf("name='%s'", filename)))
	if err != nil {
		return nil, err
	}

	if len(listRes.Files) == 0 {
		return nil, providers.ErrFileNotFound
	}

	getRes, err := g.srv.Files.Export(listRes.Files[0].Id, "text/plain").Download()
	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(getRes.Body)
	if err != nil {
		return nil, err
	}

	bytes = bytes[3:]

	if err = os.WriteFile(filename, bytes, 0666); err != nil {
		return nil, err
	}

	return bytes, nil
}

func (g *googleDrive) SaveFile(filename string, data []byte) {
	if g.srv == nil {
		g.connect()
	}

	err := os.WriteFile(filename, data, 0666)
	if err != nil {
		log.Fatal("error while writing file: ", err)
	}

	open, err := os.Open(filename)
	if err != nil {
		log.Fatal("error while open file: ", err)
	}

	listRes, err := g.srv.Files.List().Do(googleapi.QueryParameter("q", fmt.Sprintf("name='%s'", filename)))
	if err != nil {
		log.Fatal("error while getting file id: ", err)
	}

	if len(listRes.Files) == 0 {
		_, err = g.srv.Files.Create(&drive.File{Name: filename}).Media(open).Do()
		if err != nil {
			log.Fatal("error while updating file: ", err)
		}

		return
	}

	_, err = g.srv.Files.Update(listRes.Files[0].Id, &drive.File{}).Media(open).Do()
	if err != nil {
		log.Fatal("error while updating file: ", err)
	}
}

func (g *googleDrive) connect() {
	ctx := context.Background()
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveFileScope, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	g.srv = srv
}

func NewGoogleDriveProvider() interfaces.DataProvider {
	return &googleDrive{}
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	codeCh := make(chan string, 1)

	go runServer(codeCh)

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	authCode := <-codeCh

	fmt.Println()

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}

	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func runServer(codeCh chan string) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Any("/api", func(c *gin.Context) {
		c.String(http.StatusOK, "You can now go back to the app.")
		c.Abort()

		codeCh <- c.Query("code")
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
