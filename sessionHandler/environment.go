package sessionHandler

import (
	"io/ioutil"
	"os"
	"strings"
)

// Kopiere eine Datei an eine andere Stelle.
func copyFile(src string, dst string) {
	data, err := ioutil.ReadFile(src)
	HandleError(err)

	err = ioutil.WriteFile(dst, data, 0644)
	HandleError(err)
}

// Erstelle ein Verzeichnis für Tickets, sofern dieses nicht existiert.
func CheckEnvironment() {
	assetsDir := GetAssetsDir()

	if _, err := os.Stat(assetsDir + "tickets"); os.IsNotExist(err) {
		err = os.Mkdir(assetsDir+"tickets", 0755)
		HandleError(err)
		if _, err := os.Stat(assetsDir + "ticketId_resource.json"); os.IsNotExist(err) {
			srcFile := strings.Join([]string{assetsDir, "rollback/default/ticketId_resource.json"}, "")
			dstFile := strings.Join([]string{assetsDir, "ticketId_resource.json"}, "")
			copyFile(srcFile, dstFile)
		}
	}

	if _, err := os.Stat(assetsDir + "users.json"); os.IsNotExist(err) {
		srcFile := strings.Join([]string{assetsDir, "rollback/default/users.json"}, "")
		dstFile := strings.Join([]string{assetsDir, "users.json"}, "")
		copyFile(srcFile, dstFile)
	}
}

// Setze den Webserver zurück, indem alle Tickets und Nutzerdaten gelöscht werden.
// Dabei wird zunächst jeweils geprüft, ob das Objekt existiert.
func ResetData() {
	assetsDir := GetAssetsDir()

	if _, err := os.Stat(assetsDir + "tickets"); os.IsNotExist(err) == false {
		err := os.RemoveAll(assetsDir + "tickets")
		HandleError(err)
	}

	if _, err := os.Stat(assetsDir + "users.json"); os.IsNotExist(err) == false {
		err = os.Remove(assetsDir + "users.json")
		HandleError(err)
	}

	if _, err := os.Stat(assetsDir + "ticketId_resource.json"); os.IsNotExist(err) == false {
		err = os.Remove(assetsDir + "ticketId_resource.json")
		HandleError(err)
	}
}

// Setze den Webserver zurück und installiere Testdaten.
func DemoMode() {
	assetsDir := GetAssetsDir()
	ResetData()

	// Erstelle Verzeichnis für Tickets, falls dieses nicht existiert.
	if _, err := os.Stat(assetsDir + "tickets"); os.IsNotExist(err) {
		err = os.Mkdir(assetsDir+"tickets", 0755)
		HandleError(err)
	}

	// Kopiere alle Tickets aus dem Demo-Ordner in den Zielordner.
	src := strings.Join([]string{assetsDir, "rollback/demo/tickets"}, "")
	files, err := ioutil.ReadDir(src)
	HandleError(err)

	for _, file := range files {
		srcFile := strings.Join([]string{assetsDir, "rollback/demo/tickets/", file.Name()}, "")
		dstFile := strings.Join([]string{assetsDir, "tickets/", file.Name()}, "")
		copyFile(srcFile, dstFile)
	}

	// Kopiere ID-Ressource für Tickets aus dem Demo-Ordner in den Zielordner.
	srcFile := strings.Join([]string{assetsDir, "rollback/demo/ticketId_resource.json"}, "")
	dstFile := strings.Join([]string{assetsDir, "ticketId_resource.json"}, "")
	copyFile(srcFile, dstFile)

	// Kopiere Nutzerdaten aus dem Demo-Ordner in den Zielordner.
	srcFile = strings.Join([]string{assetsDir, "rollback/demo/users.json"}, "")
	dstFile = strings.Join([]string{assetsDir, "users.json"}, "")
	copyFile(srcFile, dstFile)
}
