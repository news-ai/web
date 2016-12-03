package encrypt_test

import (
	"fmt"

	"github.com/news-ai/web/encrypt"
)

func ExampleEncryptString() {
	encryptMessage, _ := encrypt.EncryptString("hello")
	decryptMessage, _ := encrypt.DecryptString(encryptMessage)
	fmt.Println(decryptMessage)
	// Output: hello
}
