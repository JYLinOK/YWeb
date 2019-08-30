package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// _____________________________________________________________________________________
// _____________________________________________________________________________________
// Get all dirs and files

// Define the global variables
// Deifne the slices to save the data.
var slicDirs = []string{}
var staticFsDirs = []string{}
var slicFiles = []string{}

// Define the WebROOT
var yWebROOT string = ""
var yWorkROOT string = ""

func getAllDirAndFiles(pathName string) ([]string, []string) {

	// Use the ioutil.ReadDir() to get the detail data of the specified dir.
	dirData, err := ioutil.ReadDir(pathName)

	if err != nil {
		fmt.Println("ioutil.ReadDir Errror: ", err)
	} else {
		// Traversal the dir
		for _, fi := range dirData {
			// Define the the full-dir-name of the variable
			savePath := pathName
			pathName := pathName + "/" + fi.Name()

			// ________________________________________________
			fmt.Println("Top pathName =", pathName)
			fmt.Println("Top fi.Name() =", fi.Name())
			// ________________________________________________

			// Judget the content of the value of the variable is a path or not.
			if fi.IsDir() {

				// Record the dir data with slice.
				slicDirs = append(slicDirs, pathName)
				// Recall the traversal function
				getAllDirAndFiles(pathName)

			} else {
				// yWebROOT
				// Get the pure files name and extension name from the fi.Name()
				// sp[0] = pure-file-name, sp[1] = pure-file-extension

				sp := strings.Split(fi.Name(), ".")

				fmt.Println("sp[0] = ", sp[0])
				fmt.Println("sp[1] = ", sp[1])

				// pure_file_name := sp[0]
				pureFileExtension := sp[1]

				// Add handlers to files
				pageFilesExtension := [...]string{"html"}
				for _, v := range pageFilesExtension {

					fmt.Println("v = ", v)

					// Handle the html type files
					if v == (pureFileExtension) {
						// yWebROOT = E:/GoWeb/YWeb/YDOOK
						// yWorkROOT = E:/GoWeb/YWeb

						fmt.Println("The file type is OK!")

						spRoot := strings.Split(yWebROOT, yWorkROOT+"/")

						fmt.Println("spRoot[1] = ", spRoot[1])

						// fi_Name := "H1.html"
						// yWorkROOT = yWorkROOT + "/"
						fmt.Println("yWorkROOT = ", yWorkROOT)
						spFile := strings.Split(pathName, yWorkROOT+"/")
						thisFilePath := spFile[1]

						fmt.Println("spFile[1] = ", spFile[1])
						fmt.Println("thisFilePath = ", thisFilePath)

						// spRoot[1] =  YDOOK
						dotN := strings.LastIndex(thisFilePath, ".")
						strDir := thisFilePath[:dotN]
						thisDirPath := strings.Split(strDir, spRoot[1])

						// ("/", "YDOOK/index.html")

						fmt.Println("thisDirPath[1] = ", thisDirPath[1])
						fmt.Println("thisFilePath = ", thisFilePath)

						if thisDirPath[1] == "/index" {
							thisDirPath[1] = "/"
							fmt.Println("/index => ", thisDirPath[1])
						}

						// AH("/", "YDOOK/index.html")
						go AH(string(thisDirPath[1]), string(thisFilePath))

					} else {
						// Handle other kinds of files
						// AHFS("E:/GoWeb/YWeb/YDOOK/PhotoHost")
						// savePath =  E:/GoWeb/YWeb/YDOOK/PhotoHost
						go AHFS(savePath)
						fmt.Println("savePath = ", savePath)

					}
				}

				// _______________________________________________________________________
				// _______________________________________________________________________
				// Handle the html type files

				fmt.Println("_____________________________________________________________________")

				slicFiles = append(slicFiles, pathName)
			}
		}
	}

	return slicDirs, slicFiles

}

// H0 struct
// Add Handlers
type H0 struct {
	dir string
}

func (h H0) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles(h.dir)
	t.Execute(w, "HelloWord!!! Host Page!")
}

// AH Add handlers to the files
// AH("/", "YDOOK/index.html")
func AH(dirHandler string, dirFile string) {
	handlerStruct := H0{
		dir: dirFile,
	}

	http.Handle(dirHandler, handlerStruct)

	// AH("/", "YDOOK/index.html")
	// AH("/H1", "YDOOK/H1.html")
	// AH("/H2", "YDOOK/H2.html")
	// fmt.Printf("http.Handle(%v, %v) T1:%T, T2:%T \n", dirHandler, handlerStruct, dirHandler, handlerStruct)
	// fmt.Printf("dirHandler = %v, dirFile = %v, T1:%T, T2:%T \n", dirHandler, dirFile, dirHandler, dirFile)

}

