package tui

import (
	"context"
	"fmt"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"torque/internal/ory"
)

type listEntry struct {
	title   string
	details []string
}

type listResultMsg struct {
	title   string
	entries []listEntry
}

func fetchList(selection int, cfg Config) tea.Cmd {
	return func() tea.Msg {
		switch selection {
		case 0:
			return fetchUsers(cfg.KratosAdminURL)
		case 1:
			return fetchDealerships(cfg.KratosAdminURL, cfg.KetoReadURL)
		case 2:
			return fetchRoles(cfg.KetoReadURL)
		}
		return listResultMsg{}
	}
}

func fetchUsers(kratosAdminURL string) listResultMsg {
	kratos := ory.NewKratosClient(kratosAdminURL)
	identities, _, err := kratos.IdentityAPI.ListIdentities(context.Background()).Execute()
	if err != nil {
		return listResultMsg{title: "Users", entries: []listEntry{{title: "error: " + err.Error()}}}
	}

	var entries []listEntry
	for _, identity := range identities {
		traits, _ := identity.Traits.(map[string]any)
		email, _ := traits["email"].(string)
		name, _ := traits["name"].(map[string]any)
		first, _ := name["first"].(string)
		last, _ := name["last"].(string)
		entries = append(entries, listEntry{
			title:   email,
			details: []string{first + " " + last},
		})
	}
	return listResultMsg{title: "Users", entries: entries}
}

func fetchDealerships(kratosAdminURL, ketoReadURL string) listResultMsg {
	keto := ory.NewKetoReadClient(ketoReadURL)
	kratos := ory.NewKratosClient(kratosAdminURL)
	ctx := context.Background()

	result, _, err := keto.RelationshipApi.GetRelationships(ctx).
		Namespace(ory.DealershipsNS).Execute()
	if err != nil || result == nil {
		return listResultMsg{title: "Dealerships", entries: []listEntry{{title: "no data"}}}
	}

	// group by dealership
	grouped := map[string][]string{}
	for _, t := range result.RelationTuples {
		identity, _, err := kratos.IdentityAPI.GetIdentity(ctx, *t.SubjectId).Execute()
		email := *t.SubjectId
		if err == nil {
			traits, _ := identity.Traits.(map[string]any)
			email, _ = traits["email"].(string)
		}
		grouped[t.Object] = append(grouped[t.Object], email+" → "+t.Relation)
	}

	var dealerships []string
	for d := range grouped {
		dealerships = append(dealerships, d)
	}
	sort.Strings(dealerships)

	var entries []listEntry
	for _, d := range dealerships {
		entries = append(entries, listEntry{title: d, details: grouped[d]})
	}
	return listResultMsg{title: "Dealerships", entries: entries}
}

func fetchRoles(ketoReadURL string) listResultMsg {
	keto := ory.NewKetoReadClient(ketoReadURL)

	result, _, err := keto.RelationshipApi.GetRelationships(context.Background()).
		Namespace(ory.PermissionsNS).Execute()
	if err != nil || result == nil {
		return listResultMsg{title: "Roles", entries: []listEntry{{title: "no data"}}}
	}

	grouped := map[string][]string{}
	for _, t := range result.RelationTuples {
		if t.SubjectId != nil {
			grouped[*t.SubjectId] = append(grouped[*t.SubjectId], t.Object)
		}
	}

	var roles []string
	for r := range grouped {
		roles = append(roles, r)
	}
	sort.Strings(roles)

	var entries []listEntry
	for _, r := range roles {
		perms := grouped[r]
		sort.Strings(perms)
		entries = append(entries, listEntry{title: r, details: perms})
	}
	return listResultMsg{title: fmt.Sprintf("Roles (%d)", len(entries)), entries: entries}
}
