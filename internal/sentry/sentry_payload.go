package sentry

import (
	"encoding/json"
	"fmt"
	"time"
)

type WebhookPayload struct {
	Event  Event
	Action EventAction `json:"action"`
	Actor  struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"actor"`
	Installation struct {
		Uuid string `json:"uuid"`
	} `json:"installation"`
	Data interface{}
}

func (w *WebhookPayload) UnmarshalJSON(data []byte) error {
	// Correctly set the Data interface type
	switch w.Event {
	case InstallationEvent:
	case UninstallationEvent:
		w.Data = new(InstallationData)
	case IssueAlertEvent:
		w.Data = new(IssueAlertData)
	case MetricAlertEvent:
		w.Data = new(MetricAlertData)
	case IssueEvent:
		w.Data = new(IssueData)
	case ErrorEvent:
		w.Data = new(ErrorData)
	default:
		return fmt.Errorf("unknown event %s", w.Event)
	}

	type tmp WebhookPayload // avoids infinite recursion
	return json.Unmarshal(data, (*tmp)(w))
}

type InstallationData struct {
	Installation struct {
		Status       string `json:"status"`
		Organization struct {
			Slug string `json:"slug"`
		} `json:"organization"`
		App struct {
			Uuid string `json:"uuid"`
			Slug string `json:"slug"`
		} `json:"app"`
		Code string `json:"code"`
		Uuid string `json:"uuid"`
	} `json:"installation"`
}

type IssueAlertData struct {
	Event struct {
		Ref        int `json:"_ref"`
		RefVersion int `json:"_ref_version"`
		Contexts   struct {
			Browser struct {
				Name    string `json:"name"`
				Type    string `json:"type"`
				Version string `json:"version"`
			} `json:"browser"`
			Os struct {
				Name    string `json:"name"`
				Type    string `json:"type"`
				Version string `json:"version"`
			} `json:"os"`
		} `json:"contexts"`
		Culprit   string      `json:"culprit"`
		Datetime  time.Time   `json:"datetime"`
		Dist      interface{} `json:"dist"`
		EventId   string      `json:"event_id"`
		Exception struct {
			Values []struct {
				Mechanism struct {
					Data struct {
						Message string `json:"message"`
						Mode    string `json:"mode"`
						Name    string `json:"name"`
					} `json:"data"`
					Description interface{} `json:"description"`
					Handled     bool        `json:"handled"`
					HelpLink    interface{} `json:"help_link"`
					Meta        interface{} `json:"meta"`
					Synthetic   interface{} `json:"synthetic"`
					Type        string      `json:"type"`
				} `json:"mechanism"`
				Stacktrace struct {
					Frames []struct {
						AbsPath     string  `json:"abs_path"`
						Colno       int     `json:"colno"`
						ContextLine *string `json:"context_line"`
						Data        struct {
							OrigInApp int `json:"orig_in_app"`
						} `json:"data"`
						Errors          interface{} `json:"errors"`
						Filename        string      `json:"filename"`
						Function        interface{} `json:"function"`
						ImageAddr       interface{} `json:"image_addr"`
						InApp           bool        `json:"in_app"`
						InstructionAddr interface{} `json:"instruction_addr"`
						Lineno          int         `json:"lineno"`
						Module          *string     `json:"module"`
						Package         interface{} `json:"package"`
						Platform        interface{} `json:"platform"`
						PostContext     interface{} `json:"post_context"`
						PreContext      interface{} `json:"pre_context"`
						RawFunction     interface{} `json:"raw_function"`
						Symbol          interface{} `json:"symbol"`
						SymbolAddr      interface{} `json:"symbol_addr"`
						Trust           interface{} `json:"trust"`
						Vars            interface{} `json:"vars"`
					} `json:"frames"`
				} `json:"stacktrace"`
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"values"`
		} `json:"exception"`
		Fingerprint    []string `json:"fingerprint"`
		GroupingConfig struct {
			Enhancements string `json:"enhancements"`
			Id           string `json:"id"`
		} `json:"grouping_config"`
		Hashes   []string `json:"hashes"`
		IssueUrl string   `json:"issue_url"`
		KeyId    string   `json:"key_id"`
		Level    string   `json:"level"`
		Location string   `json:"location"`
		Logger   string   `json:"logger"`
		Message  string   `json:"message"`
		Metadata struct {
			Filename string `json:"filename"`
			Type     string `json:"type"`
			Value    string `json:"value"`
		} `json:"metadata"`
		Platform string      `json:"platform"`
		Project  int         `json:"project"`
		Received float64     `json:"received"`
		Release  interface{} `json:"release"`
		Request  struct {
			Cookies             interface{}   `json:"cookies"`
			Data                interface{}   `json:"data"`
			Env                 interface{}   `json:"env"`
			Fragment            interface{}   `json:"fragment"`
			Headers             [][]string    `json:"headers"`
			InferredContentType interface{}   `json:"inferred_content_type"`
			Method              interface{}   `json:"method"`
			QueryString         []interface{} `json:"query_string"`
			Url                 string        `json:"url"`
		} `json:"request"`
		Sdk struct {
			Integrations []string `json:"integrations"`
			Name         string   `json:"name"`
			Packages     []struct {
				Name    string `json:"name"`
				Version string `json:"version"`
			} `json:"packages"`
			Version string `json:"version"`
		} `json:"sdk"`
		Tags      [][]string  `json:"tags"`
		TimeSpent interface{} `json:"time_spent"`
		Timestamp float64     `json:"timestamp"`
		Title     string      `json:"title"`
		Type      string      `json:"type"`
		Url       string      `json:"url"`
		User      struct {
			IpAddress string `json:"ip_address"`
		} `json:"user"`
		Version string `json:"version"`
		WebUrl  string `json:"web_url"`
	} `json:"event"`
	TriggeredRule string `json:"triggered_rule"`
}

