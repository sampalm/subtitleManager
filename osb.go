package main

import (
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
var c = &Controller{}

func encodeOSB(e Encoder) string {
	es := e.Encode()
	rep := strings.NewReplacer("=", "-", "&", "/")
	return rep.Replace(es)
}

func (c *Controller) osbRequest(method string, params url.Values) error {
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

func DownloadHashed(root, hash string, size int64) {
	var params = make(url.Values)
	c.RootFolder = root
	c.DefaultLanguage = "eng"
	c.Filter = false
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

func HashFile(file *os.File) (hash string, size int64, err error) {
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
