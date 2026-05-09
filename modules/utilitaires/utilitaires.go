/*
Copyright ou © ou Copr. Didier Hoizé, (20 avril 2026)

aquarium[@]mailo[.]com

Ce logiciel est un programme informatique servant à l'analyse des collectes
traçologiques effectuées avec le logiciel DFIR-ORC.

Ce logiciel est régi par la licence CeCILL soumise au droit français et
respectant les principes de diffusion des logiciels libres. Vous pouvez
utiliser, modifier et/ou redistribuer ce programme sous les conditions
de la licence CeCILL telle que diffusée par le CEA, le CNRS et l'INRIA
sur le site "http://www.cecill.info".

En contrepartie de l'accessibilité au code source et des droits de copie,
de modification et de redistribution accordés par cette licence, il n'est
offert aux utilisateurs qu'une garantie limitée.  Pour les mêmes raisons,
seule une responsabilité restreinte pèse sur l'auteur du programme,  le
titulaire des droits patrimoniaux et les concédants successifs.

A cet égard  l'attention de l'utilisateur est attirée sur les risques
associés au chargement,  à l'utilisation,  à la modification et/ou au
développement et à la reproduction du logiciel par l'utilisateur étant
donné sa spécificité de logiciel libre, qui peut le rendre complexe à
manipuler et qui le réserve donc à des développeurs et des professionnels
avertis possédant  des  connaissances  informatiques approfondies.  Les
utilisateurs sont donc invités à charger  et  tester  l'adéquation  du
logiciel à leurs besoins dans des conditions permettant d'assurer la
sécurité de leurs systèmes et ou de leurs données et, plus généralement,
à l'utiliser et l'exploiter dans les mêmes conditions de sécurité.

Le fait que vous puissiez accéder à cet en-tête signifie que vous avez
pris connaissance de la licence CeCILL, et que vous en avez accepté les
termes.
*/

package utilitaires

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func GetCheminBaseApplication() (string, error) {
	// Retourne le dossier contenant l'executable de l'application.
	emplacementExecutable, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(emplacementExecutable), nil
}

/***************************************************************************************/
/************************* LOGGER FUNCTIONS **********************************/
/***************************************************************************************/
const (
	dossierLogs    = "logs"
	prefixeFichier = "aquarium"
)

var (
	loggerMu     sync.RWMutex
	loggerActif  bool
	loggerGlobal *zap.Logger
)

func InitLogger(cheminBase string, debugActif bool) error {
	// Initialise le logger central une seule fois pour la session.
	// Si le debug est desactive, le logger devient inactif (no-op).
	loggerMu.Lock()
	defer loggerMu.Unlock()

	loggerActif = debugActif
	if !loggerActif {
		if loggerGlobal != nil {
			_ = loggerGlobal.Sync()
			loggerGlobal = nil
		}
		return nil
	}

	cheminLogs := filepath.Join(cheminBase, dossierLogs)
	if err := os.MkdirAll(cheminLogs, 0o755); err != nil {
		return fmt.Errorf("creation dossier logs: %w", err)
	}

	rotation := &lumberjack.Logger{
		Filename:   filepath.Join(cheminLogs, fmt.Sprintf("%s_%s.log", prefixeFichier, time.Now().Format("20060102_150405"))),
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}

	config := zap.NewProductionEncoderConfig()
	config.TimeKey = "ts"
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.MessageKey = "msg"
	config.LevelKey = "level"
	config.EncodeLevel = zapcore.LowercaseLevelEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		zapcore.AddSync(rotation),
		zapcore.DebugLevel,
	)

	if loggerGlobal != nil {
		_ = loggerGlobal.Sync()
	}
	loggerGlobal = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return nil
}

func LoggerActif() bool {
	// Indique si le logger est actif et pret a ecrire. .Helper de vérification d’état (notamment pour tests).
	loggerMu.RLock()
	defer loggerMu.RUnlock()
	return loggerActif && loggerGlobal != nil
}

func LogEvent(niveau string, evenement string, attributs map[string]interface{}, message string) {
	// Point d'entree unique de journalisation (niveau, evenement, attributs, message libre).
	// Si message est vide, l'evenement est utilise comme message.
	loggerMu.RLock()
	defer loggerMu.RUnlock()
	if !loggerActif || loggerGlobal == nil {
		return
	}

	champs := []zap.Field{
		zap.String("event", evenement),
	}
	if attributs != nil {
		champs = append(champs, zap.Any("attrs", attributs))
	}

	if strings.TrimSpace(message) == "" {
		message = evenement
	}

	niveau = strings.ToLower(strings.TrimSpace(niveau))
	switch niveau {
	case "debug":
		loggerGlobal.Debug(message, champs...)
	case "info":
		loggerGlobal.Info(message, champs...)
	case "warn":
		loggerGlobal.Warn(message, champs...)
	case "error":
		loggerGlobal.Error(message, champs...)
	default:
		loggerGlobal.Info(message, champs...)
	}
}

func SyncLogger() {
	// Force la vidange des buffers vers le fichier de log (utile a l'arret).
	loggerMu.RLock()
	defer loggerMu.RUnlock()
	if loggerGlobal != nil {
		_ = loggerGlobal.Sync()
	}
}
