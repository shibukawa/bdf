package docflow

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf/print/ipp"
)

// An app sends with credentials of its own: a token of the client
// credentials grant, or Basic over IPP, to the endpoints that list it,
// within its scopes.
func TestApps(t *testing.T) {
	e := newEnv(t, nil)
	ctx := context.Background()
	e.display("wh-a", Node{})
	// registered with its secret shown once
	e.admin("PUT", "/api/apps/bad", App{Scopes: []string{"fly"}}, http.StatusBadRequest)
	var made struct {
		App
		Secret string `json:"secret"`
	}
	json.Unmarshal(e.admin("PUT", "/api/apps/erp", App{Title: "ERP", Scopes: []string{ScopePrint, ScopePrintTemplate}}, http.StatusOK), &made)
	if !strings.HasPrefix(made.Secret, "bdfa_") || made.SecretHash != "" || !slices.Equal(made.Scopes, []string{ScopePrint, ScopePrintTemplate}) || made.Created.IsZero() {
		t.Fatalf("app %+v", made)
	}
	var again struct {
		Secret string `json:"secret"`
	}
	json.Unmarshal(e.admin("PUT", "/api/apps/erp", App{Title: "ERP 2", Scopes: []string{ScopePrint, ScopePrintTemplate}}, http.StatusOK), &again)
	if again.Secret != "" {
		t.Fatal("the secret was shown twice")
	}
	var apps []App
	json.Unmarshal(e.admin("GET", "/api/apps", nil, http.StatusOK), &apps)
	if len(apps) != 1 || apps[0].Title != "ERP 2" || apps[0].SecretHash != "" {
		t.Fatalf("apps %+v", apps)
	}

	// a token for the credentials: Basic, or in the form; not for others
	token := func(id, secret, form string) (int, map[string]any) {
		t.Helper()
		res, out := e.doWith(func(r *http.Request) {
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			if id != "" {
				r.SetBasicAuth(id, secret)
			}
		}, "POST", OAuthTokenPath, []byte(form))
		var body map[string]any
		json.Unmarshal(out, &body)
		return res.StatusCode, body
	}
	code, body := token("erp", made.Secret, "grant_type=client_credentials")
	tok, _ := body["access_token"].(string)
	if code != http.StatusOK || tok == "" || body["token_type"] != "Bearer" || body["scope"] != "print print.template" || body["expires_in"].(float64) != 600 {
		t.Fatalf("token: %d %+v", code, body)
	}
	if code, body := token("", "", "grant_type=client_credentials&client_id=erp&client_secret="+made.Secret+"&scope=print"); code != http.StatusOK || body["scope"] != "print" {
		t.Fatalf("token in the form: %d %+v", code, body)
	}
	if code, body := token("erp", "wrong", "grant_type=client_credentials"); code != http.StatusUnauthorized || body["error"] != "invalid_client" {
		t.Fatalf("a wrong secret: %d %+v", code, body)
	}
	if code, body := token("nobody", "x", "grant_type=client_credentials"); code != http.StatusUnauthorized || body["error"] != "invalid_client" {
		t.Fatalf("an unknown client: %d %+v", code, body)
	}
	if code, body := token("erp", made.Secret, "grant_type=password"); code != http.StatusBadRequest || body["error"] != "unsupported_grant_type" {
		t.Fatalf("another grant: %d %+v", code, body)
	}
	if code, body := token("erp", made.Secret, "grant_type=client_credentials&scope=admin"); code != http.StatusBadRequest || body["error"] != "invalid_scope" {
		t.Fatalf("a scope not the app's: %d %+v", code, body)
	}
	var me struct {
		Kind, ID string
		Admin    bool
	}
	_, out := e.do(nil, "GET", "/api/me", tok, nil)
	json.Unmarshal(out, &me)
	if me.Kind != PrincipalApp || me.ID != "erp" || me.Admin {
		t.Fatalf("me %+v", me)
	}

	// an endpoint that lists the app takes its documents; one that does
	// not refuses; so does a display
	e.node("intake", Node{Kind: KindEndpoint, To: []string{"wh-a"}, Apps: []string{"nobody"}}, http.StatusBadRequest)
	e.node("intake", Node{Kind: KindEndpoint, To: []string{"wh-a"}, Apps: []string{"erp"}}, http.StatusOK)
	e.node("other", Node{Kind: KindEndpoint, To: []string{"wh-a"}}, http.StatusOK)
	doc := testDoc(t, 1)
	res, out := e.do(nil, "POST", "/in/intake?name=from-erp", tok, doc)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("the app at its endpoint: %d %s", res.StatusCode, out)
	}
	var j Job
	json.Unmarshal(out, &j)
	if j.SubmittedBy == nil || *j.SubmittedBy != (Principal{Kind: PrincipalApp, ID: "erp", Name: "ERP 2", Via: ViaClientCredentials}) {
		t.Fatalf("submitted by %+v", j.SubmittedBy)
	}
	if res, _ := e.do(nil, "POST", "/in/other?name=x", tok, doc); res.StatusCode != http.StatusForbidden {
		t.Fatalf("the app at another endpoint: %d", res.StatusCode)
	}
	if res, _ := e.do(nil, "POST", "/api/destinations/wh-a/jobs?name=x", tok, doc); res.StatusCode != http.StatusForbidden {
		t.Fatalf("the app at a display: %d", res.StatusCode)
	}
	if res, _ := e.do(nil, "POST", "/api/destinations/intake/jobs?name=x", tok, doc); res.StatusCode != http.StatusCreated {
		t.Fatalf("the app at its endpoint through the API: %d", res.StatusCode)
	}
	var ds []Destination
	_, out = e.do(nil, "GET", "/api/destinations", tok, nil)
	json.Unmarshal(out, &ds)
	if len(ds) != 1 || ds[0].Name != "intake" {
		t.Fatalf("the app's destinations %+v", ds)
	}
	// over IPP, with Basic client_id:secret
	c := &ipp.Client{URL: e.srv.URL + IPPPath + "/intake", User: "erp", Username: "erp", Password: made.Secret}
	if js, err := c.PrintJob(ctx, "ipp", "application/x-bdf", nil, bytes.NewReader(doc)); err != nil || js.State != ipp.JobCompleted {
		t.Fatalf("ipp as the app: %+v %v", js, err)
	}
	c.URL = e.srv.URL + IPPPath + "/wh-a"
	if _, err := c.PrintJob(ctx, "ipp", "application/x-bdf", nil, bytes.NewReader(doc)); !ipp.IsStatus(err, ipp.StatusClientErrorNotAuthorized) {
		t.Fatalf("ipp as the app to a display: %v", err)
	}
	c.Password = "wrong"
	if _, err := c.PrintJob(ctx, "ipp", "application/x-bdf", nil, bytes.NewReader(doc)); err == nil {
		t.Fatal("ipp with a wrong secret")
	}

	// it sees its own jobs, cancels them with the scope only
	var jobs []Job
	_, out = e.do(nil, "GET", "/api/my/jobs", tok, nil)
	json.Unmarshal(out, &jobs)
	if len(jobs) != 3 {
		t.Fatalf("the app's jobs %+v", jobs)
	}
	if res, _ := e.do(nil, "POST", "/api/jobs/"+j.ID+"/cancel", tok, nil); res.StatusCode != http.StatusForbidden {
		t.Fatalf("cancel without the scope: %d", res.StatusCode)
	}
	e.admin("PUT", "/api/apps/erp", App{Scopes: []string{ScopePrint, ScopeJobsCancel}, Group: "warehouse-a"}, http.StatusOK)
	_, body = token("erp", made.Secret, "grant_type=client_credentials")
	tok2 := body["access_token"].(string)
	if res, _ := e.do(nil, "POST", "/api/jobs/"+j.ID+"/cancel", tok2, nil); res.StatusCode != http.StatusOK {
		t.Fatalf("cancel with the scope: %d", res.StatusCode)
	}
	// the token granted print.template no longer has it; sending data (for
	// a template node) needs it
	e.node("forms", Node{Kind: KindTemplate, Template: "slip", To: []string{"wh-a"}}, http.StatusOK)
	e.node("data-in", Node{Kind: KindEndpoint, Accept: AcceptData, To: []string{"forms"}, Apps: []string{"erp"}}, http.StatusOK)
	if res, _ := e.do(nil, "POST", "/in/data-in", tok2, map[string]any{"x": 1}); res.StatusCode != http.StatusForbidden {
		t.Fatalf("data without the scope: %d", res.StatusCode)
	}
	e.admin("PUT", "/api/apps/erp", App{Scopes: []string{ScopePrint, ScopePrintTemplate, ScopeJobsCancel}, Group: "warehouse-a"}, http.StatusOK)
	_, body = token("erp", made.Secret, "grant_type=client_credentials")
	tok3 := body["access_token"].(string)
	if res, out := e.do(nil, "POST", "/in/data-in", tok3, map[string]any{"x": 1}); res.StatusCode != http.StatusCreated {
		t.Fatalf("data with the scope: %d %s", res.StatusCode, out)
	}
	e.admin("DELETE", "/api/nodes/data-in", nil, http.StatusNoContent)
	e.admin("DELETE", "/api/nodes/forms", nil, http.StatusNoContent)
	// its jobs are its group's now
	res, out = e.do(nil, "POST", "/in/intake?name=grouped", tok2, doc)
	json.Unmarshal(out, &j)
	if res.StatusCode != http.StatusCreated || j.Group != "warehouse-a" {
		t.Fatalf("the app's group: %d %+v", res.StatusCode, j)
	}

	// the scopes tick and admin
	var cron struct {
		Secret string `json:"secret"`
	}
	json.Unmarshal(e.admin("PUT", "/api/apps/cron", App{Scopes: []string{ScopeTick}}, http.StatusOK), &cron)
	_, body = token("cron", cron.Secret, "grant_type=client_credentials")
	cronTok := body["access_token"].(string)
	if res, _ := e.do(nil, "POST", "/api/tick", cronTok, nil); res.StatusCode != http.StatusOK {
		t.Fatalf("tick as cron: %d", res.StatusCode)
	}
	if res, _ := e.do(nil, "POST", "/api/tick", tok2, nil); res.StatusCode != http.StatusForbidden {
		t.Fatalf("tick as erp: %d", res.StatusCode)
	}
	if res, _ := e.do(nil, "GET", "/api/nodes", cronTok, nil); res.StatusCode != http.StatusForbidden {
		t.Fatalf("the API as cron: %d", res.StatusCode)
	}
	var ops struct {
		Secret string `json:"secret"`
	}
	json.Unmarshal(e.admin("PUT", "/api/apps/ops", App{Scopes: []string{ScopeAdmin}}, http.StatusOK), &ops)
	_, body = token("ops", ops.Secret, "grant_type=client_credentials")
	if res, _ := e.do(nil, "GET", "/api/nodes", body["access_token"].(string), nil); res.StatusCode != http.StatusOK {
		t.Fatalf("the API as ops: %d", res.StatusCode)
	}

	// disabled: no token, and the tokens out stop; rotated: the old secret stops
	e.admin("PUT", "/api/apps/erp", App{Scopes: []string{ScopePrint}, Disabled: true}, http.StatusOK)
	if code, _ := token("erp", made.Secret, "grant_type=client_credentials"); code != http.StatusUnauthorized {
		t.Fatalf("a token of a disabled app: %d", code)
	}
	if res, _ := e.do(nil, "GET", "/api/me", tok2, nil); res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("a disabled app's token: %d", res.StatusCode)
	}
	e.admin("PUT", "/api/apps/erp", App{Scopes: []string{ScopePrint}}, http.StatusOK)
	var rot map[string]string
	json.Unmarshal(e.admin("POST", "/api/apps/erp/secret", nil, http.StatusOK), &rot)
	if code, _ := token("erp", made.Secret, "grant_type=client_credentials"); code != http.StatusUnauthorized {
		t.Fatalf("the old secret: %d", code)
	}
	if code, _ := token("erp", rot["secret"], "grant_type=client_credentials"); code != http.StatusOK {
		t.Fatalf("the new secret: %d", code)
	}
	// a token ends after AppTokenFor
	e.clock.add(AppTokenFor + time.Second)
	if res, _ := e.do(nil, "GET", "/api/me", cronTok, nil); res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("an expired token: %d", res.StatusCode)
	}

	// an app an endpoint lists is not deleted; the audit log has it all
	e.admin("DELETE", "/api/apps/erp", nil, http.StatusConflict)
	e.node("intake", Node{Kind: KindEndpoint, To: []string{"wh-a"}}, http.StatusOK)
	e.admin("DELETE", "/api/apps/erp", nil, http.StatusNoContent)
	e.admin("GET", "/api/apps/erp", nil, http.StatusNotFound)
	var log []AuditEntry
	json.Unmarshal(e.admin("GET", "/api/audit?limit=100", nil, http.StatusOK), &log)
	have := map[string]int{}
	for _, en := range log {
		have[en.Action]++
	}
	if have["put-app"] < 4 || have["rotate-app-secret"] != 1 || have["delete-app"] != 1 || have["token-refused"] < 3 {
		t.Fatalf("audit %v", have)
	}
}
