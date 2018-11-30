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

	err = ioutil.WriteFile(dst, data, 0744)
	HandleError(err)
}

// Erstelle ein Verzeichnis für Tickets, sofern dieses nicht existiert.
func CheckEnvironment() {
	assetsDir := GetAssetsDir()

	if _, err := os.Stat(assetsDir + "tickets"); os.IsNotExist(err) {
		err = os.Mkdir(assetsDir+"tickets", 0744)
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
	if _, err := os.Stat(assetsDir + "tickets"); os.IsNotExist(err) == false {
		src := strings.Join([]string{assetsDir, "tickets"}, "")
		files, err := ioutil.ReadDir(src)
		HandleError(err)

		// Wenn Tickets vorhanden sind, werden alle gesichert.
		if len(files) > 0 {
			err = os.Mkdir(assetsDir+"rollback/backup/tickets", 0744)
			HandleError(err)

			for _, file := range files {
				srcFile := strings.Join([]string{assetsDir, "tickets/", file.Name()}, "")
				dstFile := strings.Join([]string{assetsDir, "rollback/backup/tickets/", file.Name()}, "")
				copyFile(srcFile, dstFile)
			}
		}
	}

	// Sichere die Nutzerdaten, sofern diese existieren.
	if _, err := os.Stat(assetsDir + "users.json"); os.IsNotExist(err) == false {
		srcFile := strings.Join([]string{assetsDir, "users.json"}, "")
		dstFile := strings.Join([]string{assetsDir, "rollback/backup/users.json"}, "")
		copyFile(srcFile, dstFile)
	}

	// Sichere die Daten für Ticket-IDs, sofern diese existieren.
	if _, err := os.Stat(assetsDir + "ticketId_resource.json"); os.IsNotExist(err) == false {
		srcFile := strings.Join([]string{assetsDir, "ticketId_resource.json"}, "")
		dstFile := strings.Join([]string{assetsDir, "rollback/backup/ticketId_resource.json"}, "")
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
		if _, err := os.Stat(assetsDir + "rollback/backup/tickets"); os.IsNotExist(err) == false {
			src := strings.Join([]string{assetsDir, "rollback/backup/tickets"}, "")
			files, err := ioutil.ReadDir(src)
			HandleError(err)

			// Wenn Tickets vorhanden sind, werden alle geladen.
			if len(files) > 0 {
				err = os.Mkdir(assetsDir+"tickets", 0700)
				HandleError(err)

				for _, file := range files {
					srcFile := strings.Join([]string{assetsDir, "rollback/backup/tickets/", file.Name()}, "")
					dstFile := strings.Join([]string{assetsDir, "tickets/", file.Name()}, "")
					copyFile(srcFile, dstFile)
				}
			}
		}

		// Lade die Nutzerdaten, sofern diese existieren.
		if _, err := os.Stat(assetsDir + "rollback/backup/users.json"); os.IsNotExist(err) == false {
			srcFile := strings.Join([]string{assetsDir, "rollback/backup/users.json"}, "")
			dstFile := strings.Join([]string{assetsDir, "users.json"}, "")
			copyFile(srcFile, dstFile)
		}

		// Lade die Daten für Ticket-IDs, sofern diese existieren.
		if _, err := os.Stat(assetsDir + "rollback/backup/ticketId_resource.json"); os.IsNotExist(err) == false {
			srcFile := strings.Join([]string{assetsDir, "rollback/backup/ticketId_resource.json"}, "")
			dstFile := strings.Join([]string{assetsDir, "ticketId_resource.json"}, "")
			copyFile(srcFile, dstFile)
		}

		// Lösche abschließend das Backup.
		err := os.RemoveAll(assetsDir + "rollback/backup")
		HandleError(err)
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
		err = os.Mkdir(assetsDir+"tickets", 0744)
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