// AHFS Add handles to the FileServers
func AHFS(dir string) {
	fmt.Println("Enter the AHFS")
	flag := false

	for _, v := range staticFsDirs {
		fmt.Println("Enter the for loop")
		if v == dir {
			flag = true
			fmt.Println("flag = ", flag)
		}
	}

	if flag == false {
		// Add the new item to the staticFsDirs slice
		staticFsDirs = append(staticFsDirs, dir)

		fs := http.FileServer(http.Dir(dir))
		sp := strings.Split(dir, "/")
		sDir := "/" + sp[len(sp)-1] + "/"
		http.Handle(sDir, http.StripPrefix(sDir, fs))

		fmt.Println("AHFS() = ", dir)
	}

}

// _____________________________________________________________________________________
// _____________________________________________________________________________________
// Define the get json data function

func getConfig(dir string) map[string]string {

	webConfig := map[string]string{
		// "HostIp": "127.0.0.1",
		// "EnPort": "8000",
	}

	// Define the struct for the config of Yweb
	type Post struct {
		IP     string `json:"host_ip"`
		PORT   string `json:"en_port"`
		ROOT   string `json:"web_root"`
		YWROOT string `json:"ywork_root"`
	}

	// Open and read the json file
	JSONF, err := os.Open(dir)
	if err != nil {
		fmt.Println("os.Open() ERROR: ", err)
	} else {
		fmt.Println("Openning the JSON file successfully!")
	}

	JSONFData, err := ioutil.ReadAll(JSONF)
	if err != nil {
		fmt.Println("ioutil.ReadAll() ERROR ：", err)
	} else {
		fmt.Println("Reading the JSON file successfully!")
	}

	var readPost Post
	json.Unmarshal(JSONFData, &readPost)

	webConfig["HostIp"] = readPost.IP
	webConfig["EnPort"] = readPost.PORT
	webConfig["WebROOT"] = readPost.ROOT
	webConfig["YWorkROOT"] = readPost.YWROOT

	fmt.Println("HostIp = ", webConfig["HostIp"])
	fmt.Println("EnPort = ", webConfig["EnPort"])
	fmt.Println("WebROOT = ", webConfig["WebROOT"])
	fmt.Println("YWorkROOT = ", webConfig["YWorkROOT"])

	defer JSONF.Close()

	return webConfig

}

// _____________________________________________________________________________________
// _____________________________________________________________________________________
func main() {

	// slicDirs = nil
	// slicFiles = nil

	fmt.Println("slicDirs = ", slicDirs)
	fmt.Println("slicFiles = ", slicFiles)

	// Get the detail configure
	confConfig := "Config/config.json"
	confConfigData := getConfig(confConfig)

	hostIP := confConfigData["HostIp"]
	enPort := confConfigData["EnPort"]
	yWebROOT = confConfigData["WebROOT"]
	yWorkROOT = confConfigData["YWorkROOT"]

	fmt.Println("hostIP = ", hostIP)
	fmt.Println("enPort = ", enPort)
	fmt.Println("yWorkROOT = ", yWorkROOT)
	fmt.Println("yWebROOT = ", yWebROOT)

	// _____________________________________________________________________________________________________
	// _____________________________________________________________________________________________________
	fmt.Println("______________________________________________________")
	pathName := yWebROOT
	a, b := getAllDirAndFiles(pathName)

	for i, c := range a {
		fmt.Printf("Dirs Items NO. %d is %v \n", i, c)
	}
	for i, c := range b {
		fmt.Printf("Files Items NO. %d is %v \n", i, c)
	}

	// _____________________________________________________________________________________________________
	// _____________________________________________________________________________________________________

	// AHFS("E:/GoWeb/YWeb/YDOOK/PhotoHost")
	// AHFS("E:/GoWeb/YWeb/YDOOK/CSS")

	// ____________________________________________________
	// ____________________________________________________
	// ____________________________________________________
	// ____________________________________________________

	// Set the erver and run it
	server := http.Server{
		Addr: hostIP + ":" + enPort,
	}
	server.ListenAndServe()

}
