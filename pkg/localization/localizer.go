package localization

import (
	"borrow_book/internal/config"
	"borrow_book/pkg/logger"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

// Localizer manages localized messages.
type Localizer struct {
	messages map[string]interface{}
	mu       sync.RWMutex
	log      *logger.Logger
	path     string
	language string
}

// NewLocalizer creates a new Localizer instance.
func NewLocalizer(cfg *config.Config, log *logger.Logger) (*Localizer, error) {
	lang := cfg.Server.Language
	if lang == "" {
		lang = "en" // default language
	}

	i18nPath := cfg.Server.I18NPath
	if i18nPath == "" {
		return nil, fmt.Errorf("I18N_PATH is not set in the configuration")
	}

	// Ensure the path ends with a slash
	if !strings.HasSuffix(i18nPath, "/") {
		i18nPath += "/"
	}

	localizer := &Localizer{
		messages: make(map[string]interface{}),
		log:      log,
		path:     i18nPath,
		language: lang,
	}

	// Load the default language messages
	if err := localizer.loadYAML(lang); err != nil {
		return nil, err
	}

	return localizer, nil
}

// loadYAML loads the YAML file for the specified language.
func (l *Localizer) loadYAML(lang string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	filePath := fmt.Sprintf("%s%s.yaml", l.path, lang)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		l.log.Errorf("Error reading YAML file '%s': %v", filePath, err)
		return err
	}

	var yamlData map[string]interface{}
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		l.log.Errorf("Error unmarshaling YAML data from '%s': %v", filePath, err)
		return err
	}

	l.messages = yamlData
	l.language = lang
	l.log.Infof("Loaded messages for language '%s' from '%s'", lang, filePath)
	return nil
}

// SetLanguage switches the localization to the specified language.
func (l *Localizer) SetLanguage(lang string) error {
	return l.loadYAML(lang)
}

// Message retrieves the localized message using the key defined in LocalizedString.
func (l *Localizer) Message(key string) string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	keys := strings.Split(key, ".")
	var result interface{} = l.messages

	for _, k := range keys {
		switch value := result.(type) {
		case map[interface{}]interface{}:
			result = value[k]
		case map[string]interface{}:
			result = value[k]
		default:
			l.log.Errorf("Key type mismatch or not found: %s", k)
			return "Message not found"
		}

		if result == nil {
			l.log.Errorf("Key not found: %s", k)
			return "Message not found"
		}
	}

	if msg, ok := result.(string); ok {
		return msg
	}

	l.log.Errorf("Message for key '%s' is not a string", key)
	return "Message not found"
}
