package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

type Controller struct {
	RootFolder      string
	DefaultLanguage string
	MultiLanguage   []string
	RatingScore     int
	Filter          bool
	Subtitles       []Subtitle
	*http.Client
}

type Subtitle struct {
	Rating       string `json:"SubRating"`
	FileName     string `json:"SubFileName"`
	LanguageID   string `json:"SubLanguageID"`
	LanguageName string
	DownloadLink string `json:"SubDownloadLink"`
}

type Encoder interface {
	Encode() string
}

const chunkSize = 65536

var client = &http.Client{}
var uri = "https://rest.opensubtitles.org/search/"

func encodeOSB(e Encoder) string {
	es := e.Encode()
	rep := strings.NewReplacer("=", "-", "&", "/")
	return rep.Replace(es)
}

func (c *Controller) osbRequest(method string, params url.Values) error {
	if c.MultiLanguage == nil {
		params.Add("sublanguageid", c.DefaultLanguage)
	}
	uri += encodeOSB(params)

	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "TemporaryUserAgent")
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		panic(err) // FATAL ERROR - TRY TO RECOVER ?
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&c.Subtitles)
	if err != nil {
		return err
	}
	fmt.Println("**** REQUEST HEAD => ", uri)
	return nil
}

func (c *Controller) download(link, filename string) error {
	res, err := http.Get(link)
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			panic(err)
		}
	}()
	b, err := gzip.NewReader(res.Body)
	if err != nil {
		if err == gzip.ErrHeader {
			return fmt.Errorf("gzipExtract: couldn't download the file, try again later")
		}
		return err
	}
	defer func() {
		if err := b.Close(); err != nil {
			panic(err)
		}
	}()
	file, err := os.OpenFile(filepath.Join(c.RootFolder, filename), syscall.O_CREAT|syscall.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("%s: %s", filename, err.Error())
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	if _, err := io.Copy(file, b); err != nil {
		return fmt.Errorf("Could't download your file. Err: %s", err.Error())
	}
	return nil
}

func DownloadHashed(c *Controller, hash string, size int64) {
	var params = make(url.Values)
	params.Add("moviebytesize", fmt.Sprint(size))
	params.Add("moviehash", hash)
	err := c.osbRequest(http.MethodGet, params)
	if err != nil {
		panic(err) // raising a flag
	}
	for i := range c.Subtitles {
		err := c.download(c.Subtitles[i].DownloadLink, c.Subtitles[i].FileName)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("File %s downloaded successfully.\n", c.Subtitles[i].FileName)
	}
}

func DownloadQuery(c *Controller, params url.Values) {
	err := c.osbRequest(http.MethodGet, params)
	if err != nil {
		panic(err) // raising a flag
	}
	for i := range c.Subtitles {
		fmt.Printf("[%d] Subtitle: %s\nRating: %s - Language: %s/%s\n",
			i,
			c.Subtitles[i].FileName,
			c.Subtitles[i].Rating,
			c.Subtitles[i].LanguageID,
			c.Subtitles[i].LanguageName,
		)
	}
	var ss string
	var reader = bufio.NewReader(in)
	fmt.Println("\n-------------------------------\nSELECT SUBTITLES TO DOWNLOAD BY NUMBER:\n(Enter 'd' to end task or 'c' to cancel)\n-------------------------------")
	for {
		fmt.Print("-> ")
		ss, _ = reader.ReadString('\n')
		ss = strings.TrimSpace(strings.Replace(ss, "\n", "", -1))
		if ss == "d" {
			break
		}
		if ss == "c" {
			fmt.Println("Task canceled.")
			return
		}

		n, err := strconv.Atoi(ss)
		if err != nil {
			fmt.Println("Invalid number")
			continue
		}
		if n > len(c.Subtitles) {
			fmt.Println("Invalid subtitle")
			continue
		}
		go func(s Subtitle) {
			fmt.Printf("Starting to download %s...\n-> ", s.FileName)
			if gerr := c.download(s.DownloadLink, s.FileName); gerr != nil {
				fmt.Printf("Error while downloading: %s\n-> ", gerr.Error())
				return
			}
		}(c.Subtitles[n])
	}
	fmt.Println("Everything worked out... cya ;)")
}

func GetHashFiles(c *Controller, path string, p PullFiles) {
	var fl []File
	exts = map[string]bool{".mp4": true, ".mkv": true, ".avi": true, ".wmv": true}
	if err := p(path, "", &fl); err != nil {
		panic(err)
	}
	for _, f := range fl {
		file, err := os.Open(f.Name)
		if err != nil {
			fmt.Printf("Could not open file %s: %s\n", f.Name, err.Error())
			continue
		}
		hash, size, err := hashFile(file)
		if err != nil {
			fmt.Printf("Could not hash file %s: %s\n", f.Name, err.Error())
			continue
		}
		DownloadHashed(c, hash, size)
	}
}

func hashFile(file *os.File) (hash string, size int64, err error) {
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		return
	}
	if fi.Size() < chunkSize {
		return "", 0, fmt.Errorf("hashFile: Size: file is too small")
	}

	// Read head and tail blocks.
	buf := make([]byte, chunkSize*2)
	err = readChunk(file, 0, buf[:chunkSize])
	if err != nil {
		return
	}
	err = readChunk(file, fi.Size()-chunkSize, buf[chunkSize:])
	if err != nil {
		return
	}

	// Convert to uint64, and sum.
	var h uint64
	var nums [(chunkSize * 2) / 8]uint64
	reader := bytes.NewReader(buf)
	err = binary.Read(reader, binary.LittleEndian, &nums)
	if err != nil {
		return "", 0, fmt.Errorf("hashFile: BinaryRead: %v", err)
	}
	for _, num := range nums {
		h += num
	}
	return fmt.Sprintf("%x", h+uint64(fi.Size())), fi.Size(), nil
}

// Read a chunk of a file at `offset` so as to fill `buf`.
func readChunk(file *os.File, offset int64, buf []byte) (err error) {
	n, err := file.ReadAt(buf, offset)
	if err != nil {
		return
	}
	if n != chunkSize {
		return fmt.Errorf("Invalid read %v", n)
	}
	return
}
