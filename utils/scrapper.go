package utils 

import(
	"net/http"
	"fmt"
	"io"
	"bytes"
	
)

var baseUrl = "https://maplestorywiki.net"

var characters = "https://maplestorywiki.net/w/Characters_and_Skills"



// retrieves all class icons from each maplestory class
func GetAllIcons(){
	
}

// retrives the list of all of the characters 
func GetCharacters(){
	resp,err := http.Get(characters)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var output bytes.Buffer 
	io.Copy(&output, resp.Body)
		
	var characterSection bytes.Buffer
	
	data := output.Bytes()
	inCharacterSection := false
	for i,b := range data {
		if inCharacterSection == true {
			if b =='D'{
				if bytes.HasSuffix(data[i:i+12], []byte("Discontinued")){
					inCharacterSection = false 
					break
				}else {
					characterSection.WriteByte(b)
				}
			} else {
				characterSection.WriteByte(b)			
			}
		} else {
			if b == 'b' {
				if bytes.HasSuffix(data[i:i+5], []byte("below"))  == true {
					inCharacterSection = true
				}
			}
		}
	}
	
	characterUrls := characterSection.Bytes()
	inCharacterSection = false
	var urlExtension bytes.Buffer

	urls := make([]string, 10)
	i := 0
	for i < len(characterUrls){
		if inCharacterSection == true {
			if characterUrls[i]== '"'{
				inCharacterSection = false
				url := urlExtension.String()
				urls = append(urls, url)
				urlExtension.Reset()
			}else {
				urlExtension.WriteByte(characterUrls[i])
			}
			i += 1
		} else {
			if characterUrls[i]== 'h' {
				if i < len(characterUrls){
					if bytes.HasSuffix(characterUrls[i:i+6], []byte("href=\"")){
						inCharacterSection = true
						i += 6
					} else {
						i += 1
					}
				} else {
					break
				}
			} else {
				i += 1
			}
		}
	}
	
	for _,url := range urls{
		if url == ""{
			continue
		}else {
			GetIcons(url)
		}
	}


}

//retrieves all class icons from the specified maplestory class
func GetIcons(class string){
	url := baseUrl + class + "/Skills"
	fmt.Println("url is: ", url)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var output bytes.Buffer

	io.Copy(&output, resp.Body)

	data := output.Bytes()

	inIconSection := false
	var iconSections bytes.Buffer
	for i,b := range data {
		if inIconSection == true{
			if b == 'd'{
				if bytes.HasSuffix(data[i:i+8], []byte("decoding")){
					inIconSection = false
				} else {
					iconSections.WriteByte(b)
				}
			} else {
				iconSections.WriteByte(b)
			}
		} else {
			if b == 'A'{
				if bytes.HasSuffix(data[i:i+6], []byte("Active")){
					inIconSection = true
				}
			}
		}
	}

	iconUrls := make([]string, 20)
	iconData := iconSections.Bytes()	
	var iconUrl bytes.Buffer
	inUrlSection := false	
	i := 0 
	for i < len(iconData){
		if inUrlSection{
			if iconData[i] == '"'{
				url := iconUrl.String()
				iconUrls = append(iconUrls, url)
				iconUrl.Reset()
				inUrlSection = false
				i += 1
			} else {
				iconUrl.WriteByte(iconData[i])
				i += 1
			}
		} else{
			if iconData[i] == 's' {
				if bytes.HasSuffix(iconData[i:i+5], []byte("src=\"")){
					i += 5
					inUrlSection = true
				} else {
					i += 1
				}
			} else {
				i += 1
			}
		}
	}

	for _, url := range iconUrls{
		if url == "" {
			continue
		} else {
			// do i request to each image url and save response to a png file in the corresponding class file
			fmt.Println("image url:", url)
			saveIconImage(url, class[3:])
		}
	}
}

func saveIconImage(url, class string){

	
}