type MetricAlertData struct {
	DescriptionText  string `json:"description_text"`
	DescriptionTitle string `json:"description_title"`
	MetricAlert      struct {
		AlertRule struct {
			Aggregate          string        `json:"aggregate"`
			CreatedBy          interface{}   `json:"created_by"`
			Dataset            string        `json:"dataset"`
			DateCreated        time.Time     `json:"date_created"`
			DateModified       time.Time     `json:"date_modified"`
			Environment        interface{}   `json:"environment"`
			Id                 string        `json:"id"`
			IncludeAllProjects bool          `json:"include_all_projects"`
			Name               string        `json:"name"`
			OrganizationId     string        `json:"organization_id"`
			Projects           []string      `json:"projects"`
			Query              string        `json:"query"`
			Resolution         int           `json:"resolution"`
			ResolveThreshold   interface{}   `json:"resolve_threshold"`
			Status             int           `json:"status"`
			ThresholdPeriod    int           `json:"threshold_period"`
			ThresholdType      int           `json:"threshold_type"`
			TimeWindow         int           `json:"time_window"`
			Triggers           []interface{} `json:"triggers"`
		} `json:"alert_rule"`
		DateClosed     interface{} `json:"date_closed"`
		DateCreated    time.Time   `json:"date_created"`
		DateDetected   time.Time   `json:"date_detected"`
		DateStarted    time.Time   `json:"date_started"`
		Id             string      `json:"id"`
		Identifier     string      `json:"identifier"`
		OrganizationId string      `json:"organization_id"`
		Projects       []string    `json:"projects"`
		Status         int         `json:"status"`
		StatusMethod   int         `json:"status_method"`
		Title          string      `json:"title"`
		Type           int         `json:"type"`
	} `json:"metric_alert"`
	WebUrl string `json:"web_url"`
}

type IssueData struct {
	Issue struct {
		Annotations  []interface{} `json:"annotations"`
		AssignedTo   interface{}   `json:"assignedTo"`
		Count        string        `json:"count"`
		Culprit      string        `json:"culprit"`
		FirstSeen    time.Time     `json:"firstSeen"`
		HasSeen      bool          `json:"hasSeen"`
		Id           string        `json:"id"`
		IsBookmarked bool          `json:"isBookmarked"`
		IsPublic     bool          `json:"isPublic"`
		IsSubscribed bool          `json:"isSubscribed"`
		LastSeen     time.Time     `json:"lastSeen"`
		Level        string        `json:"level"`
		Logger       interface{}   `json:"logger"`
		Metadata     struct {
			Filename string `json:"filename"`
			Type     string `json:"type"`
			Value    string `json:"value"`
		} `json:"metadata"`
		NumComments int    `json:"numComments"`
		Permalink   string `json:"permalink"`
		Platform    string `json:"platform"`
		Project     struct {
			Id       string `json:"id"`
			Name     string `json:"name"`
			Platform string `json:"platform"`
			Slug     string `json:"slug"`
		} `json:"project"`
		ShareId       interface{} `json:"shareId"`
		ShortId       string      `json:"shortId"`
		Status        string      `json:"status"`
		StatusDetails struct {
		} `json:"statusDetails"`
		SubscriptionDetails interface{} `json:"subscriptionDetails"`
		Title               string      `json:"title"`
		Type                string      `json:"type"`
		UserCount           int         `json:"userCount"`
	} `json:"issue"`
}

