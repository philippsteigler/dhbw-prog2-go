package sessionHandler

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Kopiere eine Datei an eine andere Stelle.
func copyFile(src string, dst string) {
	data, err := ioutil.ReadFile(src)
	HandleError(err)

	err = ioutil.WriteFile(dst, data, 0744)
	HandleError(err)
}

// Erstelle ein Verzeichnis für Tickets und Mails, sofern dieses nicht existiert.
func CheckEnvironment() {
	assetsDir := GetAssetsDir()

	if _, err := os.Stat(assetsDir + "tickets"); os.IsNotExist(err) {
		err = os.Mkdir(assetsDir+"tickets", 0744)
		HandleError(err)
	}

	if _, err := os.Stat(assetsDir + "mails"); os.IsNotExist(err) {
		err = os.Mkdir(assetsDir+"mails", 0744)
		HandleError(err)
	}

	if _, err := os.Stat(assetsDir + "users.json"); os.IsNotExist(err) {
		srcFile := strings.Join([]string{assetsDir, "rollback/default/users.json"}, "")
		dstFile := strings.Join([]string{assetsDir, "users.json"}, "")
		copyFile(srcFile, dstFile)
	}
}

//typeOfFiles entweder ticket oder mails
func checkIfFilesExistAndBackup(typeOfFiles string) {
	assetsDir := GetAssetsDir()
	rollbackBackupPath := "rollback/backup/" + typeOfFiles

	// Überprüfe, ob Tickets/Mails existieren.
	if _, err := os.Stat(assetsDir + typeOfFiles); os.IsNotExist(err) == false {
		src := strings.Join([]string{assetsDir, typeOfFiles}, "")
		files, err := ioutil.ReadDir(src)
		HandleError(err)

		// Wenn Tickets/Mails vorhanden sind, werden alle gesichert.
		if len(files) > 0 {
			err = os.Mkdir(assetsDir+rollbackBackupPath, 0744)
			HandleError(err)

			for _, file := range files {
				srcFile := strings.Join([]string{assetsDir, typeOfFiles + "/", file.Name()}, "")
				dstFile := strings.Join([]string{assetsDir, rollbackBackupPath + "/", file.Name()}, "")
				copyFile(srcFile, dstFile)
			}
		}
	}
}

//typeOfFiles entweder ticket oder mails
func checkIfBackupExistAndLoadFiles(typeOfFiles string) {
	assetsDir := GetAssetsDir()
	rollbackBackupPath := "rollback/backup/" + typeOfFiles
	// Überprüfe, ob das Backup Tickets/Mails enthält.
	if _, err := os.Stat(assetsDir + rollbackBackupPath); os.IsNotExist(err) == false {
		src := strings.Join([]string{assetsDir, rollbackBackupPath}, "")
		files, err := ioutil.ReadDir(src)
		HandleError(err)

		// Wenn Tickets/Mails vorhanden sind, werden alle geladen.
		if len(files) > 0 {

			err = os.Mkdir(assetsDir+typeOfFiles, 0700)
			fmt.Println(err)

			for _, file := range files {
				srcFile := strings.Join([]string{assetsDir, rollbackBackupPath + "/", file.Name()}, "")
				dstFile := strings.Join([]string{assetsDir, typeOfFiles + "/", file.Name()}, "")
				copyFile(srcFile, dstFile)
			}
		}
	}
}

// Hilfsfunktion für Unit-Tests.
// Sichere alle produktiven Daten des nach ./assets/rollback/backup.
func BackupEnvironment() {
	assetsDir := GetAssetsDir()

	// Lösche alte Backups, sofern diese vorhanden sind.
	if _, err := os.Stat(assetsDir + "rollback/backup"); os.IsNotExist(err) {
		err = os.Mkdir(assetsDir+"rollback/backup", 0744)
		HandleError(err)
	} else {
		err := os.RemoveAll(assetsDir + "rollback/backup")
		HandleError(err)
		err = os.Mkdir(assetsDir+"rollback/backup", 0744)
		HandleError(err)
	}

	// Überprüfe, ob Tickets existieren.
	checkIfFilesExistAndBackup("tickets")

	// Überprüfe, ob Mails existieren.
	checkIfFilesExistAndBackup("mails")

	// Sichere die Nutzerdaten, sofern diese existieren.
	if _, err := os.Stat(assetsDir + "users.json"); os.IsNotExist(err) == false {
		srcFile := strings.Join([]string{assetsDir, "users.json"}, "")
		dstFile := strings.Join([]string{assetsDir, "rollback/backup/users.json"}, "")
		copyFile(srcFile, dstFile)
	}
}

