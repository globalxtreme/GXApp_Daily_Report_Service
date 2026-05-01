package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	xtrememodel "github.com/globalxtreme/go-core/v2/model"
	xtremerabbitmq "github.com/globalxtreme/go-core/v2/rabbitmq"
	"github.com/joho/godotenv"
	"io"
	"mime/multipart"
	"service/internal/pkg/config"
	"service/internal/pkg/constant"
	"service/internal/pkg/model"
	"strings"
	"time"
)

type testing struct {
	Test *xtrememodel.ArrayMapInterfaceColumn
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}

	Close := config.InitDB()
	defer Close()

	var session model.ReportSession
	config.PgSQL.First(&session, "id = ?", 5)

	date := time.Now()
	fmt.Println("time.Now():", date)
	fmt.Println("timezone:", date.Location())

	session.CompletedAt = &date
	config.PgSQL.Save(&session)
}

func getCase1() interface{} {
	return map[string]interface{}{
		"name": "Dedi",
	}
}

func sendRabbitMQ() {
	RabbitMQClose := config.InitRabbitMQ()
	defer RabbitMQClose()

	dialRabbitMQConnClose := config.InitRabbitMQConnection()
	defer dialRabbitMQConnClose()

	senderId := "1"
	senderType := "messages"
	push := xtremerabbitmq.RabbitMQ{
		Connection: xtremerabbitmq.RABBITMQ_CONNECTION_GLOBAL,
		Queue:      constant.RABBITMQ_QUEUE_SERVICE_DOMAIN_FEATURE_ACTION,
		SenderId:   &senderId,
		SenderType: &senderType,
		Data: map[string]interface{}{
			"name": "Testing",
			"subs": []string{"test"},
		},
	}
	push.OnDelivery("services", true)
	push.Push()
}

func padKey(key string, length int) []byte {
	if len(key) >= length {
		return []byte(key[:length])
	}
	return append([]byte(key), bytes.Repeat([]byte("0"), length-len(key))...)
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	if length == 0 || length%blockSize != 0 {
		return nil, fmt.Errorf("invalid padded data")
	}
	pad := int(data[length-1])
	if pad == 0 || pad > blockSize {
		return nil, fmt.Errorf("invalid padding")
	}
	return data[:length-pad], nil
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	pad := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, pad...)
}

func DecryptLaravelHash(key string, id float64, encrypted string) (string, error) {
	if strings.HasPrefix(encrypted, "_open:") {
		encrypted = strings.TrimPrefix(encrypted, "_open:")
	}

	raw, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("base64 decode error: %w", err)
	}
	if len(raw) < aes.BlockSize {
		return "", fmt.Errorf("data too short")
	}

	iv := raw[:aes.BlockSize]
	ciphertextBase64 := raw[aes.BlockSize:]

	ciphertext, err := base64.StdEncoding.DecodeString(string(ciphertextBase64))
	if err != nil {
		return "", fmt.Errorf("decode nested base64 error: %w", err)
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return "", fmt.Errorf("invalid ciphertext block size")
	}

	// 4. Generate AES key
	rawKey := padKey(fmt.Sprintf("%s%f", key, id), 16)
	block, err := aes.NewCipher(rawKey)
	if err != nil {
		return "", fmt.Errorf("aes cipher error: %w", err)
	}

	// 5. Decrypt
	decrypted := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decrypted, ciphertext)

	// 6. Unpad
	result, err := pkcs7Unpad(decrypted, aes.BlockSize)
	if err != nil {
		return "", fmt.Errorf("unpad error: %w", err)
	}

	return string(result), nil
}

func EncryptLaravelStyle(key string, id float64, plaintext []byte) (string, error) {
	// AES Key: key + id
	rawKey := padKey(fmt.Sprintf("%s%f", key, id), 16)

	block, err := aes.NewCipher(rawKey)
	if err != nil {
		return "", fmt.Errorf("cipher error: %w", err)
	}

	// Generate IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("IV error: %w", err)
	}

	// Pad plaintext
	plainBytes := pkcs7Pad(plaintext, aes.BlockSize)

	// Encrypt
	ciphertext := make([]byte, len(plainBytes))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plainBytes)

	// Inner base64 (like PHP's openssl_encrypt default)
	ciphertextBase64 := base64.StdEncoding.EncodeToString(ciphertext)

	// Outer base64: IV (raw) + inner ciphertext (string)
	full := append(iv, []byte(ciphertextBase64)...)
	finalEncoded := base64.StdEncoding.EncodeToString(full)

	return "_open:" + finalEncoded, nil
}

func SendAPI() {
	//// Open file
	//file, err := os.Open("storage/app/public/tmp/CeBxK2EWtKIeqtOy5bsw1698052181068410000.png")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//defer file.Close()
	//
	//fileStat, err := file.Stat()
	//if err != nil {
	//	fmt.Println("Error getting file stats:", err)
	//	return
	//}
	//
	//// Membuat FileHeader manual (ganti sesuai kebutuhan)
	//fileHeader := &multipart.FileHeader{
	//	Filename: "CeBxK2EWtKIeqtOy5bsw1698052181068410000.png",
	//	Size:     fileStat.Size(),
	//	//// ContentType, biasanya kamu bisa tentukan berdasarkan ekstensi file atau MIME type
	//	//Header: map[string][]string{
	//	//	"Content-Type": {mimeType},
	//	//},
	//}

	//data := Payload{
	//	User: User{
	//		Name:  "Yuswa GPT",
	//		Email: "yuswa@example.com",
	//	},
	//	Meta: Meta{
	//		Location: Location{
	//			Lat:  10.123,
	//			Long: 106.456,
	//		},
	//	},
	//	//Document: UploadFile{
	//	//	File:        file,
	//	//	FileHandler: fileHeader,
	//	//},
	//}

	//api := xtremeapi.NewXtremeAPI()
	//
	//res := api.Post("http://127.0.0.1:8000/testing", Payload{})
	//fmt.Println(res)
}

type Payload struct {
	User     User       `json:"user"`
	Meta     Meta       `json:"meta"`
	Document UploadFile `json:"document"`
}

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Meta struct {
	Location Location `json:"location"`
}

type Location struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type UploadFile struct {
	File        multipart.File
	FileHandler *multipart.FileHeader
}
