package tui

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	ketoclient "github.com/ory/keto-client-go"
	kratosclient "github.com/ory/kratos-client-go"
	"torque/internal/ory"
)

var rolePermissions = map[string][]string{
	"admin": {
		"conversation:read", "conversation:send_message", "conversation:pause",
		"vehicle:read", "vehicle:write", "vehicle:update", "vehicle:delete",
		"customer:read", "customer:write", "customer:update", "customer:delete",
	},
	"support": {
		"conversation:read", "conversation:send_message", "conversation:pause",
		"vehicle:read",
		"customer:read",
	},
	"mechanical": {
		"conversation:read",
		"vehicle:read", "vehicle:write", "vehicle:update",
		"customer:read",
	},
}

type seedUserData struct {
	Email     string
	FirstName string
	LastName  string
	Roles     []struct {
		Dealership string
		Role       string
	}
}

var seedUsersData = []seedUserData{
	{
		Email: "admin@torque.com", FirstName: "Admin", LastName: "Torque",
		Roles: []struct {
			Dealership string
			Role       string
		}{{Dealership: "sp-001", Role: "admin"}},
	},
	{
		Email: "support@torque.com", FirstName: "Support", LastName: "Torque",
		Roles: []struct {
			Dealership string
			Role       string
		}{{Dealership: "sp-001", Role: "support"}},
	},
	{
		Email: "mechanical@torque.com", FirstName: "Mechanical", LastName: "Torque",
		Roles: []struct {
			Dealership string
			Role       string
		}{{Dealership: "sp-001", Role: "mechanical"}},
	},
}

func startSeed(selection int, ch chan<- tea.Msg, cfg Config) tea.Cmd {
	return func() tea.Msg {
		emit := func(section, item string, ok bool) {
			ch <- progressMsg{section: section, item: item, ok: ok}
		}

		keto := ory.NewKetoWriteClient(cfg.KetoWriteURL)

		switch selection {
		case 0:
			execSeedPermissions(keto, emit)
		case 1:
			execSeedUsers(ory.NewKratosClient(cfg.KratosAdminURL), keto, emit)
		}

		ch <- runDoneMsg{}
		return nil
	}
}

func execSeedPermissions(keto *ketoclient.APIClient, emit func(string, string, bool)) {
	ctx := context.Background()

	_, err := keto.RelationshipApi.DeleteRelationships(ctx).
		Namespace(ory.PermissionsNS).Execute()
	emit("Clear permissions", ory.PermissionsNS, err == nil)

	for role, perms := range rolePermissions {
		for _, perm := range perms {
			_, _, err := keto.RelationshipApi.CreateRelationship(ctx).
				CreateRelationshipBody(ketoclient.CreateRelationshipBody{
					Namespace: ketoclient.PtrString(ory.PermissionsNS),
					Object:    ketoclient.PtrString(perm),
					Relation:  ketoclient.PtrString("allowed"),
					SubjectId: ketoclient.PtrString(role),
				}).Execute()
			emit("Seed permissions", role+" → "+perm, err == nil)
		}
	}
}

func execSeedUsers(kratos *kratosclient.APIClient, keto *ketoclient.APIClient, emit func(string, string, bool)) {
	ctx := context.Background()
	var created []struct {
		id    string
		email string
		roles []struct{ Dealership, Role string }
	}

	// clear
	identities, _, err := kratos.IdentityAPI.ListIdentities(ctx).Execute()
	if err == nil && len(identities) > 0 {
		cleared := 0
		for _, identity := range identities {
			_, e := kratos.IdentityAPI.DeleteIdentity(ctx, identity.Id).Execute()
			if e == nil {
				cleared++
			}
		}
		emit("Clear", fmt.Sprintf("%d user(s) removed", cleared), true)
	}

	_, errD := keto.RelationshipApi.DeleteRelationships(ctx).
		Namespace(ory.DealershipsNS).Execute()
	emit("Clear", ory.DealershipsNS, errD == nil)

	// create
	for _, u := range seedUsersData {
		identity, _, err := kratos.IdentityAPI.CreateIdentity(ctx).
			CreateIdentityBody(kratosclient.CreateIdentityBody{
				SchemaId: "default",
				Traits: map[string]any{
					"email": u.Email,
					"name":  map[string]any{"first": u.FirstName, "last": u.LastName},
				},
			}).Execute()
		if err != nil {
			emit("Create users", u.Email, false)
			continue
		}
		emit("Create users", u.Email, true)

		var roles []struct{ Dealership, Role string }
		for _, r := range u.Roles {
			roles = append(roles, struct{ Dealership, Role string }{r.Dealership, r.Role})
		}
		created = append(created, struct {
			id    string
			email string
			roles []struct{ Dealership, Role string }
		}{identity.Id, u.Email, roles})
	}

	// assign roles
	for _, u := range created {
		for _, r := range u.roles {
			_, _, err := keto.RelationshipApi.CreateRelationship(ctx).
				CreateRelationshipBody(ketoclient.CreateRelationshipBody{
					Namespace: ketoclient.PtrString(ory.DealershipsNS),
					Object:    ketoclient.PtrString(r.Dealership),
					Relation:  ketoclient.PtrString(r.Role),
					SubjectId: ketoclient.PtrString(u.id),
				}).Execute()
			emit("Assign roles", u.email+" → "+r.Role+" @ "+r.Dealership, err == nil)
		}
	}
}