type ErrorData struct {
	Error struct {
		Ref        int `json:"_ref"`
		RefVersion int `json:"_ref_version"`
		Contexts   struct {
			Browser struct {
				Name    string `json:"name"`
				Type    string `json:"type"`
				Version string `json:"version"`
			} `json:"browser"`
			Os struct {
				Name    string `json:"name"`
				Type    string `json:"type"`
				Version string `json:"version"`
			} `json:"os"`
		} `json:"contexts"`
		Culprit   string      `json:"culprit"`
		Datetime  time.Time   `json:"datetime"`
		Dist      interface{} `json:"dist"`
		EventId   string      `json:"event_id"`
		Exception struct {
			Values []struct {
				Mechanism struct {
					Data struct {
						Message string `json:"message"`
						Mode    string `json:"mode"`
						Name    string `json:"name"`
					} `json:"data"`
					Description interface{} `json:"description"`
					Handled     bool        `json:"handled"`
					HelpLink    interface{} `json:"help_link"`
					Meta        interface{} `json:"meta"`
					Synthetic   interface{} `json:"synthetic"`
					Type        string      `json:"type"`
				} `json:"mechanism"`
				Stacktrace struct {
					Frames []struct {
						AbsPath     string `json:"abs_path"`
						Colno       int    `json:"colno"`
						ContextLine string `json:"context_line"`
						Data        struct {
							OrigInApp int `json:"orig_in_app"`
						} `json:"data"`
						Errors          interface{} `json:"errors"`
						Filename        string      `json:"filename"`
						Function        interface{} `json:"function"`
						ImageAddr       interface{} `json:"image_addr"`
						InApp           bool        `json:"in_app"`
						InstructionAddr interface{} `json:"instruction_addr"`
						Lineno          int         `json:"lineno"`
						Module          string      `json:"module"`
						Package         interface{} `json:"package"`
						Platform        interface{} `json:"platform"`
						PostContext     []string    `json:"post_context"`
						PreContext      []string    `json:"pre_context"`
						RawFunction     interface{} `json:"raw_function"`
						Symbol          interface{} `json:"symbol"`
						SymbolAddr      interface{} `json:"symbol_addr"`
						Trust           interface{} `json:"trust"`
						Vars            interface{} `json:"vars"`
					} `json:"frames"`
				} `json:"stacktrace"`
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"values"`
		} `json:"exception"`
		Fingerprint    []string `json:"fingerprint"`
		GroupingConfig struct {
			Enhancements string `json:"enhancements"`
			Id           string `json:"id"`
		} `json:"grouping_config"`
		Hashes   []string `json:"hashes"`
		IssueUrl string   `json:"issue_url"`
		KeyId    string   `json:"key_id"`
		Level    string   `json:"level"`
		Location string   `json:"location"`
		Logger   string   `json:"logger"`
		Message  string   `json:"message"`
		Metadata struct {
			Filename string `json:"filename"`
			Type     string `json:"type"`
			Value    string `json:"value"`
		} `json:"metadata"`
		Platform string      `json:"platform"`
		Project  int         `json:"project"`
		Received float64     `json:"received"`
		Release  interface{} `json:"release"`
		Request  struct {
			Cookies             interface{}   `json:"cookies"`
			Data                interface{}   `json:"data"`
			Env                 interface{}   `json:"env"`
			Fragment            interface{}   `json:"fragment"`
			Headers             [][]string    `json:"headers"`
			InferredContentType interface{}   `json:"inferred_content_type"`
			Method              interface{}   `json:"method"`
			QueryString         []interface{} `json:"query_string"`
			Url                 string        `json:"url"`
		} `json:"request"`
		Sdk struct {
			Integrations []string `json:"integrations"`
			Name         string   `json:"name"`
			Packages     []struct {
				Name    string `json:"name"`
				Version string `json:"version"`
			} `json:"packages"`
			Version string `json:"version"`
		} `json:"sdk"`
		Tags      [][]string  `json:"tags"`
		TimeSpent interface{} `json:"time_spent"`
		Timestamp float64     `json:"timestamp"`
		Title     string      `json:"title"`
		Type      string      `json:"type"`
		Url       string      `json:"url"`
		User      struct {
			IpAddress string `json:"ip_address"`
		} `json:"user"`
		Version string `json:"version"`
		WebUrl  string `json:"web_url"`
	} `json:"error"`
}