// Hilfsfunktion für Unit-Tests.
// Lösche alle produktiven Daten und lade ein Backup.
func RestoreEnvironment() {
	assetsDir := GetAssetsDir()

	// Produktive Daten werden nur durch ein Backup ersetzt, wenn eins vorhanden ist.
	if _, err := os.Stat(assetsDir + "rollback/backup"); os.IsNotExist(err) == false {
		ResetData()

		// Überprüfe, ob das Backup Tickets enthält.
		checkIfBackupExistAndLoadFiles("tickets")

		//Überprüfe, ob das Backup Mails enthält.
		checkIfBackupExistAndLoadFiles("mails")

		// Lade die Nutzerdaten, sofern diese existieren.
		if _, err := os.Stat(assetsDir + "rollback/backup/users.json"); os.IsNotExist(err) == false {
			srcFile := strings.Join([]string{assetsDir, "rollback/backup/users.json"}, "")
			dstFile := strings.Join([]string{assetsDir, "users.json"}, "")
			copyFile(srcFile, dstFile)
		}

		// Lösche abschließend das Backup.
		err := os.RemoveAll(assetsDir + "rollback/backup")
		HandleError(err)
	}
}

// Setze den Webserver zurück, indem alle Tickets, Mails und Nutzerdaten gelöscht werden.
// Dabei wird zunächst jeweils geprüft, ob das Objekt existiert.
func ResetData() {
	assetsDir := GetAssetsDir()

	if _, err := os.Stat(assetsDir + "tickets"); os.IsNotExist(err) == false {
		err := os.RemoveAll(assetsDir + "tickets")
		HandleError(err)
	}

	if _, err := os.Stat(assetsDir + "mails"); os.IsNotExist(err) == false {
		err := os.RemoveAll(assetsDir + "mails")
		HandleError(err)
	}

	if _, err := os.Stat(assetsDir + "users.json"); os.IsNotExist(err) == false {
		err = os.Remove(assetsDir + "users.json")
		HandleError(err)
	}

	CheckEnvironment()
}

//typeOfFiles entweder ticket oder mails
func createDemoFiles(typeOfFile string) {
	assetsDir := GetAssetsDir()
	rollbackDemoPath := "rollback/demo/" + typeOfFile
	// Erstelle Verzeichnis für Tickets/Mails, falls dieses nicht existiert.
	if _, err := os.Stat(assetsDir + typeOfFile); os.IsNotExist(err) {
		err = os.Mkdir(assetsDir+typeOfFile, 0744)
		HandleError(err)
	}

	// Kopiere alle Tickets/Mails aus dem Demo-Ordner in den Zielordner.
	src := strings.Join([]string{assetsDir, rollbackDemoPath}, "")
	files, err := ioutil.ReadDir(src)
	HandleError(err)

	for _, file := range files {
		srcFile := strings.Join([]string{assetsDir, rollbackDemoPath + "/", file.Name()}, "")
		dstFile := strings.Join([]string{assetsDir, typeOfFile + "/", file.Name()}, "")
		copyFile(srcFile, dstFile)
	}
}

// Setze den Webserver zurück und installiere Testdaten.
func DemoMode() {
	assetsDir := GetAssetsDir()
	ResetData()

	// Kopiere alle Tickets aus dem Demo-Ordner in den Zielordner.
	createDemoFiles("tickets")

	// Kopiere alle Mails aus dem Demo-Ordner in den Zielordner.
	createDemoFiles("mails")

	// Kopiere Nutzerdaten aus dem Demo-Ordner in den Zielordner.
	srcFile := strings.Join([]string{assetsDir, "rollback/demo/users.json"}, "")
	dstFile := strings.Join([]string{assetsDir, "users.json"}, "")
	copyFile(srcFile, dstFile)
}
