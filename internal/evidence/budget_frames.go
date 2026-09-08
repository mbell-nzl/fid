package evidence

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/mbell-nzl/fid/internal/workflow"
)

type committedBudgetStop struct {
	stopID, attentionID, details string
}

// Only Stops committed by this execution's two Events can produce a control.
// Counters come from the immutable crossing fact, not a later derived total.
func appendCommittedBudgetStops(conn *sql.DB, live *LiveResult, runID, startedEventID, receiptEventID string) error {
	stops, err := loadCommittedBudgetStops(conn, startedEventID, receiptEventID)
	if err != nil || len(stops) == 0 {
		return err
	}
	chains, err := loadBudgetChains(conn, runID)
	if err != nil {
		return err
	}
	policies := make(map[string]workflow.Budget, len(chains))
	for _, chain := range chains {
		policies[chain.Budget.Name] = chain.Budget
	}
	controls := make([]map[string]any, 0, len(stops))
	for _, stop := range stops {
		var detail struct {
			Name      string `json:"budget_name"`
			Consumed  *int64 `json:"consumed"`
			Remaining *int64 `json:"remaining"`
		}
		if err := json.Unmarshal([]byte(stop.details), &detail); err != nil {
			return fmt.Errorf("budget control: decode crossing fact: %w", err)
		}
		budget, exists := policies[detail.Name]
		if !exists || detail.Consumed == nil || detail.Remaining == nil {
			return fmt.Errorf("budget control: crossing fact has no bound known budget")
		}
		budget.Consumed, budget.Remaining, budget.State = detail.Consumed, detail.Remaining, "exhausted"
		controls = append(controls, map[string]any{"kind": "budget-stop", "stop_ids": []string{stop.stopID}, "attention_ids": []string{stop.attentionID}, "budget": budget})
	}
	sort.SliceStable(controls, func(i, j int) bool {
		return controls[i]["budget"].(workflow.Budget).Name < controls[j]["budget"].(workflow.Budget).Name
	})
	for _, control := range controls {
		live.appendControl(control)
	}
	return nil
}

func loadCommittedBudgetStops(conn *sql.DB, startedEventID, receiptEventID string) ([]committedBudgetStop, error) {
	rows, err := conn.Query(`SELECT s.id, l.attention_id, s.details_json
FROM stops AS s JOIN attention_stop_links AS l ON l.stop_id = s.id
WHERE s.code = 'budget-exhausted' AND s.created_event_id IN (?, ?)
ORDER BY s.id, l.attention_id`, startedEventID, receiptEventID)
	if err != nil {
		return nil, fmt.Errorf("budget control: read committed Stops: %w", err)
	}
	defer func() { _ = rows.Close() }()
	stops := []committedBudgetStop{}
	for rows.Next() {
		var stop committedBudgetStop
		if err := rows.Scan(&stop.stopID, &stop.attentionID, &stop.details); err != nil {
			return nil, err
		}
		stops = append(stops, stop)
	}
	return stops, rows.Err()
}
