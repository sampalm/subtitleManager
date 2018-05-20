package osb

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Subtitle struct {
	Rating       string `json:"SubRating"`
	FileName     string `json:"SubFileName"`
	LanguageID   string `json:"SubLanguageID"`
	LanguageName string
	DownloadLink string `json:"SubDownloadLink"`
	InfoFormat   string
	Format       string `json:"SubFormat"`
	Body         *bytes.Buffer
}

const ChunkSize = 65536

var client = &http.Client{}

func SearchHashSub(hash, size string) (subs []Subtitle, err error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://rest.opensubtitles.org/search/moviebytesize-%s/moviehash-%s", size, hash), nil)
	if err != nil {
		err = fmt.Errorf("searchSub: NewRequest: %s", err.Error())
	}
	req.Header.Set("User-Agent", "TemporaryUserAgent")
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("searchSub: ClientDo: %s", err.Error())
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&subs)
	if err != nil {
		err = fmt.Errorf("searchSub: Decoding: %s", err.Error())
	}
	return
}

func DownloadSub(sub *Subtitle) error {
	var b = &bytes.Buffer{}
	res, err := http.Get(sub.DownloadLink)
	if err != nil {
		return fmt.Errorf("downloadSub: downloadSubtitle: %s", err.Error())
	}
	defer res.Body.Close()
	rd, err := gzip.NewReader(res.Body)
	if err != nil {
		return fmt.Errorf("downloadSub: gzipExtract: %s", err.Error())
	}
	defer rd.Close()
	if _, err = io.Copy(b, rd); err != nil {
		return fmt.Errorf("downloadSub: copySubtitle: %s", err.Error())
	}
	sub.Body = b
	return nil
}

func HashFile(file *os.File) (hash uint64, size int64, err error) {
	defer file.Close()
	fi, err := file.Stat()
	log.Println("FILE: ", fi.Name())
	if err != nil {
		return
	}
	if fi.Size() < ChunkSize {
		return 0, 0, fmt.Errorf("hashFile: Size: file is too small")
	}

	// Read head and tail blocks.
	buf := make([]byte, ChunkSize*2)
	err = readChunk(file, 0, buf[:ChunkSize])
	if err != nil {
		return
	}
	err = readChunk(file, fi.Size()-ChunkSize, buf[ChunkSize:])
	if err != nil {
		return
	}

	// Convert to uint64, and sum.
	var nums [(ChunkSize * 2) / 8]uint64
	reader := bytes.NewReader(buf)
	err = binary.Read(reader, binary.LittleEndian, &nums)
	if err != nil {
		return 0, 0, fmt.Errorf("hashFile: BinaryRead: %v", err)
	}
	for _, num := range nums {
		hash += num
	}
	return hash + uint64(fi.Size()), fi.Size(), nil
}

// Read a chunk of a file at `offset` so as to fill `buf`.
func readChunk(file *os.File, offset int64, buf []byte) (err error) {
	n, err := file.ReadAt(buf, offset)
	if err != nil {
		return
	}
	if n != ChunkSize {
		return fmt.Errorf("Invalid read %v", n)
	}
	return
}
